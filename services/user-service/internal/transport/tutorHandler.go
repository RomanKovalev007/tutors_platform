package transport

import (
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/user"
	"context"
	"user-service/internal/models"
)

func (h *ApiServer) CreateTutorProfile(ctx context.Context, req *pb.CreateTutorProfileRequest) (*pb.TutorProfileResponse, error){
	tutorReq := models.TutorProfile{
		UserID: req.UserId,
		Bio: req.Bio,
		Specialization: req.Specialization,
		Experience: req.ExperienceYears,
	}

	tutor, err := h.tutorService.CreateTutorProfile(ctx, &tutorReq)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.TutorProfileResponse{
		Profile: &pb.TutorProfile{
			UserId: tutor.UserID,
			Bio: tutor.Bio,
			Specialization: tutor.Specialization,
			ExperienceYears: tutor.Experience,
		},
	}, nil
}


func (h *ApiServer) GetTutorProfile(ctx context.Context, req *pb.GetTutorProfileRequest) (*pb.TutorProfileResponse, error){
 
	tutor, err := h.tutorService.GetTutorProfile(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.TutorProfileResponse{
		Profile: &pb.TutorProfile{
			UserId: tutor.UserID,
			Bio: tutor.Bio,
			Specialization: tutor.Specialization,
			ExperienceYears: tutor.Experience,
		},
	}, nil

}


func (h *ApiServer) UpdateTutorProfile(ctx context.Context, req *pb.UpdateTutorProfileRequest) (*pb.TutorProfileResponse, error){
	tutorReq := models.TutorProfile{
		UserID: req.UserId,
		Bio: req.Bio,
		Specialization: req.Specialization,
		Experience: req.ExperienceYears,
	}

	tutor, err := h.tutorService.UpdateTutorProfile(ctx, &tutorReq)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.TutorProfileResponse{
		Profile: &pb.TutorProfile{
			UserId: tutor.UserID,
			Bio: tutor.Bio,
			Specialization: tutor.Specialization,
			ExperienceYears: tutor.Experience,
		},
	}, nil
}


func (h *ApiServer) DeleteTutorProfile(ctx context.Context, req *pb.DeleteTutorProfileRequest) (*pb.EmptyResponse, error){
	err := h.tutorService.DeleteTutorProfile(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return nil, nil
}

func (h *ApiServer) ValidateTutor(ctx context.Context, req *pb.ValidateTutorRequest) (*pb.ValidateTutorResponse, error){
	status, err := h.tutorService.ValidateTutor(ctx, req.UserId)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.ValidateTutorResponse{IsValidTutor: status}, nil
}