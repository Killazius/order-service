package middleware

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"order-service/internal/logger"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
	guid := uuid.New().String()
	ctx = context.WithValue(ctx, logger.RequestID, guid)

	logger.GetLogger().Info(ctx, "incoming request",
		zap.String("method", info.FullMethod),
	)

	res, err := next(ctx, req)
	if err != nil {
		logger.GetLogger().Error(ctx, "request failed",
			zap.String("method", info.FullMethod),
			zap.String("request_id", guid),
			zap.Error(err),
		)
	}

	return res, err
}
