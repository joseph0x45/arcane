package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/joseph0x45/arcane/handlers"
	"github.com/joseph0x45/arcane/repository"
	_ "github.com/lib/pq"
)

const localFrontendURL = "http://localhost:5173/"

func serveFrontend(env string) http.Handler {
	if env == "dev" {
		target, _ := url.Parse(localFrontendURL)
		proxy := httputil.NewSingleHostReverseProxy(target)
		log.Println("Proxying frontend to Vite on port 8080")
		return proxy
	}
	//serve from embedded files
	log.Println("Serving frontend from embedded files")
	fs := http.FileServer(http.FS(embeddedFiles))
	return http.StripPrefix("/", fs)
}

func connectToDB() *sqlx.DB {
	dbURL := os.Getenv("DB_URL")
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}
	log.Println("Connected to database")
	return db
}

var env = "prod"

func main() {
	godotenv.Load()
	mux := http.NewServeMux()

	if env == "prod" {
		log.Println("Running in prod mode")
		return
	}
	log.Println("Running in dev mode")

	db := connectToDB()

	//repositories
	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewSessionRepo(db)

	//handlers
	authHandler := handlers.NewAuthHandler(userRepo, sessionRepo)

	mux.Handle("/", serveFrontend(env))

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
