package tests

import (
	"net"
	"testing"

	"github.com/zeromake/http2"
	"github.com/zeromake/http2/tests/utils"
)

func handle(conn net.Conn) error {
	return conn.Close()
}

func TestHttp2(t *testing.T) {
	server, err := utils.NewServer(8888, handle)
	if err != nil {
		panic(err)
	}
	go server.Accept()
	frame := http2.FrameHeader{}
	t.Error(frame)
}
