package middleware

import (
	"net/http"

	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
)

type AuthMiddleware struct {
	cfg *config.AuthConfig
}

func NewAuthMiddleware(cfg *config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
	}
}
func (auth *AuthMiddleware) WithPerms(next http.Handler, perms []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.JsonError(w, `{"error":"Unauthenticated"}`, http.StatusUnauthorized)
		//next.ServeHTTP(w, r)
	})
}
func (auth *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.JsonError(w, `{"error":"Unauthenticated"}`, http.StatusUnauthorized)
		//next.ServeHTTP(w, r)
	})
}

// Rejects requests with token
func (auth *AuthMiddleware) Reject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.JsonError(w, `{"error":"Authenticated"}`, http.StatusForbidden)
		//next.ServeHTTP(w, r)
	})
}
