package presentation

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		logger.Info(
			"grpc request started",
			zap.String("method", info.FullMethod),
		)

		res, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			logger.Error(
				"grpc request failed",
				zap.String("method", info.FullMethod),
				zap.Int64("duration_ms", duration.Milliseconds()),
				zap.Error(err),
			)
		} else {
			logger.Info(
				"grpc request completed",
				zap.String("method", info.FullMethod),
				zap.Int64("duration_ms", duration.Milliseconds()),
			)
		}

		return res, err
	}
}

func UnaryClientInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		logger.Info(
			"grpc client request started",
			zap.String("method", method),
		)

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			logger.Error(
				"grpc client request failed",
				zap.String("method", method),
				zap.Int64("duration_ms", duration.Milliseconds()),
				zap.Error(err),
			)
		} else {
			logger.Info(
				"grpc client request completed",
				zap.String("method", method),
				zap.Int64("duration_ms", duration.Milliseconds()),
			)
		}

		return err
	}
}
