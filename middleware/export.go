package middleware

import (
	"Fuse/core"
)

func Defaults() []core.HandlerFunc {
	return []core.HandlerFunc{
		Recovery(),
		RequestID(),
		Logger(),
	}
}
