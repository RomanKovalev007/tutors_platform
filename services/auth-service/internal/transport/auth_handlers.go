package transport

import (
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"
	"auth_service/internal/models"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *ApiServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error){
	user, tokens, err := h.authService.Register(ctx, req.Email, req.Password, req.ConfirmPassword)
	if err != nil{
		return nil, err.Message
	}
	userResp := pb.User{
		Id: user.ID,
		Email: user.Email,
		IsActive: user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	return &pb.RegisterResponse{
		User: &userResp,
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken}, nil
}

func (h *ApiServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error){
	user, tokens, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	userResp := pb.User{
		Id: user.ID,
		Email: user.Email,
		IsActive: user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	return &pb.LoginResponse{
		User: &userResp,
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken}, nil

}

func (h *ApiServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.EmptyResponse, error){
	err := h.authService.Logout(ctx, req.RefreshToken)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return nil, nil
}

func (h *ApiServer) RefreshTokens(ctx context.Context, req *pb.RefreshTokensRequest) (*pb.RefreshTokenResponse, error){
	tokens, err := h.authService.Refresh(ctx, req.RefreshToken)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.RefreshTokenResponse{
		RefreshToken: tokens.RefreshToken,
		AccessToken: tokens.AccessToken,
		ExpiresIn: tokens.ExpiresAt,
	}, nil
}

func (h *ApiServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error){
	userID, exp, err := h.authService.ValidateAccessToken(ctx, req.Token)
	if err != nil{
		if err.Code == models.TOKENEXPIRED{
			return &pb.ValidateTokenResponse{
				IsValid: false,
				UserId: userID,
				ExpiresIn: int32(exp),
			}, nil
		}
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.ValidateTokenResponse{
		IsValid: true,
		UserId: userID,
		ExpiresIn: int32(exp),
	}, nil
}