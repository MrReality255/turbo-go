package broker

import (
	"sync"
	"time"

	"github.com/MrReality255/turbo-go/tg/utils"
)

type IMember[Command ICommand] interface {
	Request(receiver Handle, cmd Command) (Command, error)
	RequestMultiple(receiver Handle, cmd Command, handler RequestHandler[Command])
	Send(receiver Handle, cmd Command)
	Subscribe(cmdType ...uint32)
	Close()
}

type memberWrapper[Command ICommand] struct {
	id             Handle
	descriptor     CommandDescriptor[Command]
	messageHandler MemberMessageHandler[Command]
	broker         *controller[Command]
	requestTimeout time.Duration

	reqManager map[Handle]IRequestManager[Command]
	mx         sync.Mutex
}

func (m *memberWrapper[Command]) Close() {
	m.broker.removeMember(m.id)
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
	if receiver == HandleAny {
		panic("request must have a receiver")
	}
	m.getReqManager(receiver).RequestMultiple(cmd, handler)
}

func (m *memberWrapper[Command]) Send(receiver Handle, cmd Command) {
	go m.broker.send(m.id, receiver, cmd)
}

func (m *memberWrapper[Command]) Subscribe(cmdType ...uint32) {
	m.broker.subscribe(m.id, cmdType...)
}

func (m *memberWrapper[Command]) getReqManager(receiver Handle) IRequestManager[Command] {
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.reqManager[receiver] == nil {
		m.reqManager[receiver] = NewRequestManager[Command](
			m.descriptor,
			func(cmd Command) error {
				m.Send(receiver, cmd)
				return nil
			},
			func(cmd Command) {
				m.messageHandler(receiver, cmd, m)
			},
			m.requestTimeout,
		)
	}
	return m.reqManager[receiver]
}

func (m *memberWrapper[Command]) handleMessage(sender Handle, cmd Command) {
	m.getReqManager(sender).Accept(cmd)
}
