package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"messanger/internal/lib/logger/sl"
	"messanger/internal/user/repository"
	"messanger/internal/user/usecase"
	"net/http"
)

type UserHandler struct {
	log  *slog.Logger
	auth usecase.AuthUC
}

func NewUserHandler(log *slog.Logger, auth usecase.AuthUC) *UserHandler {
	return &UserHandler{
		log:  log,
		auth: auth,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "user.http.handler.RegisterHandler"

	log := h.log.With(
		slog.String("op", op),
	)

	var registerDTO RegisterReqDTO

	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, repository.ErrUserAlreadyExist) {
			log.Warn(errDTO.String(), sl.Err(err))
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	if err := registerDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		log.Error("validation error", sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	uid, err := h.auth.Register(r.Context(), registerDTO.Name, registerDTO.Login, registerDTO.Password)
	if err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, repository.ErrUserAlreadyExist) {
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := RegisterRes{
		ID: uid,
	}

	json.NewEncoder(w).Encode(resp)
}
