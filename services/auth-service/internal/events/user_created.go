package events

import "time"

type UserRegisteredEvent struct {
	EventType  string    `json:"event_type"`
	EventID    string    `json:"event_id"`
	OccurredAt time.Time `json:"occurred_at"`
	Payload    UserRegisteredPayload `json:"payload"`
}

type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}