package endpoint

import (
	"fmt"
	. "httpstore/datasource"
	"io/ioutil"
	"net/http"
	"strings"
)

func Store(datasource *Datasource) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		switch request.Method {
		case http.MethodGet:
			fmt.Println("Processing GET request")
		case http.MethodPut:
			fmt.Println("Processing PUT request")
			contentHeader := request.Header.Get("Content-Type")
			fmt.Println("contentHeader: ", contentHeader)
			if contentHeader != "" {
				if contentHeader != "text/plain; charset=utf-8" {
					msg := "Content-Type header is not text/plain; charset=utf-8"
					http.Error(responseWriter, msg, http.StatusUnsupportedMediaType)
					return
				}
			}
			fmt.Printf("Datasource currently has %d stored value", datasource.Size())

			fmt.Println("request: ", &request)
			fmt.Println("request.RemoteAddr", request.RemoteAddr)
			fmt.Println("request.Method", request.Method)
			fmt.Println("request.URL", request.URL)

			key := strings.TrimPrefix(request.URL.Path, "/store/")
			fmt.Println("key: ", key)

			bytes, err := ioutil.ReadAll(request.Body)
			defer request.Body.Close()
			if err != nil || len(bytes) == 0 {
				fmt.Println("err: ", err)
				fmt.Println("len(bytes): ", len(bytes))
			} else {
				fmt.Printf("\nbytes: %s", bytes)
			}

			//value := http.MaxBytesReader(responseWriter, request.Body, 1048576)
			putErr := datasource.Put(Key(key), NewData(request.Header.Get("authorization"), "testValueForNow"))
			if putErr != nil {
				http.Error(responseWriter, "error", http.StatusForbidden)
			}
			fmt.Printf("Datasource currently has %d stored value", datasource.Size())

		case http.MethodDelete:
			fmt.Println("Processing DELETE request")
		default:
			fmt.Println("Method not found")
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}
