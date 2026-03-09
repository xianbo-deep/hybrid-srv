package middleware

import "github.com/xianbo-deep/Fuse/core"

func DistributedTrace() core.HandlerFunc {
	return func(c core.Ctx) core.Result {

		return c.Next()
	}
}
