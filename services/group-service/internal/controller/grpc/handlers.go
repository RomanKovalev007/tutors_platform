package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"group_service/internal/models"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"
)

func errorResponse(code, message string) *pb.Error {
	return &pb.Error{
		Code:    code,
		Message: message,
	}
}

func (s *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	if req.TutorId == "" || req.Name == "" {
		return &pb.CreateGroupResponse{
			Result: &pb.CreateGroupResponse_Error{
				Error: errorResponse("INVALID_ARGUMENT", "tutor_id and name are required"),
			},
		}, status.Error(codes.InvalidArgument, "tutor_id and name are required")
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if userID != req.TutorId {
		return &pb.CreateGroupResponse{
			Result: &pb.CreateGroupResponse_Error{
				Error: errorResponse("PERMISSION_DENIED", "you can only create groups for yourself"),
			},
		}, status.Error(codes.PermissionDenied, "you can only create groups for yourself")
	}

	group, err := s.groupsUsecase.CreateGroup(ctx, req.TutorId, req.Name, req.Description)
	if err != nil {
		if err == models.ErrTutorIsNotValid {
			return &pb.CreateGroupResponse{
				Result: &pb.CreateGroupResponse_Error{
					Error: errorResponse("PERMISSION_DENIED", "invalid or non-tutor user"),
				},
			}, status.Error(codes.PermissionDenied, "invalid or non-tutor user")
		}
		return &pb.CreateGroupResponse{
			Result: &pb.CreateGroupResponse_Error{
				Error: errorResponse("INTERNAL", err.Error()),
			},
		}, status.Error(codes.Internal, "failed to create group: "+err.Error())
	}

	return &pb.CreateGroupResponse{
		Result: &pb.CreateGroupResponse_Group{
			Group: convertGroup(group),
		},
	}, nil
}

func (s *Server) ListGroups(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	var groups []*models.Group
	var err error

	switch filter := req.Filter.(type) {
	case *pb.ListGroupsRequest_TutorId:
		if filter.TutorId == "" {
			return &pb.ListGroupsResponse{
				Error: errorResponse("INVALID_ARGUMENT", "tutor_id cannot be empty"),
			}, status.Error(codes.InvalidArgument, "tutor_id cannot be empty")
		}
		groups, err = s.groupsUsecase.ListGroupsByTutor(ctx, filter.TutorId, req.IncludeMembers)

	case *pb.ListGroupsRequest_StudentId:
		if filter.StudentId == "" {
			return &pb.ListGroupsResponse{
				Error: errorResponse("INVALID_ARGUMENT", "student_id cannot be empty"),
			}, status.Error(codes.InvalidArgument, "student_id cannot be empty")
		}
		groups, err = s.groupsUsecase.ListGroupsByStudent(ctx, filter.StudentId, req.IncludeMembers)

	default:
		return &pb.ListGroupsResponse{
			Error: errorResponse("INVALID_ARGUMENT", "filter (tutor_id or student_id) is required"),
		}, status.Error(codes.InvalidArgument, "filter (tutor_id or student_id) is required")
	}

	if err != nil {
		if err == models.ErrTutorIsNotValid {
			return &pb.ListGroupsResponse{
				Error: errorResponse("PERMISSION_DENIED", "invalid tutor"),
			}, status.Error(codes.PermissionDenied, "invalid tutor")
		}
		return &pb.ListGroupsResponse{
			Error: errorResponse("INTERNAL", err.Error()),
		}, status.Error(codes.Internal, "failed to list groups: "+err.Error())
	}

	pbGroups := make([]*pb.Group, 0, len(groups))
	for _, g := range groups {
		pbGroups = append(pbGroups, convertGroup(g))
	}

	return &pb.ListGroupsResponse{Groups: pbGroups}, nil
}

func (s *Server) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	if req.Id == "" {
		return &pb.GetGroupResponse{
			Result: &pb.GetGroupResponse_Error{
				Error: errorResponse("INVALID_ARGUMENT", "group_id is required"),
			},
		}, status.Error(codes.InvalidArgument, "group_id is required")
	}

	group, err := s.groupsUsecase.GetGroup(ctx, req.Id, req.IncludeMembers)
	if err != nil {
		return &pb.GetGroupResponse{
			Result: &pb.GetGroupResponse_Error{
				Error: errorResponse("NOT_FOUND", "group not found"),
			},
		}, status.Error(codes.NotFound, "group not found")
	}

	return &pb.GetGroupResponse{
		Result: &pb.GetGroupResponse_Group{
			Group: convertGroup(group),
		},
	}, nil
}

