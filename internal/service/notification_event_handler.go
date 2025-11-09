package service

import (
	"encoding/json"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/sirupsen/logrus"
)

// LikeEvent represents the structure of the message payload
type LikeEvent struct {
	ID          string `json:"id"`
	PostID      string `json:"post_id"`
	RecipientID string `json:"recipient_id"`
	CreatedAt   string `json:"created_at"`
}

// HandleRabbitEvent handles incoming RabbitMQ events for the notification service
func HandleRabbitEvent(eventType string, data []byte) {
	logger := logging.GetLogger()

	var evt LikeEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		logger.Errorf("❌ Failed to unmarshal notification event (%s): %v", eventType, err)
		return
	}

	switch eventType {
	case "like.created":
		handleLikeCreated(evt, logger)
	default:
		logger.Warnf("⚠️ Unrecognized event type: %s (payload: %s)", eventType, string(data))
	}
}

// --- Individual event-specific handlers ---

func handleLikeCreated(evt LikeEvent, logger *logrus.Logger) {
	logger.Infof("Like Created: ID=%s User=%s Post=%s CreatedAt=%s", evt.ID, evt.RecipientID, evt.PostID, evt.CreatedAt)
	// TODO: Add business logic, e.g.:
	// - Store in database
	// - Trigger sending pipeline (email/SMS)
}
