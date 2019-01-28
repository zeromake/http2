package utils

import (
	"fmt"
	"net"
)

type handleType func(conn net.Conn) error

// Server http2 server
type Server struct {
	Port   int
	l      net.Listener
	handle handleType
}

// New 构造函数
func NewServer(port int, handle handleType) (*Server, error) {
	listen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}
	return &Server{
		Port:   port,
		l:      listen,
		handle: handle,
	}, nil
}

// Accept 开启连接读取
func (s *Server) Accept() {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			break
		}
		go func() {
			err := s.handle(conn)
			if err != nil {
				panic(err)
			}
		}()
	}
}

// Close 关闭
func (s *Server) Close() error {
	return s.l.Close()
}
