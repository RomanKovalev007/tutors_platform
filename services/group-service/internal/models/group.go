package models

import (
	"time"
)

type Group struct {
	ID          string         `json:"id" db:"id"`
	TutorID     string         `json:"tutor_id" db:"tutor_id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	Members     []*GroupMember `json:"members" db:"-"`
}
