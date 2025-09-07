# 博客系统后端（go-zero + etcd）

本项目基于 go-zero 框架，从 0-1 搭建了一个“用户（user）+ 博客（blog）+ HTTP 网关（gateway）”的后端系统，使用 etcd 进行服务注册发现，完成用户注册/登录、文章 CRUD、评论的基础能力，并提供端到端测试示例。

# 原理
goctl 生成 API+RPC 框架；zrpc 用 etcd 做服务注册发现；
gateway 使用 zrpc.RpcClientConf 连接 user/blog；
错误处理：handler 用 httpx.ErrorCtx 统一返回；logic 返回 error 触发 HTTP 4xx/5xx。
<!-- 简化的认证流程：token=“userId|timestamp”；Verify 解析 token 前段为 userId；
内存数据库代替真实存储；结构清晰，可平滑替换为 MySQL/Redis； -->

## 环境与依赖
- Go >= 1.23
- protoc、protoc-gen-go、protoc-gen-go-grpc、goctl
- etcd（用于注册发现）

安装：
```bash
# 安装 goctl + protoc 插件（GOBIN 建议在 PATH 中）
go install github.com/zeromicro/go-zero/tools/goctl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# 验证（确保这些都在 PATH 中）
which goctl && which protoc &&  which protoc-gen-go && which protoc-gen-go-grpc
ls -l /Users/xuchenyao/go/bin | cat

# 初始化 task4_gozero 子模块
go mod init gotask/task4_gozero
# 修正环境变量
export PATH=/Users/xuchenyao/go/bin:$PATH
export GOBIN=/Users/xuchenyao/go/bin && export GO111MODULE=on

#goctl 生成 user/blog RPC 和 gateway API 服务骨架
mkdir -p rpc gateway 
cd rpc
goctl rpc new user
goctl rpc new blog
cd gateway
goctl api new gateway

#更新 proto 与 api 定义，并用 goctl 重新生成代码
cd rpc/user && goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
cd rpc/blog && goctl rpc protoc blog.proto --go_out=. --go-grpc_out=. --zrpc_out=.
cd gateway && goctl api go -api gateway.api -dir .
实现 user/blog 的内存逻辑与 JWT，并配置 etcd 和端口
```

## 服务启动
```bash
# 安装并启动 etcd（macOS brew 示例）
brew install etcd
mkdir -p /tmp/etcd-data
nohup etcd --data-dir /tmp/etcd-data \
           --listen-client-urls http://127.0.0.1:2379 \
           --advertise-client-urls http://127.0.0.1:2379 \
           >/tmp/etcd.log 2>&1 & 
pgrep etcd
pkill etcd        # 或者
kill <PID>

# 确保 etcd 已启动后，分别启动 user、blog、gateway：
cd rpc/user && go run user.go -f etc/user.yaml >/tmp/user.log 2>&1 &
cd rpc/blog && go run blog.go -f etc/blog.yaml >/tmp/blog.log 2>&1 &
cd gateway && ./gateway -f etc/gateway-api.yaml >/tmp/gw.log 2>&1 &
lsof -iTCP:8081 -sTCP:LISTEN
lsof -iTCP:8082 -sTCP:LISTEN
lsof -iTCP:8888 -sTCP:LISTEN
```

## API 与测试（curl）
- 注册用户
```bash
curl -s -X POST http://127.0.0.1:8888/api/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"u3","password":"p3","email":"u3@a.com"}'
# 响应示例：
# {"userId":1003}
```

- 登录获取 token
```bash
curl -s -X POST http://127.0.0.1:8888/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"u3","password":"p3"}'
# 响应示例：
# {"token":"1003|1757186387","userId":1003}
```

- 创建文章（需要 Authorization）
```bash
curl -s -X POST http://127.0.0.1:8888/api/posts \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer 1003|1757186387' \
  -d '{"title":"t5","content":"c5"}'
# 响应示例：
# {"id":2001,"title":"t5","content":"c5","userId":1003,"createdAt":1757186387,"updatedAt":1757186387}
```

