package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	fmt.Println("Received client request to Generate a Token")
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("client data decoded successfully")
	// validate the user and get corresponding secret key
	resp := getUserInfo(creds.Username)
	if resp.Err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	fmt.Println("Got the Secret")

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(resp.SecretKey))
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

	user := r.URL.Query().Get("user")

	// validate the user and get corresponding secret key
	resp := getUserInfo(user)
	if resp.Err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	fmt.Println("Got the Secret")

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(resp.SecretKey), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid", "username": claims.Username})
}

/*func getUserInfo(user string) models.UserwithError {
	Resp := models.UserwithError{}

	url := fmt.Sprintf("http://db-handler:8082/users/%s", user)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		Resp.Err = err
		return Resp
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		Resp.Err = err
		return Resp
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		Resp.Err = err
		return Resp
	}

	fmt.Println("Response:", string(body))

	err = json.Unmarshal(body, &Resp)
	if err != nil {
		Resp.Err = err
		return Resp
	}

	temp, _ := helper.DecodeString(Resp.SecretKey)
	Resp.SecretKey = string(temp)

	return Resp
}*/
