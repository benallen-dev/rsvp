package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/store"
	"rsvp/web"
)

const PORT = ":8080"

func chanConsumer(c chan string) {
	for msg := range c {
		log.Infof("Received: %s", msg)
	}
	log.Info("Channel closed")
}

func main() {
	log.Info("Starting")

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

	// Experiment chan listening
	myChan := make(chan string)
	go chanConsumer(myChan)

	myChan <- "hello"
	myChan <- "world"
	close(myChan)

	time.Sleep(5000)

	// Start web server
	mux := web.NewMux(s)

	fmt.Println("listening on " + PORT)
	http.ListenAndServe(PORT, mux)
}
