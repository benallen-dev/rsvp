package store

import (
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			day BOOLEAN,
			evening BOOLEAN
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS rsvps (
            id TEXT PRIMARY KEY,
            invite_id TEXT NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            attending_day BOOLEAN,
            attending_evening BOOLEAN,
            diet_notes TEXT,
            message TEXT,
            FOREIGN KEY (invite_id) REFERENCES invites(id)
        )
    `)
    if err != nil {
        log.Fatal(err)
    }
}
