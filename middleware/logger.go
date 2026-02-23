package middleware

import (
	"hybrid-srv/core"
	"log"
	"time"
)

func Logger() core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(c *core.Ctx) {
			start := time.Now()
			log.Printf("[request] id=%s start", c.RequestID)
			next(c)
			log.Printf("[req] id=%s done cost=%s aborted=%v", c.RequestID, time.Since(start), c.Aborted())
		}
	}
}
