package service

import (
	"context"
	"log"
	"time"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"

	"task_service/internal/models"
	"google.golang.org/grpc/codes"
)

type Repository interface {
	// Task operations
	CreateTask(ctx context.Context, task models.AssignedTask) (*models.AssignedTask, error)
	UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.AssignedTask, error)
	SoftDeleteTask(ctx context.Context, userID, taskID string) error
	GetTaskByID(ctx context.Context, taskID string) (*models.AssignedTask, error)
	GetTasks(ctx context.Context, filter models.TaskFilter) ([]*models.AssignedTaskShort, int32, error)

	// Submission operations
	CreateSubmission(ctx context.Context, submission models.SubmittedTask) (*models.SubmittedTask, error)
	UpdateSubmission(ctx context.Context, req models.UpdateSubmissionRequest) (*models.SubmittedTask, error)
	DeleteSubmission(ctx context.Context, userID, submissionID string) error
	GradeSubmission(ctx context.Context, grade models.SubmissionGrade) (*models.SubmittedTask, error)
	ResetGrade(ctx context.Context, userID, submissionID string) error
	GetSubmissionByID(ctx context.Context, submissionID string) (*models.SubmittedTask, error)
	GetSubmissions(ctx context.Context, filter models.SubmissionFilter) ([]*models.SubmittedTaskShort, int32, error)

	// Utility
	MarkExpiredTasks(ctx context.Context) error
}

type GroupClient interface {
	GetGroupInfo(ctx context.Context, groupID string) (*pb.Group, error)
	GetGroupMembers(ctx context.Context, groupID string) ([]*pb.GroupMember, error)
}

type Service struct {
	repo        Repository
	groupClient GroupClient
	// logger      *log.Logger
}

func NewService(repo Repository, client GroupClient) *Service {
	return &Service{
		repo:        repo,
		groupClient: client,
	}
}

func (s *Service) StartExpiryChecker(ctx context.Context) {
	go func() {
		const defaultInterval = 1 * time.Minute
		ticker := time.NewTicker(defaultInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Expiry checker stopped")
				return
			case <-ticker.C:
				s.checkExpiredTasks(ctx)
			}
		}
	}()
}

func (s *Service) checkExpiredTasks(ctx context.Context) {
	err := s.repo.MarkExpiredTasks(ctx)
	if err != nil {
		log.Printf("[ERROR] MarkExpiredTasks failed: %v", err)
	}
}

func (s *Service) CreateTask(ctx context.Context, task models.AssignedTask) (*models.AssignedTask, *models.Error) {
	if task.Deadline.Before(time.Now()) {
		log.Printf("[VALIDATION] CreateTask: deadline in the past, taskID: %s", task.ID)
		return nil, &models.Error{
			Code:    codes.InvalidArgument,
			Message: "deadline must be in the future",
		}
	}

	if task.MaxScore < 0 {
		log.Printf("[VALIDATION] CreateTask: negative maxScore %d, taskID: %s", task.MaxScore, task.ID)
		return nil, &models.Error{
			Code:    codes.InvalidArgument,
			Message: "max_score must be positive",
		}
	}

	group, err := s.groupClient.GetGroupInfo(ctx, task.GroupId)
	if err != nil {
		log.Printf("[GROUP_SERVICE] GetGroupInfo failed: %v, groupID: %s", err, task.GroupId)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "group_service error",
		}
	}

	if group.TutorId != task.TutorId {
		log.Printf("[VALIDATION] CreateTask: tutor mismatch, taskTutor: %s, groupTutor: %s",
			task.TutorId, group.TutorId)
		return nil, &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "user is not tutor in this group",
		}
	}

	task.Status = models.TaskStatusActive
	task.CreatedAt = time.Now()

	createdTask, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		log.Printf("[REPOSITORY] CreateTask failed: %v, taskID: %s", err, task.ID)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "failed to create task",
		}
	}

	log.Printf("[SUCCESS] Task created: %s for group: %s", createdTask.ID, createdTask.GroupId)
	return createdTask, nil
}

