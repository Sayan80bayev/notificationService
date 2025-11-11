package events

import "github.com/google/uuid"

const (
	SubscriptionCreated = "subscription.created"
)

type SubscriptionCreatedPayload struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FolloweeID uuid.UUID `json:"followee_id"`
	CreatedAt  int64     `json:"created_at_unix"`
}
