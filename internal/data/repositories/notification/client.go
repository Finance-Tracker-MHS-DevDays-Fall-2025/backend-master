package notification

import (
	"context"
	"fmt"

	pb "backend-master/internal/api-gen/proto/notification"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client pb.NotificationServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func NewClient(
	address string,
	logger *zap.Logger,
) (*NotificationClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}

	return &NotificationClient{
		client: pb.NewNotificationServiceClient(conn),
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *NotificationClient) SendNotification(
	ctx context.Context,
	req *pb.SendNotificationRequest,
) (*pb.SendNotificationResponse, error) {
	resp, err := c.client.SendNotification(ctx, req)
	if err != nil {
		c.logger.Error("failed to send notification", zap.Error(err))
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	return resp, nil
}

func (c *NotificationClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

