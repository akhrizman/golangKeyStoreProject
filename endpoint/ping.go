package endpoint

import (
	"httpstore/log4g"
	"net/http"
)

var PingEndpoint = "/ping"

func Ping(responseWriter http.ResponseWriter, request *http.Request) {
	log4g.Request.Println(log4g.NewRequestLogEntry(request))
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)
	responseWriter.Write([]byte("pong"))
}
