package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/models"
	"server/pkg"
	"server/store"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	Users  *store.Users
	Teams  *store.Teams
	Logger *logrus.Logger
}

func NewAuthHandler(u *store.Users, t *store.Teams, l *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		Users:  u,
		Teams:  t,
		Logger: l,
	}
}

func (h *AuthHandler) RequestGithubAuth(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_AUTH_REDIRECT_URI"),
	)
	data, err := json.Marshal(
		map[string]string{
			"url": redirectURL,
		},
	)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *AuthHandler) GithubAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	accessToken, err := pkg.GetGitHubAccessToken(code, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := pkg.GetGithubData(accessToken, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, ok := data["id"].(float64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	strId := fmt.Sprintf("%.f", id)
	email, ok := data["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	username, ok := data["name"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	avatar, ok := data["avatar_url"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbUser, err := h.Users.GetByGithubId(strId)
	_ = dbUser
	if err != nil {
		if err == sql.ErrNoRows {
			//User is signing in for the first time
			userId := uuid.NewString()
			newUser := &models.User{
				Id:        userId,
				GithubId:  strId,
				Email:     email,
				Username:  username,
				AvatarURL: avatar,
			}
			err = h.Users.Insert(newUser)
			if err != nil {
				h.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//Generate some sort of session id
			http.Redirect(w, r, os.Getenv("ARCANE_WEB_URL"), http.StatusTemporaryRedirect)
			return
		}
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Generate some sort of session id
	http.Redirect(w, r, os.Getenv("ARCANE_WEB_URL"), http.StatusTemporaryRedirect)
	return
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/auth/github", h.RequestGithubAuth)
	r.Get("/auth/github/callback", h.GithubAuthCallback)
}
