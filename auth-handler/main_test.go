package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bkojha74/mocro-service/auth-handler/models"
	"github.com/bkojha74/mocro-service/auth-handler/router"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// Mock data
var validCredentials = models.Credentials{
	Username: "user",
	Password: "pass",
}

var invalidCredentials = models.Credentials{
	Username: "invalid",
	Password: "invalid",
}

func TestGenerateToken(t *testing.T) {
	router := router.GetRouter()

	t.Run("Generate Token - Success", func(t *testing.T) {
		requestBody, _ := json.Marshal(validCredentials)
		req, _ := http.NewRequest("POST", "/generate-token", bytes.NewBuffer(requestBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["token"])

		// Validate the token
		token, err := jwt.ParseWithClaims(response["token"], &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return models.JwtKey, nil
		})
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.Valid)
	})

	t.Run("Generate Token - Invalid Credentials", func(t *testing.T) {
		requestBody, _ := json.Marshal(invalidCredentials)
		req, _ := http.NewRequest("POST", "/generate-token", bytes.NewBuffer(requestBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Generate Token - Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/generate-token", bytes.NewBuffer([]byte("invalid body")))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestMain(m *testing.M) {
	// Set up environment variable for JWT secret key
	os.Setenv("JWT_SECRET_KEY", "my_secret_key")

	code := m.Run()

	// Clean up
	os.Unsetenv("JWT_SECRET_KEY")

	os.Exit(code)
}
