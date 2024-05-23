package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bkojha74/micro-service/db-handler/models"
	"github.com/gorilla/mux"
)

type Handler struct {
	UserModel models.UserModel
}

// CreateUserHandler creates a new user.
// @Summary Create a new user
// @Description Create a new user with username and password
// @Tags CRUD Operation Managament
// @Accept json
// @Produce json
// @Param user body models.User true "User info"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Error creating user"
// @Router /users [post]
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.UserModel.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// ReadUserHandler reads a user by username.
// @Summary Read a user by username
// @Description Get user details by username
// @Tags CRUD Operation Managament
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error reading user"
// @Router /users/{username} [get]
func (h *Handler) ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := h.UserModel.ReadUser(username)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler updates a user's password.
// @Summary Update a user's password
// @Description Update a user's password by username
// @Tags CRUD Operation Managament
// @Accept json
// @Produce json
// @Param user body models.User true "User info"
// @Success 200 {string} string "User updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Error updating user"
// @Router /users [put]
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.UserModel.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User updated successfully"))
}

// DeleteUserHandler deletes a user by username.
// @Summary Delete a user by username
// @Description Delete a user by username
// @Tags CRUD Operation Managament
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {string} string "User deleted successfully"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error deleting user"
// @Router /users/{username} [delete]
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/users/")
	err := h.UserModel.DeleteUser(username)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
