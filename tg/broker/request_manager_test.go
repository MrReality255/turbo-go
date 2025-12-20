package broker

import (
	"fmt"
	"github.com/MrReality255/turbo-go/tg/utils"
	"testing"
	"time"
)

func TestRequestManager(t *testing.T) {
	l := utils.NewStringList(1024, true)

	var sender func(cmd uint32) error

	m := NewRequestManager[uint32](
		CommandDescriptor[uint32]{
			GetID: func(cmd uint32) Handle {
				return Handle(cmd & 0xFF)
			},
			GetRef: func(cmd uint32) Handle {
				return Handle(cmd & 0xFF)
			},
		},
		func(cmd uint32) error {
			return sender(cmd)
		},
		func(Cmd uint32) {
			l.Addf("Receiving %v", Cmd)
		},
		time.Second,
	)

	sender = func(cmd uint32) error {
		l.Addf("Sending %v", cmd)
		go m.Accept(0x3400 | (cmd & 0xFF))
		return nil
	}

	r, err := m.Request(0x3333)
	utils.TestAsString(t, 1, "request", "13363 <nil>", fmt.Sprintf("%v %v", r, err))

}
