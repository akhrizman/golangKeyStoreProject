package endpoint

import (
	"encoding/json"
	. "httpstore/datasource"
	"httpstore/log4g"
	"httpstore/server"
	"net/http"
	"strings"
)

var (
	ListEndpoint    = "/list/"
	jsonContentType = "application/json; charset=utf-8"
)

// List Handler to return detailed information about key-value pairs in the datastore
func List(datasource *Datasource) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		log4g.Request.Println(log4g.NewRequestLogEntry(request))

		user := server.AuthorizeUser(responseWriter, request)
		if user == "" {
			//Responses handled during Authorization
			log4g.Info.Println("Unable to process list request: Failed Authorization")
			return
		}

		responseWriter.Header().Set(contentTypeHeaderKey, jsonContentType)

		key := strings.TrimPrefix(request.URL.Path, ListEndpoint)

		switch request.Method {
		case http.MethodGet:
			if key == "" {
				entriesJson, errGetAll := json.Marshal(datasource.GetAllEntries())
				if errGetAll != nil {
					log4g.Error.Println("Unable to convert Detail Entries to JSON", errGetAll)
					responseWriter.WriteHeader(http.StatusInternalServerError)
				} else {
					responseWriter.WriteHeader(http.StatusOK)
					_, err := responseWriter.Write(entriesJson)
					if err != nil {
						responseWriter.WriteHeader(http.StatusInternalServerError)
					}
				}
			} else {
				entry, errGetOne := datasource.GetEntry(Key(key))
				if errGetOne != nil {
					responseWriter.WriteHeader(http.StatusNotFound)
					_, writeErr := responseWriter.Write([]byte(keyNotFoundRespText))
					if writeErr != nil {
						responseWriter.WriteHeader(http.StatusInternalServerError)
					}
				} else {
					entryJson, errJson := json.Marshal(entry)
					if errJson != nil {
						log4g.Error.Println("Unable to convert Detail Entry to JSON", errJson)
						responseWriter.WriteHeader(http.StatusInternalServerError)
					} else {
						responseWriter.WriteHeader(http.StatusOK)
						_, err := responseWriter.Write(entryJson)
						if err != nil {
							responseWriter.WriteHeader(http.StatusInternalServerError)
						}
					}
				}
			}
		default:
			responseWriter.WriteHeader(http.StatusNotFound)
			_, err := responseWriter.Write([]byte(`{"message": "not found"}`))
			if err != nil {
				responseWriter.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}
