package jwt

import (
	"messanger/internal/user"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user user.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["name"] = user.Name
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseTokenWithClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token.Raw, claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	return claims, nil
}

func ValidateToken(w http.ResponseWriter, r *http.Request) *jwt.Token {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	tokenString := strings.Split(authHeader, " ")[1]
	if tokenString == "" {
		return nil
	}

	return VerifyToken(tokenString)
}

func VerifyToken(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil
	}

	if !token.Valid {
		return nil
	}

	return token
}
