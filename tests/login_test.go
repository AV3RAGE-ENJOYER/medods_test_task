package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/handlers"

	"github.com/stretchr/testify/assert"
)

func TestLoginRoute(t *testing.T) {
	PATH := "/api/v1/auth/login"
	w := httptest.NewRecorder()

	exampleUser := handlers.UserRequest{
		Email:    "test1@gmail.com",
		Password: "test",
	}

	req := sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Non existing user", func(t *testing.T) {
		assert.Equal(t, 404, w.Code)
	})

	w = httptest.NewRecorder()
	exampleUser = handlers.UserRequest{
		Email:    "admin@gmail.com",
		Password: "admin1",
	}

	req = sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Incorrect password", func(t *testing.T) {
		assert.Equal(t, 401, w.Code)
	})

	w = httptest.NewRecorder()
	exampleUser = handlers.UserRequest{
		Email:    "admin@gmail.com",
		Password: "admin",
	}

	req = sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Correct password", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}
