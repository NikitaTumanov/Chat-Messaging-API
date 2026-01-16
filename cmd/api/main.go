package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/database/repository"
	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/logger"
	"github.com/NikitaTumanov/Chat-Messaging-API/internal/pkg/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultPort = ":8080"
)

func connectDB(log *slog.Logger) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Error("DATABASE_URL in .env is not exists")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect to DB", "error", err)
	}
	return db
}

func main() {
	log := logger.New()

	db := connectDB(log)

	chatRepo := repository.NewChatRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	var handler server.HTTPHandler
	handler = server.NewHandler(log, chatRepo, msgRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/chats/", handler.Route)

	loggedMux := server.LoggingMiddleware(log, mux)

	log.Info("server started", "port", defaultPort)
	err := http.ListenAndServe(defaultPort, loggedMux)
	if err != nil {
		log.Error("server stopped", "error", err)
	}
}
