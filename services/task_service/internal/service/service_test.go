package service

import (
	"context"
	"testing"
	"time"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"
	"task_service/internal/models"
	"google.golang.org/grpc/codes"
)

type mockRepository struct {
	tasks       map[string]*models.AssignedTask
	submissions map[string]*models.SubmittedTask
	createErr   error
	getErr      error
	updateErr   error
	deleteErr   error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		tasks:       make(map[string]*models.AssignedTask),
		submissions: make(map[string]*models.SubmittedTask),
	}
}

func (m *mockRepository) CreateTask(ctx context.Context, task models.AssignedTask) (*models.AssignedTask, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.tasks[task.ID] = &task
	return &task, nil
}

func (m *mockRepository) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.AssignedTask, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	task, exists := m.tasks[req.TaskID]
	if !exists {
		return nil, models.ErrNotFound
	}
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.MaxScore != nil {
		task.MaxScore = *req.MaxScore
	}
	if req.Deadline != nil {
		task.Deadline = *req.Deadline
	}
	now := time.Now()
	task.UpdatedAt = &now
	return task, nil
}

func (m *mockRepository) SoftDeleteTask(ctx context.Context, userID, taskID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.tasks, taskID)
	return nil
}

func (m *mockRepository) GetTaskByID(ctx context.Context, taskID string) (*models.AssignedTask, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	task, exists := m.tasks[taskID]
	if !exists {
		return nil, models.ErrNotFound
	}
	return task, nil
}

func (m *mockRepository) GetTasks(ctx context.Context, filter models.TaskFilter) ([]*models.AssignedTaskShort, int32, error) {
	var result []*models.AssignedTaskShort
	for _, task := range m.tasks {
		if filter.GroupID != "" && task.GroupId != filter.GroupID {
			continue
		}
		if filter.UserID != "" && task.TutorId != filter.UserID {
			continue
		}
		result = append(result, &models.AssignedTaskShort{
			ID:       task.ID,
			GroupID:  task.GroupId,
			TutorID:  task.TutorId,
			Title:    task.Title,
			Deadline: task.Deadline,
			Status:   task.Status,
		})
	}
	return result, int32(len(result)), nil
}

func (m *mockRepository) CreateSubmission(ctx context.Context, submission models.SubmittedTask) (*models.SubmittedTask, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.submissions[submission.ID] = &submission
	return &submission, nil
}

func (m *mockRepository) UpdateSubmission(ctx context.Context, req models.UpdateSubmissionRequest) (*models.SubmittedTask, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	submission, exists := m.submissions[req.SubmussionID]
	if !exists {
		return nil, models.ErrNotFound
	}
	if req.Content != nil {
		submission.Content = *req.Content
	}
	now := time.Now()
	submission.UpdatedAt = &now
	return submission, nil
}

func (m *mockRepository) DeleteSubmission(ctx context.Context, userID, submissionID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.submissions, submissionID)
	return nil
}

func (m *mockRepository) GradeSubmission(ctx context.Context, grade models.SubmissionGrade) (*models.SubmittedTask, error) {
	submission, exists := m.submissions[grade.SubmissionId]
	if !exists {
		return nil, models.ErrNotFound
	}
	submission.Score = grade.Score
	submission.Feedback = grade.Feedback
	submission.Status = models.SubmissionStatusVerified
	return submission, nil
}

func (m *mockRepository) ResetGrade(ctx context.Context, userID, submissionID string) error {
	submission, exists := m.submissions[submissionID]
	if !exists {
		return models.ErrNotFound
	}
	submission.Score = nil
	submission.Feedback = nil
	submission.Status = models.SubmissionStatusPending
	return nil
}

func (m *mockRepository) GetSubmissionByID(ctx context.Context, submissionID string) (*models.SubmittedTask, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	submission, exists := m.submissions[submissionID]
	if !exists {
		return nil, models.ErrNotFound
	}
	return submission, nil
}

func (m *mockRepository) GetSubmissions(ctx context.Context, filter models.SubmissionFilter) ([]*models.SubmittedTaskShort, int32, error) {
	var result []*models.SubmittedTaskShort
	for _, sub := range m.submissions {
		if filter.TaskID != "" && sub.TaskID != filter.TaskID {
			continue
		}
		result = append(result, &models.SubmittedTaskShort{
			ID:          sub.ID,
			TaskID:      sub.TaskID,
			StudentID:   sub.StudentID,
			Score:       sub.Score,
			Status:      sub.Status,
			SubmittedAt: sub.CreatedAt,
		})
	}
	return result, int32(len(result)), nil
}

func (m *mockRepository) MarkExpiredTasks(ctx context.Context) error {
	now := time.Now()
	for _, task := range m.tasks {
		if task.Deadline.Before(now) && task.Status == models.TaskStatusActive {
			task.Status = models.TaskStatusExpired
		}
	}
	return nil
}

type mockGroupClient struct {
	groups  map[string]*pb.Group
	members map[string][]*pb.GroupMember
	getErr  error
}

