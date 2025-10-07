package middleware

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
)

type AuthMiddleware struct {
	cfg *config.AuthConfig
	db  *sql.DB
}

func NewAuthMiddleware(cfg *config.AuthConfig, db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
		db:  db,
	}
}

func (auth *AuthMiddleware) parseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.cfg.Secret), nil
	})
}
func (auth *AuthMiddleware) getIDFromToken(token *jwt.Token) (int, error) {
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims type")
	}
	id, ok := claims["sub"].(float64) //All numbers are stored as float64
	if !ok {
		return 0, errors.New("sub claim isn't an int")
	}

	// TODO: optionally check if user exists in DB
	return int(id), nil
}

func (auth *AuthMiddleware) validToken(tokenStr string) error {
	token, err := auth.parseToken(tokenStr)

	if err != nil {
		return err
	}

	if _, err := auth.getIDFromToken(token); err != nil {
		return err
	}
	return nil
}
func (auth *AuthMiddleware) WithPerms(next http.Handler, perms []string) http.Handler {
	return auth.Require(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _ := auth.parseToken(r.Header.Get("Authorization")) //This token was checked by Require
			id, _ := auth.getIDFromToken(token)

			query := `
				SELECT COUNT(DISTINCT p.name)
				FROM user_permission up
				JOIN permission p ON up.permission_id = p.id
				WHERE up.user_id = $1 AND p.name = ANY($2)
			`

			var count int
			err := auth.db.QueryRow(query, id, pq.Array(perms)).Scan(&count)
			if err != nil {
				slog.Error("WithPerms query error", "err", err)
				util.JsonError(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
				return
			}

			// Check if user has all required permissions
			if count < len(perms) {
				util.JsonError(w, `{"error":"Forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}),
	)
}
func (auth *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			util.JsonError(w, `{"error":"Missing authentication header"}`, http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(reqToken, "Bearer ") {
			util.JsonError(w, `{"error":"Invalid Token"}`, http.StatusUnauthorized)
			return
		}
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		err := auth.validToken(reqToken)
		if err != nil {
			util.JsonError(w, `{"error":"Invalid Token"}`, http.StatusUnauthorized)
			return
		}
		r.Header.Set("Authorization", reqToken) //Easy to get after this middleware
		next.ServeHTTP(w, r)
	})
}

// Rejects requests with token
func (auth *AuthMiddleware) Reject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString != "" {
			util.JsonError(w, `{"error":"Authenticated users cannot access this endpoint"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
