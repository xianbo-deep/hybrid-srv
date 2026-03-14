// Package middleware 提供了常用的中间件，帮助你进行便捷开发。
//
//   - [Recovery]: 异常捕获，用于捕获程序中的异常避免panic导致系统崩溃。
//   - [RequestID]: 请求ID生成。为每个请求基于时间戳和随机数创建唯一ID。
//   - [Logger]: 日志打印。
//   - [RateLimit]: 限流器。基于LRU缓存对每个IP进行限流。
//   - [DistributedTrace]: 链路追踪。
//
// 当然我们无法覆盖所有常用的中间件，因此 fuse 也支持自定义中间件，示例如下：
//
//	func AuthMiddleware() fuse.HandlerFunc {
//		return func(c fuse.Ctx) fuse.Result {
//			/* 在这里编写你的中间件 */
//		}
//	}
package middleware

import (
	"github.com/xianbo-deep/Fuse/core"
)

// Defaults 导出常用的中间件。
func Defaults() []core.HandlerFunc {
	return []core.HandlerFunc{
		Recovery(),
		RequestID(),
		Logger(),
		RateLimit(),
		DistributedTrace(),
	}
}
