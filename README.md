# hybrid-srv

一个轻量的 Go 服务端骨架，目标是把 **HTTP** 和未来的 **gRPC** 接入统一到一套“上下文 + 中间件链”的模型里。

当前仓库已实现：

- `core`：框架抽象（`Ctx`、`HandlerFunc`、常量、响应写入包装）
- `httpx`：基于 `net/http` 的 Engine / Router / Context
- `middleware`：内置中间件（Recovery / RequestID / Logger）
- `grpcx`：目录预留（当前为空）

> 说明：目前路由是 **method + path 精确匹配**（map），尚不支持参数路由/通配符。

---

## 快速开始（HTTP）

### 1) 安装

这是一个 Go module：`module hybrid-srv`。

你可以在项目根目录执行：

```powershell
go test ./...
```

### 2) 最小示例

仓库里提供了一个最小可运行示例：`test/main.go`。

运行：

```powershell
go run .\test
```

然后访问：

- `GET http://127.0.0.1:8080/ping`  -> `{"message":"pong"}`

---

## 使用方式

### 注册路由

当前支持的快捷方法：

- `e.Get(path, handler)`
- `e.Post(path, handler)`
- `e.Put(path, handler)`
- `e.Delete(path, handler)`




## 内置中间件

`middleware.Defaults()` 顺序为：

1. `middleware.Recovery()`：捕获 panic，记录成 `c.Err(err)`，并 `Abort()`
2. `middleware.RequestID()`：生成 `request_id` 并放入 `Ctx`
3. `middleware.Logger()`：打印请求开始/结束日志，包含 method/path/rid/cost/aborted/err


---

## 目录结构

```text
core/        框架抽象：Ctx 接口、HandlerFunc、常量、ResponseWriter 包装
httpx/       HTTP 实现：Engine、Router、HTTP Context（实现 core.Ctx）
middleware/  内置中间件：Recovery、RequestID、Logger
grpcx/       gRPC 预留目录（当前为空）
test/        最小可运行示例（go run .\\test）
docs/        文档目录（如后续补充设计说明/示例）
```

---

## 限制与后续方向

- Router 目前是 map 精确匹配；源码里已标注 `// TODO 改trie树`，后续可升级为 Trie/radix tree。
- `grpcx` 目录目前为空，可以作为未来的 gRPC Engine/Interceptor/Ctx 适配层。

---

