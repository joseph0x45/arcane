package middleware

import (
	"context"
	"net/http"

	"github.com/joseph0x45/arcane/server/httputils"
	"github.com/joseph0x45/arcane/server/repository"
	"github.com/joseph0x45/arcane/server/sharedconsts"
)

type AuthMiddleware struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthMiddleware(
	userRepo *repository.UserRepo,
	sessionRepo *repository.SessionRepo,
) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Authorization")
		if sessionID == "" {
			httputils.WriteError(w, sharedconsts.HTTP_401_ERROR, http.StatusUnauthorized)
			return
		}
		session, err := m.sessionRepo.GetByID(sessionID)
		if err != nil || session == nil {
			httputils.WriteError(w, sharedconsts.HTTP_401_ERROR, http.StatusUnauthorized)
			return
		}
		user, err := m.userRepo.GetByID(session.ID)
		if err != nil || user == nil {
			httputils.WriteError(w, sharedconsts.HTTP_401_ERROR, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