func (s *Server) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.UpdateGroupResponse, error) {
	if req.Id == "" {
		return &pb.UpdateGroupResponse{
			Result: &pb.UpdateGroupResponse_Error{
				Error: errorResponse("INVALID_ARGUMENT", "group_id is required"),
			},
		}, status.Error(codes.InvalidArgument, "group_id is required")
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var name, desc *string
	if req.Name != nil {
		name = req.Name
	}
	if req.Description != nil {
		desc = req.Description
	}

	updatedGroup, err := s.groupsUsecase.UpdateGroup(ctx, req.Id, userID, name, desc)
	if err != nil {
		if err == models.ErrTutorIsNotValid {
			return &pb.UpdateGroupResponse{
				Result: &pb.UpdateGroupResponse_Error{
					Error: errorResponse("PERMISSION_DENIED", "you do not have permission to modify this group"),
				},
			}, status.Error(codes.PermissionDenied, "you do not have permission to modify this group")
		}
		return &pb.UpdateGroupResponse{
			Result: &pb.UpdateGroupResponse_Error{
				Error: errorResponse("INTERNAL", err.Error()),
			},
		}, status.Error(codes.Internal, "failed to update group: "+err.Error())
	}

	return &pb.UpdateGroupResponse{
		Result: &pb.UpdateGroupResponse_Group{
			Group: convertGroup(updatedGroup),
		},
	}, nil
}

func (s *Server) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteGroupResponse, error) {
	if req.Id == "" {
		return &pb.DeleteGroupResponse{
			Error: errorResponse("INVALID_ARGUMENT", "group_id is required"),
		}, status.Error(codes.InvalidArgument, "group_id is required")
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.groupsUsecase.DeleteGroup(ctx, req.Id, userID); err != nil {
		if err == models.ErrTutorIsNotValid {
			return &pb.DeleteGroupResponse{
				Error: errorResponse("PERMISSION_DENIED", "you do not have permission to delete this group"),
			}, status.Error(codes.PermissionDenied, "you do not have permission to delete this group")
		}
		return &pb.DeleteGroupResponse{
			Error: errorResponse("INTERNAL", err.Error()),
		}, status.Error(codes.Internal, "failed to delete group: "+err.Error())
	}

	return &pb.DeleteGroupResponse{}, nil
}

func (s *Server) ListGroupMembers(ctx context.Context, req *pb.ListGroupMembersRequest) (*pb.ListGroupMembersResponse, error) {
	if req.GroupId == "" {
		return &pb.ListGroupMembersResponse{
			Error: errorResponse("INVALID_ARGUMENT", "group_id is required"),
		}, status.Error(codes.InvalidArgument, "group_id is required")
	}

	members, err := s.groupsUsecase.ListGroupMembers(ctx, req.GroupId)
	if err != nil {
		return &pb.ListGroupMembersResponse{
			Error: errorResponse("INTERNAL", err.Error()),
		}, status.Error(codes.Internal, "failed to list members: "+err.Error())
	}

	pbMembers := make([]*pb.GroupMember, len(members))
	for i, m := range members {
		pbMembers[i] = &pb.GroupMember{
			GroupId:   m.GroupID,
			StudentId: m.StudentID,
			JoinedAt:  timestamppb.New(m.CreatedAt),
		}
	}

	return &pb.ListGroupMembersResponse{Members: pbMembers}, nil
}

func (s *Server) AddGroupMembers(ctx context.Context, req *pb.AddGroupMembersRequest) (*pb.AddGroupMembersResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GroupId == "" {
		return &pb.AddGroupMembersResponse{
			Error: errorResponse("INVALID_ARGUMENT", "group_id is required"),
		}, status.Error(codes.InvalidArgument, "group_id is required")
	}
	if len(req.StudentIds) == 0 {
		return &pb.AddGroupMembersResponse{
			Error: errorResponse("INVALID_ARGUMENT", "student_ids cannot be empty"),
		}, status.Error(codes.InvalidArgument, "student_ids cannot be empty")
	}

	addedCount, err := s.groupsUsecase.AddGroupMembers(ctx, req.GroupId, userID, req.StudentIds)
	if err != nil {
		if err == models.ErrTutorIsNotValid {
			return &pb.AddGroupMembersResponse{
				Error: errorResponse("PERMISSION_DENIED", "you do not have permission to modify this group"),
			}, status.Error(codes.PermissionDenied, "you do not have permission to modify this group")
		}
		return &pb.AddGroupMembersResponse{
			Error: errorResponse("INTERNAL", err.Error()),
		}, status.Error(codes.Internal, "failed to add members: "+err.Error())
	}

	return &pb.AddGroupMembersResponse{AddedCount: int32(addedCount)}, nil
}

func (s *Server) RemoveGroupMembers(ctx context.Context, req *pb.RemoveGroupMembersRequest) (*pb.RemoveGroupMembersResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GroupId == "" {
		return &pb.RemoveGroupMembersResponse{
			Error: errorResponse("INVALID_ARGUMENT", "group_id is required"),
		}, status.Error(codes.InvalidArgument, "group_id is required")
	}
	if len(req.StudentIds) == 0 {
		return &pb.RemoveGroupMembersResponse{
			Error: errorResponse("INVALID_ARGUMENT", "student_ids cannot be empty"),
		}, status.Error(codes.InvalidArgument, "student_ids cannot be empty")
	}

	removedCount, err := s.groupsUsecase.RemoveGroupMembers(ctx, req.GroupId, userID, req.StudentIds)
	if err != nil {
		if err == models.ErrTutorIsNotValid {
			return &pb.RemoveGroupMembersResponse{
				Error: errorResponse("PERMISSION_DENIED", "you do not have permission to modify this group"),
			}, status.Error(codes.PermissionDenied, "you do not have permission to modify this group")
		}
		return &pb.RemoveGroupMembersResponse{
			Error: errorResponse("INTERNAL", err.Error()),
		}, status.Error(codes.Internal, "failed to remove members: "+err.Error())
	}

	return &pb.RemoveGroupMembersResponse{RemovedCount: int32(removedCount)}, nil
}
