package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"$appRepo/internal/core"
	"$appRepo/internal/users"
)

func TestAuthHandler_GetLogin(t *testing.T) {
	handler := AuthHandler{router: core.Router{}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	handler.GetLogin(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetRegister(t *testing.T) {
	handler := AuthHandler{router: core.Router{}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	handler.GetRegister(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_PostLogin_Success(t *testing.T) {
	handler := AuthHandler{
        router: core.Router{},
        UserService: users.MockUserService{},
        PasswordService: MockPasswordService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	handler.PostLogin(w, r)

    assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestAuthHandler_PostLogin_Error(t *testing.T) {
	handler := AuthHandler{
        router: core.Router{},
        UserService: users.MockUserService{ShouldError: true},
        PasswordService: MockPasswordService{ShouldError: true},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	handler.PostLogin(w, r)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_PostRegister_Success(t *testing.T) {
	handler := AuthHandler{
        router: core.Router{},
        UserService: users.MockUserService{},
        PasswordService: MockPasswordService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	handler.PostRegister(w, r)

    assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestAuthHandler_PostRegister_UserError(t *testing.T) {
	handler := AuthHandler{
        router: core.Router{},
        UserService: users.MockUserService{ShouldError: true},
        PasswordService: MockPasswordService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	handler.PostRegister(w, r)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}


func TestAuthHandler_PostRegister_PasswordError(t *testing.T) {
	handler := AuthHandler{
        router: core.Router{},
        UserService: users.MockUserService{},
        PasswordService: MockPasswordService{ShouldError: true},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	handler.PostRegister(w, r)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}
