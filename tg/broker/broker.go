package broker

import (
	"errors"
	"sync"
	"time"

	"github.com/MrReality255/turbo-go/tg/utils"
)

const (
	queueSize = 64

	HandleAny = 0
)

var (
	ErrRequestTimeout = errors.New("request timeout")
)

type Handle uint64

type ICommand any

/*
type ICommand interface {
	GetID() Handle
	GetRefID() Handle
}
*/

type RequestHandler[Cmd ICommand] func(cmd Cmd, err error) bool
type MemberMessageHandler[Cmd ICommand] func(sender Handle, msg Cmd, member IMember[Cmd])

/*
type IMember[Command ICommand] interface {
	GetID() Handle
	HandleMessage(sender Handle, msg Command)
}
*/

type IBroker[Cmd ICommand] interface {
	AddMember(id Handle, messageHandler MemberMessageHandler[Cmd]) IMember[Cmd]
	Stop()
}

type IMember[Command ICommand] interface {
	Request(receiver Handle, cmd Command) (Command, error)
	RequestMultiple(receiver Handle, cmd Command, handler RequestHandler[Command])
	Send(receiver Handle, cmd Command)
	Subscribe(cmdType ...uint32)
}

type messageRequest[Command ICommand] struct {
	mx          sync.Mutex
	handler     RequestHandler[Command]
	nextTimeout time.Time
	closed      bool
}

type memberWrapper[Command ICommand] struct {
	id             Handle
	messageHandler MemberMessageHandler[Command]
	broker         *controller[Command]
	requestTimeout time.Duration

	activeRequests map[Handle]*messageRequest[Command]
	mx             sync.Mutex
}

type messageWrapper[Command ICommand] struct {
	cmd      Command
	sender   Handle
	receiver Handle
}

type subscribersMap map[uint32]map[Handle]bool

type controller[Command ICommand] struct {
	requestTimeout time.Duration
	msgFct         func(cmd Command) Handle

	mx sync.Mutex

	chQueue     chan *messageWrapper[Command]
	closed      bool
	members     map[Handle]*memberWrapper[Command]
	subscribers subscribersMap
}

func New[Command ICommand](
	getMessageID func(cmd Command) Handle,
	timeout time.Duration,
) IBroker[Command] {
	c := &controller[Command]{
		msgFct:         getMessageID,
		requestTimeout: timeout,
		chQueue:        make(chan *messageWrapper[Command], queueSize),
		members:        make(map[Handle]*memberWrapper[Command]),
		subscribers:    make(subscribersMap),
	}
	go c.startLoop()
	return c
}

func (c *controller[Command]) AddMember(
	handle Handle, messageHandler MemberMessageHandler[Command],
) IMember[Command] {
	c.mx.Lock()
	defer c.mx.Unlock()
	wrapper := &memberWrapper[Command]{
		id:             handle,
		messageHandler: messageHandler,
		broker:         c,
		requestTimeout: c.requestTimeout,
		activeRequests: make(map[Handle]*messageRequest[Command]),
	}
	c.members[handle] = wrapper
	return wrapper
}

func (c *controller[Command]) Stop() {
	utils.ExecLocked(&c.mx, func() {
		c.closed = true
	})
	close(c.chQueue)
}

func (c *controller[Command]) startLoop() {
	for {
		msg, ok := <-c.chQueue
		if !ok {
			break
		}
		c.dispatchMessage(msg)
	}
}

func (c *controller[Command]) dispatchMessage(msg *messageWrapper[Command]) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// the message has an exact receiver: pass it to the receiver
	if msg.receiver.GetSeqID() != HandleAny {
		go c.members[msg.receiver].handleMessage(msg.sender, msg.cmd)
		return
	}

	// pass the message to everyone who is subscribed to the message type
	var (
		msgType     = msg.cmd.GetID().GetTypeID()
		handled     = make(map[Handle]bool)
		handledType = make(map[uint32]bool)
	)

	for _, subID := range []uint32{msgType, 0} {
		for subscriber, ok := range c.subscribers[subID] {
			if !ok || handled[subscriber] {
				continue
			}
			handled[subscriber] = true
			handledType[subscriber.GetTypeID()] = true
			go c.members[subscriber].handleMessage(msg.sender, msg.cmd)
		}
	}

	// the message has receiver type: send it to one receiver of this type
	if recTypeID := msg.receiver.GetTypeID(); recTypeID != 0 && !handledType[recTypeID] {
		for m, h := range c.members {
			if m.GetTypeID() == recTypeID {
				go h.handleMessage(msg.sender, msg.cmd)
				return
			}
		}
	}
}

