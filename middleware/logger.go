package middleware

import (
	"hybrid-srv/core"
	"log"
	"time"
)

func Logger() core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(c core.Ctx) {
			start := time.Now()
			method, _ := c.Get("method")
			path, _ := c.Get("path")
			rid, _ := c.Get("request_id")
			log.Printf("[req] start method=%v path=%v rid=%v", method, path, rid)
			next(c)
			log.Printf("[req] done  method=%v path=%v rid=%v cost=%s aborted=%v",
				method, path, rid, time.Since(start), c.Aborted())
		}
	}
}
