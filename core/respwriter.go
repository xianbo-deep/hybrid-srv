package core

import (
	"bufio"
	"net"
	"net/http"
)

type HeadWriter interface {
	WriteHeader(int)
	Written() bool
}

/*
装饰器模式
包装原有的ResponseWriter
*/
type ResponseWriterWrapper struct {
	http.ResponseWriter
	status  int
	written bool
	size    int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{ResponseWriter: w, status: http.StatusOK}
}

// 设置HTTP头部和状态码
func (w *ResponseWriterWrapper) WriteHeader(code int) {
	if w.written {
		return
	}
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// 设置HTTP响应体
func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	// 未设置头部 默认使用200
	if !w.written {
		w.status = http.StatusOK
		w.written = true
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

// 流式/SSE
func (w *ResponseWriterWrapper) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// WebSocket
func (w *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

func (w *ResponseWriterWrapper) Status() int   { return w.status }
func (w *ResponseWriterWrapper) Written() bool { return w.written }
func (w *ResponseWriterWrapper) Size() int     { return w.size }
