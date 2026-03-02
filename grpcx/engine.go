package grpcx

import (
	"Fuse/core"
	"Fuse/middleware"
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Engine struct {
	server *grpc.Server
	mws    []core.HandlerFunc
}

func New() *Engine {
	e := &Engine{
		mws: make([]core.HandlerFunc, 0),
	}
	if e.server == nil {
		e.server = grpc.NewServer(grpc.UnaryInterceptor(e.unaryInterceptor()))
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

// 将中间件包装成grpc的拦截器
func (e *Engine) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 获取新的上下文
		c := NewCtx(ctx, req)

		// 补充基础元信息
		c.Set(core.CtxKeyProtocol, core.ProtocolGRPC)
		c.Set(core.CtxKeyMethod, core.MethodRPC)
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
