package main

import (
	"io"
	"log/slog"
	"messanger/internal/lib/logger/handlers/slogpretty"
	"os"

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

	log := setupLogger(envLocal, os.Stdout)

	log.Info("TEST", slog.String("something", "Hello World"))

	log.Debug("debug")

	log.Error("something")
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
