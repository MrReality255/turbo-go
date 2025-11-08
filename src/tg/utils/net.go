package utils

import (
	"fmt"
	"net"
	"net/http"
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

func GetHttpFrom(url string, target interface{}) error {
	var cl http.Client
	r, err := cl.Get(url)
	if err != nil {
		return err
	}
	return FromJSON(r.Body, target)
}
