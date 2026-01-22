package models

type GetGroupTasksParams struct {
	GroupID string
	Offset  int32
	Limit   int32
	Status  *TaskStatus
}

type GetAssignedTasksParams struct {
	UserID string
	Offset int32
	Limit  int32
	Status *TaskStatus
}

type CreatedByMeParams struct {
	UserID string
	Offset int32
	Limit  int32
	Status *TaskStatus
}

type GetSubmissionsParams struct {
	TaskID string
	Offset int32
	Limit  int32
	Status *SubmissionStatus
}
