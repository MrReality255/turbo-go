package net

import (
	"io"
	"sync"

	"github.com/MrReality255/turbo-go/tg/log"
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
	RequestHandlerLoop(requestHandler func(msg T) (T, error))
	Wait()
	Write(data T) error
}

type TypedSocket[T any] struct {
	conn     io.Closer
	mx       *sync.Mutex
	chClose  chan bool
	isClosed bool

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
		mx:   &sync.Mutex{},
		conn: src,
		onRead: func() (T, error) {
			return tsf.onRead(src)
		},
		onWrite: func(data T) error {
			return tsf.onWrite(src, data)
		},
	}
}

func (m *TypedSocket[T]) abort(hint string, err error) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.isClosed {
		return
	}

	log.LogError(hint, err)
	if err = m.closeLocked(); err != nil {
		log.LogError("error while closing the socket: %v", err)
	}
}

func (m *TypedSocket[T]) closeLocked() error {
	if m.isClosed {
		return nil
	}
	m.isClosed = true
	return m.conn.Close()
}

func (m *TypedSocket[T]) Close() error {
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.chClose != nil {
		m.chClose <- true
		close(m.chClose)
		m.chClose = nil
	}
	return m.closeLocked()
}

func (m *TypedSocket[T]) Read() (T, error) {
	return m.onRead()
}

func (m *TypedSocket[T]) RequestHandlerLoop(requestHandler func(msg T) (T, error)) {
	for {
		msg, err := m.Read()
		if err != nil {
			m.abort("error while reading the message", err)
		}

		response, err := requestHandler(msg)
		if err != nil {
			m.abort("error while handling the request", err)
			return
		}

		if err = m.Write(response); err != nil {
			m.abort("error while writing the response", err)
		}
	}
}

func (m *TypedSocket[T]) Wait() {
	utils.ExecLocked(m.mx, func() {
		if m.chClose != nil {
			panic("socket is already waiting")
		}
		m.chClose = make(chan bool, 1)
	})
	<-m.chClose
}

func (m *TypedSocket[T]) Write(data T) error {
	return m.onWrite(data)
}
