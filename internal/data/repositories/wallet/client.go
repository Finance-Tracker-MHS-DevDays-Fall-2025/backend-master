package wallet

import (
	"context"
	"fmt"

	pb "backend-master/internal/api-gen/proto/wallet"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type WalletClient struct {
	client pb.WalletServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func NewClient(
	address string,
	logger *zap.Logger,
	opts ...grpc.DialOption,
) (*WalletClient, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to wallet service: %w", err)
	}

	return &WalletClient{
		client: pb.NewWalletServiceClient(conn),
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *WalletClient) GetAccounts(
	ctx context.Context,
	req *pb.GetAccountsRequest,
) (*pb.GetAccountsResponse, error) {
	resp, err := c.client.GetAccounts(ctx, req)
	if err != nil {
		c.logger.Error("failed to get accounts", zap.Error(err))
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return resp, nil
}

func (c *WalletClient) GetTransactions(
	ctx context.Context,
	req *pb.GetTransactionsRequest,
) (*pb.GetTransactionsResponse, error) {
	resp, err := c.client.GetTransactions(ctx, req)
	if err != nil {
		c.logger.Error("failed to get transactions", zap.Error(err))
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return resp, nil
}

func (c *WalletClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
