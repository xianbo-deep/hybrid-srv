package core

import (
	"context"
)

type Ctx interface {
	// 返回上下文
	Context() context.Context
	// 更改上下文
	WithContext(context.Context)
	// 终止
	Abort()
	// 是否终止
	Aborted() bool

	Set(key string, val any)
	Get(key string) (any, bool)

	// 记录错误
	Err(err error)
	// 最后一个错误
	Error() error
	// 全部错误
	Errors() []error
}
