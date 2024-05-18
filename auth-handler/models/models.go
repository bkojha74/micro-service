package models

import (
	"github.com/bkojha74/mocro-service/auth-handler/helper"
	"github.com/dgrijalva/jwt-go"
)

// JWT key used to create the signature
var JwtKey = []byte(helper.GetEnv("JWT_SECRET_KEY"))

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
