package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/roberveral/gophercises/urlshort2/urlshort/bolt"
	log "github.com/sirupsen/logrus"
	bbolt "go.etcd.io/bbolt"
)

func main() {
	port := flag.Int("port", 8080, "The port where the server is listening")
	dbPath := flag.String("db", "my.db", "Path to the Bolt database file to use")
	flag.Parse()

	db, err := bbolt.Open(*dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler, err := bolt.New(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Starting UrlShortener server in port: %d. Using Bolt database: %s", *port, *dbPath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
