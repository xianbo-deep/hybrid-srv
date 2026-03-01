package grpcx

import (
	"Fuse/core"
	"Fuse/middleware"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Engine struct {
	server *grpc.Server
	mws    []core.HandlerFunc
}

func New() *Engine {
	return &Engine{
		server: grpc.NewServer(),
		mws:    make([]core.HandlerFunc, 0),
	}
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

// 将中间件包装成grpc的拦截器
func (e *Engine) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 获取新的上下文
		c := NewCtx(context.Background(), req)

		// 补充基础元信息
		c.Set(core.CtxKeyProtocol, core.ProtocolGRPC)
		c.Set(core.CtxKeyMethod, core.MethodRPC)
		c.Set(core.CtxKeyPath, info.FullMethod)

		// 组装调用链
		hs := make([]core.HandlerFunc, 0)
		hs = append(hs, e.mws...)

		// 获取业务函数
		grpcCodeHandler := func(c core.Ctx) core.Result {
			realResp, realErr := handler(c.Context(), req)
			if realErr != nil {
				return c.Fail(core.CodeInternal, realErr.Error())
			}
			return c.Success(realResp)
		}
		hs = append(hs, grpcCodeHandler)
		// 将调用链挂载到上下文执行
		c.handlers = hs
		c.index = -1
		res := c.Next()

		if res.Code != core.CodeSuccess {
			grpcCode := res.GetGrpcStatus()
			if grpcCode == 0 {
				grpcCode = grpcCodeFromBizCode(res.Code)
			}
			return nil, status.Error(grpcCode, res.Msg)
		}

		return res.Data, nil

	}
}
