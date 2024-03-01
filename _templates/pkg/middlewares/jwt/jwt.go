package jwt

// TODO: Complete refactor is likely needed

import (
	"errors"
	"net/http"
	"time"

	"$appRepo/pkg/core"
	"$appRepo/pkg/services"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
    Filter       func(w http.ResponseWriter, r *http.Request) bool
    Decode       func(w http.ResponseWriter, r *http.Request) (*jwt.MapClaims, error)
    Unauthorized func(w http.ResponseWriter, r *http.Request) error
    Secret       string
    Expiry       int
    UserService  services.UserService
}

func NewConfig(ctx *core.RouterContext, secret string) Config {
    return Config{
        Filter:       defaultFilter,
        Unauthorized: defaultUnauthorized,
        Decode:       makeDecoder(secret),
        Secret:       secret,
        Expiry:       int(time.Now().Add(time.Hour * 24).Unix()),
        UserService:  services.NewUserService(ctx),
    }
}

func defaultFilter(w http.ResponseWriter, r *http.Request) bool {
    excluded := []string{"/login", "/register", "/logout"}
    for _, path := range excluded {
        if path == r.URL.Path {
            return true
        }
    }

    return false
}

func defaultUnauthorized(w http.ResponseWriter, r *http.Request) error {
    return w.Redirect("/login")
}

func makeDecoder(secret string) func(w http.ResponseWriter, r *http.Request) (*jwt.MapClaims, error) {
    return func(w http.ResponseWriter, r *http.Request) (*jwt.MapClaims, error) {
        // TODO: Extract token from request cookie
        var token string

        claims := jwt.MapClaims{}
        _, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil {
            return nil, err
        }

        return &claims, nil
    }
}

func New(config Config) func (w http.ResponseWriter, r *http.Request) error {
    return func(w http.ResponseWriter, r *http.Request) error {
        if config.Filter(w, r) {
            return http.Cookie
        }

        claims, err := config.Decode(c)
        if err != nil || claims.Valid() != nil {
            return config.Unauthorized(c)
        }

        userID, err := extractUserID(claims)
        if err != nil {
            return config.Unauthorized(c)
        }

        user, err := config.UserService.GetUserByID(userID)
        if user.ID == 0 || err != nil {
            return config.Unauthorized(c)
        }

        c.Locals("user", user)
        return c.Next()
    }
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
