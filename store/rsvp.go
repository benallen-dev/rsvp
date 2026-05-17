package store

import (
	"errors"

	"rsvp/invite"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func (s *Store) CreateRSVP(inv *invite.RSVP) (*invite.RSVP, error) {
	log.Infof("Inserting RSVP into DB: id=%s, name=%s, type=%s", inv.Id, inv.Name, inv.RsvpType)
	
	result, err := s.db.Exec(`
		INSERT INTO rsvps (id, rsvp_type, name, attending_day, attending_evening, diet_notes, message)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, inv.Id.String(), inv.RsvpType, inv.Name, inv.AttendingDay, inv.AttendingEvening, inv.DietNotes, inv.Message)
	
	if err != nil {
		log.Errorf("Failed to insert RSVP: %v", err)
		return nil, err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		log.Errorf("Failed to get rows affected: %v", err)
		return nil, err
	}
	
	log.Infof("RSVP inserted successfully, rows affected: %d", rows)
	return inv, nil
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

func (s *Store) ReadAllRSVPs() ([]*invite.RSVP, error) {
	rows, err := s.db.Query(`
		SELECT id, rsvp_type, name, timestamp, attending_day, attending_evening, diet_notes, message
		FROM rsvps
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rsvps []*invite.RSVP

	for rows.Next() {
		var id, rsvpType, name, timestamp, dietNotes, message string
		var attendingDay, attendingEvening bool

		if err := rows.Scan(&id, &rsvpType, &name, &timestamp, &attendingDay, &attendingEvening, &dietNotes, &message); err != nil {
			log.Errorf("scan error: %v", err)
			return nil, err
		}

		// Parse timestamp
		ts, err := parseTimestamp(timestamp)
		if err != nil {
			log.Errorf("timestamp parse error: %v", err)
			return nil, err
		}

		// Parse UUID
		rsvpID, err := uuid.Parse(id)
		if err != nil {
			log.Errorf("uuid parse error: %v", err)
			return nil, err
		}

		rsvp := &invite.RSVP{
			Id:               rsvpID,
			Timestamp:        ts,
			RsvpType:         rsvpType,
			Name:             name,
			AttendingDay:     attendingDay,
			AttendingEvening: attendingEvening,
			DietNotes:        dietNotes,
			Message:          message,
		}

		rsvps = append(rsvps, rsvp)
	}

	return rsvps, rows.Err()
}

func (s *Store) ReadRSVPs() ([]*invite.RSVP, error) {
	return nil, errors.New("Not implemented")
}

// func (s *Store) getRSVPsForInvite(inviteID string) ([]*invite.RSVP, error) {
// 	// TODO: Prepared statement
// 	rows, err := s.db.Query(`
// 		SELECT id, invite_id, timestamp, attending_day, attending_evening, diet_notes, message 
// 		FROM rsvps 
// 		WHERE invite_id = ?
// 		ORDER BY timestamp
// 	`, inviteID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var rsvps []*invite.RSVP

// 	for rows.Next() {
// 		var id, inviteIDStr, timestamp, dietNotes, message string
// 		var attendingDay, attendingEvening bool

// 		if err := rows.Scan(&id, &inviteIDStr, &timestamp, &attendingDay, &attendingEvening, &dietNotes, &message); err != nil {
// 			return nil, err
// 		}

// 		// Parse UUIDs
// 		rsvpID, err := uuid.Parse(id)
// 		if err != nil {
// 			return nil, err
// 		}

// 		inviteUUID, err := uuid.Parse(inviteIDStr)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Parse timestamp
// 		ts, err := parseTimestamp(timestamp)
// 		if err != nil {
// 			return nil, err
// 		}

// 		rsvp := &invite.RSVP{
// 			Id:               rsvpID,
// 			Timestamp:        ts,
// 			AttendingDay:     attendingDay,
// 			AttendingEvening: attendingEvening,
// 			DietNotes:        dietNotes,
// 			Message:          message,
// 		}

// 		rsvps = append(rsvps, rsvp)
// 	}

// 	return rsvps, rows.Err()
// }
