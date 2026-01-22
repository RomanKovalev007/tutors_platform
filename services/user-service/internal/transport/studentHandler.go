package transport

import (
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/user"
	"context"
	"user-service/internal/models"
)

func (h *ApiServer) CreateStudentProfile(ctx context.Context, req *pb.CreateStudentProfileRequest) (*pb.StudentProfileResponse, error){
	studentReq := models.StudentProfile{
		UserID: req.UserId,
		Bio: req.Bio,
		Grade: req.GradeLevel,
	}

	student, err := h.studentService.CreateStudentProfile(ctx, &studentReq)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.StudentProfileResponse{
		Profile: &pb.StudentProfile{
			UserId: student.UserID,
			Bio: student.Bio,
			GradeLevel: student.Grade,
		},
	}, nil
}


func (h *ApiServer) GetStudentProfile(ctx context.Context, req *pb.GetStudentProfileRequest) (*pb.StudentProfileResponse, error){
 
	student, err := h.studentService.GetStudentProfile(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return &pb.StudentProfileResponse{
		Profile: &pb.StudentProfile{
			UserId: student.UserID,
			Bio: student.Bio,
			GradeLevel: student.Grade,
		},
	}, nil

}


func (h *ApiServer) UpdateStudentProfile(ctx context.Context, req *pb.UpdateStudentProfileRequest) (*pb.StudentProfileResponse, error){
	studentReq := models.StudentProfile{
		UserID: req.UserId,
		Bio: req.Bio,
		Grade: req.GradeLevel,
	}

	student, err := h.studentService.UpdateStudentProfile(ctx, &studentReq)
	if err != nil{
		st := parseError(err)
		return nil, st.Err()
	}

	return &pb.StudentProfileResponse{
		Profile: &pb.StudentProfile{
			UserId: student.UserID,
			Bio: student.Bio,
			GradeLevel: student.Grade,
		},
	}, nil
}


func (h *ApiServer) DeleteStudentProfile(ctx context.Context, req *pb.DeleteStudentProfileRequest) (*pb.EmptyResponse, error){
	err := h.studentService.DeleteStudentProfile(ctx, req.UserId)
	if err != nil{
		st :=  parseError(err)
		return nil, st.Err()
	}

	return nil, nil
}