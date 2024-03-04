package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var tokenName = "authorization"
var UserIdCookieName = "user_id"

func AuthticateMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Cookie(tokenName)
        claims, err := decodeToken(token)
        if err != nil {
            unauthorized(w, r)
            return
        }

        userId, err := extractUserID(claims)
        if err != nil {
            unauthorized(w, r)
            return
        }

        cookie := http.Cookie{
            Name:     UserIdCookieName,
            Value:    userId,
        }

        r.AddCookie(&cookie)

        next.ServeHTTP(w, r)
    })
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/login", http.StatusUnauthorized)
}

func decodeToken(token string) (jwt.MapClaims, error) {
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token, &claims, paseJwt)
    if err != nil {
        return nil, err
    }

    return claims, claims.Valid()
}

func parseJwt(token *jwt.Token) (interface{}, error) {
    secret := []byte(os.Getenv("JWT_SECRET"))
    return secret, nil
}

func extractUserID(claims *jwt.MapClaims) (int, error) {
    userIDValue, ok := (*claims)["uid"]
    if !ok {
        return 0, errors.New("user ID not found in claims")
    }

    userID, ok := userIDValue.(float64) // JWT decodes numbers as float64
    if !ok {
        return 0, errors.New("user ID is not a valid number")
    }

    return int(userID), nil
}