func newMockGroupClient() *mockGroupClient {
	return &mockGroupClient{
		groups:  make(map[string]*pb.Group),
		members: make(map[string][]*pb.GroupMember),
	}
}

func (m *mockGroupClient) GetGroupInfo(ctx context.Context, groupID string) (*pb.Group, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	group, exists := m.groups[groupID]
	if !exists {
		return nil, models.ErrNotFound
	}
	return group, nil
}

func (m *mockGroupClient) GetGroupMembers(ctx context.Context, groupID string) ([]*pb.GroupMember, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	members, exists := m.members[groupID]
	if !exists {
		return []*pb.GroupMember{}, nil
	}
	return members, nil
}

func newTestService() (*Service, *mockRepository, *mockGroupClient) {
	repo := newMockRepository()
	groupClient := newMockGroupClient()
	svc := NewService(repo, groupClient)
	return svc, repo, groupClient
}

func TestService_CreateTask_Success(t *testing.T) {
	svc, _, groupClient := newTestService()
	ctx := context.Background()

	groupClient.groups["group-1"] = &pb.Group{
		Id:      "group-1",
		TutorId: "tutor-1",
	}

	task := models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
	}

	result, err := svc.CreateTask(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected task, got nil")
	}
	if result.ID != task.ID {
		t.Errorf("expected ID %s, got %s", task.ID, result.ID)
	}
	if result.Status != models.TaskStatusActive {
		t.Errorf("expected status ACTIVE, got %s", result.Status)
	}
}

func TestService_CreateTask_DeadlineInPast(t *testing.T) {
	svc, _, groupClient := newTestService()
	ctx := context.Background()

	groupClient.groups["group-1"] = &pb.Group{
		Id:      "group-1",
		TutorId: "tutor-1",
	}

	task := models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(-24 * time.Hour),
	}

	result, err := svc.CreateTask(ctx, task)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil task")
	}
}

func TestService_CreateTask_NegativeMaxScore(t *testing.T) {
	svc, _, groupClient := newTestService()
	ctx := context.Background()

	groupClient.groups["group-1"] = &pb.Group{
		Id:      "group-1",
		TutorId: "tutor-1",
	}

	task := models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: -10,
		Deadline: time.Now().Add(24 * time.Hour),
	}

	result, err := svc.CreateTask(ctx, task)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil task")
	}
}

func TestService_CreateTask_WrongTutor(t *testing.T) {
	svc, _, groupClient := newTestService()
	ctx := context.Background()

	groupClient.groups["group-1"] = &pb.Group{
		Id:      "group-1",
		TutorId: "tutor-2",
	}

	task := models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
	}

	result, err := svc.CreateTask(ctx, task)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.FailedPrecondition {
		t.Errorf("expected code FailedPrecondition, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil task")
	}
}

func TestService_GetTask_Success(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	result, err := svc.GetTask(ctx, task.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected task, got nil")
	}
	if result.Title != task.Title {
		t.Errorf("expected title %s, got %s", task.Title, result.Title)
	}
}

func TestService_GetTask_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	result, err := svc.GetTask(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.NotFound {
		t.Errorf("expected code NotFound, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil task")
	}
}

func TestService_UpdateTask_Success(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Original Title",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	newTitle := "Updated Title"
	req := models.UpdateTaskRequest{
		TutorID: "tutor-1",
		TaskID:  "task-1",
		Title:   &newTitle,
	}

	result, err := svc.UpdateTask(ctx, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected task, got nil")
	}
	if result.Title != newTitle {
		t.Errorf("expected title %s, got %s", newTitle, result.Title)
	}
}

