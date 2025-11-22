package comm

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/MrReality255/turbo-go/tg/utils"
)

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
			err := utils.CloseAfter(conn, func() error {
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

func GetHttpFrom(url string, target interface{}) error {
	var cl http.Client
	r, err := cl.Get(url)
	if err != nil {
		return err
	}
	return utils.FromJSON(r.Body, target)
}
