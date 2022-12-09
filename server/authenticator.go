package server

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"httpstore/log4g"
	"net/http"
	"strings"
	"time"
)

var AuthorizationHeaderKey = "Authorization"

var users = map[string]string{
	"user_a": "$2a$08$GLVOud3QynSlvqon6qZOTeGLI37RimXVrBChQ1cn3qZWyQBI414ty",
	"user_b": "$2a$08$q1YsCnUlI2gI6hBsnhscBOZmKaUSm9mujQB2mQkcFKBlExhkm0ZOu",
	"user_c": "$2a$08$OCBcwJnsWspac0cqSMyEB.yAho0kqDX54GDvTt9LlJL/fdi6wiwUG",
	"admin":  "$2a$08$ofxzK.j/k8kC5eJ/cjOTAe9z5E1d5Hxj5bch6A60sj5XEAiVa2YvW",
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("bird_person") // USE REAL KEY IN PRODUCTION!!!

// AuthorizeUser Return user provided by bearer token if request is authenticated
func AuthorizeUser(responseWriter http.ResponseWriter, request *http.Request) string {
	auth := request.Header.Get(AuthorizationHeaderKey)
	if auth == "" {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return ""
	}
	tokenString := strings.Replace(auth, "Bearer ", "", 1)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log4g.Info.Printf("Error signature invalid %v\n", err)
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return ""
		}
		log4g.Error.Printf("Error processing JWT token %v\n", err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return ""
	}
	if !token.Valid {
		log4g.Error.Printf("Error invalid token %v\n", err)
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
		log4g.Error.Println("Error parsing basic auth", ok)
		return ""
	}
	hashedPassword, ok := users[username]

	if !ok {
		log4g.Warning.Printf("Unknown user: %s", username)
		return ""
	}
	if !PasswordValidated(password, hashedPassword) {
		log4g.Warning.Printf("Incorrect Password for user: %s\n", username)
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
		log4g.Error.Printf("Error creating the token: %v", err)
		return ""
	}
	return tokenString
}

// PasswordValidated Confirm a password matches its hash
func PasswordValidated(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
