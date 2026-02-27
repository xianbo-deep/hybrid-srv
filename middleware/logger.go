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
		method := getString(c, core.CtxKeyMethod)
		path := getString(c, core.CtxKeyPath)
		protocol := getString(c, core.CtxKeyProtocol)

		log.Printf("[FUSE] ---> [%s] %s %s", protocol, method, path)

		res := c.Next()

		cost := time.Since(start)
		// 错误
		err := c.Error()
		if err != nil {
			log.Printf("[FUSE] <--- [ERROR] [%s] %s %s | rid=%s | cost=%v | aborted=%v | err=%v",
				protocol, method, path, rid, cost, c.Aborted(), err)
		} else if res.Code != core.CodeSuccess {
			log.Printf("[FUSE] <--- [FAIL] [%s] %s %s | rid=%s | cost=%v | code=%d | msg=%s",
				protocol, method, path, rid, cost, res.Code, res.Msg)
		} else {
			log.Printf("[FUSE] <--- [OK] [%s] %s %s | rid=%s | cost=%v",
				protocol, method, path, rid, cost)
		}
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
