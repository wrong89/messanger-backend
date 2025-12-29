package repository

import (
	"context"
	"errors"
	"messanger/internal/user"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type UserRepo interface {
	Create(ctx context.Context, user user.User) (uint64, error)
	GetByID(ctx context.Context, id uint64) (user.User, error)
	GetByLogin(ctx context.Context, login string) (user.User, error)
	Update(ctx context.Context, user user.User) error
	Delete(ctx context.Context, id uint64) error
}
