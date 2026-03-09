package core

import (
	"bufio"
	"net"
	"net/http"
)

// HeadWriter 是一个接口，定义了响应写入器的基本功能。
//
// 它主要用于中间件和处理器中检测响应状态，避免重复写入头部。
type HeadWriter interface {
	WriteHeader(int)
	Written() bool
}

// ResponseWriterWrapper 是 HTTP 响应写入器的包装器，采用装饰器模式。
//
// 它包装了原生的 [http.ResponseWriter]，并添加了状态码跟踪、写入状态检测和响应大小统计功能。
type ResponseWriterWrapper struct {
	http.ResponseWriter
	status  int
	written bool
	size    int
}

// NewResponseWriter 创建一个新的 [ResponseWriterWrapper] 实例。
//
// w: 底层原生的 [http.ResponseWriter]，通常是 net/http 包提供的。
//
// 返回值: 一个新的 [ResponseWriterWrapper] 指针，初始状态码为 200，写入状态为 false。
func NewResponseWriter(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{ResponseWriter: w, status: http.StatusOK}
}

// WriteHeader 设置 HTTP 响应状态码并标记响应头已写入。
//
// 如果响应头已经写入（w.written 为 true），则直接返回，避免重复设置。
func (w *ResponseWriterWrapper) WriteHeader(code int) {
	if w.written {
		return
	}
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Write 写入响应体数据到客户端。
//
// 如果响应头尚未写入，会自动设置状态码为 200 并标记为已写入。
//
// 这个方法会统计已写入的字节数，并转发到底层 ResponseWriter。
//
// 返回值: 写入的字节数和可能的错误。
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

// Flush 实现 [http.Flusher] 接口，用于流式响应。
func (w *ResponseWriterWrapper) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack 实现 [http.Hijacker] 接口，用于 WebSocket 等功能。
//
// 返回值: 原始的 TCP 连接、带缓冲的读写器，以及可能的错误。
//
// 如果底层 ResponseWriter 不支持 Hijack，返回 http.ErrNotSupported。
func (w *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

// Status 返回当前的 HTTP 响应状态码。
func (w *ResponseWriterWrapper) Status() int { return w.status }

// Written 返回响应头是否已经被写入。
func (w *ResponseWriterWrapper) Written() bool { return w.written }

// Size 返回写入的响应体的字节数。
func (w *ResponseWriterWrapper) Size() int { return w.size }
