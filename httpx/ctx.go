package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Ctx struct {
	ctx context.Context

	// http相关
	Writer  http.ResponseWriter
	Request *http.Request

	values  map[string]any
	aborted bool

	// 错误
	errs []error

	// 锁
	mu sync.RWMutex
}

// 实现core.Ctx接口
func (c *Ctx) Context() context.Context {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ctx
}

func (c *Ctx) WithContext(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	c.mu.Lock()
	c.ctx = ctx
	c.mu.Unlock()
}

func (c *Ctx) Set(key string, value any) {
	c.mu.Lock()
	c.values[key] = value
	c.mu.Unlock()
}

func (c *Ctx) Get(key string) (any, bool) {
	c.mu.RLock()
	v, ok := c.values[key]
	c.mu.RUnlock()
	return v, ok
}

func (c *Ctx) Abort() {
	c.mu.Lock()
	c.aborted = true
	c.mu.Unlock()
}

func (c *Ctx) Aborted() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.aborted
}

func (c *Ctx) Err(err error) {
	if err == nil {
		return
	}
	c.mu.Lock()
	c.errs = append(c.errs, err)
	c.mu.Unlock()
}

func (c *Ctx) Error() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.errs) == 0 {
		return nil
	}
	return c.errs[len(c.errs)-1]
}

func (c *Ctx) Errors() []error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// 返回拷贝，防止外部篡改
	out := make([]error, len(c.errs))
	copy(out, c.errs)
	return out
}

func NewCtx(ctx context.Context) *Ctx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Ctx{ctx: ctx, values: make(map[string]any)}
}

// 设置状态码（不写入header）
func (c *Ctx) Status(code int) {
	// 类型断言
	if rw, ok := c.Writer.(*responseWriter); ok {
		if !rw.written {
			rw.status = code
			return
		}
		return
	}
}

// 设置text/plain
func (c *Ctx) String(code int, s string) {
	h := c.Writer.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "text/plain; charset=utf-8")
	}
	c.Writer.WriteHeader(code)
	_, _ = c.Writer.Write([]byte(s))
}

// 设置json
func (c *Ctx) JSON(code int, v any) {
	h := c.Writer.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "application/json; charset=utf-8")
	}
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Set("Content-Length", strconv.Itoa(len(b)))
	c.Writer.WriteHeader(code)
	_, _ = c.Writer.Write(b)

}
