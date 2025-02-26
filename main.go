package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
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

var env = "prod"

func main() {
	mux := http.NewServeMux()

	//handlers
	if env == "prod" {
		log.Println("Running in prod mode")
		return
	}
	log.Println("Running in dev mode")

	mux.Handle("/", serveFrontend(env))

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
