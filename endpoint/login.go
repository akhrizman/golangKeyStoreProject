package endpoint

import (
	"fmt"
	. "httpstore/logging"
	. "httpstore/server"
	"net/http"
)

var LoginEndpoint = "/login"

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	InfoLogger.Println("Processing login request")
	RequestLogger.Println(NewRequestLogEntry(request))
	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

	username := Authenticate(request)
	if username == "" {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		bearerToken := GenerateBearerToken(username)
		if bearerToken == "" {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseWriter.Header().Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
		InfoLogger.Printf("Authenticated User %s", username)
	}
}
