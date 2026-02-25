package middleware

import (
	"fmt"
	"hybrid-srv/core"
)

func Recovery() core.HandlerFunc {
	return func(c core.Ctx) {
		// 捕获panic
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic: %v", r)
				c.Err(err)

				c.Abort()
				return
			}
		}()

		c.Next()
	}
}