func TestService_UpdateTask_NotOwner(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Original Title",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	newTitle := "Updated Title"
	req := models.UpdateTaskRequest{
		TutorID: "tutor-2",
		TaskID:  "task-1",
		Title:   &newTitle,
	}

	result, err := svc.UpdateTask(ctx, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.PermissionDenied {
		t.Errorf("expected code PermissionDenied, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil task")
	}
}

func TestService_UpdateTask_ExpiredTask(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Original Title",
		MaxScore: 100,
		Deadline: time.Now().Add(-24 * time.Hour),
		Status:   models.TaskStatusExpired,
	}
	repo.tasks[task.ID] = task

	newTitle := "Updated Title"
	req := models.UpdateTaskRequest{
		TutorID: "tutor-1",
		TaskID:  "task-1",
		Title:   &newTitle,
	}

	result, err := svc.UpdateTask(ctx, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.FailedPrecondition {
		t.Errorf("expected code FailedPrecondition, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil task")
	}
}

func TestService_DeleteTask_Success(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	err := svc.DeleteTask(ctx, "tutor-1", "task-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestService_DeleteTask_NotOwner(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	err := svc.DeleteTask(ctx, "tutor-2", "task-1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.PermissionDenied {
		t.Errorf("expected code PermissionDenied, got %v", err.Code)
	}
}

func TestService_GetGroupTasks_Success(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task1 := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Task 1",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	task2 := &models.AssignedTask{
		ID:       "task-2",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Task 2",
		MaxScore: 50,
		Deadline: time.Now().Add(48 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task1.ID] = task1
	repo.tasks[task2.ID] = task2

	params := models.GetGroupTasksParams{
		GroupID: "group-1",
		Limit:   10,
		Offset:  0,
	}

	tasks, total, err := svc.GetGroupTasks(ctx, params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestService_CreateSubmission_Success(t *testing.T) {
	svc, repo, groupClient := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	groupClient.members["group-1"] = []*pb.GroupMember{
		{StudentId: "student-1"},
	}

	submission := models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "task-1",
		StudentID: "student-1",
		Content:   "My submission",
	}

	result, err := svc.CreateSubmission(ctx, submission)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected submission, got nil")
	}
	if result.Status != models.SubmissionStatusPending {
		t.Errorf("expected status PENDING, got %s", result.Status)
	}
}

func TestService_CreateSubmission_TaskNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	submission := models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "nonexistent-task",
		StudentID: "student-1",
		Content:   "My submission",
	}

	result, err := svc.CreateSubmission(ctx, submission)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.NotFound {
		t.Errorf("expected code NotFound, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil submission")
	}
}

func TestService_CreateSubmission_StudentNotInGroup(t *testing.T) {
	svc, repo, groupClient := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	groupClient.members["group-1"] = []*pb.GroupMember{
		{StudentId: "student-2"},
	}

	submission := models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "task-1",
		StudentID: "student-1",
		Content:   "My submission",
	}

	result, err := svc.CreateSubmission(ctx, submission)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil submission")
	}
}

func TestService_GradeSubmission_Success(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	submission := &models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "task-1",
		StudentID: "student-1",
		Content:   "My submission",
		Status:    models.SubmissionStatusPending,
	}
	repo.submissions[submission.ID] = submission

	score := int32(85)
	feedback := "Good work!"
	grade := models.SubmissionGrade{
		SubmissionId: "submission-1",
		TutorId:      "tutor-1",
		Score:        &score,
		Feedback:     &feedback,
	}

	result, err := svc.GradeSubmission(ctx, grade)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected submission, got nil")
	}
	if *result.Score != score {
		t.Errorf("expected score %d, got %d", score, *result.Score)
	}
	if result.Status != models.SubmissionStatusVerified {
		t.Errorf("expected status VERIFIED, got %s", result.Status)
	}
}

func TestService_GradeSubmission_NotTutor(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	submission := &models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "task-1",
		StudentID: "student-1",
		Content:   "My submission",
		Status:    models.SubmissionStatusPending,
	}
	repo.submissions[submission.ID] = submission

	score := int32(85)
	grade := models.SubmissionGrade{
		SubmissionId: "submission-1",
		TutorId:      "tutor-2",
		Score:        &score,
	}

	result, err := svc.GradeSubmission(ctx, grade)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.PermissionDenied {
		t.Errorf("expected code PermissionDenied, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil submission")
	}
}

func TestService_GradeSubmission_ScoreExceedsMax(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	task := &models.AssignedTask{
		ID:       "task-1",
		GroupId:  "group-1",
		TutorId:  "tutor-1",
		Title:    "Test Task",
		MaxScore: 100,
		Deadline: time.Now().Add(24 * time.Hour),
		Status:   models.TaskStatusActive,
	}
	repo.tasks[task.ID] = task

	submission := &models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "task-1",
		StudentID: "student-1",
		Content:   "My submission",
		Status:    models.SubmissionStatusPending,
	}
	repo.submissions[submission.ID] = submission

	score := int32(150)
	grade := models.SubmissionGrade{
		SubmissionId: "submission-1",
		TutorId:      "tutor-1",
		Score:        &score,
	}

	result, err := svc.GradeSubmission(ctx, grade)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.InvalidArgument {
		t.Errorf("expected code InvalidArgument, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil submission")
	}
}

func TestService_GetTaskSubmission_Success(t *testing.T) {
	svc, repo, _ := newTestService()
	ctx := context.Background()

	submission := &models.SubmittedTask{
		ID:        "submission-1",
		TaskID:    "task-1",
		StudentID: "student-1",
		Content:   "My submission",
		Status:    models.SubmissionStatusPending,
	}
	repo.submissions[submission.ID] = submission

	result, err := svc.GetTaskSubmission(ctx, "submission-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected submission, got nil")
	}
	if result.Content != submission.Content {
		t.Errorf("expected content %s, got %s", submission.Content, result.Content)
	}
}

func TestService_GetTaskSubmission_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	result, err := svc.GetTaskSubmission(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Code != codes.NotFound {
		t.Errorf("expected code NotFound, got %v", err.Code)
	}
	if result != nil {
		t.Error("expected nil submission")
	}
}
