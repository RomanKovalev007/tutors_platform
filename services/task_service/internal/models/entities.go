package models

import "time"

type AssignedTask struct {
	ID          string
	GroupId     string
	TutorId     string
	Title       string
	Description *string
	MaxScore    int32
	Deadline    time.Time
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type AssignedTaskShort struct {
	ID       string
	GroupID  string
	TutorID  string
	Title    string
	Deadline time.Time
	Status   TaskStatus
}

type SubmittedTask struct {
	ID        string
	TaskID    string
	StudentID string
	Content   string
	Status    SubmissionStatus
	Score     *int32
	Feedback  *string
	CreatedAt time.Time
	UpdatedAt *time.Time
	OverdueBy *time.Duration
}

type SubmittedTaskShort struct {
	ID          string
	TaskID      string
	StudentID   string
	Score       *int32
	Status      SubmissionStatus
	SubmittedAt time.Time
}

type SubmissionGrade struct {
	SubmissionId string
	TutorId      string
	Score        *int32
	Feedback     *string
	Status       *SubmissionStatus
}

type UpdateTaskRequest struct {
	TutorID     string
	TaskID      string
	Title       *string
	Description *string
	MaxScore    *int32
	Deadline    *time.Time
}

type UpdateSubmissionRequest struct {
	UserID       string
	SubmussionID string
	Content      *string
}
