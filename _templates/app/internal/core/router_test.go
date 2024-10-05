package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapRoutes(t *testing.T) {
	router := &Router{
        Mux: http.NewServeMux(),
        Routes: []Route{
            {
                path: "/test",
                handler: func(w http.ResponseWriter, r *http.Request){},
            },
        },
    }

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

	router.MapRoutes()
    router.Mux.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Result().StatusCode)
}

var called bool
type middleware struct {}
func (mw middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    called = true
}

func TestMapRoutes_CallsMiddleware(t *testing.T) {
	router := &Router{
        Mux: http.NewServeMux(),
        Routes: []Route{
            {
                path: "/test",
                handler: func(w http.ResponseWriter, r *http.Request){},
                middleware: []http.Handler{middleware{}},
            },
        },
    }

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

	router.MapRoutes()
    router.Mux.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Result().StatusCode)
    assert.True(t, called)
}
