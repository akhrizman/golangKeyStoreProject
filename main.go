package main

import (
	"fmt"
	"httpstore/datasource"
	. "httpstore/endpoint"
	"httpstore/log4g"
	"httpstore/server"
	"net/http"
)

var port int

func main() {
	log4g.SetupLogFiles()
	defer log4g.CloseLogFiles()

	port = server.ValidatePort()

	log4g.Info.Println("STARTING HTTP DATASTORE")
	var datasource = datasource.NewDatasource(10)

	log4g.Info.Println("Setting Up Route Handlers")
	http.HandleFunc(PingEndpoint, Ping)
	http.HandleFunc(LoginEndpoint, Login)
	http.HandleFunc(DatastoreEndpoint, Datastore(&datasource))
	http.HandleFunc(ListEndpoint, List(&datasource))
	http.HandleFunc(ShutdownEndpoint, Shutdown)

	apiHost := server.ApiHost(port)
	log4g.Info.Printf("Server available, see - http://%s", apiHost)
	err := http.ListenAndServe(fmt.Sprint(apiHost), nil)
	if err != nil {
		server.ExitOnErrors(port, err)
	}
}
