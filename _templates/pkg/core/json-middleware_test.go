package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMiddleware(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        middleware := JsonMiddleware{}
        middleware.ServeHTTP(w, r)
        w.WriteHeader(http.StatusOK) 
    })

    handler.ServeHTTP(rec, req)
    res := rec.Result()

    assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}
