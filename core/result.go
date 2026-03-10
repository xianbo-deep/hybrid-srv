package core

// H 是 map[string]any 的简短类型别名，用于方便地构建 JSON 响应数据。
type H map[string]any

// Result 是 Fuse 框架的跨协议统一返回值结构体。
//
// 它封装了业务响应数据，可以在 HTTP、gRPC 等不同协议间保持一致的数据结构。
type Result struct {
	// 0表示成功 其它表示失败
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`

	// 跨协议元信息
	Meta map[string]string `json:"meta,omitempty"`

	// http状态码
	httpStatus int

	// grpc状态码
	grpcStatus int
}

// Success 创建一个成功的 [Result] 实例。
//
// data: 成功时返回的业务数据，可以是任意类型（结构体、切片、基础类型等）。
//
// 返回值: Code 为 0 的成功 [Result]，Msg 为空字符串。
func Success(data any) Result {
	return Result{
		Code: 0,
		Data: data,
	}
}

// Fail 创建一个失败的 [Result] 实例。
//
// code: 业务错误码，非零值，建议在业务层统一定义。
//
// msg: 错误描述信息，用于客户端展示和日志记录。
//
// 返回值: 包含指定错误码和信息的 [Result]，Data 为 nil。
func Fail(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
	}
}

// WithMsg 设置返回 [Result] 的 Msg，支持链式调用。
func (r Result) WithMsg(msg string) Result {
	r.Msg = msg
	return r
}

// WithMeta 设置返回 [Result] 的元数据，支持链式调用。
func (r Result) WithMeta(k, v string) Result {
	if r.Meta == nil {
		r.Meta = map[string]string{}
	}
	r.Meta[k] = v
	return r
}

// WithData 设置返回 [Result] 携带的信息，支持链式调用。
func (r Result) WithData(data any) Result {
	r.Data = data
	return r
}

// WithHttpStatus 设置 [Result] 的 HTTP 状态码。
func (r Result) WithHttpStatus(status int) Result {
	r.httpStatus = status
	return r
}

// GetHttpStatus 获取 [Result] 的 HTTP 状态码。
func (r Result) GetHttpStatus() int {
	return r.httpStatus
}

// WithGrpcStatus 设置 [Result] 的 GRPC 状态码。
func (r Result) WithGrpcStatus(grpcStatus int) Result {
	r.grpcStatus = grpcStatus
	return r
}

// GetGrpcStatus 获取 [Result] 的 GRPC 状态码。
func (r Result) GetGrpcStatus() int {
	return r.grpcStatus
}
