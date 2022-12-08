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

var (
	DatastoreEndpoint    = "/datastore/"
	textContentType      = "text/plain; charset=utf-8"
	forbiddenRespText    = "Forbidden"
	keyNotFoundRespText  = "404 key not found"
	okResponseText       = "OK"
	contentTypeHeaderKey = "Content-Type"
)

func Store(ds *Datasource) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		log4g.Request.Println(log4g.NewRequestLogEntry(request))

		user := server.AuthorizedUser(responseWriter, request)
		if user == "" {
			log4g.Info.Printf("Unable to process request: Failed Authorization", request.Method)
			return
		}

		responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

		switch request.Method {
		case http.MethodGet:
			key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)
			data, getErr := ds.Get(Key(key))
			if getErr != nil {
				responseWriter.WriteHeader(http.StatusNotFound)
				responseWriter.Write([]byte(keyNotFoundRespText))
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte(data.GetValue()))
			}

		case http.MethodPut:
			contentHeader := request.Header.Get(contentTypeHeaderKey)
			if contentHeader != "" {
				if contentHeader != textContentType {
					log4g.Error.Printf("%s header is not %s", contentTypeHeaderKey, textContentType)
					responseWriter.WriteHeader(http.StatusUnsupportedMediaType)
					return
				}
			}

			key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)

			bytes, err := io.ReadAll(request.Body)
			defer request.Body.Close()
			newValue := string(bytes)
			//fmt.Println("request body received: ", newValue)
			// TODO Not sure how to handle scenario where user does not provide a value
			if err != nil || len(bytes) == 0 {
				log4g.Warning.Println("request body empty, setting value")
			} else {
				log4g.Info.Printf("request body value found, setting for key %s", key)
			}

			putErr := ds.Put(Key(key), NewData(user, newValue))
			if putErr != nil {
				log4g.Info.Printf("Unauthorized update to %s attempted by user: %s", key, user)
				responseWriter.WriteHeader(http.StatusForbidden)
				responseWriter.Write([]byte(forbiddenRespText))
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte(okResponseText))
			}

		case http.MethodDelete:
			key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)
			delErr := ds.Delete(Key(key), user)

			switch {
			case errors.Is(delErr, ErrKeyNotFound):
				responseWriter.WriteHeader(http.StatusNotFound)
				responseWriter.Write([]byte(keyNotFoundRespText))
			case errors.Is(delErr, ErrValueDeleteForbidden):
				log4g.Info.Printf("Unauthorized deletion of %s attempted by user: %s", key, user)
				responseWriter.WriteHeader(http.StatusForbidden)
				responseWriter.Write([]byte(forbiddenRespText))
			default:
				log4g.Info.Printf("%s deleted successfully", key)
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte(okResponseText))
			}

		default:
			log4g.Warning.Println("Method not found")
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}
