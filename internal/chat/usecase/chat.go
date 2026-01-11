package usecase

import (
	"context"
	"log/slog"
	"messanger/internal/chat"
	"messanger/internal/chat/repository"
)

type ChatUC interface {
	CreateChannel(ctx context.Context, address string) (uint64, error)
	CreateGroup(ctx context.Context, address string) (uint64, error)
	CreatePrivate(ctx context.Context, address string) (uint64, error)
}

type Chat struct {
	log      *slog.Logger
	chatRepo repository.ChatRepo
}

func NewChat(log *slog.Logger, chatRepo repository.ChatRepo) *Chat {
	return &Chat{
		log:      log,
		chatRepo: chatRepo,
	}
}

func (c *Chat) CreateChannel(ctx context.Context, address string) (uint64, error) {
	newChat := chat.Chat{
		Type:    "channel",
		Address: address,
	}

	chatID, err := c.chatRepo.Create(ctx, newChat)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (c *Chat) CreateGroup(ctx context.Context, address string) (uint64, error) {
	newChat := chat.Chat{
		Type:    "group",
		Address: address,
	}

	chatID, err := c.chatRepo.Create(ctx, newChat)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (c *Chat) CreatePrivate(ctx context.Context, address string) (uint64, error) {
	newChat := chat.Chat{
		Type:    "private",
		Address: address,
	}

	chatID, err := c.chatRepo.Create(ctx, newChat)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
