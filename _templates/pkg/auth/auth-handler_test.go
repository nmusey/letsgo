package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"$appRepo/pkg/core"
	"$appRepo/pkg/users"
)

func TestAuthHandler_GetLogin(t *testing.T) {
	handler := AuthHandler{ctx: &core.RouterContext{}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	err := handler.GetLogin(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetRegister(t *testing.T) {
	handler := AuthHandler{ctx: &core.RouterContext{}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	err := handler.GetRegister(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_PostLogin_Success(t *testing.T) {
	handler := AuthHandler{
        ctx: &core.RouterContext{},
        UserService: users.MockUserService{},
        PasswordService: MockPasswordService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	err := handler.PostLogin(w, r)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestAuthHandler_PostLogin_Error(t *testing.T) {
	handler := AuthHandler{
        ctx: &core.RouterContext{},
        UserService: users.MockUserService{ShouldError: true},
        PasswordService: MockPasswordService{ShouldError: true},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	err := handler.PostLogin(w, r)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_PostRegister_Success(t *testing.T) {
	handler := AuthHandler{
        ctx: &core.RouterContext{},
        UserService: users.MockUserService{},
        PasswordService: MockPasswordService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	err := handler.PostRegister(w, r)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestAuthHandler_PostRegister_UserError(t *testing.T) {
	handler := AuthHandler{
        ctx: &core.RouterContext{},
        UserService: users.MockUserService{ShouldError: true},
        PasswordService: MockPasswordService{},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	err := handler.PostRegister(w, r)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}


func TestAuthHandler_PostRegister_PasswordError(t *testing.T) {
	handler := AuthHandler{
        ctx: &core.RouterContext{},
        UserService: users.MockUserService{},
        PasswordService: MockPasswordService{ShouldError: true},
    }

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)
	err := handler.PostRegister(w, r)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}
