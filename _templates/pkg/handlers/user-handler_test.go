package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"$appRepo/pkg/core"
	"$appRepo/pkg/services"
)

func TestUserHandler_GetUsers(t *testing.T) {
	handler := UserHandler{
        ctx: &core.RouterContext{},
        UserService: services.MockUserService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	err := handler.GetUsers(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

    ct := w.Header().Get("Content-Type")
    assert.Equal(t, "application/json", ct)
}

