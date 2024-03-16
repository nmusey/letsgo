package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var tokenName = "authorization"
var UserIdCookieName = "user_id"

type AuthenticateMiddleware struct {
    handler http.Handler
}

func NewAuthenticatedMiddleware(handler http.Handler) *AuthenticateMiddleware {
    return &AuthenticateMiddleware{handler: handler}
}

func (m AuthenticateMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tokenCookie, err := r.Cookie(tokenName)
    if err != nil {
        unauthorized(w, r)
        return
    }

    userId, err := decodeToken(tokenCookie.Value)
    if err != nil {
        unauthorized(w, r)
        return
    }

    cookie := http.Cookie{
        Name:     UserIdCookieName,
        Value:    userId,
    }

    r.AddCookie(&cookie)
    m.handler.ServeHTTP(w, r)
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/login", http.StatusUnauthorized)
}

func decodeToken(token string) (string, error) {
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token, &claims, parseJwt)
    if err != nil {
        return "", err
    }

    userId, err := extractUserID(claims)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("%d", userId), nil
}

func parseJwt(token *jwt.Token) (interface{}, error) {
    secret := []byte(os.Getenv("JWT_SECRET"))
    return secret, nil
}

func extractUserID(claims jwt.MapClaims) (int, error) {
    userIDValue, ok := (claims)["uid"]
    if !ok {
        return 0, errors.New("user ID not found in claims")
    }

    userID, ok := userIDValue.(float64) // JWT decodes numbers as float64
    if !ok {
        return 0, errors.New("user ID is not a valid number")
    }

    return int(userID), nil
}
