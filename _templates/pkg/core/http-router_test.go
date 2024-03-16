package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapRoutes(t *testing.T) {
	router := &HttpRouter{Mux: http.NewServeMux()}

	routes := map[string]HttpHandler{
		"/test": func(w http.ResponseWriter, r *http.Request) error {
			return nil
		},
	}

	router.MapRoutes(routes)

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    router.Mux.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Result().StatusCode)
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"message": "hello"}
	err := WriteJSON(w, payload)
	if err != nil {
		t.Errorf("WriteJSON failed: %v", err)
	}

    result := w.Result()
    assert.Equal(t, 200, result.StatusCode)

    header := result.Header.Get("Content-Type")
    assert.Equal(t, "application/json", header)

    body := make([]byte, w.Body.Len())
    if _, err = w.Body.Read(body); err != nil {
        t.Errorf("Error reading response body: %v", err)
    }

    var responseJSON map[string]interface{}
    if err = json.Unmarshal(body, &responseJSON); err != nil {
        t.Errorf("Error unmarshalling JSON: %v", err)
    }

    assert.Equal(t, "hello", responseJSON["message"])
}

func TestReadJSON(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"message": "hello"}`))
    payload := make(map[string]string)
	err := ReadJSON(r, &payload)
	if err != nil {
		t.Errorf("ReadJSON failed: %v", err)
	}

    assert.Equal(t, "hello", payload["message"])
}
