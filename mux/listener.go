package mux

import (
	"net"
)

type FakeListener struct {
	addr     net.Addr
	connChan chan net.Conn
	done     chan struct{}
}

func NewFakeListener(addr net.Addr) *FakeListener {
	return &FakeListener{
		addr:     addr,
		connChan: make(chan net.Conn, 128),
		done:     make(chan struct{}),
	}
}

// 实现net.Listener接口

// Serve会调用这个方法等待连接
func (l *FakeListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connChan:
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *FakeListener) Addr() net.Addr {
	return l.addr
}

func (l *FakeListener) Close() error {
	close(l.done)
	return nil
}

// 暴露给multiplexer 供其push连接进来
func (l *FakeListener) Push(conn net.Conn) {
	select {
	case l.connChan <- conn:
	case <-l.done:
		conn.Close()
	}
}
