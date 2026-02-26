package middleware

import (
	"Fuse/core"
	"fmt"
	"log"
	"time"
)

func Logger() core.HandlerFunc {
	return func(c core.Ctx) core.Result {
		start := time.Now()
		rid := getString(c, core.CtxKeyRequestID)
		log.Printf("[req] start rid=%v", rid)
		res := c.Next()
		method := getString(c, core.CtxKeyMethod)
		path := getString(c, core.CtxKeyPath)
		rid = getString(c, core.CtxKeyRequestID)

		cost := time.Since(start)
		// 错误
		err := c.Error()
		if err != nil {
			log.Printf("[Hybrid-Srv] req done  method=%v path=%v rid=%v cost=%v aborted=%v err=%v",
				method, path, rid, cost, c.Aborted(), err)
		}
		log.Printf("[Hybrid-Srv] req done  method=%s path=%s rid=%s cost=%v aborted=%v",
			method, path, rid, cost, c.Aborted())
		return res
	}

}

func getString(c core.Ctx, key string) string {
	v, ok := c.Get(key)
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
