package utils

import (
	"fmt"
	"github.com/MrReality255/turbo-go/tg/log"
	"io"
	"net"
	"net/http"
	"sync"
)

type IAbstractSocket interface {
	io.Closer
	Read() ([]byte, error)
}

type TypedSocket[T any] struct {
	conn     io.Closer
	mx       *sync.Mutex
	isClosed bool

	Read  func() (T, error)
	Write func(data T) error
}

type TypedSocketFactory[T any] struct {
	onRead  func(conn IAbstractSocket) (T, error)
	onWrite func(conn IAbstractSocket, data T) error
}

func GetServerAddr(port int) string {
	return fmt.Sprintf("0.0.0.0:%d", port)
}

func ServeLocal(port int, handler func(conn net.Conn) error, errHandler func(error)) error {
	return Serve(GetServerAddr(port), handler, errHandler)
}
func Serve(addr string, handler func(conn net.Conn) error, errHandler func(err error)) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()

		if err != nil {
			return err
		}

		go func() {
			err := CloseAfter(conn, func() error {
				return handler(conn)
			})
			if err != nil && errHandler != nil {
				errHandler(err)
			}
		}()
	}
}

func NewTcpClient(addr string, port int) (net.Conn, error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))
	return c, err
}

func NewTcpServer(port int, handler func(c net.Conn), errHandler func(error)) (io.Closer, error) {
	listener, err := net.Listen("tcp4", GetServerAddr(port))
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			c, err := listener.Accept()
			if err != nil && errHandler != nil {
				errHandler(err)
			}
			handler(c)
		}

	}()
	return listener, nil
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

func GetHttpFrom(url string, target interface{}) error {
	var cl http.Client
	r, err := cl.Get(url)
	if err != nil {
		return err
	}
	return FromJSON(r.Body, target)
}

func (tsf *TypedSocketFactory[T]) NewTcpClient(addr string, port int) (*TypedSocket[T], error) {
	c, err := NewTcpClient(addr, port)
	if err != nil {
		return nil, err
	}
	return tsf.New(c), nil
}

func (tsf *TypedSocketFactory[T]) New(src IAbstractSocket) *TypedSocket[T] {
	return &TypedSocket[T]{
		mx:   &sync.Mutex{},
		conn: src,
		Read: func() (T, error) {
			return tsf.onRead(src)
		},
		Write: func(data T) error {
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
	return m.closeLocked()
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
