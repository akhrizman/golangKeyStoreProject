package main

import (
	"fmt"
	. "httpstore/datasource"
	"httpstore/endpoint"
	. "httpstore/logging"
	"httpstore/server"
	"net/http"
)

var port int

func main() {
	SetupLogFiles()
	defer CloseLogFiles()

	port = server.ValidatePort()

	InfoLogger.Println("STARTING KEYSTORE API")
	var datasource = NewDatasource()

	http.HandleFunc("/ping", endpoint.Ping)
	http.HandleFunc("/store/", endpoint.Store(&datasource))

	LogAppStart(port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
	if err != nil {
		server.ExitOnErrors(port, err)
	}
}
