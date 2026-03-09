package core

// 请求上下文中的键名
const (
	CtxKeyRequestID = "request_id"
	CtxKeyMethod    = "method"
	CtxKeyPath      = "path"
	CtxKeyProtocol  = "protocol"
)

// 协议名称
const (
	ProtocolHTTP = "HTTP"
	ProtocolGRPC = "GRPC"
	ProtocolCRON = "CRON"
	ProtocolWS   = "WEBSOCKET"
	ProtocolSSE  = "SSE"
)

// HTTP 请求方法与 GRPC 拦截器类型
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"

	MethodUnary  = "unary"
	MethodStream = "stream"
)

// 状态码
const (
	CodeSuccess      = 0
	CodeBadRequest   = 1001
	CodeUnauthorized = 2001
	CodeForbidden    = 3001
	CodeNotFound     = 4004
	CodeInternal     = 9001
)
