package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bkojha74/micro-service/auth-handler/helper"
	"github.com/bkojha74/micro-service/auth-handler/models"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken godoc
// @Summary Generate JWT token
// @Description Generate a JWT token for authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "User credentials"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /generate-token [post]
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate credentials (this is just an example, replace with actual validation)
	if creds.Username != "user" || creds.Password != "pass" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(helper.GetEnv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// VerifyToken godoc
// @Summary Verify a JWT token
// @Description Verify the provided JWT token
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Invalid token"
// @Router /verify-token [post]
// @Security BearerAuth
func VerifyToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(helper.GetEnv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid", "username": claims.Username})
}
