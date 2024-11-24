package auth

import (
	"testing"

	"$appRepo/internal/users"
)

type serviceTest[T interface{}] struct {
	message         string
	shouldError     bool
	shouldUserError bool
	password        string
	expectedResult  T
	expectedError   bool
}

func TestAuthService_StorePassword(t *testing.T) {
	tests := []serviceTest[interface{}]{
		serviceTest[interface{}]{
			message:         "valid",
			shouldError:     false,
			shouldUserError: false,
			password:        "password",
			expectedError:   false,
		},
		serviceTest[interface{}]{
			message:         "store error",
			shouldError:     true,
			shouldUserError: false,
			password:        "password",
			expectedError:   true,
		},
	}

	for _, test := range tests {
		service := &AuthService{
			Store: LocalAuthStore{ShouldError: test.shouldError},
		}

		err := service.StorePassword(test.password, 1)

		if err == nil && test.expectedError {
			t.Errorf("%s: no error returned", test.message)
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s: error returned: %v", test.message, err)
		}
	}
}

func TestAuthService_CheckPassword(t *testing.T) {
	tests := []serviceTest[bool]{
		serviceTest[bool]{
			message:         "wrong password",
			shouldError:     false,
			shouldUserError: false,
			password:        "password",
			expectedResult:  false,
			expectedError:   false,
		},
		serviceTest[bool]{
			message:         "correct password",
			shouldError:     false,
			shouldUserError: false,
			password:        "test",
			expectedResult:  true,
			expectedError:   false,
		},
		serviceTest[bool]{
			message:         "store error",
			shouldError:     true,
			shouldUserError: false,
			password:        "password",
			expectedError:   true,
		},
	}

	for _, test := range tests {
		service := &AuthService{
			Store: LocalAuthStore{ShouldError: test.shouldError},
		}

		correct, err := service.CheckPassword(test.password, 1)

		if err == nil && test.expectedError {
			t.Errorf("%s: no error returned", test.message)
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s: error returned: %v", test.message, err)
		}

        if correct != test.expectedResult {
            t.Errorf("%s: incorrect result: %t", test.message, correct)
        }
	}
}

func TestAuthService_GetAuthCookie(t *testing.T) {
    tests := []serviceTest[string]{
        serviceTest[string]{
            message: "happy path",
            expectedError: false,
        },
    }

    for _, test := range tests {
        service := &AuthService{
            Store: LocalAuthStore{ShouldError: test.shouldError},
        }

        user := users.User{}
        cookie, err := service.GetAuthCookie(&user)

		if err == nil && test.expectedError {
			t.Errorf("%s: no error returned", test.message)
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s: error returned: %v", test.message, err)
		}

        if cookie.Name != "Authorization" {
            t.Errorf("%s: incorrect cookie name: %s", test.message, cookie.Name)
        }

        if len(cookie.Value) < 5 { // 2 bytes are ., 1 byte for header, 1 byte for payload, 1 byte for signature
            t.Errorf("%s: incorrect cookie value: %s", test.message, cookie.Value)
        }
    }
}
