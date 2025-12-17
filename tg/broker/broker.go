package broker

import (
	"sync"
	"time"

	"github.com/MrReality255/turbo-go/tg/utils"
)

const (
	queueSize = 64
)

type RequestHandler[Cmd ICommand] func(cmd Cmd, err error) bool
type MemberMessageHandler[Cmd ICommand] func(sender Handle, msg Cmd, member IMember[Cmd])

type IBroker[Cmd ICommand] interface {
	AddMember(id Handle, messageHandler MemberMessageHandler[Cmd]) IMember[Cmd]
	Stop()
}

type messageWrapper[Command ICommand] struct {
	cmd      Command
	sender   Handle
	receiver Handle
}

type subscribersMap map[uint32]map[Handle]bool

type controller[Command ICommand] struct {
	requestTimeout time.Duration
	descriptor     CommandDescriptor[Command]

	mx sync.Mutex

	chQueue     chan *messageWrapper[Command]
	closed      bool
	members     map[Handle]*memberWrapper[Command]
	subscribers subscribersMap
}

func New[Command ICommand](
	descriptor CommandDescriptor[Command],
	timeout time.Duration,
) IBroker[Command] {
	c := &controller[Command]{
		descriptor:     descriptor,
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
		descriptor:     c.descriptor,
		messageHandler: messageHandler,
		broker:         c,
		requestTimeout: c.requestTimeout,
		reqManager:     make(map[Handle]*requestManager[Command]),
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
		msgType     = c.descriptor.GetID(msg.cmd).GetTypeID()
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
