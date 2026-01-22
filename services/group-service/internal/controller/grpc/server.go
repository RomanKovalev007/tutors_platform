package grpc

import (
	"context"
	"fmt"
	"group_service/internal/models"
	"net"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"

	"google.golang.org/grpc"
)

type GroupsUsecase interface {
	// Управление группами
	CreateGroup(ctx context.Context, tutorID, name, desc string) (*models.Group, error)
	GetGroup(ctx context.Context, id string, includeMembers bool) (*models.Group, error)
	UpdateGroup(ctx context.Context, groupIdStr, userIdStr string, name, desc *string) (*models.Group, error)
	DeleteGroup(ctx context.Context, groupIdStr, userIdStr string) error

	// Получение списков групп
	ListGroupsByTutor(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error)
	ListGroupsByStudent(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error)

	// Управление участниками
	ListGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error)
	AddGroupMembers(ctx context.Context, groupIDStr, userIdStr string, studentIDStrs []string) (int, error)
	RemoveGroupMembers(ctx context.Context, groupIDStr, userIdStr string, studentIDStrs []string) (int, error)
}

type Server struct {
	pb.GroupsServiceServer
	srv *grpc.Server

	groupsUsecase GroupsUsecase
}

func NewServer(groupsUsecase GroupsUsecase) *Server {
	grpcSrv := grpc.NewServer()

	server := &Server{
		srv:           grpcSrv,
		groupsUsecase: groupsUsecase,
	}

	pb.RegisterGroupsServiceServer(grpcSrv, server)

	return server
}

func (s *Server) Run(grpcPort int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen tcp: %w", err)
	}

	if err := s.srv.Serve(lis); err != nil {
		return fmt.Errorf("error grpc server serving: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
