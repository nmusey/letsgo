package jwt

import (
    "errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"$appRepo/pkg/core"
    "$appRepo/pkg/models"
    "$appRepo/pkg/services"
)

type IMiddleware interface {
    Next() error
    Redirect(path string) error
    Path() string
    Cookies(name string) string
    Locals(key string, value interface{})
}

type IRedirecter interface {
    Redirect(path string) error
}

type IPathGetter interface {
    Path() string
}

type ICookieHolder interface {
    Cookies(name string) string
}

type IUserFetcher interface {
    GetUserByID(id int) (models.User, error)
}

type Config struct {
    Filter       func(c IPathGetter) bool
    Decode       func(c ICookieHolder) (*jwt.MapClaims, error)
    Unauthorized func(c IRedirecter) error
    Secret       string
    Expiry       int
    UserFetcher  IUserFetcher
}

func NewConfig(ctx *core.RouterContext, secret string) Config {
    return Config{
        Filter:       defaultFilter,
        Unauthorized: defaultUnauthorized,
        Decode:       makeDecoder(secret),
        Secret:       secret,
        Expiry:       int(time.Now().Add(time.Hour * 24).Unix()),
        UserFetcher:  services.NewUserService(ctx),
    }
}

func defaultFilter(c IPathGetter) bool {
    excluded := []string{"/login", "/register", "/logout"}
    for _, path := range excluded {
        if path == c.Path() {
            return true
        }
    }

    return false
}

func defaultUnauthorized(c IRedirecter) error {
    return c.Redirect("/login")
}

func makeDecoder(secret string) func(c ICookieHolder) (*jwt.MapClaims, error) {
    return func(c ICookieHolder) (*jwt.MapClaims, error) {
        token := c.Cookies("jwt")
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

func New(config Config) func (c IMiddleware) error {
    return func(c IMiddleware) error {
        if config.Filter(c) {
            return c.Next()
        }

        claims, err := config.Decode(c)
        if err != nil || claims.Valid() != nil {
            return config.Unauthorized(c)
        }

        userID, err := extractUserID(claims)
        if err != nil {
            return config.Unauthorized(c)
        }

        user, err := config.UserFetcher.GetUserByID(userID)
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
