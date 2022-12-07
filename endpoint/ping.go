package endpoint

import (
	. "httpstore/logging"
	"net/http"
)

func Ping(response http.ResponseWriter, request *http.Request) {
	RequestLogger.Println(NewRequestLogEntry(request))
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "text/plain charset=utf-8")
	response.Write([]byte("pong"))
}
