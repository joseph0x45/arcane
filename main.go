package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joseph0x45/arcane/handlers"
)

func main() {

	//handlers
	uiHandler := handlers.NewUIHandler()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("we good"))
	})

	uiHandler.RegisterRoutes(r)

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}

	log.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
