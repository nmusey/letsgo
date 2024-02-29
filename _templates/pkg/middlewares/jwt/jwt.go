package jwt

import (
	"errors"
	"time"

	"$appRepo/pkg/core"
	"$appRepo/pkg/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
    Filter       func(c *fiber.Ctx) bool
    Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
    Unauthorized func(c *fiber.Ctx) error
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

func defaultFilter(c *fiber.Ctx) bool {
    excluded := []string{"/login", "/register", "/logout"}
    for _, path := range excluded {
        if path == c.Path() {
            return true
        }
    }

    return false
}

func defaultUnauthorized(c *fiber.Ctx) error {
    return c.Redirect("/login")
}

func makeDecoder(secret string) func(c *fiber.Ctx) (*jwt.MapClaims, error) {
    return func(c *fiber.Ctx) (*jwt.MapClaims, error) {
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

func New(config Config) func (c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
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
