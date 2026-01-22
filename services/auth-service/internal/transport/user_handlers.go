package transport

import (
	"context"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func (h *ApiServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error){
	users, err := h.userService.GetAllUsers(ctx, req.Limit, req.Offset, req.IsActive)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	var usersResp []*pb.User
	for _,user := range users{
		usersResp = append(usersResp,
			&pb.User{
				Id: user.ID, 
				Email: user.Email, 
				IsActive: user.IsActive, 
				CreatedAt: timestamppb.New(user.CreatedAt)})
	}

	return &pb.ListUsersResponse{Users: usersResp}, nil
}


func (h *ApiServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserResponse, error){
 
	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id: user.ID, 
			Email: user.Email, 
			IsActive: user.IsActive, 
			CreatedAt: timestamppb.New(user.CreatedAt)}}, nil
}


func (h *ApiServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error){
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id: user.ID, 
			Email: user.Email, 
			IsActive: user.IsActive, 
			CreatedAt: timestamppb.New(user.CreatedAt)}}, nil
}


func (h *ApiServer) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*pb.EmptyResponse, error){
	err := h.userService.UpdateIsActiveStatus(ctx, req.UserId, req.IsActive)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return nil, nil
}


func (h *ApiServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.EmptyResponse, error){
	err := h.userService.DeleteUser(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return nil, nil
}