func (s *Service) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.AssignedTask, *models.Error) {
	if req.MaxScore != nil && *req.MaxScore < 0 {
		log.Printf("[VALIDATION] CreateTask: negative maxScore %d, taskID: %s", req.MaxScore, req.TaskID)
		return nil, &models.Error{
			Code:    codes.InvalidArgument,
			Message: "max_score must be positive",
		}
	}

	currentTask, err := s.repo.GetTaskByID(ctx, req.TaskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s", err, req.TaskID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if currentTask == nil {
		log.Printf("[VALIDATION] Task not found for update: %s", req.TaskID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if currentTask.Status == models.TaskStatusExpired {
		log.Printf("[VALIDATION] Cannot update expired task: %s", req.TaskID)
		return nil, &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot update expired task",
		}
	}

	if req.TutorID != currentTask.TutorId {
		log.Printf("[PERMISSION] UpdateTask denied: reqTutor=%s, taskTutor=%s",
			req.TutorID, currentTask.TutorId)
		return nil, &models.Error{
			Code:    codes.PermissionDenied,
			Message: "only task creator can update it",
		}
	}

	if req.Deadline != nil && req.Deadline.Before(time.Now()) {
		log.Printf("[VALIDATION] New deadline in the past: %v, taskID: %s", req.Deadline, req.TaskID)
		return nil, &models.Error{
			Code:    codes.InvalidArgument,
			Message: "new deadline must be in the future",
		}
	}

	updatedTask, err := s.repo.UpdateTask(ctx, req)
	if err != nil {
		log.Printf("[REPOSITORY] UpdateTask failed: %v, taskID: %s", err, req.TaskID)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "failed to update task",
		}
	}

	log.Printf("[SUCCESS] Task updated: %s", req.TaskID)
	return updatedTask, nil
}

func (s *Service) DeleteTask(ctx context.Context, userID, taskID string) *models.Error {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s", err, taskID)
		return &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if task == nil {
		log.Printf("[VALIDATION] Task not found for deletion: %s", taskID)
		return &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if task.TutorId != userID {
		log.Printf("[PERMISSION] DeleteTask denied: user=%s, taskTutor=%s", userID, task.TutorId)
		return &models.Error{
			Code:    codes.PermissionDenied,
			Message: "only task creator can delete it",
		}
	}

	submissions, _, err := s.repo.GetSubmissions(ctx, models.SubmissionFilter{
		TaskID: taskID,
		Limit:  1,
	})
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissions failed: %v, taskID: %s", err, taskID)
		// Продолжаем, это не критично
	}

	if len(submissions) > 0 {
		log.Printf("[VALIDATION] Cannot delete task with submissions: %s, count: %d",
			taskID, len(submissions))
		return &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot delete task with submissions",
		}
	}

	if err := s.repo.SoftDeleteTask(ctx, userID, taskID); err != nil {
		log.Printf("[REPOSITORY] SoftDeleteTask failed: %v, taskID: %s, userID: %s",
			err, taskID, userID)
		return &models.Error{
			Code:    codes.Internal,
			Message: "failed to delete task",
		}
	}

	log.Printf("[SUCCESS] Task soft deleted: %s by user: %s", taskID, userID)
	return nil
}

