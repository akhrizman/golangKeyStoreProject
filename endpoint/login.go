package endpoint

import (
	"fmt"
	"httpstore/log4g"
	"httpstore/server"
	"net/http"
)

var LoginEndpoint = "/login"

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	log4g.Info.Println("Processing login request")
	log4g.Request.Println(log4g.NewRequestLogEntry(request))
	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

	username := server.Authenticate(request)
	if username == "" {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		bearerToken := server.GenerateBearerToken(username)
		if bearerToken == "" {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseWriter.Header().Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
		log4g.Info.Printf("Authenticated User %s", username)
	}
}
