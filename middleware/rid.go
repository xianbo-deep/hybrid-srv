package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hybrid-srv/core"
	"time"
)

func RequestID() core.HandlerFunc {
	return func(c core.Ctx) {
		// 已存在则跳过
		if v, ok := c.Get(core.CtxKeyRequestID); ok {
			if s, ok := v.(string); ok && s != "" {
				c.Next()
				return
			}
		}

		requestID := generateRequestID()
		c.Set(core.CtxKeyRequestID, requestID)
		c.Next()
	}
}

func generateRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	// 时间戳-随机数
	return fmt.Sprintf("%x-%s", time.Now().UnixNano(), hex.EncodeToString(b[:]))
}
