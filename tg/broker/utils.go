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

func NewHandle(typeID, seqID uint32) Handle {
	return Handle((uint64(typeID) << 32) | uint64(seqID))
}

func NewHandleType(typeID uint32) Handle {
	return NewHandle(typeID, 0)
}
