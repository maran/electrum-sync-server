package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ianschenck/envflag"
	"log"
	"net/http"
	"os"
)

const (
	Version = "0.0.3"
)

var dbType = envflag.String("DBTYPE", "sqlite3", "Database to use, valid options sqlite3 or postgres")
var dbDir = envflag.String("DBDIR", defaultDbDir(), "Sqlite3 - Path where the application should look for the database file.")

var user = envflag.String("PGUSER", "root", "postgres - Postgres user name")
var password = envflag.String("PGPASSWORD", "''", "postgres - Postgres password")
var host = envflag.String("PGHOST", "127.0.0.1", "postgres - Postgres hostname")
var db = envflag.String("PGDATABASE", "electrum-label-sync", "postgres - Postgres database name")

var listenPort = envflag.String("LISTENPORT", "0.0.0.0:8080", "Port where the json api should listen at in host:port format.")

var useTLS = envflag.Bool("useTls", false, "Serve json api conncetions over TLS.")
var certPath = envflag.String("certPath", "cert.pem", "Path to TLS certificate")
var keyPath = envflag.String("keyPath", "key.pem", "Path to Keyfile")

func main() {
	envflag.Parse()
	var sm SyncMaster

	if *dbType == "sqlite3" {
		var opts DbOpts
		opts.DbType = *dbType
		opts.DbPath = *dbDir
		sm = newSyncMaster(opts)
	} else if *dbType == "postgres" {
		var opts DbOpts
		opts.DbType = *dbType
		opts.User = *user
		opts.Password = *password
		opts.Host = *host
		opts.Dbname = *db
		sm = newSyncMaster(opts)
	} else {
		log.Fatal("Please define which database to use, sqlite3 or postgres")
	}

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

	// Special exception to make Heroku like deployment easy
	if os.Getenv("PORT") != "" {
		*listenPort = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	}

	sm.logger.Info("Server started and listening on %s", *listenPort)
	if *useTLS {
		sm.logger.Info("Using SSL with certificate '%s' and keyfile '%s'", *certPath, *keyPath)
		log.Fatal(http.ListenAndServeTLS(*listenPort, *certPath, *keyPath, api.MakeHandler()))
	} else {
		log.Fatal(http.ListenAndServe(*listenPort, api.MakeHandler()))
	}
}
