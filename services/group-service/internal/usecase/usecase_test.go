package usecase_test

import (
	"context"
	"errors"
	"testing"

	"group_service/internal/models"
	"group_service/internal/usecase"
)

type mockRepo struct {
	createGroupCalled   bool
	getGroupCalled      bool
	updateGroupCalled   bool
	deleteGroupCalled   bool
	addMembersCalled    bool
	removeMembersCalled bool

	getGroupFirstResult  *models.Group // первая загрузка (до обновления)
	getGroupSecondResult *models.Group // вторая загрузка (после обновления)
	getGroupErr          error
	updateGroupErr       error
	deleteGroupErr       error
	addMembersCount      int
	addMembersErr        error
	removeMembersCount   int
	removeMembersErr     error

	// Счётчик вызовов GetGroup
	getGroupCallCount int

	// Поля для переопределения методов
	getGroupMembersFunc   func(ctx context.Context, groupID string) ([]*models.GroupMember, error)
	listTutorGroupsFunc   func(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error)
	listStudentGroupsFunc func(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error)
}

func (m *mockRepo) CreateGroup(ctx context.Context, group *models.Group) error {
	m.createGroupCalled = true
	return nil
}

func (m *mockRepo) GetGroup(ctx context.Context, id string, includeMembers bool) (*models.Group, error) {
	m.getGroupCalled = true
	m.getGroupCallCount++
	if m.getGroupCallCount == 1 {
		return m.getGroupFirstResult, m.getGroupErr
	}
	return m.getGroupSecondResult, m.getGroupErr
}

func (m *mockRepo) UpdateGroup(ctx context.Context, id string, name, desc *string) error {
	m.updateGroupCalled = true
	return m.updateGroupErr
}

func (m *mockRepo) DeleteGroup(ctx context.Context, id string) error {
	m.deleteGroupCalled = true
	return m.deleteGroupErr
}

func (m *mockRepo) AddMembers(ctx context.Context, groupID string, studentIDs []string) (int, error) {
	m.addMembersCalled = true
	return m.addMembersCount, m.addMembersErr
}

func (m *mockRepo) RemoveMembers(ctx context.Context, groupID string, studentIDs []string) (int, error) {
	m.removeMembersCalled = true
	return m.removeMembersCount, m.removeMembersErr
}

func (m *mockRepo) ListTutorGroups(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error) {
	if m.listTutorGroupsFunc != nil {
		return m.listTutorGroupsFunc(ctx, tutorID, includeMembers)
	}
	return nil, nil
}

func (m *mockRepo) ListStudentGroups(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error) {
	if m.listStudentGroupsFunc != nil {
		return m.listStudentGroupsFunc(ctx, studentID, includeMembers)
	}
	return nil, nil
}

func (m *mockRepo) GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	if m.getGroupMembersFunc != nil {
		return m.getGroupMembersFunc(ctx, groupID)
	}
	return nil, nil
}

type mockUserClient struct {
	validateCalled bool
	validateResult bool
	validateErr    error
}

func (m *mockUserClient) ValidateTutor(ctx context.Context, tutorId string) (bool, error) {
	m.validateCalled = true
	return m.validateResult, m.validateErr
}

func TestCreateGroup_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	user := &mockUserClient{validateResult: true}
	u := usecase.NewGroupsUsecase(repo, user)

	tutorID := "tutor123"
	group, err := u.CreateGroup(ctx, tutorID, "Math 10A", "Algebra")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if group == nil {
		t.Fatal("group is nil")
	}
	if group.Name != "Math 10A" {
		t.Errorf("wrong name: got %s, want Math 10A", group.Name)
	}
	if !repo.createGroupCalled {
		t.Error("CreateGroup not called")
	}
	if !user.validateCalled {
		t.Error("ValidateTutor not called")
	}
}

func TestCreateGroup_NotTutor(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	user := &mockUserClient{validateResult: false}
	u := usecase.NewGroupsUsecase(repo, user)

	_, err := u.CreateGroup(ctx, "tutor456", "Test", "")

	if err == nil {
		t.Error("expected error")
	}
	if !errors.Is(err, models.ErrTutorIsNotValid) {
		t.Errorf("wrong error: %v", err)
	}
	if repo.createGroupCalled {
		t.Error("CreateGroup called when tutor invalid")
	}
}

func TestUpdateGroup_Success(t *testing.T) {
	ctx := context.Background()
	groupID := "group123"
	tutorID := "tutor123"

	repo := &mockRepo{
		getGroupFirstResult: &models.Group{
			ID:      groupID,
			TutorID: tutorID,
			Name:    "Old Name",
		},
		getGroupSecondResult: &models.Group{
			ID:      groupID,
			TutorID: tutorID,
			Name:    "New Name",
		},
	}
	u := usecase.NewGroupsUsecase(repo, nil)

	newName := "New Name"
	updated, err := u.UpdateGroup(ctx, groupID, tutorID, &newName, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("wrong name: got %s, want New Name", updated.Name)
	}
	if repo.getGroupCallCount != 2 {
		t.Errorf("GetGroup called %d times, want 2", repo.getGroupCallCount)
	}
	if !repo.updateGroupCalled {
		t.Error("UpdateGroup not called")
	}
}

func TestUpdateGroup_NotOwner(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		getGroupFirstResult: &models.Group{
			ID:      "group123",
			TutorID: "tutor123",
		},
	}
	wrongUser := "tutor456"
	u := usecase.NewGroupsUsecase(repo, nil)

	_, err := u.UpdateGroup(ctx, repo.getGroupFirstResult.ID, wrongUser, nil, nil)

	if err == nil {
		t.Error("expected error")
	}
	if !errors.Is(err, models.ErrTutorIsNotValid) {
		t.Errorf("wrong error: %v", err)
	}
	if !repo.getGroupCalled {
		t.Error("GetGroup not called")
	}
	if repo.updateGroupCalled {
		t.Error("UpdateGroup called when not owner")
	}
}

