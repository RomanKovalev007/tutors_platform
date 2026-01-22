package models

type TaskStatus string

const (
	TaskStatusActive  TaskStatus = "ACTIVE"
	TaskStatusExpired TaskStatus = "EXPIRED"
)

type SubmissionStatus string

const (
	SubmissionStatusPending  SubmissionStatus = "PENDING"
	SubmissionStatusVerified SubmissionStatus = "VERIFIED"
)
