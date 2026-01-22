package transport

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/task"

	"gitlab.crja72.ru/aisavelev-edu.hse.ru/internal/models"
	"google.golang.org/grpc"
)

type TaskService interface {
	// Task operation
	CreateTask(ctx context.Context, task models.AssignedTask) (*models.AssignedTask, *models.Error)
	UpdateTask(ctx context.Context, task models.UpdateTaskRequest) (*models.AssignedTask, *models.Error)
	DeleteTask(ctx context.Context, userID, taskID string) *models.Error
	GetTask(ctx context.Context, taskID string) (*models.AssignedTask, *models.Error)
	GetGroupTasks(ctx context.Context, params models.GetGroupTasksParams) ([]*models.AssignedTaskShort, int32, *models.Error)
	GetCreatedByMeTasks(ctx context.Context, params models.CreatedByMeParams) ([]*models.AssignedTaskShort, int32, *models.Error)

	// Submission operation
	CreateSubmission(ctx context.Context, task models.SubmittedTask) (*models.SubmittedTask, *models.Error)
	UpdateSubmission(ctx context.Context, task models.UpdateSubmissionRequest) (*models.SubmittedTask, *models.Error)
	DeleteSubmission(ctx context.Context, UserID, SubmissionID string) *models.Error
	GradeSubmission(ctx context.Context, grade models.SubmissionGrade) (*models.SubmittedTask, *models.Error)
	ResetGrade(ctx context.Context, userID, gradeID string) *models.Error
	GetTaskSubmission(ctx context.Context, submissionID string) (*models.SubmittedTask, *models.Error)
	GetTaskSubmissions(ctx context.Context, params models.GetSubmissionsParams) ([]*models.SubmittedTaskShort, int32, *models.Error)
}

type Server struct {
	service    TaskService
	grpcServer *grpc.Server
	listener   net.Listener
	pb.UnimplementedTaskServiceServer
}

func NewServer(port string, service TaskService) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to open port %s: %v", port, err)
	}

	const defaultMaxMsgSize = 1024 * 1024 * 20
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(defaultMaxMsgSize),
		grpc.MaxSendMsgSize(defaultMaxMsgSize),
	}

	grpcServer := grpc.NewServer(opts...)

	server := &Server{
		service:    service,
		grpcServer: grpcServer,
		listener:   lis,
	}

	pb.RegisterTaskServiceServer(grpcServer, server)

	return server, nil
}

func (s *Server) Start() {
	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	}()
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