func (s *Service) GetTask(ctx context.Context, taskID string) (*models.AssignedTask, *models.Error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s", err, taskID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if task == nil {
		log.Printf("[VALIDATION] Task not found: %s", taskID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	return task, nil
}

func (s *Service) GetGroupTasks(ctx context.Context, params models.GetGroupTasksParams) ([]*models.AssignedTaskShort, int32, *models.Error) {
	filter := models.TaskFilter{
		GroupID: params.GroupID,
		Type:    models.AssignedType,
		Offset:  params.Offset,
		Limit:   params.Limit,
	}

	tasks, total, err := s.repo.GetTasks(ctx, filter)
	if err != nil {
		log.Printf("[REPOSITORY] GetTasks failed: %v, groupID: %s", err, params.GroupID)
		return nil, 0, &models.Error{
			Code:    codes.Internal,
			Message: "failed to get group tasks",
		}
	}

	log.Printf("[INFO] GetGroupTasks: group=%s, found=%d, total=%d",
		params.GroupID, len(tasks), total)
	return tasks, total, nil
}

func (s *Service) GetCreatedByMeTasks(ctx context.Context, params models.CreatedByMeParams) ([]*models.AssignedTaskShort, int32, *models.Error) {
	filter := models.TaskFilter{
		UserID: params.UserID,
		Type:   models.CreatedType,
		Offset: params.Offset,
		Limit:  params.Limit,
	}

	tasks, total, err := s.repo.GetTasks(ctx, filter)
	if err != nil {
		log.Printf("[REPOSITORY] GetTasks failed: %v, userID: %s", err, params.UserID)
		return nil, 0, &models.Error{
			Code:    codes.Internal,
			Message: "failed to get created tasks",
		}
	}

	log.Printf("[INFO] GetCreatedByMeTasks: user=%s, found=%d, total=%d",
		params.UserID, len(tasks), total)
	return tasks, total, nil
}

func (s *Service) CreateSubmission(ctx context.Context, submission models.SubmittedTask) (*models.SubmittedTask, *models.Error) {
	task, err := s.repo.GetTaskByID(ctx, submission.TaskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s", err, submission.TaskID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if task == nil {
		log.Printf("[VALIDATION] Task not found for submission: %s", submission.TaskID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "task not found",
		}
	}

	if task.Status == models.TaskStatusExpired {
		log.Printf("[VALIDATION] Cannot submit to expired task: %s", submission.TaskID)
		return nil, &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot submit to expired task",
		}
	}

	if time.Now().After(task.Deadline) {
		log.Printf("[VALIDATION] Task deadline passed: %s, deadline: %v",
			submission.TaskID, task.Deadline)
		return nil, &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "task deadline has passed",
		}
	}

	members, err := s.groupClient.GetGroupMembers(ctx, task.GroupId)
	if err != nil {
		log.Printf("[GROUP_SERVICE] GetGroupMembers failed: %v, groupID: %s", err, task.GroupId)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "group_service error",
		}
	}

	isStudent := false
	for _, member := range members {
		if member.StudentId == submission.StudentID {
			isStudent = true
			break
		}
	}

	if !isStudent {
		log.Printf("[VALIDATION] User not in group: student=%s, group=%s",
			submission.StudentID, task.GroupId)
		return nil, &models.Error{
			Code:    codes.InvalidArgument,
			Message: "user is not student in this group",
		}
	}

	submission.Status = models.SubmissionStatusPending
	submission.CreatedAt = time.Now()

	createdSubmission, err := s.repo.CreateSubmission(ctx, submission)
	if err != nil {
		log.Printf("[REPOSITORY] CreateSubmission failed: %v, taskID: %s, studentID: %s",
			err, submission.TaskID, submission.StudentID)

		if err == models.ErrAlreadyGraded {
			return nil, &models.Error{
				Code:    codes.FailedPrecondition,
				Message: "cannot update verified submission",
			}
		}
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "failed to create submission",
		}
	}

	log.Printf("[SUCCESS] Submission created: %s for task: %s by student: %s",
		createdSubmission.ID, submission.TaskID, submission.StudentID)
	return createdSubmission, nil
}

