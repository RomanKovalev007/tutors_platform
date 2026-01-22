package repository

import (
	"context"
	"database/sql"
	"user-service/internal/models"
	"user-service/pkg/postgres"
)

type studentProfileRepository struct {
    db *sql.DB
}

func NewStudentProfileRepository(db *sql.DB) *studentProfileRepository {
    return &studentProfileRepository{db: db}
}

func (r *studentProfileRepository) CreateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	setStatusQuery := `
        UPDATE user_profiles 
        SET is_student = $1
        WHERE user_id = $2
    `
    
    _, err = tx.ExecContext(ctx, setStatusQuery, true, student.UserID)
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

	createQuery := `
        INSERT INTO student_profiles
		(user_id, bio, grade)
        VALUES ($1, $2, $3)
		RETURNING user_id, bio, grade
    `

	err = tx.QueryRowContext(ctx, createQuery,
		student.UserID, 
        student.Bio,
		student.Grade).
		Scan(
		&student.UserID, 
        &student.Bio,
		student.Grade)
	

	if err != nil{
		if postgres.IsDuplicateKeyError(err){
			return nil, ErrUserExists
		}
		return nil, err
	}

	
	if err = tx.Commit(); err != nil {
		return nil, err
	}

    return student, nil
}

func (r *studentProfileRepository) GetStudentProfileByID(ctx context.Context, id string) (*models.StudentProfile, error) {
    query := `
        SELECT user_id, bio, grade
        FROM student_profiles 
        WHERE user_id = $1
    `
    
    var student models.StudentProfile
    err := r.db.QueryRowContext(ctx, query, id).
		Scan(
		&student.UserID, 
        &student.Bio,
		student.Grade)
    
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

    
    return &student, nil
}


func (r *studentProfileRepository) UpdateStudentProfile(ctx context.Context, student *models.StudentProfile) (*models.StudentProfile, error) {
	query := `
        UPDATE student_profiles
        SET bio = $1, grade = $2
        WHERE user_id = $3 
        RETURNING user_id, bio, grade
    `

    err := r.db.QueryRowContext(ctx, query, student.Bio, student.Grade, student.UserID).
		Scan(
		&student.UserID, 
        &student.Bio,
		student.Grade)
    
	if err != nil{
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err 
	}

    return student, nil
}

func (r *studentProfileRepository) DeleteStudentProfie(ctx context.Context, id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	setStatusQuery := `
        UPDATE user_profiles
        SET is_student = $1
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
        DELETE FROM student_profiles
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