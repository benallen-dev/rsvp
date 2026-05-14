package main

import (
	_ "embed"
	"encoding/csv"
	"os"
	"strings"

	"rsvp/invite"
	"rsvp/store"

	"github.com/charmbracelet/log"
)

//go:embed invites.csv
var guestlistCSV string

func main() {
	log.Info("Initialising guestlist")


//  ── setup db ─────────────────────────────────────────────────────────────
	s, closedb, err := store.Init("./rsvp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer closedb()

	s.LogJournalMode()
	err = s.DeleteAll()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}


//  ── Read from csv ────────────────────────────────────────────────────────
	reader := csv.NewReader(strings.NewReader(guestlistCSV))
	records, err := reader.ReadAll()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Loop through them and insert into the DB
	for _, record := range records[1:] {
		name := record[0]
		timeSlot := record[1]

		log.Infof("Read line: %s - %s", name, timeSlot)
		err := s.CreateInvite(invite.NewInvite(name, timeSlot == "day"))
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}


	}

}
