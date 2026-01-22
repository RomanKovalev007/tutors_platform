package models

import (
	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/task"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateTaskFromProto(proto *pb.CreateTaskRequest) *AssignedTask {
	if proto == nil {
		return nil
	}

	task := &AssignedTask{
		GroupId:     proto.GetGroupId(),
		TutorId:     proto.GetTutorId(),
		Title:       proto.GetTitle(),
		Description: proto.Description,
		MaxScore:    proto.GetMaxScore(),
	}

	if proto.GetDeadline() != nil {
		task.Deadline = proto.GetDeadline().AsTime()
	}

	return task
}

func TaskToProto(task *AssignedTask) *pb.AssignedTask {
	if task == nil {
		return &pb.AssignedTask{}
	}

	protoTask := &pb.AssignedTask{
		TaskId:      task.ID,
		GroupId:     task.GroupId,
		TutorId:     task.TutorId,
		Title:       task.Title,
		Description: task.Description,
		MaxScore:    task.MaxScore,
		Status:      taskStatusToProto(task.Status),
		CreatedAt:   timestamppb.New(task.CreatedAt),
	}

	if !task.Deadline.IsZero() {
		protoTask.Deadline = timestamppb.New(task.Deadline)
	}

	return protoTask
}

func UpdateTaskFromProto(proto *pb.UpdateTaskRequest) *UpdateTaskRequest {
	if proto == nil {
		return nil
	}

	req := &UpdateTaskRequest{
		TutorID:     proto.GetTutorId(),
		TaskID:      proto.GetTaskId(),
		Title:       proto.Title,
		Description: proto.Description,
	}

	if proto.MaxScore != nil {
		req.MaxScore = proto.MaxScore
	}

	if proto.Deadline != nil {
		deadline := proto.GetDeadline().AsTime()
		req.Deadline = &deadline
	}

	return req
}

func UpdateTaskToProto(task *AssignedTask) *pb.UpdateTaskResponse {
	if task == nil {
		return &pb.UpdateTaskResponse{}
	}

	return &pb.UpdateTaskResponse{
		Task: TaskToProto(task),
	}
}

func TasksListToProto(tasks []*AssignedTaskShort) []*pb.AssignedTaskShort {
	if tasks == nil {
		return []*pb.AssignedTaskShort{}
	}

	proto := make([]*pb.AssignedTaskShort, len(tasks))
	for i, task := range tasks {
		proto[i] = ShortTaskToProto(task)
	}
	return proto
}

func ShortTaskToProto(task *AssignedTaskShort) *pb.AssignedTaskShort {
	if task == nil {
		return &pb.AssignedTaskShort{}
	}

	protoTask := &pb.AssignedTaskShort{
		TaskId:  task.ID,
		GroupId: task.GroupID,
		TutorId: task.TutorID,
		Title:   task.Title,
		Status:  taskStatusToProto(task.Status),
	}

	if !task.Deadline.IsZero() {
		protoTask.Deadline = timestamppb.New(task.Deadline)
	}

	return protoTask
}

func CreateSubmissionFromProto(proto *pb.CreateSubmissionRequest) *SubmittedTask {
	if proto == nil {
		return nil
	}

	return &SubmittedTask{
		TaskID:    proto.GetTaskId(),
		StudentID: proto.GetStudentId(),
		Content:   proto.GetContent(),
	}
}

func SubmissionToProto(submission *SubmittedTask) *pb.SubmittedTask {
	if submission == nil {
		return &pb.SubmittedTask{}
	}

	protoSubmission := &pb.SubmittedTask{
		SubmissionId: submission.ID,
		TaskId:       submission.TaskID,
		StudentId:    submission.StudentID,
		Content:      submission.Content,
		Status:       submissionStatusToProto(submission.Status),
		Score:        submission.Score,
		Feedback:     submission.Feedback,
		CreatedAt:    timestamppb.New(submission.CreatedAt),
	}

	// Проверяем указатели перед разыменованием
	if submission.UpdatedAt != nil {
		protoSubmission.UpdatedAt = timestamppb.New(*submission.UpdatedAt)
	}

	if submission.OverdueBy != nil {
		protoSubmission.OverdueBy = durationpb.New(*submission.OverdueBy)
	}

	return protoSubmission
}

func taskStatusToProto(status TaskStatus) pb.AssignedTaskStatus {
	switch status {
	case TaskStatusActive:
		return pb.AssignedTaskStatus_ACTIVE
	case TaskStatusExpired:
		return pb.AssignedTaskStatus_EXPIRED
	default:
		return pb.AssignedTaskStatus_ACTIVE
	}
}

func UpdateSubmissionFromProto(proto *pb.UpdateSubmissionRequest) *UpdateSubmissionRequest {
	if proto == nil {
		return nil
	}

	return &UpdateSubmissionRequest{
		UserID:       proto.GetUserId(),
		SubmussionID: proto.GetSubmissionId(),
		Content:      proto.Content,
	}
}

func GradeSubmissionFromProto(proto *pb.GradeSubmissionRequest) *SubmissionGrade {
	if proto == nil {
		return nil
	}

	return &SubmissionGrade{
		TutorId:      proto.GetTutorId(),
		SubmissionId: proto.GetSubmissionId(),
		Score:        proto.Score,
		Feedback:     proto.Feedback,
	}
}

func SubmissionShortToProto(submission *SubmittedTaskShort) *pb.SubmittedTaskShort {
	if submission == nil {
		return &pb.SubmittedTaskShort{}
	}

	protoSub := &pb.SubmittedTaskShort{
		SubmissionId: submission.ID,
		TaskId:       submission.TaskID,
		StudentId:    submission.StudentID,
		Score:        submission.Score,
		Status:       submissionStatusToProto(submission.Status),
	}

	if !submission.SubmittedAt.IsZero() {
		protoSub.CreatedAt = timestamppb.New(submission.SubmittedAt)
	}

	return protoSub
}

func SubmissionsListToProto(submissions []*SubmittedTaskShort) []*pb.SubmittedTaskShort {
	if submissions == nil {
		return []*pb.SubmittedTaskShort{}
	}

	result := make([]*pb.SubmittedTaskShort, len(submissions))
	for i, submission := range submissions {
		result[i] = SubmissionShortToProto(submission)
	}
	return result
}

func submissionStatusToProto(status SubmissionStatus) pb.SubmittedTaskStatus {
	switch status {
	case SubmissionStatusPending:
		return pb.SubmittedTaskStatus_PENDING
	case SubmissionStatusVerified:
		return pb.SubmittedTaskStatus_VERIFIED
	default:
		return pb.SubmittedTaskStatus_UNSPECIFIED
	}
}
