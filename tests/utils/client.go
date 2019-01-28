package utils

import (
	"fmt"
	"net"
)

// Client 客户端
type Client struct {
	Port int
	conn net.Conn
}

// New 构造
func NewClient(port int) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if conn != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		Port: port,
	}, nil
}

// Handle 回调
func (c *Client) Handle(handle func(conn net.Conn) error) error {
	return handle(c.conn)
}

// Close 关闭
func (c *Client) Close() error {
	return c.conn.Close()
}
