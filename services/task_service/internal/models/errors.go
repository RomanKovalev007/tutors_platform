package models

import (
	"errors"

	"google.golang.org/grpc/codes"
)

type Error struct {
	Code    codes.Code
	Message string
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyGraded = errors.New("submission already graded")
)