- 列表文章
```bash
curl -s http://127.0.0.1:8888/api/posts
# 响应示例：
# {"posts":[{"id":2001,"title":"t5","content":"c5","userId":1003,"createdAt":1757186387,"updatedAt":1757186387}]}

下面以调用列表接口
curl -s http://127.0.0.1:8888/api/posts
为例，按调用链逐步说明数据流与数据结构变化。
1) HTTP 网关接收请求
路由匹配: gateway/internal/handler/routes.go 将 GET /api/posts 路由到 ListPostsHandler.
参数解析: ListPostsHandler 使用 httpx.Parse 将 query 绑定为 types.ListPostsReq{Page, Size}（若未传则使用默认值，见 gateway.api）。
进入逻辑层: 构造 logic.NewListPostsLogic(ctx, svcCtx)，调用 ListPosts(req)。
2) 网关逻辑调用 RPC（Blog）
位置: gateway/internal/logic/listpostslogic.go
请求对象映射:
入参 types.ListPostsReq{Page, Size} → RPC 入参 blogclient.ListPostsRequest{Page, Size}
服务发现与调用:
svcCtx.BlogRpc 是 blogclient.Blog 客户端，初始化于 gateway/internal/svc/servicecontext.go，通过 zrpc.MustNewClient(c.BlogRpc) 读取 etc/gateway-api.yaml 的 Etcd 配置（Hosts: 127.0.0.1:2379, Key: blog.rpc），在 etcd 上发现 blog 服务并发起 gRPC 调用。
gRPC 序列化: 将 ListPostsRequest 序列化，通过 gRPC 发给 blog rpc 服务。
3) Blog RPC 服务端接收与分发
位置: rpc/blog/internal/server/blogserver.go
分发: gRPC 框架将 ListPosts 请求分发到 logic.NewListPostsLogic(ctx, svcCtx).ListPosts(in)。
4) Blog 逻辑层读取内存数据
位置: rpc/blog/internal/logic/listpostslogic.go
数据结构（只读，不修改）:
svcCtx.Posts map[int64]*PostRecord（在 rpc/blog/internal/svc/servicecontext.go 初始化）
处理流程:
加读锁 RLock
遍历 Posts，将每条 *PostRecord{Id, Title, Content, UserId, CreatedAt, UpdatedAt} 映射为 *blog.Post{...}
去锁 RUnlock
组装 blog.ListPostsResponse{Posts: []*blog.Post}
返回: gRPC 将 ListPostsResponse 序列化返回网关。
5) 网关逻辑层映射响应
位置: gateway/internal/logic/listpostslogic.go
映射:
RPC 返回 *blogclient.ListPostsResponse{Posts: []*blogclient.Post}
转为 HTTP 出参 types.ListPostsResp{Posts: []types.PostResp}，字段一一对应（Id/Title/Content/UserId/CreatedAt/UpdatedAt）。
返回给 handler。
6) 网关 handler 输出 HTTP 响应
位置: gateway/internal/handler/listpostshandler.go
输出: httpx.OkJsonCtx 将 types.ListPostsResp 编码为 JSON，HTTP 200 返回。
最终返回示例（与你实际结果一致）:
{"posts":[{"id":2001,"title":"t5","content":"c5","userId":1003,"createdAt":1757186387,"updatedAt":1757186387}]}
数据结构变化小结
网关侧: 仅做对象映射（types ↔ blogclient），不改内部状态。
Blog 服务侧: 本次为只读操作，svcCtx.Posts 不发生变更（加读锁遍历）。
etcd: 仅用于服务发现，不涉及数据变化。
涉及的关键文件/函数
路由: gateway/internal/handler/routes.go → GET /api/posts → ListPostsHandler
Handler: gateway/internal/handler/listpostshandler.go
网关逻辑: gateway/internal/logic/listpostslogic.go（调用 svcCtx.BlogRpc.ListPosts 并做映射）
RPC 客户端注入: gateway/internal/svc/servicecontext.go（blogclient.Blog）
Blog 逻辑: rpc/blog/internal/logic/listpostslogic.go（从 svcCtx.Posts 读取并返回）
Blog 内存模型: rpc/blog/internal/svc/servicecontext.go（Posts map[int64]*PostRecord）
```

> 说明：评论的创建与查询接口也已生成（/api/posts/:id/comments），可按以上模式添加 curl 测试。


## 项目结构

```
task4_gozero/
├─ task4.md                      # 需求说明
├─ README.md                     # 本文档
├─ go.mod / go.sum
├─ rpc/
│  ├─ user/
│  │  ├─ user.proto              # Register/Login/Verify
│  │  ├─ user.go                 # RPC 服务入口
│  │  ├─ etc/user.yaml           # 配置（ListenOn/Etcd/JwtSecret）
│  │  └─ internal/
│  │     ├─ config/config.go     # zrpc.RpcServerConf (+ JwtSecret)
│  │     ├─ svc/servicecontext.go# 内存用户表 Users / 互斥锁
│  │     └─ logic/               # 注册/登录/验证 实现
│  └─ blog/
│     ├─ blog.proto              # Post/Comment CRUD/列表
│     ├─ blog.go                 # RPC 服务入口
│     ├─ etc/blog.yaml           # 配置（ListenOn/Etcd）
│     └─ internal/
│        ├─ config/config.go
│        ├─ svc/servicecontext.go# 内存 Posts/Comments / 互斥锁
│        └─ logic/               # 文章/评论逻辑实现
└─ gateway/
   ├─ gateway.api                # HTTP API 定义
   ├─ gateway.go                 # HTTP 服务入口
   ├─ etc/gateway-api.yaml       # RestConf + UserRpc/BlogRpc(Etcd) + JwtSecret
   └─ internal/
      ├─ config/config.go        # RestConf + zrpc.RpcClientConf
      ├─ svc/servicecontext.go   # user/blog RPC 客户端注入
      ├─ handler/                # 路由、参数解析、Header 注入
      └─ logic/                  # 转发到 user/blog RPC 的业务逻辑
```

