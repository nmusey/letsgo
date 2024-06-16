package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapRoutes(t *testing.T) {
	router := &HttpRouter{Mux: http.NewServeMux()}

	routes := []Route{
        {
            path: "/test",
            handler: func(w http.ResponseWriter, r *http.Request){},
		},
    }

	router.MapRoutes(routes)

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    router.Mux.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Result().StatusCode)
}

var called bool
type middleware struct {}
func (mw middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    called = true
}

func TestMapRoutes_CallsMiddleware(t *testing.T) {
	router := &HttpRouter{Mux: http.NewServeMux()}
	routes := []Route{
        {
            path: "/test",
            handler: func(w http.ResponseWriter, r *http.Request){},
            middleware: []http.Handler{middleware{}},
		},
    }

	router.MapRoutes(routes)

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    router.Mux.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Result().StatusCode)
    assert.True(t, called)
}
