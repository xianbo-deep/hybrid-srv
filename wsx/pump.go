package wsx

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Pump 是 websocket 模块的消息泵，执行从服务端发送消息到客户端的任务。
//
// 用户通过 [WsContext.Send] 发送信息到 writeChan，消息泵从 writeChan（有缓冲通道）接收需要发送到客户端的信息并执行发送。
//
// 使用读锁防止发信息与心跳检测产生冲突，保证线程安全。
//
// 内部维护 done 监听用户是否断连，防止资源泄漏。
type Pump struct {
	writeChan chan []byte
	conn      *websocket.Conn
	done      chan struct{}
	mu        *sync.Mutex
}

func NewPump(conn *websocket.Conn, done chan struct{}, writeChan chan []byte, mu *sync.Mutex) *Pump {
	return &Pump{
		writeChan: writeChan,
		conn:      conn,
		done:      done,
		mu:        mu,
	}
}

func (p *Pump) WritePump() {
	for {
		select {
		case msg, ok := <-p.writeChan:
			if !ok {
				return
			}
			p.mu.Lock()
			p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // 防止客户端网络不好 消息无法发出阻塞协程
			p.conn.WriteMessage(websocket.TextMessage, msg)
			p.mu.Unlock()
		case <-p.done:
			return
		}
	}
}
