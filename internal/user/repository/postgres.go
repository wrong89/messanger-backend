package repository

import (
	"context"
	"fmt"
	"messanger/internal/user"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context, dbURL string) (*Storage, error) {
	const op = "user.repository.postgres.New"

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close(ctx)
}

func (s *Storage) CreateUser(ctx context.Context, user user.User) (uint64, error) {
	const op = "user.repository.postgres.New"

	sql := `INSERT INTO users(name, login, password_hash) VALUES(@name, @login, @password_hash) RETURNING id;`
	args := pgx.NamedArgs{
		"name":          user.Name,
		"login":         user.Login,
		"password_hash": user.PasswordHash,
	}

	var usrID uint64

	err := s.db.QueryRow(ctx, sql, args).Scan(&usrID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return usrID, nil
}
