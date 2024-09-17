package client

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetGRPCConn(address string, logger *zap.Logger, interceptors ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptors...),
	)

	if err != nil {
		logger.Fatal("failed to connect to server", zap.String("address", address), zap.Error(err))
	}

	return conn
}
