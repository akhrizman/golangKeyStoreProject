package server

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	. "httpstore/logging"
	"net/http"
	"time"
)

var users = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("bird_person")

// Authorize Return user if request is authenticated
func Authorize(request *http.Request) string {
	bearerToken := request.Header.Get("Authorization")
	fmt.Println(bearerToken)
	return ""
}

// Authenticate Validate user-provided credentials
func Authenticate(request *http.Request) string {
	username, password, ok := request.BasicAuth()
	if !ok {
		ErrorLogger.Println("Error parsing basic auth", ok)
		return ""
	}
	expectedPassword, ok := users[username]
	if !ok {
		WarningLogger.Printf("Unknown user: %s", username)
		return ""
	}
	if password != expectedPassword {
		WarningLogger.Printf("Incorrect Password for user: %s\n", username)
		return ""
	}
	return username
}

// GenerateBearerToken Create JWT claims and signed token
func GenerateBearerToken(username string) string {
	// Create JWT claims
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "httpstore service",
		},
	}
	// Create bearer token {header}.{payload}.{signature}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		ErrorLogger.Printf("Error creating the token: %v", err)
		return ""
	}
	fmt.Printf("User: %s - %s", username, tokenString)
	return tokenString
}
