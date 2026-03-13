package middleware

import (
	"github.com/xianbo-deep/Fuse/core"
)

// Defaults 导出常用的中间件。
func Defaults() []core.HandlerFunc {
	return []core.HandlerFunc{
		Recovery(),
		RequestID(),
		Logger(),
		RateLimit(),
		DistributedTrace(),
	}
}
