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
	s, closedb, err := store.Init("./test-data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer closedb()

	s.LogJournalMode()

	// data, err := s.ReadAllInvitesWithRSVPs()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, invite := range data {

	// 	fmt.Println(invite)
	// }

	// Start web server
	mux := web.NewMux(s)

	fmt.Println("listening on " + PORT)
	http.ListenAndServe(PORT, mux)
}
