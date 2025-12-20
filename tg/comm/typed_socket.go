package comm

import (
	"io"

	"github.com/MrReality255/turbo-go/tg/utils"
)

type IAbstractSocket interface {
	io.Closer
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

type ITypedSocket[T any] interface {
	io.Closer
	Read() (T, error)
	Wait()
	Write(data T) error
}

type TypedSocket[T any] struct {
	conn io.Closer
	p    utils.IRunner

	onRead  func() (T, error)
	onWrite func(data T) error
}

type TypedSocketFactory[T any] struct {
	onRead  func(conn IAbstractSocket) (T, error)
	onWrite func(conn IAbstractSocket, data T) error
}

func NewTypedSocketFactory[T any](
	onRead func(conn IAbstractSocket) (T, error),
	onWrite func(conn IAbstractSocket, data T) error,
) *TypedSocketFactory[T] {
	return &TypedSocketFactory[T]{
		onRead:  onRead,
		onWrite: onWrite,
	}
}

func (tsf *TypedSocketFactory[T]) NewTcpClient(addr string, port int) (ITypedSocket[T], error) {
	c, err := NewTcpClient(addr, port)
	if err != nil {
		return nil, err
	}
	return tsf.New(c), nil
}

func (tsf *TypedSocketFactory[T]) New(src IAbstractSocket) ITypedSocket[T] {
	return &TypedSocket[T]{
		conn: src,
		onRead: func() (T, error) {
			return tsf.onRead(src)
		},
		onWrite: func(data T) error {
			return tsf.onWrite(src, data)
		},
		p: utils.NewRuner(
			nil,
			func() error {
				return src.Close()
			},
		),
	}
}

func (m *TypedSocket[T]) Close() error {
	return m.p.Close()
}

func (m *TypedSocket[T]) Read() (T, error) {
	return m.onRead()
}

func (m *TypedSocket[T]) Wait() {
	m.p.Wait()
}

func (m *TypedSocket[T]) Write(data T) error {
	return m.onWrite(data)
}
