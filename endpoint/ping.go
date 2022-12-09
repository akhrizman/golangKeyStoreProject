package endpoint

import (
	"httpstore/log4g"
	"net/http"
)

var PingEndpoint = "/ping"

// Ping Handler to verify server is running
func Ping(responseWriter http.ResponseWriter, request *http.Request) {
	log4g.Request.Println(log4g.NewRequestLogEntry(request))
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set(ContentTypeHeaderKey, TextContentType)
	_, writeErr := responseWriter.Write([]byte("pong"))
	if writeErr != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
