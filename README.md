<div align="center">
  <h1>Fuse</h1>
  <p><strong>一个轻量级、协议无关的 Go 服务端框架骨架</strong></p>
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" />
</div>

<br />

**Fuse** 旨在提供统一的服务端开发体验。通过高度抽象的 **Context + Middleware** 模型，它将 HTTP、gRPC、Cron、WebSocket、SSE 等多种协议的处理流程归一化，让开发者能以一致的范式构建服务。

## 核心设计

**Core + Facade 模式**

Fuse 采用独立且极简的 `core` 包来定义核心接口，所有协议适配层（`httpx`, `grpcx` 等）均依赖 `core`。而 `fuse` 包作为统一门面，对核心类型进行了**类型别名 (Type Aliasing)**。

开发者**仅需引入 `fuse` 包**，即可获得完整的能力支持。你使用的 `fuse.Context` 本质上就是 `core.Ctx`，但无需直接感知内部结构，既保证了架构的解耦，又维持了 API 的整洁。

## 目录结构

```text
.
├── core/           # 核心抽象层 (Ctx, Result)，解耦依赖
├── fuse/           # 用户统一入口 (Facade)，提供极简 API
├── httpx/          # HTTP 协议适配器
├── grpcx/          # gRPC 协议适配器
├── cronx/          # 定时任务适配器
├── ssex/           # Server-Sent Events 适配器
├── wsx/            # WebSocket 适配器
├── middleware/     # 通用中间件集合
├── mux/            # 连接多路复用
└── docs/           # 详细开发文档
```

## 快速开始

### 安装

```bash
go get github.com/xianbo-deep/fuse
```

### 示例代码

```go
package main

import "Fuse/fuse"

func main() {
    // 1. 初始化引擎
    app := fuse.New()

    // 2. 注册路由 (使用 fuse.Context，无需引入 core)
    app.HTTP().Get("/ping", func(c fuse.Context) fuse.Result {
        return c.Success(fuse.H{"message": "pong"})
    })

    // 3. 启动
    if err := app.Run(":8080"); err != nil {
        panic(err)
    }
}
```

## 未来规划

- [x] 完善 gRPC 协议支持
- [x] 优化中间件链内存分配
- [x] 实现优雅停机
- [x] WebSocket 协议适配 (wsx) 及心跳检测
- [x] HTTP 路由优先级控制及 Radix Tree 实现
- [x] HTTP 路由分组
- [x] SSE 协议适配 (ssex) 及 Keep-Alive
- [x] 实现单端口协议分发 (CMUX)
- [ ] 增加更多通用中间件 (Limit, Auth, Trace)
- [x] WebSocket 消息泵 (Message Pump)
- [ ] WebSocket JSON 解析
- [ ] HTTP 参数校验 (Validator)
- [ ] 完善单元测试与 Benchmarks 对比
- [ ] 完善 GoDoc
- [ ] 补充相关文档

更多详细设计和用法请阅读 [docs](./docs) 目录下的文档。

