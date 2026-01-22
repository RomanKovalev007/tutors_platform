package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password_hash" db:"password_hash"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
