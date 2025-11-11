package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"notificationService/cmd/server/ws"
	"notificationService/internal/events"
	"notificationService/internal/model"
)

// --- Individual event-specific handlers ---

//func handleLikeCreated(evt LikeEvent, logger *logrus.Logger) {
//	logger.Infof("Like Created: ID=%s User=%s Post=%s CreatedAt=%s", evt.ID, evt.RecipientID, evt.PostID, evt.CreatedAt)
//	// TODO: Add business logic, e.g.:
//	// - Store in database
//	// - Trigger sending pipeline (email/SMS)
//}

func HandleSubscriptionCreated(svc NotificationService) func(data []byte) error {
	return func(data []byte) error {
		logger := logging.GetLogger()
		var evt events.SubscriptionCreatedPayload
		err := json.Unmarshal(data, &evt)
		if err != nil {
			return err
		}

		logger.Infof("Subscription Created: Follower=%s Followee=%s CreatedAt=%d", evt.FollowerID, evt.FolloweeID, evt.CreatedAt)
		notification := &model.Notification{
			UserID:  evt.FolloweeID,
			Message: fmt.Sprintf("You have a new follower! %s, %d", evt.FollowerID, evt.CreatedAt),
		}

		res, err := svc.CreateNotification(context.Background(), notification)
		if err != nil {
			return err
		}

		ws.SendNotification(res.UserID, res.Message)
		return nil
	}
}
