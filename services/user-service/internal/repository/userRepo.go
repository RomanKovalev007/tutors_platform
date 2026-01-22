package repository

import (
	"context"
	"database/sql"
	"user-service/internal/models"
	"user-service/pkg/postgres"
)

type userProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *userProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) SelectAllUsers(ctx context.Context, limit, offset int32) ([]models.UserProfile, error) {
	query := `
        SELECT user_id, email, name, surname, is_tutor, is_student, created_at, telegram
        FROM user_profiles
		LIMIT $1
		OFFSET $2
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserProfile
	for rows.Next() {
		var user models.UserProfile
		if err := rows.Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.IsTutor,
			&user.IsStudent,
			&user.CreatedAt,
			&user.Telegram); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}

func (r *userProfileRepository) CreateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, error) {
	query := `
        INSERT INTO user_profiles 
		(user_id, email, name, surname, telegram)
        VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
		RETURNING user_id, email, name, surname, is_tutor, is_student, created_at, telegram
    `

	err := r.db.QueryRowContext(ctx, query,
		user.UserID,
		user.Email,
		user.Name,
		user.Surname,
		user.Telegram).
		Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.IsTutor,
			&user.IsStudent,
			&user.CreatedAt,
			&user.Telegram)

	if err != nil {
		if postgres.IsDuplicateKeyError(err) {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

func (r *userProfileRepository) GetUserByID(ctx context.Context, id string) (*models.UserProfile, error) {
	query := `
        SELECT user_id, email, name, surname, is_tutor, is_student, created_at, telegram
        FROM user_profiles 
        WHERE user_id = $1
    `

	var user models.UserProfile
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.IsTutor,
			&user.IsStudent,
			&user.CreatedAt,
			&user.Telegram)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userProfileRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserProfile, error) {
	query := `
        SELECT user_id, email, name, surname, is_tutor, is_student, created_at, telegram
        FROM user_profiles 
        WHERE email = $1
    `

	var user models.UserProfile
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.IsTutor,
			&user.IsStudent,
			&user.CreatedAt,
			&user.Telegram)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userProfileRepository) UpdateUser(ctx context.Context, user *models.UserProfile) (*models.UserProfile, error) {
	query := `
        UPDATE user_profiles 
        SET name = $1, surname = $2, telegram = $3
        WHERE user_id = $4
        RETURNING user_id, email, name, surname, is_tutor, is_student, created_at, telegram
    `

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Surname, user.Telegram, user.UserID).
		Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.IsTutor,
			&user.IsStudent,
			&user.CreatedAt,
			&user.Telegram)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userProfileRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
        DELETE FROM user_profiles
		WHERE user_id = $1
    `

	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (r *userProfileRepository) GetUserTypes(ctx context.Context, id string) (*models.UserType, error) {
	query := `
        SELECT is_tutor, is_student
        FROM user_profiles 
        WHERE user_id = $1
    `

	var userType models.UserType
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&userType.IsTutor, &userType.IsStudent)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &userType, nil
}
