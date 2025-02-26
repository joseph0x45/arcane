package handlers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UIHandler struct{}

func NewUIHandler() *UIHandler {
	return &UIHandler{}
}

func (uh *UIHandler) renderHomePage(w http.ResponseWriter, _ *http.Request) {
	templ, err := template.ParseFiles("./views/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Execute(w, nil)
}

func (uh *UIHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", uh.renderHomePage)
}
