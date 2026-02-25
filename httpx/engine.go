package httpx

import (
	"hybrid-srv/core"
	"net/http"
)

type Engine struct {
	router *Router
	mws    []core.Middleware
}

func New() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

// 注入中间件
func (e *Engine) Use(middleware ...core.Middleware) {
	e.mws = append(e.mws, middleware...)
}

// 常见请求方法
func (e *Engine) Get(path string, handler core.Handler) {
	e.router.Add("GET", path, handler)
}

func (e *Engine) Post(path string, handler core.Handler) {
	e.router.Add("POST", path, handler)
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 根据请求路径匹配业务方法
	h := e.router.Match(request.Method, request.URL.Path)
	if h == nil {
		http.NotFound(writer, request)
		return
	}

	// 创建上下文
	c := NewCtx(request.Context())

	// 换成包装后的writer
	rw := newResponseWriter(writer)

	// 读取输入
	c.Writer = rw
	c.Request = request
	c.Set("method", request.Method)
	c.Set("path", request.URL.Path)

	// 执行中间件
	final := core.Chain(h, e.mws...)
	final(c)

	// 被终止
	if c.Aborted() {
		return
	}
}
