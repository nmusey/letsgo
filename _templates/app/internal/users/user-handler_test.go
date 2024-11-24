package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUsers(t *testing.T) {
	handler := UserHandler{
        UserService: &UserService{
            Store: &LocalUserStore{},
        },
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	handler.GetUsers(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

