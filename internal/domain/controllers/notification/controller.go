package notification

import (
	"context"
	"fmt"

	pb "backend-master/internal/api-gen/proto/notification"
	"backend-master/internal/data/repositories/notification"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationController interface {
	SendNotification(
		ctx context.Context,
		userID string,
		title string,
		message string,
	) (*pb.SendNotificationResponse, error)
}

type notificationControllerImpl struct {
	repo   notification.NotificationRepository
	client *notification.NotificationClient
	logger *zap.Logger
}

func NewController(
	repo notification.NotificationRepository,
	client *notification.NotificationClient,
	logger *zap.Logger,
) NotificationController {
	return &notificationControllerImpl{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (cont *notificationControllerImpl) SendNotification(
	ctx context.Context,
	userID string,
	title string,
	message string,
) (*pb.SendNotificationResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	req := &pb.SendNotificationRequest{
		UserId:  userID,
		Title:   title,
		Message: message,
	}

	resp, err := cont.client.SendNotification(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification  via client: %w", err)
	}

	_, err = cont.repo.CreateNotification(ctx, uid, title, message)
	if err != nil {
		cont.logger.Error(
			"failed to log notification to database",
			zap.Error(err),
		)
	}

	return resp, nil
}
