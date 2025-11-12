package notification

import (
	"backend-master/internal/data/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	CreateNotification(
		ctx context.Context,
		userID uuid.UUID,
		title string,
		message string,
	) (*Notification, error)

	GetNotificationsByUserID(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
	) ([]Notification, error)
}

type notificationRepositoryImpl struct {
	db     database.DBManager
	logger *zap.Logger
}

func NewRepository(
	db database.DBManager,
	logger *zap.Logger,
) NotificationRepository {
	return &notificationRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (repo *notificationRepositoryImpl) CreateNotification(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	message string,
) (*Notification, error) {
	query := `
		INSERT INTO notifications (
			id,
			user_id,
			title,
			message,
			sent_at,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, message, sent_at, created_at
	`

	notification := Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.db.GetDB().GetContext(
		ctx,
		&notification,
		query,
		notification.ID,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.SentAt,
		notification.CreatedAt,
	)
	if err != nil {
		repo.logger.Error(
			"failed to create notification",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

func (repo *notificationRepositoryImpl) GetNotificationsByUserID(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]Notification, error) {
	query := `
		SELECT 
			id,
			user_id,
			title,
			message,
			sent_at,
			created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	var notifications []Notification
	err := repo.db.GetDB().SelectContext(ctx, &notifications, query, userID, limit)
	if err != nil {
		repo.logger.Error(
			"failed to get notifications",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

