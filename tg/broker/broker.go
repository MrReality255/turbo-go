package broker

/*
type IBrokerMember[CmdTx comparable, MemberRec comparable, Command any, MsgType comparable] interface {
	About() MemberRec
	HandleMessage(msg Command)
}
*/

type IBroker[CmdTx comparable, MemberRec comparable, Command any, MsgType comparable] interface {
	AddMember(id MemberRec, )
}

type IBrokerMember[CmdTx comparable, MemberRec comparable, Command any, MsgType comparable] interface {
	Subscribe(msgType MsgType)
	Send(cmd Command)
	Request(cmd Command) (Command, error)
	RequestMultiple(cmd Command, handler func(response Command) bool)
}

type controller[T comparable, M any, C any] struct {
}

func (c controller[T, M, C]) AddMember(member IBrokerMember[C, M]) {
	//TODO implement me
	panic("implement me")
}

func (c controller[T, M, C]) Send(sender M, message C) {
	//TODO implement me
	panic("implement me")
}

func New[T comparable, M any, C any]() IBroker[T, M, C] {
	return &controller[T, M, C]{}
}
