package http

import (
	"messanger/internal/user/usecase"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := validateToken(w, r)
		if err != nil {
			errDTO := NewErrorDTO(usecase.ErrInvalidToken)
			http.Error(w, errDTO.String(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, usecase.ErrInvalidToken
	}

	tokenString := strings.Split(authHeader, " ")[1]
	if tokenString == "" {
		return nil, usecase.ErrInvalidToken
	}

	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	return token, err
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, usecase.ErrInvalidToken
	}

	return token, nil
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // или "http://localhost:5173"
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
