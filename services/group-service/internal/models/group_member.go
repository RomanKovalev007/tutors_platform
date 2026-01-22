package models

import (
	"time"
)

type GroupMember struct {
	StudentID string    `json:"student_id" db:"student_id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
