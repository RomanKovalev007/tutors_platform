package repository

import (
	"auth_service/internal/models"
	"auth_service/pkg/postgres"
	"context"
	"database/sql"
)

type userRepository struct{
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
    return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error){
	query := `
		INSERT INTO users
		(id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, is_active, created_at
	`
	
	err := r.db.QueryRowContext(ctx, query, 
		user.ID,
		user.Email,
		user.Password).
		Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt)

	if err != nil{
		if postgres.IsDuplicateKeyError(err){
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
    query := `
        SELECT id, email, is_active, created_at
        FROM users 
        WHERE id = $1
    `
    
    var user models.User
    err := r.db.QueryRowContext(ctx, query, id).
		Scan(
        &user.ID,
		&user.Email,
		&user.IsActive,
        &user.CreatedAt)
    
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

    
    return &user, nil
}


func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `
        SELECT id, email, password_hash, is_active, created_at
        FROM users 
        WHERE email = $1
    `
    
    var user models.User
    err := r.db.QueryRowContext(ctx, query, email).
		Scan(
        &user.ID,
		&user.Email,
		&user.Password,
		&user.IsActive,
        &user.CreatedAt)
    
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

    
    return &user, nil
}


func (r *userRepository) UpdatePassword(ctx context.Context, id string, password string) error {
    query := `
        UPDATE users 
        SET password_hash = $1
        WHERE id = $2
    `

    res, err := r.db.ExecContext(ctx, query, password, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrUserNotFound
    }
    
    return nil
}


func (r *userRepository) SetIsActive(ctx context.Context, id string, status bool) error {
	query := `
        UPDATE users
        SET is_active = $1
        WHERE id = $2
    `
    
	res, err := r.db.ExecContext(ctx, query, status, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrUserNotFound
    }
    
    return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
        DELETE FROM users
		WHERE id = $1
    `
	res, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrUserNotFound
    }

    return nil
}

func (r *userRepository) SelectAllUsers(ctx context.Context, limit, offset int32, isActive bool) ([]models.User, error){
	query := `
		SELECT id, email, is_active, created_at
		FROM users
		WHERE is_active = $1
		ORDER BY created_at DESC, id
		LIMIT $2
		OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, isActive, limit, offset)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next(){
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.IsActive, &user.CreatedAt); err != nil{
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return users, nil
}