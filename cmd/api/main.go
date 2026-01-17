package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

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

	handler := server.NewHandler(log, chatRepo, msgRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/chats/", handler.Route)

	loggedMux := server.LoggingMiddleware(log, mux)

	srv := &http.Server{
		Addr:              defaultPort,
		Handler:           loggedMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Info("server started", "port", defaultPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Error("server stopped", "error", err)
	}
}
