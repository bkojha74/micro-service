package models

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims struct for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Credentials represents the login credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
