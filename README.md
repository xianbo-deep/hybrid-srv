<h1 align="center">Fuse</h1>

<p align="center">
  <strong>一个轻量级、协议无关的 Go 服务端框架骨架</strong>
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" />
</p>
<br />

**Fuse** 旨在提供统一的服务端开发体验。通过高度抽象的 **Context + Middleware** 模型，它将 HTTP、gRPC、Cron、WebSocket、SSE 等多种协议的处理流程归一化，让开发者能以一致的范式构建服务。

## 核心设计

**为什么设计 Core 包？**

Fuse 采用了 **Core + Facade** 的架构模式。

-   **抽象契约 (`core`)**：`core` 包定义了框架最基础的接口（如 `Ctx`、`Result`），所有协议适配层（`httpx`、`wsx` 等）都依赖于它。这层抽象确保了不同协议之间互不干扰，实现了架构的清晰分层。
-   **统一门面 (`fuse`)**：`fuse` 包作为面向用户的统一入口，对 `core` 的类型进行了封装和别名处理。开发者在使用时只需关注 `fuse` 包提供的 API，无需关心底层的具体实现和复杂的依赖关系。

这种设计使得框架既保持了内部组件的低耦合，又为用户提供了简洁一致的开发体验。

## 支持协议

Fuse 通过连接多路复用（CMUX）技术，支持在单端口上运行多种协议。

-   **HTTP/1.1 & HTTP/2**：内置高性能路由，支持 RESTful API。
-   **gRPC**：无缝集成的 RPC 支持。
-   **WebSocket**：提供升级处理、消息泵及心跳机制。
-   **SSE**：支持 Server-Sent Events 推送。
-   **Cron**：集成的定时任务调度。

## 目录结构

```text
Fuse/
├── core/           # 核心抽象层 (Ctx, Result)
├── cronx/          # 定时任务适配器
├── docs/           # 详细开发文档
├── fuse/           # 用户统一入口 (Facade)
├── grpcx/          # gRPC 协议适配器
├── httpx/          # HTTP 协议适配器
├── middleware/     # 通用中间件集合
├── mux/            # 协议多路复用器 (CMUX)
├── ssex/           # SSE 协议适配器
└── wsx/            # WebSocket 协议适配器
```

## 快速开始

### 安装

```bash
go get github.com/xianbo-deep/Fuse
```

### 示例代码

```go
package main

import "github.com/xianbo-deep/Fuse/fuse"

func main() {
    // 1. 初始化引擎
    app := fuse.New()

    // 2. 注册路由 (使用 fuse.Context)
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
- [x] WebSocket 消息泵 (Message Pump)
- [ ] 增加更多通用中间件 (Limit, Auth, Trace)
- [ ] WebSocket JSON 解析
- [ ] HTTP 参数校验 (Validator)
- [ ] 完善单元测试与 Benchmarks 对比
- [ ] 完善 GoDoc
- [ ] 补充相关文档

更多详细设计和用法请阅读 [docs](./docs) 目录下的文档。

