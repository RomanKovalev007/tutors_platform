package grpc

import (
	"context"
	"group_service/internal/models"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no metadata in request")
	}

	values := md.Get("x-user-id")
	if len(values) == 0 {
		values = md.Get("X-User-Id")
		if len(values) == 0 {
			return "", status.Error(codes.Unauthenticated, "user id not provided by gateway")
		}
	}

	return values[0], nil
}

func convertGroup(g *models.Group) *pb.Group {
	pbG := &pb.Group{
		Id:          g.ID,
		TutorId:     g.TutorID,
		Name:        g.Name,
		Description: g.Description,
		CreatedAt:   timestamppb.New(g.CreatedAt),
		MemberCount: int32(len(g.Members)),
	}

	if g.Members != nil {
		pbG.Members = make([]*pb.GroupMember, len(g.Members))
		for i, m := range g.Members {
			pbG.Members[i] = &pb.GroupMember{
				GroupId:   m.GroupID,
				StudentId: m.StudentID,
				JoinedAt:  timestamppb.New(m.CreatedAt),
			}
		}
	}

	return pbG
}
