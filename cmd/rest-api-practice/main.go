package main

import (
	"context"
	"log"

	"github.com/ursuldaniel/brunoyam/internal/server"
	"github.com/ursuldaniel/brunoyam/internal/storage"
)

func main() {
	// store, err := storage.NewSqlite3Storage("database.db")
	store, err := storage.NewPostgresStorage(context.TODO(), "postgres://postgres:postgres@localhost:5432/brunoyam")
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(":8080", store)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
