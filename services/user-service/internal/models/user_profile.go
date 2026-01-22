package models

import "time"

type UserProfile struct {
	UserID   string `json:"user_id" db:"user_id"`
	Email    string `json:"email" db:"email"`
	Name     string `json:"name" db:"name"`
	Surname  string `json:"surname" db:"surname"`
	Telegram string `json:"telegram" db:"telegram"`
	UserType
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type UserType struct {
	IsStudent bool `json:"is_student" db:"is_student"`
	IsTutor   bool `json:"is_tutor" db:"is_tutor"`
}
