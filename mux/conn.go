package mux

import (
	"bufio"
	"net"
	"time"
)

type FuseConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewFuseConn(conn net.Conn) *FuseConn {
	return &FuseConn{
		conn:   conn,
		reader: bufio.NewReader(conn), // 相当于reader和conn绑定了 当reader的buffer不够时 会自动从conn的字节流获取
	}
}

// 实现net.conn接口
func (c *FuseConn) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *FuseConn) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c *FuseConn) Close() error {
	return c.conn.Close()
}

func (c *FuseConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *FuseConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *FuseConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *FuseConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *FuseConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// 预读字节
func (c *FuseConn) Peek(n int) ([]byte, error) {
	return c.reader.Peek(n)
}
