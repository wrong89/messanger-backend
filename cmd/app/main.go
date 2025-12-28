package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("env not found", slog.String("err", err.Error()))
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("ERROR connecting to DB:", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Проверяем соединение
	if err := conn.Ping(context.Background()); err != nil {
		fmt.Println("PING ERROR:", err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS PING")
}
