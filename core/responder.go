package core

type HttpResponder interface {
	Status(code int)
	String(code int, s string)
	JSON(code int, v any)
	Result(r Result)
}

type GrpcResponder interface {
}
