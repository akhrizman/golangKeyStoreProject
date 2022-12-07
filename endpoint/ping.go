package endpoint

import (
	"fmt"
	. "httpstore/logging"
	"net/http"
)

func Ping(responseWriter http.ResponseWriter, request *http.Request) {
	value := http.MaxBytesReader(responseWriter, request.Body, 1048576)
	fmt.Println(value)
	
	RequestLogger.Println(NewRequestLogEntry(request))
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "text/plain charset=utf-8")
	responseWriter.Write([]byte("pong"))
}
