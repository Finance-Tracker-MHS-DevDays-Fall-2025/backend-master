package market

import (
	"context"
	"fmt"

	pb "backend-master/internal/api-gen/proto/market"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MarketClient struct {
	client pb.MarketServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func NewClient(
	address string,
	logger *zap.Logger,
	opts ...grpc.DialOption,
) (*MarketClient, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to market service: %w", err)
	}

	return &MarketClient{
		client: pb.NewMarketServiceClient(conn),
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *MarketClient) GetInvestmentPositions(
	ctx context.Context,
	req *pb.GetInvestmentPositionsRequest,
) (*pb.GetInvestmentPositionsResponse, error) {
	resp, err := c.client.GetInvestmentPositions(ctx, req)
	if err != nil {
		c.logger.Error("failed to get investment positions", zap.Error(err))
		return nil, fmt.Errorf("failed to get investment positions: %w", err)
	}

	return resp, nil
}

func (c *MarketClient) GetSecurity(
	ctx context.Context,
	req *pb.GetSecurityRequest,
) (*pb.GetSecurityResponse, error) {
	resp, err := c.client.GetSecurity(ctx, req)
	if err != nil {
		c.logger.Error("failed to get security", zap.Error(err))
		return nil, fmt.Errorf("failed to get security: %w", err)
	}

	return resp, nil
}

func (c *MarketClient) GetSecuritiesPrices(
	ctx context.Context,
	req *pb.GetSecuritiesPricesRequest,
) (*pb.GetSecuritiesPricesResponse, error) {
	resp, err := c.client.GetSecuritiesPrices(ctx, req)
	if err != nil {
		c.logger.Error("failed to get securities prices", zap.Error(err))
		return nil, fmt.Errorf("failed to get securities prices: %w", err)
	}

	return resp, nil
}

func (c *MarketClient) GetSecurityPayments(
	ctx context.Context,
	req *pb.GetSecuritiesPaymentsRequest,
) (*pb.GetSecuritiesPaymentsResponse, error) {
	resp, err := c.client.GetSecurityPayments(ctx, req)
	if err != nil {
		c.logger.Error("failed to get security payments", zap.Error(err))
		return nil, fmt.Errorf("failed to get security payments: %w", err)
	}

	return resp, nil
}

func (c *MarketClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
