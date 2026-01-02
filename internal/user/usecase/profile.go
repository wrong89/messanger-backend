package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"messanger/internal/user/repository"
)

type ProfileUC interface {
	Delete(context.Context, uint64) error
}

type Profile struct {
	log      *slog.Logger
	userRepo repository.UserRepo
}

func NewProfile(log *slog.Logger, userRepo repository.UserRepo) *Profile {
	return &Profile{
		log:      log,
		userRepo: userRepo,
	}
}

func (a *Profile) Delete(ctx context.Context, id uint64) error {
	const op = "user.usecase.profile.Delete"

	log := a.log.With(
		slog.String("op", op),
		slog.Uint64("id", id),
	)

	if err := a.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("user not found")

			return fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
