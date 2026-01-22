package repository

import (
	"context"
	"database/sql"
	"user-service/internal/models"
	"user-service/pkg/postgres"
)

type tutorProfileRepository struct {
    db *sql.DB
}

func NewTutorProfileRepository(db *sql.DB) *tutorProfileRepository {
    return &tutorProfileRepository{db: db}
}

func (r *tutorProfileRepository) CreateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	setStatusQuery := `
        UPDATE user_profiles
        SET is_tutor = $1
        WHERE user_id = $2
    `
    
    _, err = tx.ExecContext(ctx, setStatusQuery, true, tutor.UserID)
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

	createQuery := `
        INSERT INTO tutor_profiles
		(user_id, bio, specialization, experience_years)
        VALUES ($1, $2, $3, $4)
		RETURNING user_id, bio, specialization, experience_years
    `

	err = tx.QueryRowContext(ctx, createQuery,
		tutor.UserID, 
        tutor.Bio,
		tutor.Specialization,
		tutor.Experience).
		Scan(
		&tutor.UserID, 
        &tutor.Bio,
		&tutor.Specialization,
		&tutor.Experience)
	

	if err != nil{
		if postgres.IsDuplicateKeyError(err){
			return nil, ErrUserExists
		}
		return nil, err
	}

	
	if err = tx.Commit(); err != nil {
		return nil, err
	}

    return tutor, nil
}

func (r *tutorProfileRepository) GetTutorProfileByID(ctx context.Context, id string) (*models.TutorProfile, error) {
    query := `
        SELECT user_id, bio, specialization, experience_years
        FROM tutor_profiles 
        WHERE user_id = $1
    `
    
    var tutor models.TutorProfile
    err := r.db.QueryRowContext(ctx, query, id).
		Scan(
		&tutor.UserID, 
        &tutor.Bio,
		&tutor.Specialization,
		&tutor.Experience)
    
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

    
    return &tutor, nil
}


func (r *tutorProfileRepository) UpdateTutorProfile(ctx context.Context, tutor *models.TutorProfile) (*models.TutorProfile, error) {
	query := `
        UPDATE tutor_profiles
        SET bio = $1, specialization = $2, experience_years = $3
        WHERE user_id = $4 
        RETURNING user_id, bio, specialization, experience_years
    `

    err := r.db.QueryRowContext(ctx, query, tutor.Bio, tutor.Specialization, tutor.Experience, tutor.UserID).
		Scan(
		&tutor.UserID, 
        &tutor.Bio,
		&tutor.Specialization,
		&tutor.Experience)
    
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

    return tutor, nil
}

func (r *tutorProfileRepository) DeleteTutorProfie(ctx context.Context, id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	setStatusQuery := `
        UPDATE user_profiles 
        SET is_tutor = $1
        WHERE user_id = $2
    `
    
    _, err = tx.ExecContext(ctx, setStatusQuery, false, id)
	if err != nil{
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err 
	}

	query := `
        DELETE FROM tutor_profiles
		WHERE user_id = $1
    `

	_, err = tx.ExecContext(ctx, query, id)
	

	if err != nil{
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

    return nil
}