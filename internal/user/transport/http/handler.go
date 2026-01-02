package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"messanger/internal/lib/logger/sl"
	"messanger/internal/user/repository"
	"messanger/internal/user/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	log     *slog.Logger
	auth    usecase.AuthUC
	profile usecase.ProfileUC
}

func NewUserHandler(log *slog.Logger, auth usecase.AuthUC, profile usecase.ProfileUC) *UserHandler {
	return &UserHandler{
		log:     log,
		auth:    auth,
		profile: profile,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "user.http.handler.Register"

	log := h.log.With(
		slog.String("op", op),
	)

	var registerDTO RegisterReqDTO

	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		errDTO := NewErrorDTO(err)
		log.Warn(errDTO.String(), sl.Err(err))
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
		if errors.Is(err, repository.ErrUserAlreadyExist) {
			errDTO := NewErrorDTO(repository.ErrUserAlreadyExist)
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := RegisterRes{
		ID: uid,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "user.http.handler.Login"

	log := h.log.With(
		slog.String("op", op),
	)

	var loginReqDTO LoginReqDTO

	if err := json.NewDecoder(r.Body).Decode(&loginReqDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	if err := loginReqDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		log.Error("validation error", sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	token, err := h.auth.Login(r.Context(), loginReqDTO.Login, loginReqDTO.Password)
	if err != nil {
		errDTO := NewErrorDTO(usecase.ErrInvalidCredentials)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	resp := LoginRes{Token: token}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "user.http.handler.Delete"

	log := h.log.With(
		slog.String("op", op),
	)

	uidStr := chi.URLParam(r, "id")
	if uidStr == "" {
		log.Warn("id not found")
		http.Error(w, "id not found", http.StatusBadRequest)
		return
	}

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		errDTO := NewErrorDTO(err)
		log.Error("incorrect id", sl.Err(err))
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	if err := h.profile.Delete(r.Context(), uint64(uid)); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			errDTO := NewErrorDTO(usecase.ErrInvalidCredentials)
			http.Error(w, errDTO.String(), http.StatusBadRequest)
			return
		}

		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
