package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var JwtKey = []byte("hello Bipin")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ContextKey string

const (
	ContextKeyUsername ContextKey = "username"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		fmt.Println("[", tokenString, "]")

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		fmt.Println(claims)
		fmt.Println(token)
		fmt.Println(err)
		fmt.Println("[", JwtKey, "]")
		fmt.Println("[", string(JwtKey), "]")

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUsername, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(ContextKeyUsername).(string)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + username})
}

func main() {
	router := mux.NewRouter()
	router.Handle("/protected", VerifyToken(http.HandlerFunc(ProtectedEndpoint))).Methods("GET")
	http.ListenAndServe(":7081", router)
}
