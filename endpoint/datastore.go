package endpoint

import (
	"errors"
	. "httpstore/datasource"
	"httpstore/log4g"
	"httpstore/server"
	"io"
	"net/http"
	"strings"
)

const (
	DatastoreEndpoint = "/datastore/"
)

// Datastore Allows users to create/read/update/delete key-value pairs in the datastore
func Datastore(ds *Datasource) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		log4g.Request.Println(log4g.NewRequestLogEntry(request))

		user := server.AuthorizeUser(responseWriter, request)
		if user == "" {
			//Responses handled during Authorization
			log4g.Info.Printf("Unable to process datastore request: Failed Authorization")
			return
		}

		responseWriter.Header().Set(ContentTypeHeaderKey, TextContentType)

		//TODO this handler assumes that an empty string is a perfectly valid key, if not
		// we would add some conditional logic here to check if it's an empty string, and
		// return a 400 BAD REQUEST ERROR regardless of which Method was used with /datastore/
		key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)

		switch request.Method {
		case http.MethodGet:
			data, getErr := ds.Get(Key(key))
			if getErr != nil {
				responseWriter.WriteHeader(http.StatusNotFound)
				_, writeErr := responseWriter.Write([]byte(KeyNotFoundRespText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				_, writeErr := responseWriter.Write([]byte(data.GetValue()))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			}

		case http.MethodPut:
			contentHeader := request.Header.Get(ContentTypeHeaderKey)
			if contentHeader != "" {
				if contentHeader != TextContentType {
					log4g.Error.Printf("%s header is not %s", ContentTypeHeaderKey, TextContentType)
					responseWriter.WriteHeader(http.StatusUnsupportedMediaType)
					return
				}
			}

			bytes, err := io.ReadAll(request.Body)
			defer request.Body.Close()
			newValue := string(bytes)
			//TODO this handler assumes that an empty string is a perfectly valid value, if not
			// we would add some conditional logic here to check if it's an empty string, and
			// return a 400 BAD REQUEST ERROR
			if err != nil || len(bytes) == 0 {
				log4g.Warning.Println("request body empty, setting value")
			} else {
				log4g.Info.Printf("request body value found, setting for key %s", key)
			}

			putErr := ds.Put(Key(key), user, newValue)
			if putErr != nil {
				log4g.Info.Printf("Unauthorized update to %s attempted by user: %s", key, user)
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
			}

		case http.MethodDelete:
			delErr := ds.Delete(Key(key), user)

			switch {
			case errors.Is(delErr, ErrKeyNotFound):
				responseWriter.WriteHeader(http.StatusNotFound)
				_, writeErr := responseWriter.Write([]byte(KeyNotFoundRespText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			case errors.Is(delErr, ErrValueDeleteForbidden):
				log4g.Info.Printf("Unauthorized deletion of %s attempted by user: %s", key, user)
				responseWriter.WriteHeader(http.StatusForbidden)
				_, writeErr := responseWriter.Write([]byte(ForbiddenRespText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			default:
				log4g.Info.Printf("%s deleted successfully", key)
				responseWriter.WriteHeader(http.StatusOK)
				_, writeErr := responseWriter.Write([]byte(OkResponseText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			}

		default:
			log4g.Warning.Println("Method not found")
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}
