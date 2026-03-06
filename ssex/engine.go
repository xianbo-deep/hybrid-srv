package ssex

import (
	"Fuse/core"
	"Fuse/httpx"
	"errors"
	"net/http"
)

type SSEHandlerFunc func(c core.Ctx, stream *Stream) error

func Upgrade(sseHandler SSEHandlerFunc) core.HandlerFunc {
	return func(c core.Ctx) core.Result {
		// 类型断言
		ctx, ok := c.(*httpx.Ctx)
		if !ok {
			return c.Fail(core.CodeBadRequest, "can not upgrade to sse without http request")
		}

		// 设置SSE响应头
		ctx.SetSSEHeader()
		ctx.Status(http.StatusOK)
		ctx.Writer.Flush()

		// 初始化stream实例
		stream := NewStream(ctx)

		// 启动守护进程 监听客户端是否断连
		go stream.startHeartPingPong()

		// 执行业务逻辑
		if err := sseHandler(ctx, stream); err != nil {
			if !errors.Is(err, errClosed) {
				c.FailWithError(err)
			}
		}

		return core.Result{}
	}
}
