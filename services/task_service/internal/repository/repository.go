package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.crja72.ru/aisavelev-edu.hse.ru/internal/models"
)

type Repository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Repository{
		pool:    pool,
		builder: builder,
	}
}

func (r *Repository) CreateTask(ctx context.Context, task models.AssignedTask) (*models.AssignedTask, error) {
	descriptionVal := interface{}(squirrel.Expr("NULL"))
	if task.Description != nil {
		descriptionVal = *task.Description
	}

	query, args, err := r.builder.
		Insert("assigned_tasks").
		Columns("group_id", "tutor_id", "title", "description",
			"max_score", "deadline", "task_status", "created_at").
		Values(task.GroupId, task.TutorId, task.Title,
			descriptionVal, task.MaxScore,
			task.Deadline, task.Status, task.CreatedAt).
		Suffix("RETURNING id, group_id, tutor_id, title, description, max_score, deadline, task_status, created_at").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var createdTask models.AssignedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&createdTask.ID, &createdTask.GroupId, &createdTask.TutorId,
		&createdTask.Title, &createdTask.Description, &createdTask.MaxScore,
		&createdTask.Deadline, &createdTask.Status, &createdTask.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &createdTask, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID string) (*models.AssignedTask, error) {
	query, args, err := r.builder.
		Select("id", "group_id", "tutor_id", "title", "description",
			"max_score", "deadline", "task_status",
			"created_at", "updated_at").
		From("assigned_tasks").
		Where(squirrel.Eq{"id": taskID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var task models.AssignedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&task.ID, &task.GroupId, &task.TutorId, &task.Title,
		&task.Description, &task.MaxScore, &task.Deadline,
		&task.Status, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.AssignedTask, error) {
	updateBuilder := r.builder.Update("assigned_tasks").
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": req.TaskID})

	if req.Title != nil {
		updateBuilder = updateBuilder.Set("title", *req.Title)
	}
	if req.Description != nil {
		updateBuilder = updateBuilder.Set("description", *req.Description)
	}
	if req.MaxScore != nil {
		updateBuilder = updateBuilder.Set("max_score", *req.MaxScore)
	}
	if req.Deadline != nil {
		updateBuilder = updateBuilder.Set("deadline", *req.Deadline)
	}

	query, args, err := updateBuilder.
		Suffix("RETURNING id, group_id, tutor_id, title, description, max_score, deadline, task_status, created_at, updated_at").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var task models.AssignedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&task.ID, &task.GroupId, &task.TutorId, &task.Title,
		&task.Description, &task.MaxScore, &task.Deadline,
		&task.Status, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &task, nil
}

func (r *Repository) SoftDeleteTask(ctx context.Context, userID, taskID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	deleteSubmissionsQuery, deleteSubmissionsArgs, err := r.builder.
		Delete("submitted_tasks").
		Where(squirrel.Eq{"task_id": taskID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build delete submissions query: %w", err)
	}

	_, err = tx.Exec(ctx, deleteSubmissionsQuery, deleteSubmissionsArgs...)
	if err != nil {
		return fmt.Errorf("delete submissions: %w", err)
	}

	deleteTaskQuery, deleteTaskArgs, err := r.builder.
		Delete("assigned_tasks").
		Where(squirrel.And{
			squirrel.Eq{"id": taskID},
			squirrel.Eq{"tutor_id": userID},
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build delete task query: %w", err)
	}

	result, err := tx.Exec(ctx, deleteTaskQuery, deleteTaskArgs...)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetTasks(ctx context.Context, filter models.TaskFilter) ([]*models.AssignedTaskShort, int32, error) {
	baseQuery := r.builder.
		Select("id", "group_id", "tutor_id", "title", "deadline", "task_status").
		From("assigned_tasks")

	if filter.GroupID != "" {
		baseQuery = baseQuery.Where(squirrel.Eq{"group_id": filter.GroupID})
	}
	if filter.Type == models.CreatedType {
		baseQuery = baseQuery.Where(squirrel.Eq{"tutor_id": filter.UserID})
	}

	countQuery := baseQuery.RemoveColumns().Column("COUNT(*)")
	countSQL, countArgs, _ := countQuery.ToSql()

	var total int32
	err := r.pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	baseQuery = baseQuery.
		OrderBy("created_at DESC").
		Offset(uint64(filter.Offset)).
		Limit(uint64(filter.Limit))

	query, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var tasks []*models.AssignedTaskShort
	for rows.Next() {
		var task models.AssignedTaskShort
		err := rows.Scan(
			&task.ID, &task.GroupID, &task.TutorID,
			&task.Title, &task.Deadline, &task.Status,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan row: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, total, nil
}

func (r *Repository) CreateSubmission(ctx context.Context, submission models.SubmittedTask) (*models.SubmittedTask, error) {
	query, args, err := r.builder.
		Insert("submitted_tasks").
		Columns("task_id", "student_id", "content",
			"status", "created_at").
		Values(submission.TaskID, submission.StudentID,
			submission.Content, submission.Status,
			submission.CreatedAt).
		Suffix(`RETURNING id, task_id, student_id, content, 
                status, score, tutor_feedback, created_at, updated_at, overdue_by`).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var createdSubmission models.SubmittedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&createdSubmission.ID, &createdSubmission.TaskID, &createdSubmission.StudentID,
		&createdSubmission.Content, &createdSubmission.Status, &createdSubmission.Score,
		&createdSubmission.Feedback, &createdSubmission.CreatedAt, &createdSubmission.UpdatedAt,
		&createdSubmission.OverdueBy,
	)

	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &createdSubmission, nil
}

func (r *Repository) GetSubmissionByID(ctx context.Context, submissionID string) (*models.SubmittedTask, error) {
	query, args, err := r.builder.
		Select("id", "task_id", "student_id", "content",
			"status", "score", "tutor_feedback", "created_at",
			"updated_at", "overdue_by").
		From("submitted_tasks").
		Where(squirrel.Eq{"id": submissionID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var submission models.SubmittedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&submission.ID, &submission.TaskID, &submission.StudentID,
		&submission.Content, &submission.Status, &submission.Score,
		&submission.Feedback, &submission.CreatedAt, &submission.UpdatedAt,
		&submission.OverdueBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &submission, nil
}

func (r *Repository) GradeSubmission(ctx context.Context, grade models.SubmissionGrade) (*models.SubmittedTask, error) {
	updateBuilder := r.builder.Update("submitted_tasks").
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": grade.SubmissionId})

	updateBuilder = updateBuilder.Set("status", models.SubmissionStatusVerified)

	if grade.Score != nil {
		updateBuilder = updateBuilder.Set("score", *grade.Score)
	}
	if grade.Feedback != nil {
		updateBuilder = updateBuilder.Set("tutor_feedback", *grade.Feedback)
	}
	if grade.Status != nil {
		updateBuilder = updateBuilder.Set("status", *grade.Status)
	}

	query, args, err := updateBuilder.
		Suffix(`RETURNING id, task_id, student_id, content, 
                status, score, tutor_feedback, created_at, updated_at, overdue_by`).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var submission models.SubmittedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&submission.ID, &submission.TaskID, &submission.StudentID,
		&submission.Content, &submission.Status, &submission.Score,
		&submission.Feedback, &submission.CreatedAt, &submission.UpdatedAt,
		&submission.OverdueBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &submission, nil
}

func (r *Repository) ResetGrade(ctx context.Context, userID, submissionID string) error {
	query, args, err := r.builder.
		Update("submitted_tasks").
		Set("status", models.SubmissionStatusPending).
		Set("score", squirrel.Expr("NULL")).
		Set("tutor_feedback", squirrel.Expr("NULL")).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": submissionID}).
		Where(squirrel.Eq{"status": models.SubmissionStatusVerified}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) UpdateSubmission(ctx context.Context, req models.UpdateSubmissionRequest) (*models.SubmittedTask, error) {
	updateBuilder := r.builder.Update("submitted_tasks").
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": req.SubmussionID}).
		Where(squirrel.NotEq{"status": models.SubmissionStatusVerified})

	if req.Content != nil {
		updateBuilder = updateBuilder.Set("content", *req.Content)
	}

	query, args, err := updateBuilder.
		Suffix(`RETURNING id, task_id, student_id, content, 
                status, score, tutor_feedback, created_at, updated_at, overdue_by`).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var submission models.SubmittedTask
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&submission.ID, &submission.TaskID, &submission.StudentID,
		&submission.Content, &submission.Status, &submission.Score,
		&submission.Feedback, &submission.CreatedAt, &submission.UpdatedAt,
		&submission.OverdueBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			var exists bool
			checkQuery := "SELECT EXISTS(SELECT 1 FROM submitted_tasks WHERE id = $1)"
			r.pool.QueryRow(ctx, checkQuery, req.SubmussionID).Scan(&exists)

			if exists {
				return nil, models.ErrAlreadyGraded
			}
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &submission, nil
}

func (r *Repository) DeleteSubmission(ctx context.Context, userID, submissionID string) error {
	query, args, err := r.builder.
		Delete("submitted_tasks").
		Where(squirrel.And{
			squirrel.Eq{"id": submissionID},
			squirrel.Eq{"student_id": userID},
			squirrel.NotEq{"status": models.SubmissionStatusVerified},
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *Repository) GetSubmissions(ctx context.Context, filter models.SubmissionFilter) ([]*models.SubmittedTaskShort, int32, error) {
	baseQuery := r.builder.
		Select("id", "task_id", "student_id", "score", "status", "created_at").
		From("submitted_tasks")

	if filter.TaskID != "" {
		baseQuery = baseQuery.Where(squirrel.Eq{"task_id": filter.TaskID})
	}
	if filter.UserID != "" {
		baseQuery = baseQuery.Where(squirrel.Eq{"student_id": filter.UserID})
	}

	countQuery := baseQuery.RemoveColumns().Column("COUNT(*)")
	countSQL, countArgs, _ := countQuery.ToSql()

	var total int32
	err := r.pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count submissions: %w", err)
	}

	baseQuery = baseQuery.
		OrderBy("created_at DESC").
		Offset(uint64(filter.Offset)).
		Limit(uint64(filter.Limit))

	query, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var submissions []*models.SubmittedTaskShort
	for rows.Next() {
		var submission models.SubmittedTaskShort
		err := rows.Scan(
			&submission.ID, &submission.TaskID, &submission.StudentID,
			&submission.Score, &submission.Status, &submission.SubmittedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan row: %w", err)
		}
		submissions = append(submissions, &submission)
	}

	return submissions, total, nil
}

func (r *Repository) MarkExpiredTasks(ctx context.Context) error {
	query, args, err := r.builder.Select("1").
		From("mark_expired_tasks()").
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %v", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("mark expired tasks: %w", err)
	}

	return nil
}

func (r *Repository) Close() {
	r.pool.Close()
}
