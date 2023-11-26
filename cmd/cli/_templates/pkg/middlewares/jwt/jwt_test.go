package jwt

import (
	"testing"
	"time"

	"$appRepo/pkg/core"
    "$appRepo/pkg/models"

	golangJwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

type MockUserService struct{}

func (s MockUserService) GetUserByID(id int) (models.User, error) {
    return models.User{
        ID: id,
    }, nil
}

func TestNewConfig(t *testing.T) {
    ctx := &core.RouterContext{}
    config := NewConfig(ctx, "testsecret")
    assert.NotNil(t, config.Filter)
    assert.NotNil(t, config.Unauthorized)
    assert.NotNil(t, config.Decode)
    assert.Equal(t, "testsecret", config.Secret)
    assert.NotNil(t, config.UserFetcher)
    assert.Equal(t, int(time.Now().Add(time.Hour * 24).Unix()), config.Expiry)
}

func TestMakeDecoder(t *testing.T) {
    decoder := makeDecoder("testsecret")
    assert.NotNil(t, decoder)
}

func TestExtractUserID(t *testing.T) {
    claims := golangJwt.MapClaims{
        "uid": float64(123),
    }

    userID, err := extractUserID(&claims)
    assert.NoError(t, err)
    assert.Equal(t, 123, userID)

    claims = golangJwt.MapClaims{
        "uid": "not-a-number",
    }

    _, err = extractUserID(&claims)
    assert.Error(t, err)
}

func TestNew(t *testing.T) {
    config := Config{
        Filter:       defaultFilter,
        Unauthorized: defaultUnauthorized,
        Decode:       makeDecoder("testsecret"),
        Secret:       "testsecret",
        UserFetcher:  &MockUserService{},
    }

    _ = New(config)
}
