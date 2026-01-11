package http

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrAddressIsEmpty = errors.New("address is empty")
	ErrChatIdIsEmpty  = errors.New("chat_id is empty")
	ErrUserIdIsEmpty  = errors.New("user_id is empty")
	ErrRoleIsEmpty    = errors.New("role is empty")
)

type CreateChatReqDTO struct {
	Address string `json:"address"`
}

func (c CreateChatReqDTO) Validate() error {
	if c.Address == "" {
		return ErrAddressIsEmpty
	}

	return nil
}

type JoinChatReqDTO struct {
	UserID uint64 `json:"user_id"`
	ChatID uint64 `json:"chat_id"`
	Role   string `json:"role"`
}

func (j JoinChatReqDTO) Validate() error {
	if j.Role == "" {
		return ErrRoleIsEmpty
	}
	if j.UserID == 0 {
		return ErrUserIdIsEmpty
	}
	if j.ChatID == 0 {
		return ErrChatIdIsEmpty
	}
	return nil
}

type LeaveChatReqDTO struct {
	UserID uint64 `json:"user_id"`
	ChatID uint64 `json:"chat_id"`
}

func (d *LeaveChatReqDTO) Validate() error {
	if d.UserID == 0 {
		return ErrUserIdIsEmpty
	}
	if d.ChatID == 0 {
		return ErrChatIdIsEmpty
	}
	return nil
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) String() string {
	b, _ := json.MarshalIndent(e, "", "    ")

	return string(b)
}
