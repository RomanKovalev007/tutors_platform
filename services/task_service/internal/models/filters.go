package models

type TaskType string

const (
	AssignedType TaskType = "ASSIGNED"
	CreatedType  TaskType = "CREATED"
)

type TaskFilter struct {
	GroupID string
	UserID  string
	Type    TaskType
	Offset  int32
	Limit   int32
}

type SubmissionFilter struct {
	TaskID string
	UserID string
	Offset int32
	Limit  int32
}
