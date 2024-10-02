package utils

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowOrigin := "http://emnaservices.online"
		if os.Getenv("IS_DEV_MODE") == "allow" {
			allowOrigin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)                        // Allow this origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Allowed headers

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r) // Call the next handler
	})
}

// ************************************ //
type AuthedUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type key int

const UserContextKey key = iota

type AuthMiddleware struct {
	db *sql.DB
}

func NewAuthMiddleware(db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{
		db: db,
	}
}

func (s *AuthMiddleware) Protect(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Split the Bearer token from the header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization format")
			return
		}

		var user AuthedUser
		query := "SELECT id, username, role FROM users WHERE access_token = $1"
		err := s.db.QueryRow(query, tokenParts[1]).Scan(&user.ID, &user.Username, &user.Role)

		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization ErrNoRows")
			return
		} else if err != nil {
			fmt.Println(err)
			RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization ErrQuery")
			return
		}

		// Pass the user info to the next handler via the request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
