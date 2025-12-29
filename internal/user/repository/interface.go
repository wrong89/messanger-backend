package repository

import (
	"context"
	"messanger/internal/user"
)

type UserRepo interface {
	UserReader
	UserWriter
}

type UserReader interface {
	GetByID(ctx context.Context, id uint64) (user.User, error)
	GetByLogin(ctx context.Context, login string) (user.User, error)
}

type UserWriter interface {
	Create(ctx context.Context, user user.User) (uint64, error)
	Update(ctx context.Context, newUser user.User) error
	Delete(ctx context.Context, id uint64) error
}
