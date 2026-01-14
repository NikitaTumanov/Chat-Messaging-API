package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/database/repository"
	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultPort = ":8080"
)

func connectDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL in .env is not exists")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := connectDB()

	chatRepo := repository.NewChatRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	var handler server.HTTPHandler
	handler = server.NewHandler(chatRepo, msgRepo)

	http.HandleFunc("/chats/", handler.Route)

	err := http.ListenAndServe(defaultPort, nil)
	if err != nil {
		log.Fatal()
	}
}
