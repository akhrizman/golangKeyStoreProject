package endpoint

import (
	"httpstore/log4g"
	"httpstore/server"
	"net/http"
	"os"
	"time"
)

const ShutdownEndpoint = "/shutdown"

// Shutdown Handler to gracefully down the http server.
func Shutdown(responseWriter http.ResponseWriter, request *http.Request) {
	log4g.Request.Println(log4g.NewRequestLogEntry(request))

	// Only admins can shut down the server via this endpoint
	user := server.AuthorizeUser(responseWriter, request)
	if user == "" {
		//Responses handled during Authorization
		log4g.Info.Println("Unable to process shutdown request: Failed Authorization")
		return
	}

	responseWriter.Header().Set(ContentTypeHeaderKey, TextContentType)

	// Only admins can shut down the server via this endpoint
	if user != "admin" {
		log4g.Warning.Printf("Unauthorized shutdown attempted by user: %s", user)
		responseWriter.WriteHeader(http.StatusForbidden)
		_, writeErr := responseWriter.Write([]byte(ForbiddenRespText))
		if writeErr != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		responseWriter.WriteHeader(http.StatusOK)
		_, writeErr := responseWriter.Write([]byte(OkResponseText))
		if writeErr != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
		log4g.Info.Printf("Shutdown initiated by user: %s", user)

		//TODO - should something else be done here??? like ensure keystore is locked,
		// all processes are finished, and/or endpoints are no longer accessible?
		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(server.ExitStatus(0))
		}()

	}
}
