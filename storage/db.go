package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/domain"
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
			person_0_first_name TEXT NOT NULL,
			person_0_last_name TEXT NOT NULL,
			person_1_first_name TEXT,
			person_1_last_name TEXT,
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
            person_0_attending_day BOOLEAN,
            person_0_attending_evening BOOLEAN,
            person_1_attending_day BOOLEAN,
            person_1_attending_evening BOOLEAN,
            has_presentation BOOLEAN,
            phone_number TEXT,
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
func (s *Store) ReadAllInvitesWithRSVPs() ([]*domain.InviteWithRSVPs, error) {
	rows, err := s.db.Query(`SELECT id, person_0_first_name, person_0_last_name, person_1_first_name, person_1_last_name, day, evening FROM invites`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.InviteWithRSVPs

	for rows.Next() {
		var id, person0FirstName, person0LastName string
		var person1FirstName, person1LastName sql.NullString
		var day, evening bool

		if err := rows.Scan(&id, &person0FirstName, &person0LastName, &person1FirstName, &person1LastName, &day, &evening); err != nil {
			return nil, err
		}

		// Parse UUID
		uuid, err := parseUUID(id)
		if err != nil {
			return nil, err
		}

		// Build People array
		people := []domain.Person{
			{FirstName: person0FirstName, LastName: person0LastName},
		}
		if person1FirstName.Valid && person1LastName.Valid {
			people = append(people, domain.Person{FirstName: person1FirstName.String, LastName: person1LastName.String})
		}

		inv := &domain.Invite{
			Id:      uuid,
			People:  people,
			Day:     day,
			Evening: evening,
		}

		// Get RSVPs for this invite
		rsvps, err := s.getRSVPsForInvite(id)
		if err != nil {
			return nil, err
		}

		result = append(result, &domain.InviteWithRSVPs{
			Invite: inv,
			RSVPs:  rsvps,
		})
	}

	return result, rows.Err()
}

// ReadInviteWithRSVPs returns a single invite with its associated RSVPs
func (s *Store) ReadInviteWithRSVPs(id string) (*domain.InviteWithRSVPs, error) {
	var person0FirstName, person0LastName string
	var person1FirstName, person1LastName sql.NullString
	var day, evening bool

	err := s.db.QueryRow(`SELECT person_0_first_name, person_0_last_name, person_1_first_name, person_1_last_name, day, evening FROM invites WHERE id = ?`, id).
		Scan(&person0FirstName, &person0LastName, &person1FirstName, &person1LastName, &day, &evening)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	// Build People array
	people := []domain.Person{
		{FirstName: person0FirstName, LastName: person0LastName},
	}
	if person1FirstName.Valid && person1LastName.Valid {
		people = append(people, domain.Person{FirstName: person1FirstName.String, LastName: person1LastName.String})
	}

	inv := &domain.Invite{
		Id:      uuid,
		People:  people,
		Day:     day,
		Evening: evening,
	}

	// Get RSVPs for this invite
	rsvps, err := s.getRSVPsForInvite(id)
	if err != nil {
		return nil, err
	}

	return &domain.InviteWithRSVPs{
		Invite: inv,
		RSVPs:  rsvps,
	}, nil
}

// getRSVPsForInvite is a helper function that fetches all RSVPs for a given invite ID
func (s *Store) getRSVPsForInvite(inviteID string) ([]*domain.RSVP, error) {
	rows, err := s.db.Query(`
		SELECT id, invite_id, timestamp, person_0_attending_day, person_0_attending_evening, person_1_attending_day, person_1_attending_evening, has_presentation, phone_number, diet_notes, message 
		FROM rsvps 
		WHERE invite_id = ?
		ORDER BY timestamp
	`, inviteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rsvps []*domain.RSVP

	for rows.Next() {
		var id, inviteIDStr, timestamp, phoneNumber, dietNotes, message string
		var person0AttendingDay, person0AttendingEvening, person1AttendingDay, person1AttendingEvening, hasPresentation bool

		if err := rows.Scan(&id, &inviteIDStr, &timestamp, &person0AttendingDay, &person0AttendingEvening, &person1AttendingDay, &person1AttendingEvening, &hasPresentation, &phoneNumber, &dietNotes, &message); err != nil {
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

		// Build Attendances array
		attendances := []domain.PersonAttendance{
			{PersonIndex: 0, AttendingDay: person0AttendingDay, AttendingEvening: person0AttendingEvening},
		}
		// Only add person 1 if they have any attendance data
		if person1AttendingDay || person1AttendingEvening {
			attendances = append(attendances, domain.PersonAttendance{
				PersonIndex:      1,
				AttendingDay:     person1AttendingDay,
				AttendingEvening: person1AttendingEvening,
			})
		}

		rsvp := &domain.RSVP{
			Id:              rsvpID,
			InviteId:        inviteUUID,
			Timestamp:       ts,
			Attendances:     attendances,
			HasPresentation: hasPresentation,
			PhoneNumber:     phoneNumber,
			DietNotes:       dietNotes,
			Message:         message,
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
