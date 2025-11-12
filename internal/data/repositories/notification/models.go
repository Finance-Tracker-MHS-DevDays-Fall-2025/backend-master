package notification

import (
	"time"

	pb "backend-master/internal/api-gen/proto/notification"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Title     string    `db:"title"`
	Message   string    `db:"message"`
	SentAt    time.Time `db:"sent_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (n *Notification) ToProto() *pb.SendNotificationRequest {
	return &pb.SendNotificationRequest{
		UserId:  n.UserID.String(),
		Title:   n.Title,
		Message: n.Message,
	}
}

