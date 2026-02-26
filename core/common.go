package core

const (
	CtxKeyRequestID = "request_id"
	CtxKeyMethod    = "method"
	CtxKeyPath      = "path"
	CtxKeyProtocal  = "protocal"
)

const (
	ProtocalHTTP = "HTTP"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

const (
	CodeOK           = 0
	CodeBadRequest   = 1001
	CodeUnauthorized = 2001
	CodeForbidden    = 3001
	CodeNotFound     = 4004
	CodeInternal     = 9001
)