func (c *controller[Command]) send(sender Handle, receiver Handle, cmd Command) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.closed {
		return
	}
	c.chQueue <- &messageWrapper[Command]{cmd: cmd, sender: sender, receiver: receiver}
}

func (c *controller[Command]) subscribe(subscriber Handle, cmdTypes ...uint32) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.closed {
		return
	}
	for _, cmdType := range cmdTypes {
		if c.subscribers[cmdType] == nil {
			c.subscribers[cmdType] = make(map[Handle]bool)
		}
		c.subscribers[cmdType][subscriber] = true
	}
}

func (m *memberWrapper[Command]) Request(receiver Handle, cmd Command) (Command, error) {
	chResponse := make(chan *utils.ItemWithErr[Command], 1)
	m.RequestMultiple(receiver, cmd, func(cmd Command, err error) bool {
		chResponse <- &utils.ItemWithErr[Command]{
			Data: cmd,
			Err:  err,
		}
		return true
	})
	result := <-chResponse
	return result.Data, result.Err
}

func (m *memberWrapper[Command]) RequestMultiple(
	receiver Handle, cmd Command, handler RequestHandler[Command],
) {
	reqHandle := cmd.GetID()
	utils.ExecLocked(&m.mx, func() {
		m.activeRequests[reqHandle] = &messageRequest[Command]{
			handler:     handler,
			nextTimeout: time.Now().Add(m.requestTimeout),
		}
	})
	m.Send(receiver, cmd)
	go m.checkRequestTimeout(reqHandle)
}

func (m *memberWrapper[Command]) Send(receiver Handle, cmd Command) {
	go m.broker.send(m.id, receiver, cmd)
}

func (m *memberWrapper[Command]) Subscribe(cmdType ...uint32) {
	m.broker.subscribe(m.id, cmdType...)
}

func (m *memberWrapper[Command]) handleMessage(sender Handle, cmd Command) {
	var (
		refID  = cmd.GetRefID()
		refReq = m.getRequest(refID)
	)
	if refReq != nil {
		if m.tryHandleResponse(refID, refReq, sender, cmd) {
			return
		}
	}
	m.messageHandler(sender, cmd, m)
}

func (m *memberWrapper[Command]) checkRequestTimeout(handle Handle) {
	time.Sleep(m.requestTimeout)
	req := m.getRequest(handle)
	if req == nil {
		return
	}
	req.mx.Lock()
	defer req.mx.Unlock()
	if req.nextTimeout.Before(time.Now()) {
		req.closed = true
		req.handler(nil, ErrRequestTimeout)
		go m.removeClosedRequest(handle)
	}
}

func (m *memberWrapper[Command]) getRequest(handle Handle) *messageRequest[Command] {
	m.mx.Lock()
	defer m.mx.Unlock()
	return m.activeRequests[handle]
}

func (m *memberWrapper[Command]) removeClosedRequest(handle Handle) {
	m.mx.Lock()
	defer m.mx.Unlock()
	req := m.activeRequests[handle]
	if req != nil && req.closed {
		delete(m.activeRequests, handle)
	}
}

func (m *memberWrapper[Command]) tryHandleResponse(
	handle Handle, req *messageRequest[Command], sender Handle, cmd Command,
) bool {
	req.mx.Lock()
	defer req.mx.Unlock()
	if req.closed {
		return false
	}
	req.nextTimeout = time.Now().Add(m.requestTimeout)
	go func() {
		isDone := req.handler(cmd, nil)
		if isDone {
			req.mx.Lock()
			defer req.mx.Unlock()
			req.closed = true
			go m.removeClosedRequest(handle)
		}
	}()
	return true
}

func (h Handle) GetTypeID() uint32 {
	return uint32(h >> 32)
}

func (h Handle) GetSeqID() uint32 {
	return uint32(h & 0xffffffff)
}
