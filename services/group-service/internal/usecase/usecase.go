package usecase

import (
	"context"
	"fmt"
	"group_service/internal/models"
	"time"

	"github.com/google/uuid"
)

type GroupsRepo interface {
	CreateGroup(ctx context.Context, group *models.Group) error
	ListTutorGroups(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error)
	ListStudentGroups(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error)
	GetGroup(ctx context.Context, id string, includeMembers bool) (*models.Group, error)
	UpdateGroup(ctx context.Context, id string, name, desc *string) error
	DeleteGroup(ctx context.Context, id string) error
	AddMembers(ctx context.Context, groupID string, studentIDs []string) (int, error)
	RemoveMembers(ctx context.Context, groupID string, studentIDs []string) (int, error)
	GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error)
}

type UserClient interface {
	ValidateTutor(ctx context.Context, tutorId string) (bool, error)
}

type GroupsUsecase struct {
	groupsRepo GroupsRepo
	userClient UserClient
}

func NewGroupsUsecase(groupsRepo GroupsRepo, userClient UserClient) *GroupsUsecase {
	return &GroupsUsecase{
		groupsRepo: groupsRepo,
		userClient: userClient,
	}
}

func (u *GroupsUsecase) CreateGroup(ctx context.Context, tutorID, name, desc string) (*models.Group, error) {
	ok, err := u.userClient.ValidateTutor(ctx, tutorID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate tutor: %w", err)
	}
	if !ok {
		return nil, models.ErrTutorIsNotValid
	}

	group := &models.Group{
		ID:          uuid.New().String(),
		TutorID:     tutorID,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
		Members:     []*models.GroupMember{},
	}

	if err := u.groupsRepo.CreateGroup(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

// ListGroupsByTutor возвращает группы репетитора
func (u *GroupsUsecase) ListGroupsByTutor(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error) {
	ok, err := u.userClient.ValidateTutor(ctx, tutorID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate tutor: %w", err)
	}
	if !ok {
		return nil, models.ErrTutorIsNotValid
	}

	groups, err := u.groupsRepo.ListTutorGroups(ctx, tutorID, includeMembers)
	if err != nil {
		return nil, fmt.Errorf("failed to list tutor groups: %w", err)
	}

	return groups, nil
}

// ListGroupsByStudent возвращает группы студента
func (u *GroupsUsecase) ListGroupsByStudent(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error) {
	groups, err := u.groupsRepo.ListStudentGroups(ctx, studentID, includeMembers)
	if err != nil {
		return nil, fmt.Errorf("failed to list student groups: %w", err)
	}

	return groups, nil
}

func (u *GroupsUsecase) GetGroup(ctx context.Context, id string, includeMembers bool) (*models.Group, error) {
	group, err := u.groupsRepo.GetGroup(ctx, id, includeMembers)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	return group, nil
}

func (u *GroupsUsecase) UpdateGroup(ctx context.Context, groupId, userId string, name, desc *string) (*models.Group, error) {
	group, err := u.groupsRepo.GetGroup(ctx, groupId, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if group.TutorID != userId {
		return nil, models.ErrTutorIsNotValid
	}

	if err := u.groupsRepo.UpdateGroup(ctx, groupId, name, desc); err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	updatedGroup, err := u.groupsRepo.GetGroup(ctx, groupId, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated group: %w", err)
	}

	return updatedGroup, nil
}

func (u *GroupsUsecase) DeleteGroup(ctx context.Context, groupId, userId string) error {
	group, err := u.groupsRepo.GetGroup(ctx, groupId, false)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}
	if group.TutorID != userId {
		return models.ErrTutorIsNotValid
	}

	if err := u.groupsRepo.DeleteGroup(ctx, groupId); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

func (u *GroupsUsecase) AddGroupMembers(ctx context.Context, groupId, userId string, studentIDs []string) (int, error) {
	group, err := u.groupsRepo.GetGroup(ctx, groupId, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get group: %w", err)
	}
	if group.TutorID != userId {
		return 0, models.ErrTutorIsNotValid
	}

	addedCount, err := u.groupsRepo.AddMembers(ctx, groupId, studentIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to add members: %w", err)
	}

	return addedCount, nil
}

func (u *GroupsUsecase) RemoveGroupMembers(ctx context.Context, groupId, userId string, studentIDs []string) (int, error) {
	group, err := u.groupsRepo.GetGroup(ctx, groupId, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get group: %w", err)
	}
	if group.TutorID != userId {
		return 0, models.ErrTutorIsNotValid
	}

	removedCount, err := u.groupsRepo.RemoveMembers(ctx, groupId, studentIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to remove members: %w", err)
	}

	return removedCount, nil
}

func (u *GroupsUsecase) ListGroupMembers(ctx context.Context, groupId string) ([]*models.GroupMember, error) {
	members, err := u.groupsRepo.GetGroupMembers(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to list group members: %w", err)
	}

	return members, nil
}
