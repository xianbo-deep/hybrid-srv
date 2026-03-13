package core

import (
	"context"
)

// Ctx 是 Fuse 中预定义的上下文接口，定义了请求处理的生命周期和共享数据。
//
// 不同协议的上下文需要实现此接口。
type Ctx interface {
	// Context 用于返回底层的上下文
	Context() context.Context
	// WithContext 用于替换底层的上下文
	WithContext(context.Context)
	// Aborted 判断当前请求的中间件链是否已被终止。
	//
	// 通常在调用 Abort 后返回 true，用于提前结束请求处理。
	Aborted() bool

	// Next 执行中间件链中的下一个处理器，支持洋葱模型的中间件架构。
	//
	// 必须在中间件的适当位置调用，否则后续中间件和业务逻辑不会执行。
	Next() Result
	// Abort 终止当前中间件链的执行
	Abort()

	// Set 在当前请求的上下文中存储一个键值对，数据仅在当前请求生命周期内有效。
	//
	// 注意：此存储非线程安全，不应在多个 goroutine 中并发访问。
	Set(key string, val any)
	// Get 获取上下文中存储的值，返回存储的值和是否存在标志。
	//
	// 如果 key 不存在，第二个返回值为 false。
	Get(key string) (any, bool)

	// Render 将 [Result] 渲染并写入具体协议的响应中。
	//
	// 在 HTTP 协议下，它会负责写入状态码和响应体；
	//
	// 在 gRPC 协议下，它会将结果映射为 gRPC 状态码。
	Render(result Result)

	// Copy 对当前上下文进行深拷贝，返回一个独立的副本。
	//
	// 重要：如果需要在新的 goroutine 中异步处理请求上下文，必须调用此方法，
	//
	// 以避免数据竞争和生命周期问题。
	Copy() Ctx

	// Success 构造一个成功的响应结果，通常包含业务数据。
	//
	// data: 成功响应中携带的业务数据，可以是任意类型。
	Success(data any) Result

	// Fail 构造一个失败的响应结果，包含错误码和错误信息。
	//
	// code: 业务错误码，用于客户端识别错误类型。
	//
	// msg: 人类可读的错误描述信息。
	Fail(code int, msg string) Result
	// FailWithError 使用 error 对象构造一个失败的响应结果。
	//
	// 这里的 error 对象在 Fuse 中为 BizError，它实现了 error 接口。
	//
	// 根据传入的 error 对象自动设置错误信息。
	FailWithError(err error) Result

	// Param 获取路由路径中的参数（例如路由 /user/:id 中的 id 值）。
	//
	// key: 路径参数的名称，如 "id"。
	//
	// 返回值: 路径参数的字符串值，如果不存在则返回空字符串。
	Param(key string) string
	// Query 获取 URL 查询字符串中的参数（如 /user?age=12 中的 age）。
	//
	// key: 查询参数的名称，如 "age"。
	//
	// 返回值: 查询参数的字符串值，如果不存在则返回空字符串。
	Query(key string) string

	// Bind 将请求体中的数据绑定到指定的结构体指针 v。
	//
	// 根据协议和 Content-Type 自动选择解析器。
	//
	// 绑定失败时会返回相应的错误。
	Bind(v any) error

	// Err 记录一个错误到上下文错误列表中，用于收集请求处理过程中的多个错误。
	//
	// 通常用于中间件或业务逻辑中记录非致命性错误。
	Err(err error)
	// Error 返回最后一个被记录的错误，如果没有错误则返回 nil。
	Error() error
	// Errors 返回当前上下文中记录的所有错误列表。
	Errors() []error

	// ClientIP 获得当前请求的IP。
	ClientIP() string
}