func (s *Service) UpdateSubmission(ctx context.Context, req models.UpdateSubmissionRequest) (*models.SubmittedTask, *models.Error) {
	currentSubmission, err := s.repo.GetSubmissionByID(ctx, req.SubmussionID)
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissionByID failed: %v, submissionID: %s",
			err, req.SubmussionID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if currentSubmission == nil {
		log.Printf("[VALIDATION] Submission not found: %s", req.SubmussionID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if currentSubmission.StudentID != req.UserID {
		log.Printf("[PERMISSION] UpdateSubmission denied: reqUser=%s, submissionStudent=%s",
			req.UserID, currentSubmission.StudentID)
		return nil, &models.Error{
			Code:    codes.PermissionDenied,
			Message: "only submission creator can update it",
		}
	}

	if currentSubmission.Status == models.SubmissionStatusVerified {
		log.Printf("[VALIDATION] Cannot update verified submission: %s", req.SubmussionID)
		return nil, &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot update verified submission",
		}
	}

	task, err := s.repo.GetTaskByID(ctx, currentSubmission.TaskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s",
			err, currentSubmission.TaskID)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "task not found",
		}
	}

	if time.Now().After(task.Deadline) {
		log.Printf("[VALIDATION] Cannot update after deadline: %s, deadline: %v",
			req.SubmussionID, task.Deadline)
		return nil, &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot update submission after deadline",
		}
	}

	updatedSubmission, err := s.repo.UpdateSubmission(ctx, req)
	if err != nil {
		log.Printf("[REPOSITORY] UpdateSubmission failed: %v, submissionID: %s",
			err, req.SubmussionID)

		if err == models.ErrAlreadyGraded {
			return nil, &models.Error{
				Code:    codes.FailedPrecondition,
				Message: "cannot update verified submission",
			}
		}
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "failed to update submission",
		}
	}

	log.Printf("[SUCCESS] Submission updated: %s", req.SubmussionID)
	return updatedSubmission, nil
}

func (s *Service) DeleteSubmission(ctx context.Context, userID, submissionID string) *models.Error {
	submission, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissionByID failed: %v, submissionID: %s", err, submissionID)
		return &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if submission == nil {
		log.Printf("[VALIDATION] Submission not found: %s", submissionID)
		return &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if submission.StudentID != userID {
		log.Printf("[PERMISSION] DeleteSubmission denied: user=%s, submissionStudent=%s",
			userID, submission.StudentID)
		return &models.Error{
			Code:    codes.PermissionDenied,
			Message: "only submission author can delete it",
		}
	}

	if submission.Status == models.SubmissionStatusVerified {
		log.Printf("[VALIDATION] Cannot delete verified submission: %s", submissionID)
		return &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot delete verified submission",
		}
	}

	if err := s.repo.DeleteSubmission(ctx, userID, submissionID); err != nil {
		log.Printf("[REPOSITORY] DeleteSubmission failed: %v, submissionID: %s, userID: %s",
			err, submissionID, userID)

		if err == models.ErrAlreadyGraded {
			return &models.Error{
				Code:    codes.FailedPrecondition,
				Message: "cannot delete verified submission",
			}
		}
		return &models.Error{
			Code:    codes.Internal,
			Message: "failed to delete submission",
		}
	}

	log.Printf("[SUCCESS] Submission deleted: %s by user: %s", submissionID, userID)
	return nil
}

