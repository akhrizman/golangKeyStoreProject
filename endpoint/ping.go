package endpoint

import "net/http"

func Ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "text/plain charset=utf-8")
	response.Write([]byte("pong"))
}
