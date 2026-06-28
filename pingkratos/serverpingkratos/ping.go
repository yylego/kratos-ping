package serverpingkratos

import (
	"context"
	"log/slog"

	pb "github.com/yylego/kratos-ping/pingkratos/clientpingkratos"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// PingService implements the Ping service with configurable logging
// PingService 实现 Ping 服务，支持可配置的日志记录
type PingService struct {
	pb.UnimplementedPingServer
	applog *slog.Logger
}

// NewPingService creates a new PingService with the provided logger
// 使用提供的 logger 创建新的 PingService
func NewPingService(applog *slog.Logger) *PingService {
	return &PingService{
		applog: applog,
	}
}

// Ping handles ping requests and returns the same message back
// Ping 处理 ping 请求并返回相同的消息
func (s *PingService) Ping(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	message := req.GetValue()
	s.applog.DebugContext(ctx, "ping service processing message", "message", message)
	return wrapperspb.String(message), nil
}
