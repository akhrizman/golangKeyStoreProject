package endpoint

import (
	"httpstore/log4g"
	"httpstore/server"
	"net/http"
	"os"
	"time"
)

var ShutdownEndpoint = "/shutdown"

func Shutdown(responseWriter http.ResponseWriter, request *http.Request) {
	log4g.Request.Println(log4g.NewRequestLogEntry(request))

	user := server.AuthorizedUser(responseWriter, request)
	if user == "" {
		//Responses handled during Authorization
		log4g.Info.Printf("Unable to process request: Failed Authorization", request.Method)
		return
	}

	responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

	if user != "admin" {
		log4g.Warning.Printf("Unauthorized shutdown attempted by user: %s", user)
		responseWriter.WriteHeader(http.StatusForbidden)
		responseWriter.Write([]byte(forbiddenRespText))
	} else {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(okResponseText))
		log4g.Info.Printf("Shutdown initiated by user: %s", user)

		//TODO - should something else be done here??? like ensure keystore is locked,
		// all processes are finished, and/or endpoints are no longer accessible?
		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(server.ExitStatus(0))
		}()

	}
}
