package middleware

import (
	"context"

	"github.com/xianbo-deep/Fuse/core"
)

// DistributedTrace 链路追踪
//
// 将 RequestID 注入标准上下文并替换掉 [core.Ctx] 的底层上下文，方便第三方服务（如 GORM，Redis）等进行调用。
func DistributedTrace() core.HandlerFunc {
	return func(c core.Ctx) core.Result {
		var traceID string
		// 从原始上下文提取请求ID
		if v, ok := c.Get(core.CtxKeyRequestID); ok {
			// 内部进行类型断言
			if s, ok := v.(string); ok {
				traceID = s
			}
		}

		if traceID == "" {
			traceID = core.DefaultTraceID
		}

		// 以原始底层上下文作为父上下文
		traceCtx := context.WithValue(c.Context(), core.CtxKeyTraceID, traceID)

		// 替换底层上下文
		c.WithContext(traceCtx)

		// 在原始上下文添加信息
		c.Set(core.CtxKeyTraceID, traceID)

		return c.Next()
	}
}
