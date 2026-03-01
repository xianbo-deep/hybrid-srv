package grpcx

import (
	"Fuse/core"
	"context"
	"sync"

	"google.golang.org/grpc/codes"
)

type Ctx struct {
	ctx context.Context

	// 原生请求对象
	request any

	// 共享数据
	values  map[string]any
	aborted bool

	// 中间件控制
	handlers []core.HandlerFunc
	index    int

	errs []error
	mu   sync.RWMutex
}

func NewCtx(ctx context.Context, request any) *Ctx {
	return &Ctx{
		ctx:     ctx,
		request: request,
		values:  make(map[string]any),
	}
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

func (c *Ctx) Next() core.Result {
	c.index++
	if c.index < len(c.handlers) {
		if c.aborted {
			return core.Result{}
		}
		return c.handlers[c.index](c)
	}
	return core.Result{}
}

func (c *Ctx) Abort() {
	c.aborted = true
	c.index = len(c.handlers)
}

func (c *Ctx) Copy() core.Ctx {

	// TODO
}

func (c *Ctx) Aborted() bool {
	return c.aborted
}

func (c *Ctx) Err(err error) {
	if err == nil {
		return
	}
	c.errs = append(c.errs, err)
}

func (c *Ctx) Error() error {
	if len(c.errs) == 0 {
		return nil
	}
	return c.errs[len(c.errs)-1]
}

func (c *Ctx) Errors() []error {
	out := make([]error, len(c.errs))
	copy(out, c.errs)
	return out
}

func (c *Ctx) Param(key string) string {
	return ""
}

func (c *Ctx) Query(key string) string {
	return ""
}

func (c *Ctx) Bind(data any) {

}

func (c *Ctx) Success(data any) core.Result {
	return core.Result{Code: core.CodeSuccess, Data: data}
}

func (c *Ctx) Fail(code int, msg string) core.Result {
	return core.Result{Code: code, Msg: msg}
}

func (c *Ctx) Render(result core.Result) {

}

func grpcCodeFromBizCode(code int) codes.Code {
	switch code {
	case core.CodeBadRequest:
		return codes.InvalidArgument
	case core.CodeUnauthorized:
		return codes.Unauthenticated
	case core.CodeForbidden:
		return codes.PermissionDenied
	case core.CodeNotFound:
		return codes.NotFound
	case core.CodeInternal:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
