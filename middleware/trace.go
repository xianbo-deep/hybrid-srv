package middleware

import "Fuse/core"

func DistributedTrace() core.HandlerFunc {
	return func(c core.Ctx) core.Result {

		return c.Next()
	}
}
