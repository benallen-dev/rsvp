package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RSVP struct {
	Id              uuid.UUID          `json:"-"`
	InviteId        uuid.UUID          `json:"-"`
	Timestamp       time.Time          `json:"-"`
	Attendances     []PersonAttendance `json:"attendances"`
	HasPresentation bool               `json:"hasPresentation"`
	PhoneNumber     string             `json:"phoneNumber"`
	DietNotes       string             `json:"dietNotes"`
	Message         string             `json:"message"`
}

func NewRsvp(inviteId string) (*RSVP, error) {
	inviteUuid, err := uuid.Parse(inviteId)
	if err != nil {
		return nil, err
	}

	return &RSVP{
		Id:        uuid.New(),
		InviteId:  inviteUuid,
		Timestamp: time.Now(),
	}, nil
}

// GetPersonAttendance returns the attendance record for a specific person index, or nil if not found
func (r *RSVP) GetPersonAttendance(personIndex int) *PersonAttendance {
	for _, a := range r.Attendances {
		if a.PersonIndex == personIndex {
			return &a
		}
	}
	return nil
}

func (r RSVP) String() string {
	var b strings.Builder
	b.WriteString("RSVP:\n")
	fmt.Fprintf(&b, "  ID:                %s\n", r.Id)
	fmt.Fprintf(&b, "  Invite ID:         %s\n", r.InviteId)
	fmt.Fprintf(&b, "  Timestamp:         %s\n", r.Timestamp.Format(time.RFC3339))
	for _, attendance := range r.Attendances {
		fmt.Fprintf(&b, "  Person %d:\n", attendance.PersonIndex)
		fmt.Fprintf(&b, "    Attending Day:     %t\n", attendance.AttendingDay)
		fmt.Fprintf(&b, "    Attending Evening: %t\n", attendance.AttendingEvening)
	}
	fmt.Fprintf(&b, "  Has Presentation:  %t\n", r.HasPresentation)
	fmt.Fprintf(&b, "  Phone Number:      %s\n", r.PhoneNumber)
	fmt.Fprintf(&b, "  Diet Notes:        %s\n", r.DietNotes)
	fmt.Fprintf(&b, "  Message:           %s\n", r.Message)
	return b.String()
}
