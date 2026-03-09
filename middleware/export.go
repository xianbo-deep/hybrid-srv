package middleware

import (
	"github.com/xianbo-deep/Fuse/core"
)

func Defaults() []core.HandlerFunc {
	return []core.HandlerFunc{
		Recovery(),
		RequestID(),
		Logger(),
	}
}
