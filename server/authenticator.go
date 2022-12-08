package server

import (
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

// Authenticate Return user if request is authenticated
func Authenticate(request *http.Request) string {
	return ""
}

func ValidateLoginCredentials(request *http.Request) string {
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

func GenerateBearerToken(username string, expirationTime time.Time) string {
	// Create JWT claims
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
	return tokenString
}
