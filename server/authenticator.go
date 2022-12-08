package server

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	. "httpstore/logging"
	"net/http"
	"strings"
	"time"
)

var AuthorizationHeaderKey = "Authorization"

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

// Authorize Return user provided by bearer token if request is authenticated
func Authorize(responseWriter http.ResponseWriter, request *http.Request) string {
	auth := request.Header.Get(AuthorizationHeaderKey)
	if auth == "" {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return ""
	}
	tokenString := strings.Replace(auth, "Bearer ", "", 1)
	fmt.Println(tokenString)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Printf("Error signature invalid %v\n", err)
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return ""
		}
		fmt.Printf("Error processing JWT token %v\n", err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return ""
	}
	if !token.Valid {
		fmt.Printf("Error invalid token %v\n", err)
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return ""
	}

	// username given in the token
	return claims.Username
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
	expirationTime := time.Now().Add(30 * time.Minute)
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
