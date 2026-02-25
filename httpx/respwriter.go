package httpx

import "net/http"

/*
装饰器模式
包装原有的ResponseWriter
*/
type responseWriter struct {
	http.ResponseWriter
	status  int
	written bool
	size    int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: 200}
}

// 设置HTTP头部和状态码
func (w *responseWriter) WriteHeader(code int) {
	if w.written {
		return
	}
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// 设置HTTP响应体
func (w *responseWriter) Write(b []byte) (int, error) {
	// 未设置头部 默认使用200
	if !w.written {
		w.WriteHeader(w.status)
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *responseWriter) Status() int   { return w.status }
func (w *responseWriter) Written() bool { return w.written }
func (w *responseWriter) Size() int     { return w.size }
