package http

import "messanger/internal/user/usecase"

type Handler struct {
	auth usecase.AuthUC
}
