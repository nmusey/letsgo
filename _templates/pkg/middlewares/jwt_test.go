package jwt

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
    name string
    expiry time.Duration
    expectedStatusCode int
}

func TestAuthenticatedMiddleware(t *testing.T) {
    testCases := []testCase{
        {
            name: "Valid token expiry",
            expiry: time.Hour * 2, 
            expectedStatusCode: http.StatusOK,
        },
        {
            name: "Past token expiry",
            expiry: -time.Hour, 
            expectedStatusCode: http.StatusUnauthorized,
        },
    }

    os.Setenv("JWT_SECRET", "secret")

    for _, tc := range testCases {
        recorder := buildHandler(tc)
        assert.Equal(t, tc.expectedStatusCode, recorder.Code)
    }

}

func buildHandler(tc testCase) httptest.ResponseRecorder {
    token := generateTestToken(tc)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: tokenName, Value: string(token)})

	recorder := httptest.NewRecorder()
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := NewAuthenticatedMiddleware(testHandler)
	middleware.ServeHTTP(recorder, req)

    return *recorder
}

func generateTestToken(tc testCase) string {
	claims := jwt.MapClaims{
        "uid": 1234,
        "exp": time.Now().Add(tc.expiry).Unix(),
    }

    token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
    return token
}
