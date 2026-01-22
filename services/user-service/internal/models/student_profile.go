package models

type StudentProfile struct {
	UserID string `json:"user_id" db:"user_id"`
	Grade  string `json:"grade" db:"grade"`
	Bio    string `json:"bio" db:"bio"`
}
