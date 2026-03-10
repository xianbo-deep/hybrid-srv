package grpcx

import (
	"context"
	"strconv"
	"sync"

	"github.com/xianbo-deep/Fuse/core"
	"github.com/xianbo-deep/Fuse/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Engine 是 Fuse 框架的 gRPC 引擎，负责管理和调度 gRPC 服务的中间件、拦截器和上下文。
//
// 支持一元拦截器和流式拦截器，支持统一的中间件链。
type Engine struct {
	// server 是底层的 gRPC 服务器实例，处理实际的 RPC 调用。
	server *grpc.Server
	// mws 是注册的中间件链。
	mws []core.HandlerFunc
	// pool 是 Ctx 对象的同步池，用于重用上下文实例，减少内存分配和垃圾回收压力。
	pool sync.Pool
}

// New 创建一个新的 gRPC 引擎实例，可选的 gRPC 服务器配置选项。
//
// 此方法会自动配置一元和流式拦截器，以支持 Fuse 的中间件机制。
//
// opts: 可选的 [grpc.ServerOption] 配置，用于自定义 gRPC 服务器行为。
func New(opts ...grpc.ServerOption) *Engine {
	e := &Engine{
		mws: make([]core.HandlerFunc, 0),
	}

	opts = append(opts,
		grpc.UnaryInterceptor(e.unaryInterceptor()),
		grpc.StreamInterceptor(e.streamInterceptor()),
	)

	if e.server == nil {
		e.server = grpc.NewServer(opts...)
	}

	e.pool.New = func() any {
		c := NewCtx(context.Background())
		return c
	}
	return e
}

// Default 创建一个带有默认中间件的 gRPC 引擎。
func Default() *Engine {
	engine := New()
	engine.Use(middleware.Defaults()...)
	return engine
}

// Use 向引擎注册一个或多个中间件。
func (e *Engine) Use(middleware ...core.HandlerFunc) {
	e.mws = append(e.mws, middleware...)
}

// Server 返回底层的 gRPC 服务器实例，允许用户注册具体的 gRPC 服务。
func (e *Engine) Server() *grpc.Server {
	return e.server
}

// unaryInterceptor 创建并返回一元 RPC 的服务器拦截器。
// 这个拦截器负责将 gRPC 的一元调用转换为 Fuse 的中间件执行流程，
// 包括上下文管理、中间件链执行、错误处理和元数据传递。
//
// 拦截器执行流程：
//  1. 从对象池获取或创建上下文
//  2. 设置协议、方法和路径元信息
//  3. 执行中间件链（包括业务中间件）
//  4. 执行业务处理函数
//  5. 处理结果，映射状态码
//  6. 重置并回收上下文
func (e *Engine) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 获取新的上下文
		c := e.pool.Get().(*Ctx)

		// 传入原生上下文
		c.WithContext(ctx)

		defer func() {
			// 清空上下文状态
			c.reset()

			// 回收上下文
			e.pool.Put(c)
		}()

		// 传入请求对象
		c.request = req

		// 补充基础元信息
		c.Set(core.CtxKeyProtocol, core.ProtocolGRPC)
		c.Set(core.CtxKeyMethod, core.MethodUnary)
		c.Set(core.CtxKeyPath, info.FullMethod)

		// 组装调用链
		c.handlers = append(c.handlers, e.mws...)

		// 业务函数
		grpcCodeHandler := func(c core.Ctx) core.Result {
			realResp, realErr := handler(c.Context(), req)
			if realErr != nil {
				return c.FailWithError(realErr)
			}
			return c.Success(realResp)
		}

		// 将调用链挂载到上下文执行
		c.handlers = append(c.handlers, grpcCodeHandler)
		c.index = -1
		res := c.Next()

		// 业务状态码写到元数据中
		trailer := metadata.Pairs("x-biz-code", strconv.Itoa(res.Code))

		// 挂载元数据
		_ = grpc.SetTrailer(ctx, trailer)

		if res.Code != core.CodeSuccess {
			grpcCode := res.GetGrpcStatus()
			var finalCode codes.Code
			if grpcCode == 0 {
				finalCode = grpcCodeFromBizCode(res.Code)
			} else {
				finalCode = codes.Code(grpcCode)
			}
			return nil, status.Error(finalCode, res.Msg)
		}

		/*
			返回的数据需要实现proto.Message接口
		*/
		return res.Data, nil
	}
}

// streamInterceptor 创建并返回流式 RPC 的服务器拦截器。
//
// 这个拦截器处理 gRPC 流式调用，支持双向流、客户端流和服务器流。
func (e *Engine) streamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		c := e.pool.Get().(*Ctx)

		// 将默认上下文更改为ServerStream提供的长连接上下文
		c.WithContext(ss.Context())

		defer func() {
			// 清空上下文状态
			c.reset()

			// 回收上下文
			e.pool.Put(c)
		}()

		c.Set(core.CtxKeyProtocol, core.ProtocolGRPC)
		c.Set(core.CtxKeyMethod, core.MethodStream)
		c.Set(core.CtxKeyPath, info.FullMethod)

		c.handlers = append(c.handlers, e.mws...)

		streamHandler := func(c core.Ctx) core.Result {
			// 执行原生流式业务逻辑
			err := handler(srv, ss)
			if err != nil {
				return c.FailWithError(err)
			}
			return c.Success(nil)
		}

		c.handlers = append(c.handlers, streamHandler)

		// 执行中间件
		c.index = -1
		res := c.Next()

		// 业务状态码写到元数据中
		trailer := metadata.Pairs("x-biz-code", strconv.Itoa(res.Code))

		// 挂载元数据
		// 这里用ss挂载trailer
		ss.SetTrailer(trailer)

		if res.Code != core.CodeSuccess {
			grpcCode := res.GetGrpcStatus()
			var finalCode codes.Code
			if grpcCode == 0 {
				finalCode = grpcCodeFromBizCode(res.Code)
			} else {
				finalCode = codes.Code(grpcCode)
			}
			// 将业务错误映射为 gRPC 标准错误
			return status.Error(finalCode, res.Msg)
		}

		return nil
	}
}
