package grpcx

import (
	"context"
	"net"
	"strconv"
	"strings"
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
	// trustedProxies 可信代理。
	trustedProxies []*net.IPNet
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
		c := NewCtx(context.Background(), e)
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
func (e *Engine) Use(mws ...core.HandlerFunc) {
	e.mws = append(e.mws, mws...)
}

// Server 返回底层的 [grpc.Server] 实例。
//
// 主要用于驱动层获取并启动服务，或者用于注册 gRPC 服务。
func (e *Engine) Server() *grpc.Server {
	return e.server
}

// unaryInterceptor 返回一元拦截器
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

// SetTrustedProxies 设置可信代理。
func (e *Engine) SetTrustedProxies(trustedProxies []string) error {
	e.trustedProxies = make([]*net.IPNet, 0, len(trustedProxies))
	for _, proxy := range trustedProxies {
		if !strings.Contains(proxy, "/") {
			if strings.Contains(proxy, ":") {
				proxy += "/128"
			} else {
				proxy += "/32"
			}
		}
		_, ipNet, err := net.ParseCIDR(proxy)
		if err != nil {
			return err
		}
		e.trustedProxies = append(e.trustedProxies, ipNet)
	}
	return nil
}

// IsTrustedProxies 判断 IP 是否是可信代理。
func (e *Engine) IsTrustedProxies(ip string) bool {
	// 默认不信任任何代理
	if len(e.trustedProxies) == 0 {
		return false
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, trustedProxy := range e.trustedProxies {
		if trustedProxy.Contains(parsedIP) {
			return true
		}
	}
	return false
}
