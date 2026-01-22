package transport

import (
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/user"
	"context"
	"user-service/internal/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *ApiServer) ListUsers(ctx context.Context, req *pb.ListUserProfilesRequest) (*pb.ListUserProfilesResponse, error){
	users, err := h.userService.GetAllUsers(ctx, req.Limit, req.Offset)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	var usersResp []*pb.UserProfile
	for _,user := range users{
		usersResp = append(usersResp,
			&pb.UserProfile{
				UserId: user.UserID, 
				Email: user.Email, 
				Name: user.Name,
				Surname: user.Surname,
				Telegram: user.Telegram,
				CreatedAt: timestamppb.New(user.CreatedAt),
			})
	}

	return &pb.ListUserProfilesResponse{Users: usersResp}, nil
}

func (h *ApiServer) CreateUserProfile(ctx context.Context, req *pb.CreateUserProfileRequest) (*pb.UserProfileResponse, error){
	userReq := models.UserProfile{
		UserID: req.UserId,
		Email: req.Email,
		Name: req.Name,
		Surname: req.Surname,
		Telegram: req.Telegram,
	}
	user, err := h.userService.CreateUser(ctx, &userReq)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.UserProfileResponse{
		Profile: &pb.UserProfile{
			UserId: user.UserID,
				Email: user.Email, 
				Name: user.Name,
				Surname: user.Surname,
				Telegram: user.Telegram,
				CreatedAt: timestamppb.New(user.CreatedAt),
			}}, nil
}


func (h *ApiServer) GetUserByID(ctx context.Context, req *pb.GetUserProfileByIDRequest) (*pb.UserProfileResponse, error){
 
	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.UserProfileResponse{
		Profile: &pb.UserProfile{
			UserId: user.UserID,
				Email: user.Email, 
				Name: user.Name,
				Surname: user.Surname,
				Telegram: user.Telegram,
				CreatedAt: timestamppb.New(user.CreatedAt),
			}}, nil
}


func (h *ApiServer) GetUserByEmail(ctx context.Context, req *pb.GetUserProfileByEmailRequest) (*pb.UserProfileResponse, error){
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.UserProfileResponse{
		Profile: &pb.UserProfile{
			UserId: user.UserID,
				Email: user.Email, 
				Name: user.Name,
				Surname: user.Surname,
				Telegram: user.Telegram,
				CreatedAt: timestamppb.New(user.CreatedAt),
			}}, nil
}


func (h *ApiServer) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UserProfileResponse, error){
	userReq := models.UserProfile{
		UserID: req.UserId,
		Name: req.Name,
		Surname: req.Surname,
		Telegram: req.Telegram,
	}

	user, err := h.userService.UpdateUser(ctx, &userReq)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.UserProfileResponse{
		Profile: &pb.UserProfile{
			UserId: user.UserID,
				Email: user.Email, 
				Name: user.Name,
				Surname: user.Surname,
				Telegram: user.Telegram,
				CreatedAt: timestamppb.New(user.CreatedAt),
			}}, nil
}


func (h *ApiServer) DeleteUser(ctx context.Context, req *pb.DeleteUserProfileRequest) (*pb.EmptyResponse, error){
	err := h.userService.DeleteUser(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return nil, nil
}

func (h *ApiServer) GetCompleteUserProfile(ctx context.Context, req *pb.GetCompliteUserProfileRequest) (*pb.GetCompliteUserProfileResponse, error){
	user, err := h.userService.GetCompleteUserProfile(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	userProfile := pb.UserProfile{
		Email: user.UserProfile.Email, 
		Name: user.UserProfile.Name,
		Surname: user.UserProfile.Surname,
		Telegram: user.UserProfile.Telegram,
		CreatedAt: timestamppb.New(user.UserProfile.CreatedAt),
	}

	tutorProfile := pb.TutorProfile{
		UserId: user.TutorProfile.UserID,
		Bio: user.TutorProfile.Bio,
		Specialization: user.TutorProfile.Specialization,
		ExperienceYears: user.TutorProfile.Experience,
	}

	stdentProfile := pb.StudentProfile{
		UserId: user.StudentProfile.UserID,
		Bio: user.StudentProfile.Bio,
		GradeLevel: user.StudentProfile.Grade,
	}

	return &pb.GetCompliteUserProfileResponse{
		UserProfile: &userProfile,
		TutorProfile: &tutorProfile,
		StudentProfile: &stdentProfile,
	}, nil
}