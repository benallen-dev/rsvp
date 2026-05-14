package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/invite"
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

	// TODO: MAKE THIS MATCH INVITE
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			day BOOLEAN,
			evening BOOLEAN
		)
	`)
	if err != nil {
		return nil, nil, err
	}

	// TODO: MAKE THIS MATCH RSVP
	_, err = s.db.Exec(`
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
		return nil, nil, err
	}

	return s, db.Close, nil
}


//  ── Helper functions used in public API ───────────────────────────────────

// getRSVPsForInvite is a helper function that fetches all RSVPs for a given invite ID
func (s *Store) getRSVPsForInvite(inviteID string) ([]*invite.RSVP, error) {
	rows, err := s.db.Query(`
		SELECT id, invite_id, timestamp, attending_day, attending_evening, diet_notes, message 
		FROM rsvps 
		WHERE invite_id = ?
		ORDER BY timestamp
	`, inviteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rsvps []*invite.RSVP

	for rows.Next() {
		var id, inviteIDStr, timestamp, dietNotes, message string
		var attendingDay, attendingEvening bool

		if err := rows.Scan(&id, &inviteIDStr, &timestamp, &attendingDay, &attendingEvening, &dietNotes, &message); err != nil {
			return nil, err
		}

		// Parse UUIDs
		rsvpID, err := parseUUID(id)
		if err != nil {
			return nil, err
		}

		inviteUUID, err := parseUUID(inviteIDStr)
		if err != nil {
			return nil, err
		}

		// Parse timestamp
		ts, err := parseTimestamp(timestamp)
		if err != nil {
			return nil, err
		}

		rsvp := &invite.RSVP{
			Id:               rsvpID,
			InviteId:         inviteUUID,
			Timestamp:        ts,
			AttendingDay:     attendingDay,
			AttendingEvening: attendingEvening,
			DietNotes:        dietNotes,
			Message:          message,
		}

		rsvps = append(rsvps, rsvp)
	}

	return rsvps, rows.Err()
}

func (s *Store) writeRsvp(rsvp invite.RSVP) {

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
