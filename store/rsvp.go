package store

import (
	"errors"

	"rsvp/invite"

	"github.com/google/uuid"
)

func (s *Store) CreateRSVP(inv *invite.RSVP) (*invite.RSVP, error) {
	return nil, errors.New("Not implemented")
}

func (s *Store) ReadRSVP(id string) (*invite.RSVP, error) {
	return nil, errors.New("Not implemented")
}

func (s *Store) UpdateRSVP(inv *invite.RSVP) (*invite.RSVP, error) {
	return nil, errors.New("Not implemented")
}

func (s *Store) DestroyRSVP(id string) error {
	return errors.New("Not implemented")
}

func (s *Store) ReadRSVPs() ([]*invite.RSVP, error) {
	return nil, errors.New("Not implemented")
}

func (s *Store) getRSVPsForInvite(inviteID string) ([]*invite.RSVP, error) {
	// TODO: Prepared statement
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
		rsvpID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		inviteUUID, err := uuid.Parse(inviteIDStr)
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
