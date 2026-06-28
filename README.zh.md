[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-ping/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-ping/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-ping)](https://pkg.go.dev/github.com/yylego/kratos-ping)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-ping/main.svg)](https://coveralls.io/github/yylego/kratos-ping?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25%2B-lightgrey.svg)](https://github.com/yylego/kratos-ping)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-ping.svg)](https://github.com/yylego/kratos-ping/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-ping)](https://goreportcard.com/report/github.com/yylego/kratos-ping)

# kratos-ping

简单快速的 Kratos 框架 Ping 服务包，同时支持 gRPC 和 HTTP 协议类型。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 功能特性

✅ **简洁设计**: 清晰直观的服务结构
✅ **双协议支持**: 同时支持 gRPC 和 HTTP 协议
✅ **原生集成**: 与 Kratos 框架无缝集成
✅ **内置测试**: 包含完整的测试覆盖
✅ **现代架构**: 使用 Protocol Buffers 定义服务
✅ **无需配置**: 开箱即用的 Ping 服务实现

## 安装

```bash
go get github.com/yylego/kratos-ping/pingkratos
```

## 使用方法

### HTTP 示例

```go
package main

import (
	"time"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/middleware/logging"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/yylego/kratos-ping/clientpingkratos"
	"github.com/yylego/kratos-ping/serverpingkratos"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
)

func main() {
	// Setup logging to show ping request logs
	// 配置日志以显示 ping 请求日志
	applog := log.NewLogger(log.NewHandler())

	// Setup HTTP service on port 8000
	// 在端口 8000 配置 HTTP 服务
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			logging.Server(applog.With("caption", "HTTP")),
		),
		http.Timeout(time.Minute),
	)

	// Setup ping service
	// 配置 ping 服务
	pingService := serverpingkratos.NewPingService(applog.With("caption", "PING"))
	clientpingkratos.RegisterPingHTTPServer(httpSrv, pingService)

	// Setup and start application
	// 配置并启动应用
	app := kratos.New(
		kratos.Name("pingkratos-http-demo"),
		kratos.Server(httpSrv),
	)

	zaplog.LOG.Info("Starting HTTP Ping service", zap.String("address", ":8000"))

	// Run application (awaiting shutdown)
	// 运行应用（等待关闭）
	must.Done(app.Run())
	defer rese.F0(app.Stop)
}
```

⬆️ **源代码:** [Source](internal/demos/demo1x/main.go)

### gRPC 示例

```go
package main

import (
	"time"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/middleware/logging"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/yylego/kratos-ping/clientpingkratos"
	"github.com/yylego/kratos-ping/serverpingkratos"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
)

func main() {
	// Setup logging to show ping request logs
	// 配置日志以显示 ping 请求日志
	applog := log.NewLogger(log.NewHandler())

	// Setup gRPC service on port 9000
	// 在端口 9000 配置 gRPC 服务
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(applog.With("caption", "GRPC")),
		),
		grpc.Timeout(time.Minute),
	)

	// Setup ping service
	// 配置 ping 服务
	pingService := serverpingkratos.NewPingService(applog.With("caption", "PING"))
	clientpingkratos.RegisterPingServer(grpcSrv, pingService)

	// Setup and start application
	// 配置并启动应用
	app := kratos.New(
		kratos.Name("pingkratos-grpc-demo"),
		kratos.Server(grpcSrv),
	)

	zaplog.LOG.Info("Starting gRPC Ping service", zap.String("address", ":9000"))

	// Run application (awaiting shutdown)
	// 运行应用（等待关闭）
	must.Done(app.Run())
	defer rese.F0(app.Stop)
}
```

⬆️ **源代码:** [Source](internal/demos/demo2x/main.go)

## 依赖项

### 核心依赖

- `github.com/go-kratos/kratos/v3` - Kratos 框架
- `google.golang.org/grpc` - gRPC 支持
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/yylego/*` - 实用工具包

## 示例

### 集成到 Kratos 项目

要将 pingkratos 集成到 Kratos 项目中，请按照以下步骤操作：

**1. 添加依赖**

```bash
go get github.com/yylego/kratos-ping/pingkratos
```

**2. 配置 Wire 提供者**

在 `internal/service/service.go` 中：

```go
import (
    "github.com/google/wire"
    "github.com/yylego/kratos-ping/serverpingkratos"
)

var ProviderSet = wire.NewSet(
    NewGreeterService,
    serverpingkratos.NewPingService,
)
```

**3. 配置 HTTP 端点**

在 `internal/server/http.go` 中：

```go
import (
    "github.com/yylego/kratos-ping/clientpingkratos"
    "github.com/yylego/kratos-ping/serverpingkratos"
)

func NewHTTPServer(
    c *conf.Server,
    greeter *service.GreeterService,
    pingService *serverpingkratos.PingService,
    applog *slog.Logger,
) *http.Server {
    srv := http.NewServer(opts...)
    v1.RegisterGreeterHTTPServer(srv, greeter)
    clientpingkratos.RegisterPingHTTPServer(srv, pingService)
    return srv
}
```

**4. 配置 gRPC 端点**

在 `internal/server/grpc.go` 中：

```go
import (
    "github.com/yylego/kratos-ping/clientpingkratos"
    "github.com/yylego/kratos-ping/serverpingkratos"
)

func NewGRPCServer(
    c *conf.Server,
    greeter *service.GreeterService,
    pingService *serverpingkratos.PingService,
    applog *slog.Logger,
) *grpc.Server {
    srv := grpc.NewServer(opts...)
    v1.RegisterGreeterServer(srv, greeter)
    clientpingkratos.RegisterPingServer(srv, pingService)
    return srv
}
```

**5. 生成 Wire 代码**

```bash
wire ./cmd/demo-app/...
```

### 演示项目

完整的可运行示例：

- [pingkratos-demos](https://github.com/yylego/kratos-ping-demos) - 完整的 Kratos 项目集成示例

单元测试示例：[TEST](serverpingkratos/ping_test.go)。

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![标星点赞](https://starchart.cc/yylego/kratos-ping.svg?variant=adaptive)](https://starchart.cc/yylego/kratos-ping)