二、数据转发流程（HTTP → Gateway → RPC）
注册:
HTTP POST /api/register -> gateway/internal/handler/registerhandler.go 解析 JSON
调用 gateway/internal/logic/registerlogic.go -> svcCtx.UserRpc.Register
user/internal/logic/registerlogic.go 更新内存 Users，返回 userId
登录:
HTTP POST /api/login -> gateway logic/loginlogic.go -> svcCtx.UserRpc.Login
user/internal/logic/loginlogic.go 校验密码，生成简单 token "userId|unix"
发文:
HTTP POST /api/posts，Header: Authorization: Bearer <token>
handler 将 Authorization 存入 context -> logic/createpostlogic.go
解析 Bearer token -> 调 user.Verify -> 获取 userId -> 调 blog.CreatePost，返回 Post
列表:
HTTP GET /api/posts -> gateway logic/listpostslogic.go -> 调 blog.ListPosts
blog/internal/logic/listpostslogic.go 读取内存 Posts 返回列表
其他（更新/删除/评论）同理（作者校验等在 blog 逻辑层完成）

三、关键数据结构变化（典型调用链）
注册:
入参 RegisterReq{username,password,email}
user 服务 ServiceContext.Users map 更新：插入 UserRecord{Id: auto-increment, ...}
返回 RegisterResp{userId}
登录:
校验 Users[username].Password
生成 token="userId|timestamp"，返回 LoginResp{token,userId}
发文:
从 ctx 读取 "Authorization" header，提取 Bearer token
调 User.Verify 返回 VerifyResponse{ok,userId,username}
blog ServiceContext.Posts[id] = PostRecord{...}，返回 CreatePostResponse{Post}
列表:
遍历 Posts -> 转换为 []blog.Post -> 返回 ListPostsResponse{posts}


## 关键设计与数据流

- 注册/登录（user.rpc）：
  - Register: 写入内存 Users（username->UserRecord），返回自增 userId。
  - Login: 校验密码，返回简单 token（格式：`userId|timestamp`）。
  - Verify: 解析 token 前半段 userId，返回认证结果。

- 文章/评论（blog.rpc）：
  - Post：Create/Get/List/Update/Delete，全在内存结构 `Posts` 中操作。
  - Comment：Create/List，存储在 `Comments[postId]` 中。

- 网关（gateway）：
  - HTTP -> logic：
    - 解析 JSON/Path；
    - 从 Header 注入 `Authorization: Bearer <token>` 到 context；
    - 调用 user.Verify 获取 userId；
    - 转发到 blog/user RPC 对应接口，返回结果。



## 配置说明
示例（关键片段）：
- rpc/user/etc/user.yaml
```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8081
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: user.rpc
JwtSecret: "dev_secret_key"
```

- rpc/blog/etc/blog.yaml
```yaml
Name: blog.rpc
ListenOn: 0.0.0.0:8082
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: blog.rpc
```

- gateway/etc/gateway-api.yaml
```yaml
Name: gateway-api
Host: 0.0.0.0
Port: 8888
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
BlogRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: blog.rpc
JwtSecret: "dev_secret_key"
```


## 数据结构（内存实现）
- 用户
```go
// svc/user
type UserRecord struct {
    Id       int64
    Username string
    Password string // 演示明文，生产请使用哈希
    Email    string
}
Users map[string]*UserRecord // username -> record
```

- 文章/评论
```go
// svc/blog
type PostRecord struct {
    Id        int64
    Title     string
    Content   string
    UserId    int64
    CreatedAt int64
    UpdatedAt int64
}
Posts map[int64]*PostRecord

type CommentRecord struct {
    Id        int64
    Content   string
    UserId    int64
    PostId    int64
    CreatedAt int64
}
Comments map[int64][]*CommentRecord // postId -> comments
```

## 原理与后续改进
- goctl 生成项目骨架 + zrpc 使用 etcd 做服务注册发现；
- 网关通过 zrpc.RpcClientConf 访问 user、blog 两个 RPC；
- 认证流程目前为简化 token（`userId|timestamp`），可替换为标准 JWT（HS256）并添加过期校验、中间件拦截；
- 内存存储用于演示，可平滑替换为 MySQL/Redis（可使用 goctl model 生成模型层）；
- 错误处理与日志：go-zero 内置 httpx.ErrorCtx、logx、trace/stat 指标，便于观测与追踪。

## 常见问题
- 8888 端口被占用：
  - 使用 `lsof -iTCP:8888 -sTCP:LISTEN` 定位进程并结束后重启网关；
- etcd 未启动/连接失败：
  - 确认 127.0.0.1:2379 可用，查看 `/tmp/etcd.log` 日志；
- 插件缺失：
  - 确保 `goctl`、`protoc-gen-go`、`protoc-gen-go-grpc` 在 PATH。

---

如需接入真实数据库/鉴权中间件/单元测试或补充 Swagger 文档，可在此基础上继续扩展。 