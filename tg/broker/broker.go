package broker

import (
	"github.com/MrReality255/turbo-go/tg/utils"
	"sync"
)

/*
	- every member has a unique id (int64), consisting of 32 bit type, 32 bit id within the type
	- every command has its own id (int64), consisting of 32 bit type, 32 bit id within the type
*/

/*
type IBrokerMember[CmdTx comparable, MemberRec comparable, Command any, MsgType comparable] interface {
	About() MemberRec
	HandleMessage(msg Command)
}
*/

const (
	queueSize = 64
)

type Handle uint64

type ICommand interface {
	GetID() Handle
}

type IMember interface {
	GetID() Handle
	HandleMessage(sender Handle, msg ICommand, ref Handle)
}

type IBroker[Cmd ICommand] interface {
	AddMember(member IMember) IBrokerMember[Cmd]
	Stop()
}

type IBrokerMember[Command ICommand] interface {
	Request(cmd Command) (Command, error)
	RequestMultiple(cmd Command, handler func(response Command) bool)
	Send(receiver Handle, cmd Command, ref Handle)
	Subscribe(cmdType ...uint32)
}

type memberWrapper[Command ICommand] struct {
	member IMember
	broker *controller[Command]
}

type messageWrapper[Command ICommand] struct {
	cmd      Command
	sender   Handle
	receiver Handle
	ref      Handle
}

type subscribersMap map[uint32]map[Handle]bool

type controller[Command ICommand] struct {
	mx sync.Mutex

	chQueue     chan *messageWrapper[Command]
	closed      bool
	members     map[Handle]*memberWrapper[Command]
	subscribers subscribersMap
}

func New[Command ICommand]() IBroker[Command] {
	c := &controller[Command]{
		chQueue:     make(chan *messageWrapper[Command], queueSize),
		members:     make(map[Handle]*memberWrapper[Command]),
		subscribers: make(subscribersMap),
	}
	go c.startLoop()
	return c
}

func (c *controller[Command]) AddMember(member IMember) IBrokerMember[Command] {
	c.mx.Lock()
	defer c.mx.Unlock()
	wrapper := &memberWrapper[Command]{member: member, broker: c}
	c.members[member.GetID()] = wrapper
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
	if msg.receiver.GetSeqID() != 0 {
		go c.members[msg.receiver].handleMessage(msg.sender, msg.cmd, msg.ref)
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
			go c.members[subscriber].handleMessage(msg.sender, msg.cmd, msg.ref)
		}
	}

	// the message has receiver type: send it to one receiver of this type
	if recTypeID := msg.receiver.GetTypeID(); recTypeID != 0 && !handledType[recTypeID] {
		for m, h := range c.members {
			if m.GetTypeID() == recTypeID {
				go h.handleMessage(msg.sender, msg.cmd, msg.ref)
				return
			}
		}
	}
}

func (c *controller[Command]) send(sender Handle, receiver Handle, cmd Command, refID Handle) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.closed {
		return
	}
	c.chQueue <- &messageWrapper[Command]{cmd: cmd, sender: sender, receiver: receiver, ref: refID}
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

func (m *memberWrapper[Command]) Request(cmd Command) (Command, error) {
	//TODO implement me
	panic("implement me")
}

func (m *memberWrapper[Command]) RequestMultiple(cmd Command, handler func(response Command) bool) {
	//TODO implement me
	panic("implement me")
}

func (m *memberWrapper[Command]) Send(receiver Handle, cmd Command, ref Handle) {
	go m.broker.send(m.member.GetID(), receiver, cmd, ref)
}

func (m *memberWrapper[Command]) Subscribe(cmdType ...uint32) {
	m.broker.subscribe(m.member.GetID(), cmdType...)
}

func (m *memberWrapper[Command]) handleMessage(sender Handle, cmd Command, ref Handle) {

}

func (h Handle) GetTypeID() uint32 {
	return uint32(h >> 32)
}

func (h Handle) GetSeqID() uint32 {
	return uint32(h & 0xffffffff)
}
