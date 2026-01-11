package repository

import (
	"context"
	"fmt"
	"messanger/internal/chat"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context, dbURL string) (*Storage, error) {
	const op = "chat.repository.postgres.New"

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close(ctx)
}

func (s *Storage) Create(ctx context.Context, chat chat.Chat) (uint64, error) {
	const op = "chat.repository.postgres.Create"

	sql := `INSERT INTO chats(type, address) VALUES(@type, @address) RETURNING id;`
	args := pgx.NamedArgs{
		"type":    chat.Type,
		"address": chat.Address,
	}

	var chatID uint64

	if err := s.db.QueryRow(ctx, sql, args).Scan(&chatID); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return chatID, nil
}

// func (s *Storage) Update(ctx context.Context, newChat chat.Chat) error {
// 	const op = "chat.repository.postgres.Update"

// 	sql := `UPDATE chats SET `
// 	args := pgx.NamedArgs{
// 		"id": id,
// 	}

// 	return nil
// }

func (s *Storage) Delete(ctx context.Context, id uint64) error {
	const op = "chat.repository.postgres.Delete"

	_, err := s.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sql := `DELETE FROM chats WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err = s.db.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, id uint64) (chat.Chat, error) {
	const op = "chat.repository.postgres.GetByID"

	sql := `SELECT id, type, address, created_at FROM chats WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	var cht chat.Chat

	err := s.db.QueryRow(ctx, sql, args).Scan(
		&cht.ID,
		&cht.Type,
		&cht.Address,
		&cht.CreatedAt,
	)
	if err != nil {
		return chat.Chat{}, fmt.Errorf("%s: %w", op, err)
	}

	return cht, nil
}

func (s *Storage) GetByAddress(ctx context.Context, address string) (chat.Chat, error) {
	const op = "chat.repository.postgres.GetByAddress"

	sql := `SELECT id, type, address, created_at FROM chats WHERE address = @address`
	args := pgx.NamedArgs{
		"address": address,
	}

	var cht chat.Chat

	err := s.db.QueryRow(ctx, sql, args).Scan(
		&cht.ID,
		&cht.Type,
		&cht.Address,
		&cht.CreatedAt,
	)
	if err != nil {
		return chat.Chat{}, fmt.Errorf("%s: %w", op, err)
	}

	return cht, nil
}

func (s *Storage) Join(ctx context.Context, role string, userID uint64, chatID uint64) error {
	const op = "chat.repository.postgres.Join"

	sql := `INSERT INTO chat_members(role, chat_id, user_id) VALUES(@role, @chat_id, @user_id)`
	args := pgx.NamedArgs{
		"role":    role,
		"chat_id": chatID,
		"user_id": userID,
	}

	if _, err := s.db.Exec(ctx, sql, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Leave(ctx context.Context, userID uint64, chatID uint64) error {
	const op = "chat.repository.postgres.Leave"

	sql := `DELETE FROM chat_members WHERE user_id = @user_id AND chat_id = @chat_id`
	args := pgx.NamedArgs{
		"user_id": userID,
		"chat_id": chatID,
	}

	if _, err := s.db.Exec(ctx, sql, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
