package grpcx

import (
	"Fuse/core"
	"Fuse/middleware"
	"context"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Engine struct {
	server *grpc.Server
	mws    []core.HandlerFunc
	pool   sync.Pool
}

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

func Default() *Engine {
	engine := New()
	engine.Use(middleware.Defaults()...)
	return engine
}

func (e *Engine) Use(middleware ...core.HandlerFunc) {
	e.mws = append(e.mws, middleware...)
}

// 获取grpc服务实例
func (e *Engine) Server() *grpc.Server {
	return e.server
}

// 一元拦截器
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
		hs := make([]core.HandlerFunc, 0)
		hs = append(hs, e.mws...)

		// 业务函数
		grpcCodeHandler := func(c core.Ctx) core.Result {
			realResp, realErr := handler(c.Context(), req)
			if realErr != nil {
				return c.FailWithError(realErr)
			}
			return c.Success(realResp)
		}

		hs = append(hs, grpcCodeHandler)
		// 将调用链挂载到上下文执行
		c.handlers = hs
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

// 流式拦截器
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

		hs := make([]core.HandlerFunc, 0)
		hs = append(hs, e.mws...)

		streamHandler := func(c core.Ctx) core.Result {
			// 执行原生流式业务逻辑
			err := handler(srv, ss)
			if err != nil {
				return c.FailWithError(err)
			}
			return c.Success(nil)
		}

		hs = append(hs, streamHandler)

		// 执行中间件
		c.handlers = hs
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
