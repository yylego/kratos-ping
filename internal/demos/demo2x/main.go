package main

import (
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/yylego/kratos-ping/pingkratos/clientpingkratos"
	"github.com/yylego/kratos-ping/pingkratos/serverpingkratos"
	"github.com/yylego/kratos-zap/zapkratos"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
)

func main() {
	// Setup logging to show ping request logs
	// 配置日志以显示 ping 请求日志
	zapKratos := zapkratos.NewZapKratos(zaplog.LOGGER, zapkratos.NewOptions())

	// Setup gRPC service on port 9000
	// 在端口 9000 配置 gRPC 服务
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(zapKratos.GetLogger("GRPC")),
		),
		grpc.Timeout(time.Minute),
	)

	// Setup ping service
	// 配置 ping 服务
	pingService := serverpingkratos.NewPingService(zapKratos.GetLogger("PING"))
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
