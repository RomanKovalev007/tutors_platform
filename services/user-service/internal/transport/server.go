package transport

import (
	"database/sql"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/kafka"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/user"
)

type ApiServer struct {
	pb.UserServiceServer

	userService    UserService
	tutorService   TutorService
	studentService StudentService

	EventHandler *service.KafkaHandler
}

func NewApiServer(pgDB *sql.DB, producer *kafka.Producer) *ApiServer {
	userRepo := repository.NewUserProfileRepository(pgDB)
	tutorRepo := repository.NewTutorProfileRepository(pgDB)
	studentRepo := repository.NewStudentProfileRepository(pgDB)

	userService := service.NewUserService(userRepo, tutorRepo, studentRepo)
	tutorService := service.NewTutorService(userRepo, tutorRepo)
	studentService := service.NewStudentService(userRepo, studentRepo)

	eventHandler := service.NewKafkaHandler(userRepo, producer)
	return &ApiServer{
		userService:    userService,
		tutorService:   tutorService,
		studentService: studentService,
		EventHandler:   eventHandler,
	}
}
