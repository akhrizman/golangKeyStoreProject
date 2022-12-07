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

	var datasource = NewDatasource()

	InfoLogger.Println("STARTING KEYSTORE API")

	http.HandleFunc("/store/", endpoint.Store(&datasource))
	http.HandleFunc("/ping", endpoint.Ping)

	LogAppStart(port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
	if err != nil {
		server.ExitOnErrors(port, err)
	}
}
