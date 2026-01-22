package adapter

import (
	"context"
	"errors"
	"fmt"
	"group_service/internal/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupsRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewGroupsRepo(db *pgxpool.Pool) *GroupsRepo {
	return &GroupsRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *GroupsRepo) CreateGroup(ctx context.Context, group *models.Group) error {
	query, args, err := r.builder.Insert("student_groups").
		Columns("id", "tutor_id", "name", "description", "created_at").
		Values(group.ID, group.TutorID, group.Name, group.Description, group.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert group: %w", err)
	}

	return nil
}

func (r *GroupsRepo) ListTutorGroups(ctx context.Context, tutorID string, includeMembers bool) ([]*models.Group, error) {
	query, args, err := r.builder.Select("id", "name", "description", "created_at").
		From("student_groups").
		Where(squirrel.Eq{"tutor_id": tutorID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tutor groups: %w", err)
	}
	defer rows.Close()

	groups := make([]*models.Group, 0)
	for rows.Next() {
		g := &models.Group{TutorID: tutorID}
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if includeMembers {
		for _, g := range groups {
			g.Members, err = r.GetGroupMembers(ctx, g.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get members for group %s: %w", g.ID, err)
			}
		}
	}

	return groups, nil
}

func (r *GroupsRepo) ListStudentGroups(ctx context.Context, studentID string, includeMembers bool) ([]*models.Group, error) {
	query, args, err := r.builder.Select("sg.id", "sg.tutor_id", "sg.name", "sg.description", "sg.created_at").
		From("student_groups sg").
		Join("group_members gm ON gm.group_id = sg.id").
		Where(squirrel.Eq{"gm.student_id": studentID}).
		OrderBy("sg.created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query student groups: %w", err)
	}
	defer rows.Close()

	groups := make([]*models.Group, 0)
	for rows.Next() {
		g := &models.Group{}
		if err := rows.Scan(&g.ID, &g.TutorID, &g.Name, &g.Description, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// только если запрошено
	if includeMembers {
		for _, g := range groups {
			g.Members, err = r.GetGroupMembers(ctx, g.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get members for group %s: %w", g.ID, err)
			}
		}
	}

	return groups, nil
}

func (r *GroupsRepo) GetGroup(ctx context.Context, id string, includeMembers bool) (*models.Group, error) {
	query, args, err := r.builder.Select("tutor_id", "name", "description", "created_at").
		From("student_groups").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	g := &models.Group{ID: id}
	err = r.db.QueryRow(ctx, query, args...).Scan(&g.TutorID, &g.Name, &g.Description, &g.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("group with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to query group: %w", err)
	}

	if includeMembers {
		g.Members, err = r.GetGroupMembers(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get members: %w", err)
		}
	} else {
		g.Members = make([]*models.GroupMember, 0)
	}

	return g, nil
}

func (r *GroupsRepo) UpdateGroup(ctx context.Context, id string, name, desc *string) error {
	updateBuilder := r.builder.Update("student_groups")

	hasUpdates := false
	if name != nil {
		updateBuilder = updateBuilder.Set("name", *name)
		hasUpdates = true
	}
	if desc != nil {
		updateBuilder = updateBuilder.Set("description", *desc)
		hasUpdates = true
	}

	if !hasUpdates {
		return nil // Нет полей для обновления
	}

	query, args, err := updateBuilder.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("group with id %s not found", id)
	}

	return nil
}

func (r *GroupsRepo) DeleteGroup(ctx context.Context, id string) error {
	if err := r.deleteGroupMembers(ctx, id); err != nil {
		return fmt.Errorf("failed to delete group members: %w", err)
	}

	query, args, err := r.builder.Delete("student_groups").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("group with id %s not found", id)
	}

	return nil
}

func (r *GroupsRepo) deleteGroupMembers(ctx context.Context, groupID string) error {
	query, args, err := r.builder.Delete("group_members").
		Where(squirrel.Eq{"group_id": groupID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete members: %w", err)
	}

	return nil
}

func (r *GroupsRepo) AddMembers(ctx context.Context, groupID string, studentIDs []string) (int, error) {
	if len(studentIDs) == 0 {
		return 0, nil
	}

	now := time.Now()
	insertBuilder := r.builder.Insert("group_members").
		Columns("student_id", "group_id", "joined_at")

	for _, sid := range studentIDs {
		insertBuilder = insertBuilder.Values(sid, groupID, now)
	}

	query, args, err := insertBuilder.
		Suffix("ON CONFLICT (group_id, student_id) DO NOTHING").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert members: %w", err)
	}

	return int(res.RowsAffected()), nil
}

func (r *GroupsRepo) RemoveMembers(ctx context.Context, groupID string, studentIDs []string) (int, error) {
	if len(studentIDs) == 0 {
		return 0, nil
	}

	query, args, err := r.builder.Delete("group_members").
		Where(squirrel.And{
			squirrel.Eq{"group_id": groupID},
			squirrel.Eq{"student_id": studentIDs},
		}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build delete query: %w", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete members: %w", err)
	}

	return int(res.RowsAffected()), nil
}

func (r *GroupsRepo) GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	query, args, err := r.builder.Select("student_id", "joined_at").
		From("group_members").
		Where(squirrel.Eq{"group_id": groupID}).
		OrderBy("joined_at ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query group members: %w", err)
	}
	defer rows.Close()

	members := make([]*models.GroupMember, 0)
	for rows.Next() {
		m := &models.GroupMember{GroupID: groupID}
		if err := rows.Scan(&m.StudentID, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return members, nil
}
