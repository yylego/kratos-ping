package serverpingkratos_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/middleware/logging"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/kratos-ping/internal/utils"
	"github.com/yylego/kratos-ping/pingkratos/clientpingkratos"
	"github.com/yylego/kratos-ping/pingkratos/serverpingkratos"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	httpPort string // Dynamic HTTP port // 动态分配的 HTTP 端口
	grpcPort string // Dynamic gRPC port // 动态分配的 gRPC 端口
)

func TestMain(m *testing.M) {
	// Create logger to show ping request logs
	// 创建 logger 以显示 ping 请求日志
	applog := log.NewLogger(log.NewHandler())

	// Create HTTP server with dynamic port (port 0 = random available port)
	// 使用动态端口创建 HTTP 服务器（端口 0 表示随机可用端口）
	httpSrv := http.NewServer(
		http.Address(":0"),
		http.Middleware(
			recovery.Recovery(),
			logging.Server(applog.With("caption", "HTTP")),
		),
		http.Timeout(time.Minute),
	)
	httpPort = utils.ExtractPort(rese.P1(httpSrv.Endpoint()))

	// Create gRPC server with dynamic port
	// 使用动态端口创建 gRPC 服务器
	grpcSrv := grpc.NewServer(
		grpc.Address(":0"),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(applog.With("caption", "GRPC")),
		),
		grpc.Timeout(time.Minute),
	)
	grpcPort = utils.ExtractPort(rese.P1(grpcSrv.Endpoint()))

	// Create ping service
	// 创建 ping 服务
	pingService := serverpingkratos.NewPingService(applog.With("caption", "PING"))
	clientpingkratos.RegisterPingHTTPServer(httpSrv, pingService)
	clientpingkratos.RegisterPingServer(grpcSrv, pingService)

	app := kratos.New(
		kratos.Name("test-ping-kratos"),
		kratos.Server(httpSrv, grpcSrv),
	)

	// Start server in background
	// 后台启动服务器
	go func() {
		must.Done(app.Run())
	}()
	defer rese.F0(app.Stop)

	// Wait a short time to ensure the server has started
	// 等待片刻以确保服务器已启动
	time.Sleep(time.Millisecond * 200)

	zaplog.LOG.Info("Starting test servers with dynamic ports",
		zap.String("http_port", httpPort),
		zap.String("grpc_port", grpcPort),
	)

	m.Run()
}

func TestPingService_Ping_HTTP(t *testing.T) {
	// Create HTTP client connecting to dynamic port
	// 创建 HTTP 客户端连接到动态端口
	conn := rese.P1(http.NewClient(
		context.Background(),
		http.WithMiddleware(recovery.Recovery()),
		http.WithEndpoint("127.0.0.1:"+httpPort),
	))
	defer rese.F0(conn.Close)

	// Test ping service via HTTP
	// 通过 HTTP 测试 ping 服务
	pingClient := clientpingkratos.NewPingHTTPClient(conn)
	ctx := context.Background()
	message := uuid.New().String()

	resp, err := pingClient.Ping(ctx, wrapperspb.String(message))
	require.NoError(t, err)
	require.Equal(t, message, resp.GetValue())
}

func TestPingService_Ping_gRPC(t *testing.T) {
	// Create gRPC client connecting to dynamic port
	// 创建 gRPC 客户端连接到动态端口
	conn := rese.P1(grpc.NewClient(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:"+grpcPort),
		grpc.WithMiddleware(recovery.Recovery()),
	))
	defer rese.F0(conn.Close)

	// Test ping service via gRPC
	// 通过 gRPC 测试 ping 服务
	pingClient := clientpingkratos.NewPingClient(conn)
	ctx := context.Background()
	message := uuid.New().String()

	resp, err := pingClient.Ping(ctx, wrapperspb.String(message))
	require.NoError(t, err)
	require.Equal(t, message, resp.GetValue())
}

// NOTE: Add test to handle blank message
// NOTE: 暂不测试空字符串的
