package fuse

import (
	"Fuse/core"
	"Fuse/grpcx"
	"Fuse/httpx"
	"net/http"
	"sync"
)

type Context = core.Ctx
type HandlerFunc = core.HandlerFunc
type Result = core.Result
type H = core.H

type Fuse struct {
	// 引擎
	httpEngine *httpx.Engine
	grpcEngine *grpcx.Engine

	// 全局中间件
	mws []core.HandlerFunc
}

func New() *Fuse {
	return &Fuse{
		httpEngine: httpx.New(),
		grpcEngine: grpcx.New(),
	}
}

func (fs *Fuse) Default() *Fuse {
	return &Fuse{
		httpEngine: httpx.Default(),
		grpcEngine: grpcx.Default(),
	}
}

// 挂载中间件
func (fs *Fuse) Use(mws ...core.HandlerFunc) {
	fs.mws = append(fs.mws, mws...)

	// 下发给底层引擎
	fs.httpEngine.Use(mws...)
	fs.grpcEngine.Use(mws...)
}

// 返回引擎
func (fs *Fuse) HTTP() *httpx.Engine {
	return fs.httpEngine
}

func (fs *Fuse) GRPC() *grpcx.Engine {
	return fs.grpcEngine
}

// 启动服务
func (fs *Fuse) Run(httpAddr string) error {
	var wg sync.WaitGroup
	if httpAddr != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = http.ListenAndServe(httpAddr, fs.httpEngine)
		}()
	}
	wg.Wait()
	return nil
}
