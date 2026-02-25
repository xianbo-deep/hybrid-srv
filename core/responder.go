package core

type Responder interface {
	Status(code int)
	String(code int, s string)
	JSON(code int, v any)
}
