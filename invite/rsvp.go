package invite

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Changing this? Don't forget to update the store tables!
type RSVP struct {
	Id               uuid.UUID `json:"-"`
	Timestamp        time.Time `json:"-"`
	RsvpType         string    `json:"rsvpType"`
	Name             string    `json:"name"`
	AttendingDay     bool      `json:"attendingDay"`
	AttendingEvening bool      `json:"attendingEvening"`
	DietNotes        string    `json:"dietNotes"`
	Message          string    `json:"message"`
}

func NewRsvp() *RSVP {
	return &RSVP{
		Id:        uuid.New(),
		Timestamp: time.Now(),
	}
}

func (r RSVP) String() string {
	return fmt.Sprintf(
		"RSVP:\n"+
		"  ID:                  %s\n"+
		"  Type:                %s\n"+
		"  Name:                %s\n"+
		"  Timestamp:           %s\n"+
		"  Attending Day:       %t\n"+
		"  Attending Evening:   %t\n"+
		"  Diet Notes:          %s\n"+
		"  Message:             %s",
		r.Id,
		r.RsvpType,
		r.Name,
		r.Timestamp.Format(time.RFC3339),
		r.AttendingDay,
		r.AttendingEvening,
		r.DietNotes,
		r.Message,
	)
}
