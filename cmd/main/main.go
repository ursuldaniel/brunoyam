package main

import (
	"context"
	"log"
	"os"

	"github.com/ursuldaniel/brunoyam/internal/server"
	"github.com/ursuldaniel/brunoyam/internal/storage/postgres"
)

func main() {
	// store, err := storage.NewSqlite3Storage("database.db")
	store, err := postgres.NewPostgresStorage(context.TODO(), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(os.Getenv("LISTEN_ADDR"), store)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
