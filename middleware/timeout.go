package middleware

import (
	"context"
	"errors"
	"hybrid-srv/core"
	"time"
)

var ErrTimeout = errors.New("request timeout")

func Timeout(d time.Duration) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(c core.Ctx) {
			ctx, cancel := context.WithTimeout(c.Context(), d)
			defer cancel()

			// 替换新的上下文
			c.WithContext(ctx)

			// 创建通道接收结束标志
			done := make(chan struct{})

			// 创建协助执行业务
			go func() {
				defer close(done)
				next(c)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				c.Err(ErrTimeout)
				c.Abort()
				return
			}
		}
	}
}
