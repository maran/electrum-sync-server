#### Electrum sync server
This is the server component for the [Bitcoin wallet Electrum](http://electrum.org)'s label synchronization feature.

There is a public node available at http://sync.bytesized-hosting.com/.

##### Installation

###### Docker
A Docker file is included in the repository. You can pull it using `docker pull maran/electrum-sync-server` and then start it using `docker run -p 0.0.0.0:8080:8080 --rm -v $PWD:/data maran/electrum-sync-server`. You can change $PWD to a folder where you want to save the database file on the host.

###### Source
This project is 'go get(able)' [install Go](http://golang.org/doc/install) and do `go get -u github.com/maran/electrum-sync-server`.

##### Configuration

Configuration happens via environment variables.

`DBTYPE`: Database to use, possible values sqlite3 or postgres.
`DBDIR`: Directory to use for sqlite3 database.
`DBUSER`: Postgres user
`DBPASSWORD`: Postgres password
`PGHOST`: Postgres hostname
`PGDATABASE`: Postgres database
`LISTENPORT`: Address to bind on. Format; ip:port
