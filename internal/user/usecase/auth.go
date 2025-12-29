package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"messanger/internal/lib/jwt"
	"messanger/internal/lib/logger/sl"
	"messanger/internal/user"
	"messanger/internal/user/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthUC interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, name, login, password string) (uint64, error)
}

type Auth struct {
	log      *slog.Logger
	userRepo repository.UserRepo
	secret   string
	tokenTTL time.Duration
}

func NewAuth(log *slog.Logger, userRepo repository.UserRepo, secret string, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		userRepo: userRepo,
		secret:   secret,
		tokenTTL: tokenTTL,
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

func (a *Auth) Login(ctx context.Context, login, password string) (string, error) {
	const op = "user.usecase.auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	usr, err := a.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(usr, a.secret, a.tokenTTL)
	if err != nil {
		log.Info("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
