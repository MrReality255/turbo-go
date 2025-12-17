package broker

const (
	HandleAny = 0
)

type Handle uint64

type ICommand any

type CommandDescriptor[Command ICommand] struct {
	GetID  func(cmd Command) Handle
	GetRef func(cmd Command) Handle
}

func (h Handle) GetTypeID() uint32 {
	return uint32(h >> 32)
}

func (h Handle) GetSeqID() uint32 {
	return uint32(h & 0xffffffff)
}