func TestDeleteGroup_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		getGroupFirstResult: &models.Group{
			ID:      "group123",
			TutorID: "tutor123",
		},
	}
	u := usecase.NewGroupsUsecase(repo, nil)

	err := u.DeleteGroup(ctx, repo.getGroupFirstResult.ID, repo.getGroupFirstResult.TutorID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !repo.deleteGroupCalled {
		t.Error("DeleteGroup not called")
	}
}

func TestAddGroupMembers_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		getGroupFirstResult: &models.Group{
			ID:      "group123",
			TutorID: "tutor123",
		},
		addMembersCount: 2,
	}
	u := usecase.NewGroupsUsecase(repo, nil)

	student1 := "student1"
	student2 := "student2"

	count, err := u.AddGroupMembers(ctx, repo.getGroupFirstResult.ID, repo.getGroupFirstResult.TutorID, []string{student1, student2})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("wrong count: got %d, want 2", count)
	}
	if !repo.addMembersCalled {
		t.Error("AddMembers not called")
	}
}

func TestAddGroupMembers_NotOwner(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		getGroupFirstResult: &models.Group{
			ID:      "group123",
			TutorID: "tutor123",
		},
	}
	wrongUser := "tutor456"
	u := usecase.NewGroupsUsecase(repo, nil)

	count, err := u.AddGroupMembers(ctx, repo.getGroupFirstResult.ID, wrongUser, []string{"student1"})

	if err == nil {
		t.Error("expected error")
	}
	if !errors.Is(err, models.ErrTutorIsNotValid) {
		t.Errorf("wrong error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
	if repo.addMembersCalled {
		t.Error("AddMembers called when not owner")
	}
}

func TestRemoveGroupMembers_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{
		getGroupFirstResult: &models.Group{
			ID:      "group123",
			TutorID: "tutor123",
		},
		removeMembersCount: 1,
	}
	u := usecase.NewGroupsUsecase(repo, nil)

	count, err := u.RemoveGroupMembers(ctx, repo.getGroupFirstResult.ID, repo.getGroupFirstResult.TutorID, []string{"student1"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("wrong count: got %d, want 1", count)
	}
	if !repo.removeMembersCalled {
		t.Error("RemoveMembers not called")
	}
}

func TestGetGroup_Success(t *testing.T) {
	ctx := context.Background()
	expectedGroup := &models.Group{
		ID:      "group123",
		TutorID: "tutor123",
		Name:    "Test Group",
	}

	repo := &mockRepo{
		getGroupFirstResult: expectedGroup,
	}
	u := usecase.NewGroupsUsecase(repo, nil)

	group, err := u.GetGroup(ctx, expectedGroup.ID, false)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if group.ID != expectedGroup.ID {
		t.Errorf("wrong group ID: got %s, want %s", group.ID, expectedGroup.ID)
	}
	if !repo.getGroupCalled {
		t.Error("GetGroup not called")
	}
}

func TestListGroupMembers_Success(t *testing.T) {
	ctx := context.Background()
	expectedMembers := []*models.GroupMember{
		{GroupID: "group123", StudentID: "student1"},
		{GroupID: "group123", StudentID: "student2"},
	}

	repo := &mockRepo{
		getGroupMembersFunc: func(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
			if groupID == "group123" {
				return expectedMembers, nil
			}
			return nil, nil
		},
	}

	u := usecase.NewGroupsUsecase(repo, nil)

	members, err := u.ListGroupMembers(ctx, "group123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(members) != 2 {
		t.Errorf("wrong number of members: got %d, want 2", len(members))
	}
}

func TestCreateGroup_UserClientError(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	user := &mockUserClient{validateErr: errors.New("user service error")}
	u := usecase.NewGroupsUsecase(repo, user)

	_, err := u.CreateGroup(ctx, "tutor123", "Test", "")

	if err == nil {
		t.Error("expected error")
	}
	if !user.validateCalled {
		t.Error("ValidateTutor not called")
	}
	if repo.createGroupCalled {
		t.Error("CreateGroup called when user client error")
	}
}

func TestListGroupsByTutor_Success(t *testing.T) {
	ctx := context.Background()
	expectedGroups := []*models.Group{
		{ID: "group1", TutorID: "tutor123", Name: "Group 1"},
		{ID: "group2", TutorID: "tutor123", Name: "Group 2"},
	}

	repo := &mockRepo{
		listTutorGroupsFunc: func(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error) {
			if tutorID == "tutor123" {
				return expectedGroups, nil
			}
			return nil, nil
		},
	}
	user := &mockUserClient{validateResult: true}
	u := usecase.NewGroupsUsecase(repo, user)

	groups, err := u.ListGroupsByTutor(ctx, "tutor123", false)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("wrong number of groups: got %d, want 2", len(groups))
	}
	if !user.validateCalled {
		t.Error("ValidateTutor not called")
	}
}

func TestListGroupsByStudent_Success(t *testing.T) {
	ctx := context.Background()
	expectedGroups := []*models.Group{
		{ID: "group1", TutorID: "tutor123", Name: "Group 1"},
	}

	repo := &mockRepo{
		listStudentGroupsFunc: func(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error) {
			if studentID == "student123" {
				return expectedGroups, nil
			}
			return nil, nil
		},
	}
	u := usecase.NewGroupsUsecase(repo, nil)

	groups, err := u.ListGroupsByStudent(ctx, "student123", false)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("wrong number of groups: got %d, want 1", len(groups))
	}
}
