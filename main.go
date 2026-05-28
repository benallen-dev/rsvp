package main

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/config"
	"rsvp/store"
	"rsvp/web"
)

const PORT = ":8080"

func main() {
	log.Info("Starting")

	logLevel, err := log.ParseLevel(config.Current.LogLevel)
	log.SetLevel(logLevel)
	log.Infof("Log level: %s", log.GetLevel().String())

	// Start database
	s, closedb, err := store.Init("./data/rsvp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer closedb()

	s.LogJournalMode()

	// Start web server
	mux := web.NewMux(s)

	fmt.Println("listening on " + PORT)
	http.ListenAndServe(PORT, mux)
}
