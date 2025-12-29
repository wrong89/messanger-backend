package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"messanger/internal/lib/logger/sl"
	"messanger/internal/user"
	"messanger/internal/user/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUC interface {
	Login(ctx context.Context, login, password string) (user.User, error)
	Register(ctx context.Context, name, login, password string) (uint64, error)
}

type Auth struct {
	log      *slog.Logger
	userRepo repository.UserRepo
}

func NewAuth(log *slog.Logger, userRepo repository.UserRepo) *Auth {
	return &Auth{
		log:      log,
		userRepo: userRepo,
	}
}

func (a *Auth) Register(ctx context.Context, name, login, password string) (uint64, error) {
	const op = "user.usecase.auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	_, err := a.userRepo.GetByLogin(ctx, login)
	if err == nil {
		log.Warn("user already exist", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, repository.ErrUserAlreadyExist)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("hash generation error", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	newUser := user.NewUser(name, login, passHash)

	id, err := a.userRepo.Create(ctx, newUser)
	if err != nil {
		log.Error("user creation error", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
