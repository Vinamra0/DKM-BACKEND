package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"dkmbackend/internal/models"
	"dkmbackend/internal/services"
)

func Auth(authSvc *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			var user *models.User
			var err error
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				// allow ADMIN_TOKEN env var as dev fallback
				adminToken := os.Getenv("ADMIN_TOKEN")
				if adminToken == "" {
					http.Error(w, "missing bearer token", http.StatusUnauthorized)
					return
				}
				// if header missing but ADMIN_TOKEN present, require exact match in header
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")
			// first attempt normal auth
			user, err = authSvc.ParseToken(token)
			if err != nil {
				// fallback: check ADMIN_TOKEN env var
				adminToken := os.Getenv("ADMIN_TOKEN")
				if adminToken != "" && token == adminToken {
					user = &models.User{Email: "admin@local", Name: "admin"}
					err = nil
				}
			}
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type ctxKey string

const userKey ctxKey = "user"

func UserFromContext(ctx context.Context) *models.User {
	if v := ctx.Value(userKey); v != nil {
		if u, ok := v.(*models.User); ok {
			return u
		}
	}
	return nil
}
