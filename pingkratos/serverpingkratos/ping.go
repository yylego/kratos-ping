package serverpingkratos

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/yylego/kratos-ping/pingkratos/clientpingkratos"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// PingService implements the Ping service with configurable logging
// PingService 实现 Ping 服务，支持可配置的日志记录
type PingService struct {
	pb.UnimplementedPingServer
	slog *log.Helper
}

// NewPingService creates a new PingService with the provided logger
// 使用提供的 logger 创建新的 PingService
func NewPingService(logger log.Logger) *PingService {
	return &PingService{
		slog: log.NewHelper(logger),
	}
}

// Ping handles ping requests and returns the same message back
// Ping 处理 ping 请求并返回相同的消息
func (s *PingService) Ping(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	message := req.GetValue()
	s.slog.WithContext(ctx).Debugf("ping service processing message: %s", message)
	return wrapperspb.String(message), nil
}
