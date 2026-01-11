package main

import (
	"context"
	"io"
	"log/slog"
	chatRepo "messanger/internal/chat/repository"
	chatHTTP "messanger/internal/chat/transport/http"
	chatUC "messanger/internal/chat/usecase"
	"messanger/internal/lib/logger/handlers/slogpretty"
	userRepo "messanger/internal/user/repository"
	userHTTP "messanger/internal/user/transport/http"
	userUC "messanger/internal/user/usecase"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/user", func(r chi.Router) {
		storage, err := userRepo.New(ctx, DATABASE_URL)
		if err != nil {
			panic(err)
		}

		auth := userUC.NewAuth(log, storage, JWT_SECRET, time.Hour*24)
		profile := userUC.NewProfile(log, storage)
		handler := userHTTP.NewUserHandler(log, auth, profile)

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		r.With(userHTTP.AuthMiddleware).Delete("/delete/{id}", handler.Delete)
	})

	r.Route("/chat", func(r chi.Router) {
		storage, err := chatRepo.New(ctx, DATABASE_URL)
		if err != nil {
			panic(err)
		}

		r.Use(userHTTP.AuthMiddleware)

		chatUc := chatUC.NewChat(log, storage)
		handler := chatHTTP.New(log, chatUc)

		r.Post("/channel", handler.CreateChannel)
		r.Post("/group", handler.CreateGroup)
		r.Post("/private", handler.CreatePrivate)

		r.Post("/join", handler.Join)
		r.Post("/leave", handler.Leave)
	})

	log.Info("trying to start server...", slog.String("addr", SERVER_ADDR))
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
