package models

type TutorProfile struct {
	UserID         string `json:"user_id" db:"user_id"`
	Bio            string `json:"bio" db:"bio"`
	Specialization string `json:"specialization" db:"specialization"`
	Experience     int32  `json:"experience_years" db:"experience_years"`
}
