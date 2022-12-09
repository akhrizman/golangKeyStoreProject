package endpoint

import (
	"fmt"
	"httpstore/log4g"
	"httpstore/server"
	"net/http"
)

const LoginEndpoint = "/login"

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	log4g.Request.Println(log4g.NewRequestLogEntry(request))
	responseWriter.Header().Set(ContentTypeHeaderKey, TextContentType)

	username := server.Authenticate(request)
	if username == "" {
		responseWriter.WriteHeader(http.StatusUnauthorized)
	} else {
		bearerToken := server.GenerateBearerToken(username)
		if bearerToken == "" {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		} else {
			responseWriter.Header().Set(server.AuthorizationHeaderKey, fmt.Sprintf("Bearer %s", bearerToken))
			log4g.Info.Printf("Authenticated User %s", username)
		}
	}
}
