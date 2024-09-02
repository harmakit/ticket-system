package interceptor

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"ticket-system/lib/logger"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger.Info("request started", zap.String("method", info.FullMethod), zap.Any("request", req))

	resp, err = handler(ctx, req)
	if err != nil {
		logger.Error("request failed", zap.String("method", info.FullMethod), zap.Error(err))
		return nil, err
	}

	logger.Info("request finished", zap.String("method", info.FullMethod), zap.Any("response", resp))

	return resp, err
}
