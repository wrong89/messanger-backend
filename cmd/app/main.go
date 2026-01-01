package main

import (
	"context"
	"io"
	"log/slog"
	"messanger/internal/lib/logger/handlers/slogpretty"
	"messanger/internal/user/repository"
	userHTTP "messanger/internal/user/transport/http"
	"messanger/internal/user/usecase"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var (
		DATABASE_URL = os.Getenv("DATABASE_URL")
		JWT_SECRET   = os.Getenv("JWT_SECRET")
		SERVER_ADDR  = os.Getenv("SERVER_ADDR")
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := setupLogger(envLocal, os.Stdout)

	storage, err := repository.New(ctx, DATABASE_URL)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/user", func(r chi.Router) {
		auth := usecase.NewAuth(log, storage, JWT_SECRET, time.Hour*24)
		handler := userHTTP.NewUserHandler(log, auth)

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})

	if err := http.ListenAndServe(SERVER_ADDR, r); err != nil {
		panic(err)
	}
}

func setupLogger(env string, out io.Writer) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog(out)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog(out io.Writer) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(out)

	return slog.New(handler)
}
