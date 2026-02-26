package middleware

import (
	"fmt"
	"hybrid-srv/core"
)

func Recovery() core.HandlerFunc {
	return func(c core.Ctx) core.Result {
		// 捕获panic
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic: %v", r)
				c.Err(err)

				c.Abort()
				c.Render(core.Fail(core.CodeInternal, "出现Panic"))
			}
		}()

		return c.Next()
	}
}
