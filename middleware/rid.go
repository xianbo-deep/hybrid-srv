package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hybrid-srv/core"
	"time"
)

func RequestID() core.HandlerFunc {
	return func(c core.Ctx) core.Result {
		// 已存在则跳过
		if v, ok := c.Get(core.CtxKeyRequestID); ok {
			if s, ok := v.(string); ok && s != "" {
				return c.Next()
			}
		}

		requestID := generateRequestID()
		c.Set(core.CtxKeyRequestID, requestID)
		return c.Next()
	}
}

func generateRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	// 时间戳-随机数
	return fmt.Sprintf("%x-%s", time.Now().UnixNano(), hex.EncodeToString(b[:]))
}