func (s *Service) GradeSubmission(ctx context.Context, grade models.SubmissionGrade) (*models.SubmittedTask, *models.Error) {
	submission, err := s.repo.GetSubmissionByID(ctx, grade.SubmissionId)
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissionByID failed: %v, submissionID: %s",
			err, grade.SubmissionId)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if submission == nil {
		log.Printf("[VALIDATION] Submission not found: %s", grade.SubmissionId)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	task, err := s.repo.GetTaskByID(ctx, submission.TaskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s", err, submission.TaskID)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "task not found",
		}
	}

	if task.TutorId != grade.TutorId {
		log.Printf("[PERMISSION] GradeSubmission denied: gradeTutor=%s, taskTutor=%s",
			grade.TutorId, task.TutorId)
		return nil, &models.Error{
			Code:    codes.PermissionDenied,
			Message: "only task tutor can grade submissions",
		}
	}

	if *grade.Score > task.MaxScore {
		log.Printf("[VALIDATION] Score exceeds max: score=%d, max=%d, submissionID: %s",
			*grade.Score, task.MaxScore, grade.SubmissionId)
		return nil, &models.Error{
			Code:    codes.InvalidArgument,
			Message: "score cannot exceed task max score",
		}
	}

	gradedSubmission, err := s.repo.GradeSubmission(ctx, grade)
	if err != nil {
		log.Printf("[REPOSITORY] GradeSubmission failed: %v, submissionID: %s, tutorID: %s",
			err, grade.SubmissionId, grade.TutorId)
		return nil, &models.Error{
			Code:    codes.Internal,
			Message: "failed to grade submission",
		}
	}

	log.Printf("[SUCCESS] Submission graded: %s, score: %d, tutor: %s",
		grade.SubmissionId, *grade.Score, grade.TutorId)
	return gradedSubmission, nil
}

func (s *Service) ResetGrade(ctx context.Context, userID, submissionID string) *models.Error {
	submission, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissionByID failed: %v, submissionID: %s", err, submissionID)
		return &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if submission == nil {
		log.Printf("[VALIDATION] Submission not found: %s", submissionID)
		return &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	task, err := s.repo.GetTaskByID(ctx, submission.TaskID)
	if err != nil {
		log.Printf("[REPOSITORY] GetTaskByID failed: %v, taskID: %s", err, submission.TaskID)
		return &models.Error{
			Code:    codes.Internal,
			Message: "task not found",
		}
	}

	if task.TutorId != userID {
		log.Printf("[PERMISSION] ResetGrade denied: user=%s, taskTutor=%s", userID, task.TutorId)
		return &models.Error{
			Code:    codes.PermissionDenied,
			Message: "only task creator can reset grade",
		}
	}

	if submission.Status != models.SubmissionStatusVerified {
		log.Printf("[VALIDATION] Cannot reset ungraded submission: %s, status: %s",
			submissionID, submission.Status)
		return &models.Error{
			Code:    codes.FailedPrecondition,
			Message: "cannot reset grade for ungraded submission",
		}
	}

	if err := s.repo.ResetGrade(ctx, userID, submissionID); err != nil {
		log.Printf("[REPOSITORY] ResetGrade failed: %v, submissionID: %s, userID: %s",
			err, submissionID, userID)
		return &models.Error{
			Code:    codes.Internal,
			Message: "failed to reset grade",
		}
	}

	log.Printf("[SUCCESS] Grade reset: %s by tutor: %s", submissionID, userID)
	return nil
}

func (s *Service) GetTaskSubmission(ctx context.Context, submissionID string) (*models.SubmittedTask, *models.Error) {
	submission, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissionByID failed: %v, submissionID: %s", err, submissionID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	if submission == nil {
		log.Printf("[VALIDATION] Submission not found: %s", submissionID)
		return nil, &models.Error{
			Code:    codes.NotFound,
			Message: "submission not found",
		}
	}

	return submission, nil
}

func (s *Service) GetTaskSubmissions(ctx context.Context, params models.GetSubmissionsParams) ([]*models.SubmittedTaskShort, int32, *models.Error) {
	filter := models.SubmissionFilter{
		TaskID: params.TaskID,
		Offset: params.Offset,
		Limit:  params.Limit,
	}

	submissions, total, err := s.repo.GetSubmissions(ctx, filter)
	if err != nil {
		log.Printf("[REPOSITORY] GetSubmissions failed: %v, taskID: %s", err, params.TaskID)
		return nil, 0, &models.Error{
			Code:    codes.Internal,
			Message: "failed to get task submissions",
		}
	}

	log.Printf("[INFO] GetTaskSubmissions: task=%s, found=%d, total=%d",
		params.TaskID, len(submissions), total)
	return submissions, total, nil
}
