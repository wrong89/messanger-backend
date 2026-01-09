package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"messanger/internal/chat/repository"
	"messanger/internal/chat/usecase"
	"messanger/internal/lib/logger/sl"
	"net/http"
)

type ChatHandler struct {
	log    *slog.Logger
	chatUC usecase.ChatUC
}

func New(log *slog.Logger, chatUC usecase.ChatUC) ChatHandler {
	return ChatHandler{
		log:    log,
		chatUC: chatUC,
	}
}

func (h *ChatHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	const op = "chat.http.handler.CreateChannel"

	log := h.log.With(
		slog.String("op", op),
	)

	var createChatDTO CreateChatReqDTO

	if err := json.NewDecoder(r.Body).Decode(&createChatDTO); err != nil {
		errDTO := NewErrorDTO(err)
		log.Warn(errDTO.String(), sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	if err := createChatDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		log.Error("validation error", sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	chatID, err := h.chatUC.CreateChannel(r.Context(), createChatDTO.Address)
	if err != nil {
		// todo: fix. Now it's isn't working
		if errors.Is(err, repository.ErrChatAlreadyExist) {
			errDTO := NewErrorDTO(repository.ErrChatAlreadyExist)
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := CreateChannelResDTO{ID: chatID}

	json.NewEncoder(w).Encode(resp)
}

func (h *ChatHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	const op = "chat.http.handler.CreateGroup"

	log := h.log.With(
		slog.String("op", op),
	)

	var createChatDTO CreateChatReqDTO

	if err := json.NewDecoder(r.Body).Decode(&createChatDTO); err != nil {
		errDTO := NewErrorDTO(err)
		log.Warn(errDTO.String(), sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	if err := createChatDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		log.Error("validation error", sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	chatID, err := h.chatUC.CreateGroup(r.Context(), createChatDTO.Address)
	if err != nil {
		if errors.Is(err, repository.ErrChatAlreadyExist) {
			errDTO := NewErrorDTO(repository.ErrChatAlreadyExist)
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := CreateChannelResDTO{ID: chatID}

	json.NewEncoder(w).Encode(resp)
}

func (h *ChatHandler) CreatePrivate(w http.ResponseWriter, r *http.Request) {
	const op = "chat.http.handler.CreatePrivate"

	log := h.log.With(
		slog.String("op", op),
	)

	var createChatDTO CreateChatReqDTO

	if err := json.NewDecoder(r.Body).Decode(&createChatDTO); err != nil {
		errDTO := NewErrorDTO(err)
		log.Warn(errDTO.String(), sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	if err := createChatDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		log.Error("validation error", sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	chatID, err := h.chatUC.CreatePrivate(r.Context(), createChatDTO.Address)
	if err != nil {
		if errors.Is(err, repository.ErrChatAlreadyExist) {
			errDTO := NewErrorDTO(repository.ErrChatAlreadyExist)
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := CreateChannelResDTO{ID: chatID}

	json.NewEncoder(w).Encode(resp)
}
