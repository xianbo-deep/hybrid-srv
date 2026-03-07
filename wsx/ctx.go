package wsx

import (
	"Fuse/core"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type WsContext struct {
	core.Ctx
	Conn    *websocket.Conn
	mu      *sync.Mutex
	MsgType int
	Data    []byte
}

func NewWsContext(c core.Ctx, conn *websocket.Conn, msgType int, data []byte, mu *sync.Mutex) *WsContext {
	return &WsContext{
		Conn:    conn,
		Ctx:     c,
		MsgType: msgType,
		Data:    data,
		mu:      mu,
	}
}

func (wsc *WsContext) Send(msgType int, data []byte) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	return wsc.Conn.WriteMessage(msgType, data)
}

func (wsc *WsContext) BindJSON(obj interface{}) error {
	return json.Unmarshal(wsc.Data, obj)
}

func (wsc *WsContext) SendJSON(obj interface{}) error {
	// 序列化成字节
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return wsc.Conn.WriteMessage(websocket.TextMessage, data)
}
