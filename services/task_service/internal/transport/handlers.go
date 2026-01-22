package transport

import (
	"context"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/task"

	"gitlab.crja72.ru/aisavelev-edu.hse.ru/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// создание и управление
func (s *Server) CreateTask(ctx context.Context, request *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if request.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id is not set")
	}

	if request.TutorId == "" {
		return nil, status.Error(codes.InvalidArgument, "tutor_id is not set")
	}

	if request.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is not set")
	}

	if request.Deadline == nil {
		return nil, status.Error(codes.InvalidArgument, "deadline is not set")
	}

	task, err := s.service.CreateTask(ctx, *models.CreateTaskFromProto(request))

	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	resp := models.TaskToProto(task)
	return &pb.CreateTaskResponse{Task: resp}, nil
}

func (s *Server) UpdateTask(ctx context.Context, request *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	if request.TutorId == "" {
		return nil, status.Error(codes.InvalidArgument, "tutor_id is not set")
	}

	if request.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	resp, err := s.service.UpdateTask(ctx, *models.UpdateTaskFromProto(request))
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return models.UpdateTaskToProto(resp), nil
}

func (s *Server) DeleteTask(ctx context.Context, request *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	if request.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is not set")
	}

	err := s.service.DeleteTask(ctx, request.GetUserId(), request.GetTaskId())
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}

func (s *Server) GetTask(ctx context.Context, request *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	if request.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	task, err := s.service.GetTask(ctx, request.GetTaskId())
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	resp := models.TaskToProto(task)
	return &pb.GetTaskResponse{Task: resp}, nil
}

func (s *Server) GetGroupTasks(ctx context.Context, request *pb.GetGroupTasksRequest) (*pb.GetGroupTasksResponse, error) {
	if request.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	params := models.GetGroupTasksParams{
		GroupID: request.GetGroupId(),
		Offset:  request.GetOffset(),
		Limit:   request.GetLimit(),
	}

	tasks, total, err := s.service.GetGroupTasks(ctx, params)
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	tasksProto := models.TasksListToProto(tasks)
	resp := &pb.GetGroupTasksResponse{
		Tasks:  tasksProto,
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  total,
	}

	return resp, nil
}

func (s *Server) GetCreatedByMeTasks(ctx context.Context, request *pb.GetCreatedByMeTasksRequest) (*pb.GetCreatedByMeTasksResponse, error) {
	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	params := models.CreatedByMeParams{
		UserID: request.GetUserId(),
		Offset: request.GetOffset(),
		Limit:  request.GetLimit(),
	}

	tasks, total, err := s.service.GetCreatedByMeTasks(ctx, params)
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	tasksProto := models.TasksListToProto(tasks)
	resp := &pb.GetCreatedByMeTasksResponse{
		Tasks:  tasksProto,
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  total,
	}

	return resp, nil
}

// решения
func (s *Server) CreateSubmission(ctx context.Context, request *pb.CreateSubmissionRequest) (*pb.CreateSubmissionResponse, error) {
	if request.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	if request.StudentId == "" {
		return nil, status.Error(codes.InvalidArgument, "student_id is not set")
	}

	if request.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is not set")
	}

	submussion, err := s.service.CreateSubmission(ctx, *models.CreateSubmissionFromProto(request))
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.CreateSubmissionResponse{Submission: models.SubmissionToProto(submussion)}, nil
}

func (s *Server) UpdateSubmission(ctx context.Context, request *pb.UpdateSubmissionRequest) (*pb.UpdateSubmissionResponse, error) {
	if request.SubmissionId == "" {
		return nil, status.Error(codes.InvalidArgument, "submission_id is not set")
	}

	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is not set")
	}

	resp, err := s.service.UpdateSubmission(ctx, *models.UpdateSubmissionFromProto(request))
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.UpdateSubmissionResponse{Submission: models.SubmissionToProto(resp)}, nil
}

func (s *Server) DeleteSubmission(ctx context.Context, request *pb.DeleteSubmissionRequest) (*pb.DeleteSubmissionResponse, error) {
	if request.SubmissionId == "" {
		return nil, status.Error(codes.InvalidArgument, "submission_id is not set")
	}

	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is not set")
	}

	err := s.service.DeleteSubmission(ctx, request.GetUserId(), request.GetSubmissionId())
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.DeleteSubmissionResponse{Success: true}, nil
}

func (s *Server) GradeSubmission(ctx context.Context, request *pb.GradeSubmissionRequest) (*pb.GradeSubmissionResponse, error) {
	if request.SubmissionId == "" {
		return nil, status.Error(codes.InvalidArgument, "submission_id is not set")
	}

	if request.TutorId == "" {
		return nil, status.Error(codes.InvalidArgument, "tutor_id is not set")
	}

	gradeReq := models.GradeSubmissionFromProto(request)
	submission, err := s.service.GradeSubmission(ctx, *gradeReq)
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.GradeSubmissionResponse{Submission: models.SubmissionToProto(submission)}, nil
}

func (s *Server) ResetGrade(ctx context.Context, request *pb.ResetGradeRequest) (*pb.ResetGradeResponse, error) {
	if request.SubmissionId == "" {
		return nil, status.Error(codes.InvalidArgument, "submission_id is not set")
	}

	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is not set")
	}

	err := s.service.ResetGrade(ctx, request.GetUserId(), request.GetSubmissionId())
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.ResetGradeResponse{Success: true}, nil
}

func (s *Server) GetSubmission(ctx context.Context, request *pb.GetSubmissionRequest) (*pb.GetSubmissionResponse, error) {
	if request.SubmissionId == "" {
		return nil, status.Error(codes.InvalidArgument, "submission_id is not set")
	}

	submission, err := s.service.GetTaskSubmission(ctx, request.GetSubmissionId())
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	return &pb.GetSubmissionResponse{Submission: models.SubmissionToProto(submission)}, nil
}

func (s *Server) GetTaskSubmissions(ctx context.Context, request *pb.GetTaskSubmissionsRequest) (*pb.GetTaskSubmissionsResponse, error) {
	if request.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id is not set")
	}

	params := models.GetSubmissionsParams{
		TaskID: request.GetTaskId(),
		Offset: request.GetOffset(),
		Limit:  request.GetLimit(),
	}

	submissions, total, err := s.service.GetTaskSubmissions(ctx, params)
	if err != nil {
		return nil, status.Error(err.Code, err.Message)
	}

	submissionsProto := models.SubmissionsListToProto(submissions)
	resp := &pb.GetTaskSubmissionsResponse{
		Submissions: submissionsProto,
		Offset:      params.Offset,
		Limit:       params.Limit,
		Total:       total,
	}

	return resp, nil
}
