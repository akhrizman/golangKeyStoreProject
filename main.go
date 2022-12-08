package main

import (
	"fmt"
	"httpstore/datasource"
	. "httpstore/endpoint"
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
	var datasource = datasource.NewDatasource()

	InfoLogger.Println("Setting Up Route Handlers")
	http.HandleFunc(PingEndpoint, Ping)
	http.HandleFunc(LoginEndpoint, Login)
	http.HandleFunc(DatastoreEndpoint, Store(&datasource))
	http.HandleFunc(ListEndpoint, List(&datasource))
	http.HandleFunc(ShutdownEndpoint, Shutdown)

	apiHost := server.ApiHost(port)
	InfoLogger.Printf("Server available, see - http://%s", apiHost)
	err := http.ListenAndServe(fmt.Sprintf(apiHost), nil)
	if err != nil {
		server.ExitOnErrors(port, err)
	}
}
