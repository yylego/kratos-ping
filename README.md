[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-ping/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-ping/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-ping)](https://pkg.go.dev/github.com/yylego/kratos-ping)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-ping/main.svg)](https://coveralls.io/github/yylego/kratos-ping?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25%2B-lightgrey.svg)](https://github.com/yylego/kratos-ping)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-ping.svg)](https://github.com/yylego/kratos-ping/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-ping)](https://goreportcard.com/report/github.com/yylego/kratos-ping)

# kratos-ping

Simple and fast Ping service package to support the Kratos framework, with both gRPC and HTTP protocols.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Features

✅ **Simple Design**: Clean and straightforward service structure
✅ **Dual-Protocol Support**: Supports both gRPC and HTTP protocols
✅ **Native Integration**: Seamless integration with the Kratos framework
✅ **Built-in Tests**: Comprehensive test coverage included
✅ **Modern Schema**: Uses Protocol Buffers to define services
✅ **Zero Config**: Out-of-the-box Ping service implementation

## Installation

```bash
go get github.com/yylego/kratos-ping/pingkratos
```

## Usage

### HTTP Example

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
	applog := log.NewLogger(log.NewHandler())

	// Setup HTTP service on port 8000
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			logging.Server(applog.With("caption", "HTTP")),
		),
		http.Timeout(time.Minute),
	)

	// Setup ping service
	pingService := serverpingkratos.NewPingService(applog.With("caption", "PING"))
	clientpingkratos.RegisterPingHTTPServer(httpSrv, pingService)

	// Setup and start application
	app := kratos.New(
		kratos.Name("pingkratos-http-demo"),
		kratos.Server(httpSrv),
	)

	zaplog.LOG.Info("Starting HTTP Ping service", zap.String("address", ":8000"))

	// Run application (awaiting shutdown)
	must.Done(app.Run())
	defer rese.F0(app.Stop)
}
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### gRPC Example

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
	applog := log.NewLogger(log.NewHandler())

	// Setup gRPC service on port 9000
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(applog.With("caption", "GRPC")),
		),
		grpc.Timeout(time.Minute),
	)

	// Setup ping service
	pingService := serverpingkratos.NewPingService(applog.With("caption", "PING"))
	clientpingkratos.RegisterPingServer(grpcSrv, pingService)

	// Setup and start application
	app := kratos.New(
		kratos.Name("pingkratos-grpc-demo"),
		kratos.Server(grpcSrv),
	)

	zaplog.LOG.Info("Starting gRPC Ping service", zap.String("address", ":9000"))

	// Run application (awaiting shutdown)
	must.Done(app.Run())
	defer rese.F0(app.Stop)
}
```

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

## Dependencies

### Core Dependencies

- `github.com/go-kratos/kratos/v3` - Kratos framework
- `google.golang.org/grpc` - gRPC support
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/yylego/*` - Utilities

## Examples

### Integration with Kratos Project

To integrate pingkratos into the Kratos project, follow these steps:

**1. Add Dependency**

```bash
go get github.com/yylego/kratos-ping/pingkratos
```

**2. Setup Wire Provider**

In `internal/service/service.go`:

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

**3. Setup HTTP Endpoint**

In `internal/server/http.go`:

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

**4. Setup gRPC Endpoint**

In `internal/server/grpc.go`:

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

**5. Generate Wire Code**

```bash
wire ./cmd/demo-app/...
```

### Demo Projects

Complete working examples:

- [pingkratos-demos](https://github.com/yylego/kratos-ping-demos) - Complete Kratos project integration

Unit test examples: [TEST](serverpingkratos/ping_test.go).

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/yylego/kratos-ping.svg?variant=adaptive)](https://starchart.cc/yylego/kratos-ping)
