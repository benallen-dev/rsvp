package store

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Returns (*Store, db.Close, error)
func Init(dbFile string) (*Store, func() error, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, nil, err
	}

	s := &Store{db: db}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, nil, err
	}

	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS rsvps (
            id TEXT PRIMARY KEY,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            rsvp_type TEXT,
            name TEXT,
            attending_ceremony BOOLEAN,
            attending_reception BOOLEAN,
            attending_dinner BOOLEAN,
            attending_party BOOLEAN,
            diet_notes TEXT,
            message TEXT,
            supersedes TEXT
        )
    `)

	return s, db.Close, err
}

// Lets you know which mode the DB is in
func (s *Store) LogJournalMode() {
	row := s.db.QueryRow("PRAGMA journal_mode")
	var mode string
	err := row.Scan(&mode)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Journal mode: %s", mode)
}

func (s *Store) DeleteAll() error {
	_, err := s.db.Exec("DELETE from invites")
	if err != nil {
		return err
	}

	_, err = s.db.Exec("DELETE from rsvps")
	return err
}

//  ── Helper functions for parsing ─────────────────────────────────────────

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func parseTimestamp(s string) (time.Time, error) {
	// Try parsing as RFC3339 first (standard ISO 8601 format)
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}

	// Try parsing as SQLite datetime format
	t, err = time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}
