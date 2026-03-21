package invite

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RSVP struct {
	Id               uuid.UUID `json:"-"`
	InviteId         uuid.UUID `json:"-"`
	Timestamp        time.Time `json:"-"`
	AttendingDay     bool      `json:"attendingDay"`
	AttendingEvening bool      `json:"attendingEvening"`
	HasPresentation  bool	   `json:"hasPresentation"`
	DietNotes        string    `json:"dietNotes"`
	Message          string    `json:"message"`
}

func NewRsvp(inviteId string) (*RSVP, error) {
	inviteUuid, err := uuid.Parse(inviteId)
	if err != nil {
		return nil, err
	}

	return &RSVP{
		Id: uuid.New(),
		InviteId: inviteUuid,
		Timestamp: time.Now(),
	}, nil
}

func (r RSVP) String() string {
	return fmt.Sprintf(
		"RSVP:\n"+
		"  ID:                  %s\n"+
		"  Invite ID:           %s\n"+
		"  Timestamp:           %s\n"+
		"  Attending Day:       %t\n"+
		"  Attending Evening:   %t\n"+
		"  Diet Notes:          %s\n"+
		"  Message:             %s",
		r.Id,
		r.InviteId,
		r.Timestamp.Format(time.RFC3339),
		r.AttendingDay,
		r.AttendingEvening,
		r.DietNotes,
		r.Message,
	)
}
