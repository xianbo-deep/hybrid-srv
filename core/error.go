package core

import "fmt"

// BizError 是 Fuse 框架定义的业务错误类型，用于在业务逻辑中统一表示和处理错误。
//
// 它包含了业务状态码、错误信息以及可选的 HTTP 和 gRPC 协议状态码。
//
// 实现了 error 接口。
type BizError struct {
	// Code 是业务状态码，用于标识具体的业务错误类型。
	Code int
	// Msg 是错误描述信息，用于客户端展示或日志记录。
	Msg string

	// HttpStatus 是 HTTP 协议对应的状态码。
	//
	// 如果不设置，框架会自动根据 Code 转换为合适的 HTTP 状态码。
	HttpStatus int
	// GrpcStatus 是 gRPC 协议对应的状态码，对应 codes.Code 枚举值。
	//
	// 如果不设置，框架会根据 Code 自动映射为合适的 gRPC 状态码。
	GrpcStatus int
}

// Error 实现了 error 接口的 Error 方法，返回格式化的错误字符串。
//
// 使得 [BizError] 可以直接作为标准 error 使用。
func (e *BizError) Error() string {
	return fmt.Sprintf("code = %d, msg = %s", e.Code, e.Msg)
}

// NewError 创建一个新的 [BizError] 实例。
//
// code: 业务错误码，建议在业务层统一定义和维护。
//
// msg: 错误描述信息，应当简洁明了，便于调试和问题定位。
//
// 返回值: 一个可继续链式调用的 *[BizError] 实例。
func NewError(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

// WithHttpStatus 为 [BizError] 设置 HTTP 协议状态码，支持链式调用。
func (e *BizError) WithHttpStatus(status int) *BizError {
	e.HttpStatus = status
	return e
}

// WithGrpcStatus 允许业务层在抛出错误时，设置 gRPC 状态码。
func (e *BizError) WithGrpcStatus(status int) *BizError {
	e.GrpcStatus = status
	return e
}
