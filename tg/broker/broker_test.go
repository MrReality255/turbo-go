package broker

/*
type testMessage struct {
	id      uint64
	refID   uint64
	payload string
}

type testSender struct{}

func (t testMessage) GetID() Handle {
	return Handle(t.id)
}

func (t testMessage) GetRefID() Handle {
	return Handle(t.refID)
}

func TestBroker(t *testing.T) {
	b := New[testMessage](time.Second)
	m := b.AddMember(
		1,
		func(sender Handle, msg testMessage, member IMember[testMessage]) {
			time.Sleep(time.Millisecond)
			response := testMessage{
				id:      msg.id + 1,
				refID:   msg.id,
				payload: fmt.Sprintf("response: %v", msg.payload),
			}
			member.Send(sender, response)

		},
	)
}

func (t *testSender) GetID() Handle {
	return 1
}

func (t *testSender) HandleMessage(sender Handle, msg testMessage) {
}
*/
