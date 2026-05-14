package store

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"rsvp/invite"
)

type Store struct {
	db *sql.DB
}

// ReadAllInvitesWithRSVPs returns all invites with their associated RSVPs
func (s *Store) ReadAllInvitesWithRSVPs() ([]*invite.InviteWithRSVP, error) {
	invites, err := s.ReadInvites()
	if err != nil {
		return nil, err
	}

	var result []*invite.InviteWithRSVP

	for _, inv := range invites {
		rsvps, err := s.getRSVPsForInvite(inv.Id.String())
		if err != nil {
			return nil, err
		}

		result = append(result, &invite.InviteWithRSVP{
			Invite: inv,
			RSVPs:  rsvps,
		})
	}

	return result, err
}

// ReadInviteWithRSVPs returns a single invite with its associated RSVPs
func (s *Store) ReadInviteWithRSVPs(id string) (*invite.InviteWithRSVP, error) {
	var name string
	var day, evening bool

	err := s.db.QueryRow(`SELECT name, day, evening FROM invites WHERE id = ?`, id).
		Scan(&name, &day, &evening)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	inv := &invite.Invite{
		Id:      uid,
		Name:    name,
		Day:     day,
		Evening: evening,
	}

	// Get RSVPs for this invite
	rsvps, err := s.getRSVPsForInvite(id)
	if err != nil {
		return nil, err
	}

	return &invite.InviteWithRSVP{
		Invite: inv,
		RSVPs:  rsvps,
	}, nil
}

