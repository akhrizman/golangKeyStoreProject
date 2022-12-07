package endpoint

import (
	"errors"
	"fmt"
	. "httpstore/datasource"
	. "httpstore/logging"
	"io/ioutil"
	"net/http"
	"strings"
)

func Store(datasource *Datasource) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		fmt.Printf("\nDatasource currently has %d stored value\n", datasource.Size())

		responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		switch request.Method {
		case http.MethodGet:
			InfoLogger.Println("Processing GET request")
			RequestLogger.Println(NewRequestLogEntry(request))

			key := strings.TrimPrefix(request.URL.Path, "/datastore/")
			data, getErr := datasource.Get(Key(key))
			if getErr != nil {
				http.Error(responseWriter, "error", http.StatusNotFound)
				responseWriter.Write([]byte("404 key not found"))
			} else {
				responseWriter.WriteHeader(http.StatusOK)
				responseWriter.Write([]byte(data.GetValue()))
			}

		case http.MethodPut:
			fmt.Println("request: ", request.Body)

			InfoLogger.Println("Processing PUT request")
			RequestLogger.Println(NewRequestLogEntry(request))

			contentHeader := request.Header.Get("Content-Type")
			if contentHeader != "" {
				if contentHeader != "text/plain; charset=utf-8" {
					msg := "Content-Type header is not text/plain; charset=utf-8"
					http.Error(responseWriter, msg, http.StatusUnsupportedMediaType)
					return
				}
			}

			key := strings.TrimPrefix(request.URL.Path, "/store/")

			bytes, err := ioutil.ReadAll(request.Body)
			//bytes, err := io.ReadAll(request.Body)
			defer request.Body.Close()
			var newValue string
			if err != nil || len(bytes) == 0 {
				fmt.Println("err: ", err)
				fmt.Println("len(bytes): ", len(bytes))
				newValue = "testValueForNow"
			} else {
				fmt.Printf("\nbytes: %s", bytes)
				newValue = string(bytes)
			}
			//value := http.MaxBytesReader(responseWriter, request.Body, 1048576)

			//valueBytes := make([]byte, 256)
			//size, err := request.Body.Read(valueBytes)
			//if err != nil || size == 0 {
			//	fmt.Println("err: ", err)
			//	fmt.Println("size: ", size)
			//} else {
			//	valueBytes = append([]byte(nil), valueBytes[:size]...)
			//	fmt.Printf("\nbytes: %s", size)
			//	fmt.Println(valueBytes)
			//}

			putErr := datasource.Put(Key(key), NewData(request.Header.Get("authorization"), newValue))
			if putErr != nil {
				http.Error(responseWriter, "error", http.StatusForbidden)
			} else {
				responseWriter.WriteHeader(http.StatusOK)
			}

		case http.MethodDelete:
			InfoLogger.Println("Processing DELETE request")
			RequestLogger.Println(NewRequestLogEntry(request))

			key := strings.TrimPrefix(request.URL.Path, "/store/")
			delErr := datasource.Delete(Key(key), request.Header.Get("authorization"))

			switch {
			case errors.Is(delErr, ErrKeyNotFound):
				responseWriter.WriteHeader(http.StatusNotFound)
				responseWriter.Write([]byte("404 key not found"))
			case errors.Is(delErr, ErrValueDeleteForbidden):
				responseWriter.WriteHeader(http.StatusForbidden)
			default:
				responseWriter.WriteHeader(http.StatusOK)
			}

		default:
			WarningLogger.Println("Method not found")
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}
