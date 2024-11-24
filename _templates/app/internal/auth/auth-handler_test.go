package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"$appRepo/internal/users"
)

type handlerTest struct {
	message         string
	shouldError     bool
	shouldUserError bool
    email           string
	password        string
	expectedStatus  int
}

func TestAuthHandler_PostLogin(t *testing.T) {
	tests := []handlerTest{
		handlerTest{
			message:         "success case",
			shouldError:     false,
			shouldUserError: false,
            email:           "test@focal.chat",
			password:        "test",
			expectedStatus:  http.StatusOK,
		},
		handlerTest{
			message:         "incorrect password",
			shouldError:     false,
			shouldUserError: false,
            email:           "test@focal.chat",
			password:        "not-correct",
			expectedStatus:  http.StatusNotFound,
		},
		handlerTest{
			message:         "auth store error",
			shouldError:     true,
			shouldUserError: false,
            email:           "test@focal.chat",
			password:        "test",
			expectedStatus:  http.StatusNotFound,
		},
	}

	for _, test := range tests {
		handler := setupHandler(test)

		var body bytes.Buffer
		req := LoginRequest{
			Email:    test.email,
			Password: test.password,
		}
		json.NewEncoder(&body).Encode(&req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", &body)
		r.Header.Set("Content-Type", "application/json")

		handler.PostLogin(w, r)
		if test.expectedStatus != w.Code {
			t.Errorf("%s failed: expected %d got %d", test.message, test.expectedStatus, w.Code)
		}
	}
}

func TestAuthHandler_PostRegister(t *testing.T) {
	tests := []handlerTest{
		handlerTest{
			message:         "success case",
			shouldError:     false,
			shouldUserError: false,
            email:           "test@focal.chat",
			password:        "testtest",
			expectedStatus:  http.StatusOK,
		},
		handlerTest{
			message:         "auth store error",
			shouldError:     true,
			shouldUserError: false,
            email:           "test@focal.chat",
			password:        "testtest",
			expectedStatus:  http.StatusInternalServerError,
		},
		handlerTest{
			message:         "password too short",
			shouldError:     false,
			shouldUserError: false,
            email:           "test@focal.chat",
			password:        "test",
			expectedStatus:  http.StatusBadRequest,
		},
		handlerTest{
			message:         "no @ in email",
			shouldError:     false,
			shouldUserError: false,
            email:           "test-at-focal.chat",
			password:        "testtest",
			expectedStatus:  http.StatusBadRequest,
		},
		handlerTest{
			message:         "invalid email name",
			shouldError:     false,
			shouldUserError: false,
            email:           "@focal.chat",
			password:        "testtest",
			expectedStatus:  http.StatusBadRequest,
		},
		handlerTest{
			message:         "invalid email domain",
			shouldError:     false,
			shouldUserError: false,
            email:           "test@",
			password:        "testtest",
			expectedStatus:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		handler := setupHandler(test)

		var body bytes.Buffer
		req := RegisterRequest{
			Email:    test.email,
			Password: test.password,
		}
		json.NewEncoder(&body).Encode(&req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/register", &body)
		r.Header.Set("Content-Type", "application/json")

		handler.PostRegister(w, r)
		if test.expectedStatus != w.Code {
			t.Errorf("%s failed: expected %d got %d", test.message, test.expectedStatus, w.Code)
		}
	}
}
func setupHandler(tc handlerTest) *AuthHandler {
	return &AuthHandler{
		service: &AuthService{
			Store: LocalAuthStore{ShouldError: tc.shouldError},
			UserService: &users.UserService{
				Store: users.LocalUserStore{ShouldError: tc.shouldUserError},
			},
		},
	}
}
