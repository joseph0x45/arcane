package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"server/models"
	"server/pkg"
	"server/store"
)

type AuthHandler struct {
	Users    *store.Users
	Teams    *store.Teams
	Sessions *store.Sessions
	Logger   *logrus.Logger
}

func NewAuthHandler(u *store.Users, t *store.Teams, s *store.Sessions, l *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		Users:    u,
		Teams:    t,
		Sessions: s,
		Logger:   l,
	}
}

func (h *AuthHandler) RequestGithubAuth(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
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
    fmt.Println("Id not detected")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	strId := fmt.Sprintf("%.f", id)
	email, ok := data["email"].(string)
	if !ok {
		email, err = pkg.GetGithubUserEmail(accessToken, h.Logger)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	username, ok := data["name"].(string)
	if !ok {
    fmt.Println("username not detected")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	avatar, ok := data["avatar_url"].(string)
	if !ok {
    fmt.Println("avatar not detected")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbUser, err := h.Users.GetByGithubId(strId)
	_ = dbUser
	if err != nil {
		if err == sql.ErrNoRows {
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
			sessionId, err := h.Sessions.Create(userId)
			if err != nil {
				h.Logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			redirectURL := fmt.Sprintf("%s/login?sessionId=%s", os.Getenv("ARCANE_WEB_URL"), sessionId)
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.Users.UpdateData(
		dbUser.Id,
		email,
		avatar,
		username,
	)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionId, err := h.Sessions.Create(dbUser.Id)
	if err != nil {
		h.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	redirectURL := fmt.Sprintf("%s/login?sessionId=%s", os.Getenv("ARCANE_WEB_URL"), sessionId)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	return
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/auth/github", h.RequestGithubAuth)
	r.Get("/auth/github/callback", h.GithubAuthCallback)
}
