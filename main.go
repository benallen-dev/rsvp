package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/store"
	"rsvp/web"
)

const PORT = ":8080"

func createNewListOfInvitees() {
	// Write to disk

}

func main() {
	db, err := sql.Open("sqlite3", "./rsvp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store.InitDb(db)
	log.Fatal("exiting early")

	mux := web.NewMux()

	fmt.Println("listening on " + PORT)
	http.ListenAndServe(PORT, mux)
}
