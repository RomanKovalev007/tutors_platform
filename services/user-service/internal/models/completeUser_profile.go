package models

type CompleteUserProfile struct {
	UserProfile    *UserProfile
	TutorProfile   *TutorProfile
	StudentProfile *StudentProfile
}
