package endpoint

import (
	. "httpstore/logging"
	"net/http"
)

func Ping(responseWriter http.ResponseWriter, request *http.Request) {
	RequestLogger.Println(NewRequestLogEntry(request))
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "text/plain charset=utf-8")
	responseWriter.Write([]byte("pong"))
}
