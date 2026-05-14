package store

import (
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/invite"
)

type Store struct {
	db *sql.DB
}

func (s *Store) ReadInvite(id string) (*invite.Invite, error) {
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

	return &invite.Invite{
		Id:      uuid,
		Name:    name,
		Day:     day,
		Evening: evening,
	}, nil
}

func (s *Store) ReadAllInvites() ([]*invite.Invite, error) {
	rows, err := s.db.Query(`SELECT id, name, day, evening FROM invites ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*invite.Invite

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
		result = append(result, inv)
	}

	return result, rows.Err()
}

// ReadAllInvitesWithRSVPs returns all invites with their associated RSVPs
func (s *Store) ReadAllInvitesWithRSVPs() ([]*invite.InviteWithRSVPs, error) {
	rows, err := s.db.Query(`SELECT id, name, day, evening FROM invites ORDER BY name ASC`)
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

// Let's you know which mode the DB is in
func (s *Store) LogJournalMode() {
	row := s.db.QueryRow("PRAGMA journal_mode")
	var mode string
	err := row.Scan(&mode)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Journal mode: %s", mode)
}

