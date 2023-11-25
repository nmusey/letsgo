package jwt

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"$appRepo/pkg/core"
	"$appRepo/pkg/services"
)

type Config struct {
	Filter          func(c *fiber.Ctx) bool
	Decode          func(c *fiber.Ctx) (*jwt.MapClaims, error)
    Unauthorized    fiber.Handler
	Secret          string
	Expiry          int
    UserService     services.UserService
}

func DefaultConfig(ctx *core.RouterContext) Config {
    secret := os.Getenv("JWT_SECRET")
    return Config{
        Filter: func(c *fiber.Ctx) bool {
            excluded := []string{"/login", "/register", "/logout"}
            for _, path := range excluded {
                if path == c.Path() {
                    return true
                }
            }

            return false
        },
        Unauthorized: func(c *fiber.Ctx) error {
            return c.Redirect("/login")
        },
        Decode: func(c *fiber.Ctx) (*jwt.MapClaims, error) {
            token := c.Cookies("jwt")
            claims := jwt.MapClaims{}
            _, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })

            if err != nil {
                return nil, err
            }

            return &claims, nil
        },
        Secret: secret,
        Expiry: int(time.Now().Add(time.Hour * 24).Unix()),
        UserService: services.NewUserService(ctx),
    }
}

func New(config Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if config.Filter(c) {
            return c.Next()
        }

        claims, err := config.Decode(c)
        if err != nil || claims.Valid() != nil {
            return config.Unauthorized(c)
        }

        var userId interface{} = (*claims)["uid"]
        if userId == nil {
            return config.Unauthorized(c)
        }

        user, err := config.UserService.GetUserByID(int(userId.(float64)))
        if user.ID == 0 || err != nil {
            return config.Unauthorized(c)
        }

        c.Locals("user", user)
        return c.Next()
    } 
}
