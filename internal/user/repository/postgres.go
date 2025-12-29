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

func (s *Storage) Create(ctx context.Context, user user.User) (uint64, error) {
	const op = "user.repository.postgres.Create"

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

func (s *Storage) Update(ctx context.Context, newUser user.User) error {
	const op = "user.repository.postgres.Update"

	_, err := s.GetByID(ctx, newUser.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sql := `UPDATE users SET name = @name, login = @login, password_hash = @password_hash WHERE id = @id`
	args := pgx.NamedArgs{
		"id":            newUser.ID,
		"name":          newUser.Name,
		"login":         newUser.Login,
		"password_hash": newUser.PasswordHash,
	}

	_, err = s.db.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id uint64) error {
	const op = "user.repository.postgres.Delete"

	_, err := s.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sql := `DELETE FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err = s.db.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, id uint64) (user.User, error) {
	const op = "user.repository.postgres.GetByID"

	sql := `SELECT id, name, login, created_at FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	var usr user.User

	err := s.db.QueryRow(ctx, sql, args).Scan(
		&usr.ID,
		&usr.Name,
		&usr.Login,
		&usr.CreatedAt,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return usr, nil
}

func (s *Storage) GetByLogin(ctx context.Context, login string) (user.User, error) {
	const op = "user.repository.postgres.GetByLogin"

	sql := `SELECT id, name, login, created_at FROM users WHERE login = @login`
	args := pgx.NamedArgs{
		"login": login,
	}

	var usr user.User

	err := s.db.QueryRow(ctx, sql, args).Scan(
		&usr.ID,
		&usr.Name,
		&usr.Login,
		&usr.CreatedAt,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return usr, nil
}
