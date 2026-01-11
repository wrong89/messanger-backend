package repository

import (
	"context"
	"messanger/internal/chat"
)

type ChatRepo interface {
	ChatReader
	ChatWriter
	ChatUserActions
}

type ChatReader interface {
	GetByID(ctx context.Context, id uint64) (chat.Chat, error)
	GetByAddress(ctx context.Context, address string) (chat.Chat, error)
}

type ChatWriter interface {
	Create(ctx context.Context, chat chat.Chat) (uint64, error)
	// Update(ctx context.Context, newChat chat.Chat) error
	Delete(ctx context.Context, id uint64) error
}

type ChatUserActions interface {
	Join(ctx context.Context, role string, userID uint64, chatID uint64) error
	Leave(ctx context.Context, userID uint64, chatID uint64) error
}
