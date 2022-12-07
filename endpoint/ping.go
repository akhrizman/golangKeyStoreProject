package endpoint

import (
	. "httpstore/logging"
	"net/http"
)

var PingEndpoint = "/ping"

func Ping(responseWriter http.ResponseWriter, request *http.Request) {
	RequestLogger.Println(NewRequestLogEntry(request))
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)
	responseWriter.Write([]byte("pong"))
}
