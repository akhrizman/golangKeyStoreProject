package endpoint

import (
	"fmt"
	. "httpstore/logging"
	. "httpstore/server"
	"net/http"
	"time"
)

var LoginEndpoint = "/login"

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	InfoLogger.Println("Processing login request")
	RequestLogger.Println(NewRequestLogEntry(request))
	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

	username := ValidateLoginCredentials(request)
	if username == "" {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		expirationTime := time.Now().Add(5 * time.Minute)
		tokenString := GenerateBearerToken(username, expirationTime)
		fmt.Println("tokenString: ", tokenString)
		if tokenString == "" {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(responseWriter, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		InfoLogger.Printf("Authenticated User %s", username)
	}
}
