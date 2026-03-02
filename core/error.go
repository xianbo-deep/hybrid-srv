package core

import "fmt"

type BizError struct {
	Code int // 业务状态码
	Msg  string

	HttpStatus int
	GrpcStatus int
}

// 实现接口才可被断言成error
func (e *BizError) Error() string {
	return fmt.Sprintf("code = %d, msg = %s", e.Code, e.Msg)
}

// 暴露接口给用户
func NewError(code int, msg string) error {
	return &BizError{Code: code, Msg: msg}
}

func (e *BizError) WithHttpStatus(status int) *BizError {
	e.HttpStatus = status
	return e
}

// WithGrpcStatus 允许业务层在抛出错误时，强行指定底层的 gRPC 状态码
func (e *BizError) WithGrpcStatus(status int) *BizError {
	e.GrpcStatus = status
	return e
}
