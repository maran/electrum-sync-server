package main

import (
	"flag"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

const (
	Version = "0.0.1"
)

var dbDir = flag.String("dbDir", defaultDbDir(), "Path where the application should look for the database file.")
var listenPort = flag.String("listenPort", ":8080", "Port where the json api should listen at in host:port format.")

func main() {
	flag.Parse()

	sm := newSyncMaster(*dbDir)
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/labels/since/:nonce/for/:mpk", sm.GetLabels},
		&rest.Route{"POST", "/label", sm.CreateLabel},
		&rest.Route{"POST", "/labels", sm.CreateLabels},
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	sm.logger.Info("Server started and listening on %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, api.MakeHandler()))
}
