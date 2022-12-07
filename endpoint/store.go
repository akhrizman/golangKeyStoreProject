package endpoint

import (
	"errors"
	"fmt"
	. "httpstore/datasource"
	. "httpstore/logging"
	"io"
	"net/http"
	"strings"
)

var (
	DatastoreEndpoint    = "/datastore/"
	textContentType      = "text/plain; charset=utf-8"
	jsonContentType      = "application/json"
	forbiddenRespText    = "Forbidden"
	keyNotFoundRespText  = "404 key not found"
	okResponseText       = "OK"
	contentTypeHeaderKey = "Content-Type"
	userHeaderKey        = "Authorization"
)

func Store(datasource *Datasource) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		InfoLogger.Printf("Processing %s request by user %s", request.Method, request.Header.Get(userHeaderKey))
		responseWriter.Header().Set(contentTypeHeaderKey, textContentType)
		switch request.Method {
		case http.MethodGet:
			RequestLogger.Println(NewRequestLogEntry(request))

			key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)
			data, getErr := datasource.Get(Key(key))
			if getErr != nil {
				http.Error(responseWriter, "error", http.StatusNotFound)
				responseWriter.Write([]byte(keyNotFoundRespText))
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte(data.GetValue()))
			}

		case http.MethodPut:
			RequestLogger.Println(NewRequestLogEntry(request))

			contentHeader := request.Header.Get(contentTypeHeaderKey)
			if contentHeader != "" {
				if contentHeader != textContentType {
					msg := fmt.Sprintf("%s header is not %s", contentTypeHeaderKey, textContentType)
					http.Error(responseWriter, msg, http.StatusUnsupportedMediaType)
					return
				}
			}

			key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)

			bytes, err := io.ReadAll(request.Body)
			defer request.Body.Close()
			newValue := string(bytes)
			fmt.Println("request body received: ", newValue)
			// TODO Not sure how to handle scenario where user does not provide a value
			if err != nil || len(bytes) == 0 {
				WarningLogger.Println("request body empty, setting value")
			} else {
				InfoLogger.Printf("request body value found, setting for key %s", key)
			}

			putErr := datasource.Put(Key(key), NewData(request.Header.Get(userHeaderKey), newValue))
			if putErr != nil {
				http.Error(responseWriter, "error", http.StatusForbidden)
				responseWriter.Write([]byte(forbiddenRespText))
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte("OK"))
			}

		case http.MethodDelete:
			RequestLogger.Println(NewRequestLogEntry(request))

			key := strings.TrimPrefix(request.URL.Path, DatastoreEndpoint)
			delErr := datasource.Delete(Key(key), request.Header.Get(userHeaderKey))

			switch {
			case errors.Is(delErr, ErrKeyNotFound):
				responseWriter.WriteHeader(http.StatusNotFound)
				responseWriter.Write([]byte(keyNotFoundRespText))
			case errors.Is(delErr, ErrValueDeleteForbidden):
				responseWriter.WriteHeader(http.StatusForbidden)
				responseWriter.Write([]byte(forbiddenRespText))
			default:
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte("OK"))
			}

		default:
			WarningLogger.Println("Method not found")
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}
