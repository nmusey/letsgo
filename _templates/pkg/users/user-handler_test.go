package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"$appRepo/pkg/core"
)

func TestUserHandler_GetUsers(t *testing.T) {
	handler := UserHandler{
        ctx: &core.RouterContext{},
        UserService: MockUserService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	handler.GetUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

    ct := w.Header().Get("Content-Type")
    assert.Equal(t, "application/json", ct)
}

