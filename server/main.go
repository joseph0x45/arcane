package main

import (
	"fmt"
	"net/http"
	"server/handlers"
	"server/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	pool := store.GetPostgresPool()

	logger := logrus.New()
	logger.SetReportCaller(true)

	usersStore := store.NewUsersStore(pool)
	teamsStore := store.NewTeamsStore(pool)

	authHandler := handlers.NewAuthHandler(usersStore, teamsStore, logger)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
  authHandler.RegisterRoutes(r)

	fmt.Println("Server launched on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
