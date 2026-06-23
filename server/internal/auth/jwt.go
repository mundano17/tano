package auth

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func jwtKey() []byte {
	return []byte(os.Getenv("JWT_PVT_KEY"))
}

func createAccessToken(userID int64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(15 * time.Minute).Unix(),
			"iss": "tanoserver",
			"sub": userID,
		})
	s, err := t.SignedString(jwtKey())
	return s, err
}

func createRefreshToken(userID int64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
			"iss": "tanoserver",
			"sub": userID,
		})
	s, err := t.SignedString(jwtKey())
	return s, err
}

func VerifyToken(tokenString string) (int64, error) {
	tok, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey(), nil
	})
	if err != nil {
		return -1, err
	}
	sub, err := tok.Claims.GetSubject()
	if err != nil {
		return -1, err
	}
	userID, err := strconv.ParseInt(sub, 10, 64)
	return userID, err
}

func IsAuthenticated(r *http.Request) (userID int64, err error) {
	authStr := r.Header.Get("Authorization")
	if authStr == "" {
		return -1, ErrUnauthorized
	}
	if !strings.HasPrefix(authStr, "Bearer ") {
		return -1, ErrUnauthorized
	}
	tokenString := strings.TrimPrefix(authStr, "Bearer ")
	userID, err = VerifyToken(tokenString)
	return userID, err
}
