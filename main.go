package main

import (
	"flag"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

const (
	Version = "0.0.2"
)

var dbDir = flag.String("dbDir", defaultDbDir(), "Path where the application should look for the database file.")
var listenPort = flag.String("listenPort", "0.0.0.0:8080", "Port where the json api should listen at in host:port format.")

var useTLS = flag.Bool("useTls", false, "Serve json api conncetions over TLS.")
var certPath = flag.String("certPath", "cert.pem", "Path to TLS certificate")
var keyPath = flag.String("keyPath", "key.pem", "Path to Keyfile")

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
	if *useTLS {
		sm.logger.Info("Using SSL with certificate '%s' and keyfile '%s'", *certPath, *keyPath)
		log.Fatal(http.ListenAndServeTLS(*listenPort, *certPath, *keyPath, api.MakeHandler()))
	} else {
		log.Fatal(http.ListenAndServe(*listenPort, api.MakeHandler()))
	}
}
