package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joseph0x45/arcane/github"
	"github.com/joseph0x45/arcane/models"
	"github.com/joseph0x45/arcane/repository"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthHandler(
	userRepo *repository.UserRepo,
	sessionRepo *repository.SessionRepo,
) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (h *AuthHandler) handleAuthRequest(w http.ResponseWriter, _ *http.Request) {
	oauthURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		os.Getenv("GH_CLIENT_ID"),
		os.Getenv("GH_REDIRECT_URI"))
	data, err := json.Marshal(
		map[string]string{
			"url": oauthURL,
		},
	)
	if err != nil {
		log.Println("Error while creating oauth url:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *AuthHandler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	error := r.URL.Query().Get("error")
	if error != "" {
		log.Println("Error while loging in:", error)
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	accessToken, err := github.ExchangeCodeWithToken(code)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	data, err := github.GetGithubUserData(accessToken)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	idStr := strconv.Itoa(data.ID)
	existingUser, err := h.userRepo.GetByGithubID(idStr)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	userID := ""
	if existingUser == nil {
		newUser := &models.User{
			ID:        ulid.Make().String(),
			GithubID:  idStr,
			Username:  data.Login,
			AvatarURL: data.AvatarURL,
			JoinedAt:  time.Now().UTC().String(),
		}
		err = h.userRepo.Insert(newUser)
		userID = newUser.ID
	} else {
		userID = existingUser.ID
		err = h.userRepo.UpdateUserData(data.Login, data.AvatarURL, existingUser.ID)
	}
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	session := &models.Session{
		ID:      ulid.Make().String(),
		UserID:  userID,
		IsValid: true,
	}
	err = h.sessionRepo.Insert(session)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/?state=error", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return
}

func (h *AuthHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session, err := h.sessionRepo.GetByID(sessionCookie.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if session == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.userRepo.GetByID(session.UserID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	data := map[string]any{
		"data": map[string]any{
			"user": user,
		},
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error while marshalling data: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", h.handleAuthRequest)
	mux.HandleFunc("/auth/callback", h.handleAuthCallback)
	mux.HandleFunc("/auth/user", h.GetUserData)
}
