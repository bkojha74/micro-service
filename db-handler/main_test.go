package main_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bkojha74/micro-service/db-handler/controller"
	"github.com/bkojha74/micro-service/db-handler/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockUserModel struct {
	CreateUserFunc func(user models.User) error
	ReadUserFunc   func(username string) (models.User, error)
	UpdateUserFunc func(user models.User) error
	DeleteUserFunc func(username string) error
}

func (m *MockUserModel) CreateUser(user models.User) error {
	return m.CreateUserFunc(user)
}

func (m *MockUserModel) ReadUser(username string) (models.User, error) {
	return m.ReadUserFunc(username)
}

func (m *MockUserModel) UpdateUser(user models.User) error {
	return m.UpdateUserFunc(user)
}

func (m *MockUserModel) DeleteUser(username string) error {
	return m.DeleteUserFunc(username)
}

func TestCreateUserHandler(t *testing.T) {
	mockUserModel := &MockUserModel{
		CreateUserFunc: func(user models.User) error {
			if user.Username == "erroruser" {
				return errors.New("error creating user")
			}
			return nil
		},
	}

	handler := &controller.Handler{UserModel: mockUserModel}

	tests := []struct {
		name         string
		input        string
		expectedCode int
		expectedBody string
	}{
		{"Valid User", `{"username": "testuser", "password": "testpassword"}`, http.StatusCreated, "User created successfully"},
		{"Invalid Request Body", `{invalid json}`, http.StatusBadRequest, "Invalid request body"},
		{"Server Error", `{"username": "erroruser", "password": "testpassword"}`, http.StatusInternalServerError, "Error creating user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(tt.input))
			rr := httptest.NewRecorder()
			handler.CreateUserHandler(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestReadUserHandler(t *testing.T) {
	mockUserModel := &MockUserModel{
		ReadUserFunc: func(username string) (models.User, error) {
			if username == "notfound" {
				return models.User{}, errors.New("user not found")
			}
			if username == "erroruser" {
				return models.User{}, errors.New("internal server error")
			}
			return models.User{Username: username, Password: "hashedpassword"}, nil
		},
	}

	handler := &controller.Handler{UserModel: mockUserModel}

	tests := []struct {
		name         string
		username     string
		expectedCode int
		expectedBody string
	}{
		{"Valid User", "testuser", http.StatusOK, `{"Username":"testuser","Password":"hashedpassword"}`},
		{"User Not Found", "notfound", http.StatusNotFound, "User not found"},
		{"Server Error", "erroruser", http.StatusInternalServerError, "Internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/users/"+tt.username, nil)
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/users/{username}", handler.ReadUserHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
			}
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	mockUserModel := &MockUserModel{
		UpdateUserFunc: func(user models.User) error {
			if user.Username == "erroruser" {
				return errors.New("internal server error")
			}
			return nil
		},
	}

	handler := &controller.Handler{UserModel: mockUserModel}

	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid User",
			requestBody:  `{"Username":"testuser","Password":"newpassword"}`,
			expectedCode: http.StatusOK,
			expectedBody: "User updated successfully",
		},
		{
			name:         "Invalid Request Body",
			requestBody:  `{"Username":"testuser","Password":}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body",
		},
		{
			name:         "Server Error",
			requestBody:  `{"Username":"erroruser","Password":"newpassword"}`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error updating user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/users", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler.UpdateUserHandler)
			handlerFunc.ServeHTTP(rr, req)

			actualBody := strings.TrimSpace(rr.Body.String())

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	mockUserModel := &MockUserModel{
		DeleteUserFunc: func(username string) error {
			if username == "erroruser" {
				return errors.New("error deleting user")
			}
			if username == "notfound" {
				return errors.New("user not found")
			}
			return nil
		},
	}

	handler := &controller.Handler{UserModel: mockUserModel}

	tests := []struct {
		name         string
		username     string
		expectedCode int
		expectedBody string
	}{
		{"Valid User", "testuser", http.StatusOK, "User deleted successfully"},
		{"User Not Found", "notfound", http.StatusNotFound, "User not found"},
		{"Server Error", "erroruser", http.StatusInternalServerError, "Error deleting user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/users/"+tt.username, nil)
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/users/{username}", handler.DeleteUserHandler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}
