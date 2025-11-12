package analyzer

import (
	"context"
	"fmt"

	pb "backend-master/internal/api-gen/proto/analyzer"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalyzerClient struct {
	client pb.AnalyzerServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func NewClient(
	address string,
	logger *zap.Logger,
) (*AnalyzerClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to analyzer service: %w", err)
	}

	return &AnalyzerClient{
		client: pb.NewAnalyzerServiceClient(conn),
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *AnalyzerClient) GetStatistics(
	ctx context.Context,
	req *pb.GetStatisticsRequest,
) (*pb.GetStatisticsResponse, error) {
	resp, err := c.client.GetStatistics(ctx, req)
	if err != nil {
		c.logger.Error("failed to get statistics", zap.Error(err))
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return resp, nil
}

func (c *AnalyzerClient) GetForecast(
	ctx context.Context,
	req *pb.GetForecastRequest,
) (*pb.GetForecastResponse, error) {
	resp, err := c.client.GetForecast(ctx, req)
	if err != nil {
		c.logger.Error("failed to get forecast", zap.Error(err))
		return nil, fmt.Errorf("failed to get forecast: %w", err)
	}

	return resp, nil
}

func (c *AnalyzerClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

