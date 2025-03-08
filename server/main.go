package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/arcane/server/handlers"
	"github.com/joseph0x45/arcane/server/logger"
	"github.com/joseph0x45/arcane/server/middleware"
	"github.com/joseph0x45/arcane/server/repository"
	_ "github.com/lib/pq"
)

func getDBPool() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		logger.Error(err)
		return nil
	}
	return db
}

func newLogger() *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}),
	)
	return logger
}

func main() {
	db := getDBPool()
	if db == nil {
		return
	}

	//repositories
	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewSessionRepo(db)

	//middlewares
	authMiddleware := middleware.NewAuthMiddleware(userRepo, sessionRepo)

	//handlers
	authHandler := handlers.NewAuthHandler(userRepo, sessionRepo, authMiddleware)
	mux := http.NewServeMux()

	//register routes
	authHandler.RegisterRoutes(mux)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}
	log.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
