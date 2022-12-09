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

		user := server.AuthorizeUser(responseWriter, request)
		if user == "" {
			log4g.Info.Printf("Unable to process datastore request: Failed Authorization")
			return
		}

		responseWriter.Header().Set(contentTypeHeaderKey, textContentType)

		key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)

		switch request.Method {
		case http.MethodGet:
			data, getErr := ds.Get(Key(key))
			if getErr != nil {
				responseWriter.WriteHeader(http.StatusNotFound)
				_, writeErr := responseWriter.Write([]byte(keyNotFoundRespText))
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
			contentHeader := request.Header.Get(contentTypeHeaderKey)
			if contentHeader != "" {
				if contentHeader != textContentType {
					log4g.Error.Printf("%s header is not %s", contentTypeHeaderKey, textContentType)
					responseWriter.WriteHeader(http.StatusUnsupportedMediaType)
					return
				}
			}

			bytes, err := io.ReadAll(request.Body)
			defer request.Body.Close()
			newValue := string(bytes)
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
				_, writeErr := responseWriter.Write([]byte(forbiddenRespText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				_, writeErr := responseWriter.Write([]byte(okResponseText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			}

		case http.MethodDelete:
			delErr := ds.Delete(Key(key), user)

			switch {
			case errors.Is(delErr, ErrKeyNotFound):
				responseWriter.WriteHeader(http.StatusNotFound)
				_, writeErr := responseWriter.Write([]byte(keyNotFoundRespText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			case errors.Is(delErr, ErrValueDeleteForbidden):
				log4g.Info.Printf("Unauthorized deletion of %s attempted by user: %s", key, user)
				responseWriter.WriteHeader(http.StatusForbidden)
				_, writeErr := responseWriter.Write([]byte(forbiddenRespText))
				if writeErr != nil {
					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			default:
				log4g.Info.Printf("%s deleted successfully", key)
				responseWriter.WriteHeader(http.StatusOK)
				_, writeErr := responseWriter.Write([]byte(okResponseText))
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
