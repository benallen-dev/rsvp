package store

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"

	"rsvp/invite"
)

func (s *Store) CreateInvite(inv *invite.Invite) error {
	// TODO: Worry about does the record already exist and that kind of stuff
	stmt, err := s.db.Prepare("INSERT INTO invites (id, name, day, evening) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(inv.Id.String(), inv.Name, inv.Day, inv.Evening)
	log.Info(res)
	return err
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
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &invite.Invite{
		Id:      uid,
		Name:    name,
		Day:     day,
		Evening: evening,
	}, nil
}

func (s *Store) UpdateInvite(inv *invite.Invite) (*invite.Invite, error) {
	return nil, errors.New("Not implemented")
}

func (s *Store) DestroyInvite(id string) error {
	return errors.New("Not implemented")
}

func (s *Store) ReadInvites() ([]*invite.Invite, error) {
	rows, err := s.db.Query(`SELECT id, name, day, evening FROM invites ORDER BY day DESC, name ASC`)
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
		result = append(result, inv)
	}

	return result, rows.Err()
}


func (s *Store) SearchInvites(q string) ([]*invite.Invite, error) {
	stmt, err := s.db.Prepare(`SELECT id, name, day, evening FROM invites WHERE name LIKE ? ORDER BY day DESC, name ASC`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + q + "%")
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
		result = append(result, inv)
	}

	return result, rows.Err()
}
