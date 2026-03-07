# 协议说明

## 单端口协议分发

**实现单端口协议分发需要在TCP层对请求进行预读，将流量分发到不同的处理器**

### 预读

TCP本质是**字节流**，没有消息边界

- 正常读取字节会导致**字节在缓冲区丢失**
- 使用`MSG_PEEK`可以实现**读取数据且防止字节丢失**
  - 从`net.Conn`读取部分字节到`bufio.Reader`，这部分字节在`net.Conn`丢失，但是会进入`bufio.Reader`的缓冲区

具体需要对`net.Conn`对象进行包装
- 读取数据时先读取`bufio.Reader`的缓存，再读取`net.Conn`的剩余字节
- 包装后的对象可以给后续应用层协议复用，因为它们只接受`net.Conn`对象
- 需要实现`net.Conn`接口

通过预读的前几个字节，可以判断连接的协议类型：

**HTTP/1.1**
- 纯文本协议。
- 报文以方法名开头：`GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `OPTIONS`, `CONNECT`, `TRACE`。
- 读取前 8 个字节即可进行匹配。

**HTTP/2 (h2c)**
- 二进制协议。
- 定义了 **固定连接前奏**：`PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n`。
- 读取到 `PRI *` 即可确定是 HTTP/2 (用于 gRPC)。

**Websocket**
- 基于 HTTP 协议握手。
- 分发器将其识别为 HTTP/1.1。
- 在 HTTP处理器 (`httpx` 或 `wsx`) 内部解析 Header，若包含 `Upgrade: websocket` 则升级连接。

**SSE (Server-Sent Events)**
- 基于 HTTP 长连接。
- 分发器将其识别为 HTTP/1.1。
- 响应头特征：
  ```
  Content-Type: text/event-stream
  Cache-Control: no-cache
  Connection: keep-alive
  ```

### 连接分发逻辑

底层使用 TCP 进行监听 (`mux.Multiplexer`)。
1.  `Accept()` 获取原始 TCP 连接。
2.  将其封装为 `FuseConn`。
3.  设置 **握手超时**（默认 3 秒），防止恶意连接耗尽资源。
4.  预读字节，进行协议匹配 (`IsHTTP1`, `IsHTTP2`)。
5.  匹配成功后，将 `FuseConn` 推送到对应的虚拟监听器 (`FakeListener`)。

### 虚拟监听器 (Fake Listener)

虚拟监听器是适配 Go 标准库 `net.Listener` 接口的关键组件。
- 它不直接监听端口，而是通过通道 (`chan net.Conn`) 接收来自 `Multiplexer` 分发的连接。
- `http.Server` 或 `grpc.Server` 使用这个虚拟监听器启动服务。
- 当应用层调用 `Serve()` 时，会从虚拟监听器的通道中取出连接进行处理。

### gRPC 集成

Fuse 通过 `grpcx` 包集成了 gRPC，并实现了与 HTTP 处理器风格一致的中间件机制。

**核心机制**
- **Interceptors (拦截器)**: 使用 `UnaryInterceptor` 和 `StreamInterceptor` 拦截所有 gRPC 请求。
- **Context 适配**: 将 `context.Context` 封装为 `core.Ctx` (具体为 `grpcx.Ctx`)，使得 gRPC 方法可以使用与 HTTP 相同的中间件和上下文操作 API。
- **协议共存**: 可以在同一个端口上运行，通过 `mux` 进行流量分发。
