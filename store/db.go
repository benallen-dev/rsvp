package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/invite"
)

type Store struct {
	db *sql.DB
}


// Returns (*Store, db.Close, error)
func Init(dbFile string) (*Store, func() error, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, nil, err
	}

	s := &Store{db: db}

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

// ReadAllInvitesWithRSVPs returns all invites with their associated RSVPs
func (s *Store) ReadAllInvitesWithRSVPs() ([]*invite.InviteWithRSVPs, error) {
	rows, err := s.db.Query(`SELECT id, name, day, evening FROM invites`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*invite.InviteWithRSVPs

	for rows.Next() {
		var id, name string
		var day, evening bool

		if err := rows.Scan(&id, &name, &day, &evening); err != nil {
			return nil, err
		}

		// Parse UUID
		uuid, err := parseUUID(id)
		if err != nil {
			return nil, err
		}

		inv := &invite.Invite{
			Id:      uuid,
			Name:    name,
			Day:     day,
			Evening: evening,
		}

		// Get RSVPs for this invite
		rsvps, err := s.getRSVPsForInvite(id)
		if err != nil {
			return nil, err
		}

		result = append(result, &invite.InviteWithRSVPs{
			Invite: inv,
			RSVPs:  rsvps,
		})
	}

	return result, rows.Err()
}

// ReadInviteWithRSVPs returns a single invite with its associated RSVPs
func (s *Store) ReadInviteWithRSVPs(id string) (*invite.InviteWithRSVPs, error) {
	var name string
	var day, evening bool

	err := s.db.QueryRow(`SELECT name, day, evening FROM invites WHERE id = ?`, id).
		Scan(&name, &day, &evening)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	inv := &invite.Invite{
		Id:      uuid,
		Name:    name,
		Day:     day,
		Evening: evening,
	}

	// Get RSVPs for this invite
	rsvps, err := s.getRSVPsForInvite(id)
	if err != nil {
		return nil, err
	}

	return &invite.InviteWithRSVPs{
		Invite: inv,
		RSVPs:  rsvps,
	}, nil
}

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

// Helper functions for parsing

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

func writeRsvp(rsvp invite.RSVP) {

}
