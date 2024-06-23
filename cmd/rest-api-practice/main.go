package main

import (
	"log"

	"github.com/ursuldaniel/brunoyam/internal/server"
	"github.com/ursuldaniel/brunoyam/internal/storage"
)

func main() {
	store, err := storage.NewStorage("database.db")
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(":8080", store)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
