package endpoint

import (
	. "httpstore/logging"
	"httpstore/server"
	"net/http"
	"os"
	"time"
)

var ShutdownEndpoint = "/shutdown"

func Shutdown(responseWriter http.ResponseWriter, request *http.Request) {
	user := request.Header.Get(userHeaderKey)
	InfoLogger.Printf("Processing %s request by user %s", request.Method, user)
	RequestLogger.Println(NewRequestLogEntry(request))
	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

	if user != "admin" {
		WarningLogger.Printf("Unauthorized shutdown attempted by user: %s", user)
		responseWriter.WriteHeader(http.StatusForbidden)
		responseWriter.Write([]byte(forbiddenRespText))
	} else {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(okResponseText))
		InfoLogger.Printf("Shutdown initiated by user: %s", user)

		//TODO - should something else be done here??? like ensure keystore is locked,
		// all processes are finished, and/or endpoints are no longer accessible?
		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(server.ExitStatus(0))
		}()

	}
}