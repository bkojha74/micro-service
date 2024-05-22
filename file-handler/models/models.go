package models

type ContextKey string

const (
	ContextKeyUsername ContextKey = "username"
)

type AuthResponse struct {
	Message string `json:"message"`
	User    string `json:"username"`
	Err     error
}
