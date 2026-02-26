# Fuse

<div align="center">
<img alt="Go" src="https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go&logoColor=white" />
</div>

**Fuse** 是一个轻量级、可扩展的 Go 服务端框架骨架。

它的核心设计理念是将 **HTTP** 和 **gRPC** (规划中) 的处理流程统一抽象为一套 **"Context + Middleware Chain"** 模型。通过定义通用的 `core` 接口，解耦底层协议与业务逻辑，让开发者能够以一致的方式编写中间件和业务处理器。

## 特性

- **统一核心抽象 (`core`)**：定义了通用的 `Ctx` 接口和 `HandlerFunc`，统一了返回值规范 `Result`。
- **模块化设计**：
  - `httpx`：基于标准库 `net/http` 的轻量级封装，提供路由、上下文管理。
  - `middleware`：开箱即用的通用中间件（Recovery, Logger, RequestID）。
  - `grpcx`：(WIP) 预留的 gRPC 扩展层。
- **中间件机制**：支持洋葱模型（Onion Architecture）的中间件链，轻松实现拦截、鉴权、日志等功能。
- **标准化响应**：内置 `Result` 结构，强制统一API返回格式（Code, Msg, Data）。

## 目录结构

```text
.
├── core/           # [核心] 框架顶层抽象 (Ctx接口, Result定义, Handler定义)
├── httpx/          # [HTTP] HTTP 协议实现 (Engine, Router, Context实现)
├── grpcx/          # [gRPC] gRPC 协议实现 (预留/开发中)
├── middleware/     # [中间件] 通用中间件集合 (Logger, Recovery, RequestID)
├── test/           # [示例] 最小可运行示例
├── go.mod          # 依赖管理
└── README.md       # 项目文档
```

## 快速开始

### 1. 环境要求

- Go 1.25+

### 2. 运行示例

本项目包含一个开箱即用的示例代码。

```bash
# 运行测试示例
go run ./test/main.go
```

服务启动后，监听 `:8080` 端口。你可以通过 curl 或浏览器访问：

```bash
curl http://localhost:8080/ping
# 输出: {"code":0,"msg":"","data":"pong"}
```

## 使用指南

### 基础用法

```go
package main

import (
    "hybrid-srv/core"
    "hybrid-srv/httpx"
    "net/http"
)

func main() {
    // 1. 创建默认引擎 (包含 Logger, Recovery, RequestID 中间件)
    app := httpx.Default()

    // 2. 注册路由
    app.Get("/hello", func(c core.Ctx) core.Result {
        // 返回统一结构
        return core.OK(map[string]string{
            "message": "Hello World",
        })
    })

    // 3. 启动服务
    http.ListenAndServe(":8080", app)
}
```

### 编写中间件

中间件遵循 `func(core.Ctx) core.Result` 签名，通过 `c.Next()` 控制执行流，支持洋葱模型。

## 目前限制与计划

- **路由系统**：目前基于 `map` 精确匹配。
  - *计划*：引入 Trie 树 (前缀树) 路由实现。
- **gRPC 支持**：`grpcx` 目录目前保留，等待实现。
  - *计划*：实现 gRPC 的拦截器适配，使其能复用 `core` 定义的中间件逻辑。
