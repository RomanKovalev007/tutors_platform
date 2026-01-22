package transport

import (
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"
	"auth_service/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	INTERNAL = status.New(codes.Internal, "internal error")
	ALREADYEXISTS = status.New(codes.AlreadyExists, "user already exists")
	INVALIDINPUT = status.New(codes.InvalidArgument, "invalid input")
	NOTFOUND = status.New(codes.NotFound, "user not found")
	OK = status.New(codes.OK, "ok")
)

func parseError(err *models.Error) *status.Status{
	e := &pb.Error{
		Code: string(err.Code), 
		Message: err.Message.Error(),
	}
	details, _ := anypb.New(e)
	
	var st *status.Status

	switch err.Code{
	case models.USERNOTFOUND:
		st, _ = NOTFOUND.WithDetails(details)

	case models.INVALIDCREDENTIALS:
		st, _ = INVALIDINPUT.WithDetails(details)

	case models.INVALIDTOKEN:
		st, _ = INVALIDINPUT.WithDetails(details)

	case models.INVALIDINPUT:
		st, _ = INVALIDINPUT.WithDetails(details)

	case models.USEREXISTS:
		st, _ = ALREADYEXISTS.WithDetails(details)

	case models.STATUS_OK:
		st, _ = OK.WithDetails(details)

	case models.INTERNALERROR:
		st, _ = INTERNAL.WithDetails(details)
	}
	return st
} 