package main

import (
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
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

	// Setup HTTP service on port 8000
	// 在端口 8000 配置 HTTP 服务
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			logging.Server(zapKratos.GetLogger("HTTP")),
		),
		http.Timeout(time.Minute),
	)

	// Setup ping service
	// 配置 ping 服务
	pingService := serverpingkratos.NewPingService(zapKratos.GetLogger("PING"))
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
