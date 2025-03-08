package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joseph0x45/arcane/server/httputils"
	"github.com/joseph0x45/arcane/server/logger"
	"github.com/joseph0x45/arcane/server/middleware"
	"github.com/joseph0x45/arcane/server/models"
	"github.com/joseph0x45/arcane/server/repository"
	"github.com/joseph0x45/arcane/server/sharedconsts"
	"github.com/joseph0x45/arcane/server/utils"
	"github.com/joseph0x45/arcane/server/validation"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	userRepo       *repository.UserRepo
	sessionRepo    *repository.SessionRepo
	authMiddleware *middleware.AuthMiddleware
}

func NewAuthHandler(
	userRepo *repository.UserRepo,
	sessionRepo *repository.SessionRepo,
	authMiddleware *middleware.AuthMiddleware,
) *AuthHandler {
	return &AuthHandler{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		authMiddleware: authMiddleware,
	}
}

func (h *AuthHandler) handleRegistration(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, "Failed to decode body", http.StatusInternalServerError)
		return
	}
	if payload.Email == "" || !validation.IsEmail(&payload.Email) {
		httputils.WriteError(w, "Invalid Email", http.StatusBadRequest)
		return
	}
	if len(payload.Password) > sharedconsts.BCRYPT_PWD_MAX_LENGH {
		httputils.WriteError(w, "Password too long. Should be less that 72 characters", http.StatusBadRequest)
		return
	}
	if payload.Password == "" {
		httputils.WriteError(w, "Invalid Password", http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetByEmail(payload.Email)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, sharedconsts.GENERIC_HTTP_500_ERROR, http.StatusConflict)
		return
	}
	if user != nil {
		httputils.WriteError(w, "Email is already in use", http.StatusConflict)
		return
	}
	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, sharedconsts.GENERIC_HTTP_500_ERROR, http.StatusInternalServerError)
		return
	}
	user = &models.User{
		ID:       ulid.Make().String(),
		Email:    payload.Email,
		Password: string(hash),
	}
	err = h.userRepo.Insert(user)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, sharedconsts.GENERIC_HTTP_500_ERROR, http.StatusInternalServerError)
		return
	}
	httputils.WriteData(w, nil, http.StatusCreated)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, "Failed to decode body", http.StatusInternalServerError)
		return
	}
	if payload.Email == "" || !validation.IsEmail(&payload.Email) {
		httputils.WriteError(w, "Invalid Email", http.StatusBadRequest)
		return
	}
	if len(payload.Password) > sharedconsts.BCRYPT_PWD_MAX_LENGH {
		httputils.WriteError(w, "Password too long. Should be less that 72 characters", http.StatusBadRequest)
		return
	}
	if payload.Password == "" {
		httputils.WriteError(w, "Invalid Password", http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetByEmail(payload.Email)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, sharedconsts.GENERIC_HTTP_500_ERROR, http.StatusInternalServerError)
		return
	}
	if user == nil {
		httputils.WriteError(w, "No user found with that Email", http.StatusBadRequest)
		return
	}
	if !utils.HashMatchesPassword(user.Password, payload.Password) {
		httputils.WriteError(w, "Wrong password", http.StatusBadRequest)
		return
	}
	session := &models.Session{
		ID:      ulid.Make().String(),
		UserID:  user.ID,
		IsValid: true,
	}
	err = h.sessionRepo.Insert(session)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, sharedconsts.GENERIC_HTTP_500_ERROR, http.StatusInternalServerError)
		return
	}
	httputils.WriteData(w, map[string]any{
		"data": map[string]string{
			"session": session.ID,
		},
	}, http.StatusOK)
}

func (h *AuthHandler) getCurrentUserInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.handleRegistration)
	mux.HandleFunc("POST /auth/login", h.handleLogin)
	mux.Handle("GET /user", h.authMiddleware.Authenticate(h.getCurrentUserInfo()))
}
