package invite

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Changing this? Don't forget to update the store tables!
type RSVP struct {
	Id                 uuid.UUID `json:"id"`
	Timestamp          time.Time `json:"timestamp"`
	Type               string    `json:"type"`
	Name               string    `json:"name"`
	AttendingCeremony  bool      `json:"attendingCeremony"`
	AttendingReception bool      `json:"attendingReception"`
	AttendingDinner    bool      `json:"attendingDinner"`
	AttendingParty     bool      `json:"attendingParty"`
	DietNotes          string    `json:"dietNotes"`
	Message            string    `json:"message"`
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
			"  Attending Ceremony:  %t\n"+
			"  Attending Reception: %t\n"+
			"  Attending Dinner:    %t\n"+
			"  Attending Party:     %t\n"+
			"  Diet Notes:          %s\n"+
			"  Message:             %s",
		r.Id,
		r.Type,
		r.Name,
		r.Timestamp.Format(time.RFC3339),
		r.AttendingCeremony,
		r.AttendingReception,
		r.AttendingDinner,
		r.AttendingParty,
		r.DietNotes,
		r.Message,
	)
}
