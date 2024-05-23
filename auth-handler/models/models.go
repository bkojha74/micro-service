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

type User struct {
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	SecretKey string `json:"secret" bson:"secret"`
	Role      string `json:"role" bson:"role"`
}

type UserwithError struct {
	User
	Err error
}
