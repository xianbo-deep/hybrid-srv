package ssex

import (
	"Fuse/httpx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Stream struct {
	ctx    *httpx.Ctx
	reqCtx context.Context
	mu     sync.Mutex
}

func NewStream(ctx *httpx.Ctx) *Stream {
	return &Stream{
		ctx:    ctx,
		reqCtx: ctx.Request.Context(),
	}
}

var errClosed = errors.New("connection closed by client")

// 发送数据
func (s *Stream) Send(event string, data any) error {
	// 进行心跳检测 客户端断连直接返回错误
	select {
	case <-s.reqCtx.Done():
		return errClosed
	default:
	}

	// 格式化数据
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		dataStr = string(b)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 写入event
	if event != "" {
		_, err := s.ctx.Writer.Write([]byte(fmt.Sprintf("event: %s\n", event)))
		if err != nil {
			return err
		}
	}

	// 写入data 多行写入
	if dataStr != "" {
		dataline := strings.Split(dataStr, "\n")
		for _, line := range dataline {
			_, err := s.ctx.Writer.Write([]byte(fmt.Sprintf("data: %s\n", line)))
			if err != nil {
				return err
			}
		}

	}
	// 补充最后的换行符
	_, err := s.ctx.Writer.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// 每次推送完数据堆缓冲区进行刷新
	s.ctx.Writer.Flush()

	return nil
}

func (s *Stream) startHeartPingPong() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.reqCtx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			s.ctx.Writer.Write([]byte(":ping\n\n"))
			s.ctx.Writer.Flush()
			s.mu.Unlock()
		}
	}
}
