package domain

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Invite struct {
	Id      uuid.UUID
	People  []Person // 1-2 elements
	Day     bool
	Evening bool
}

func NewInvite(firstName, lastName string, fullDay bool) *Invite {
	return &Invite{
		Id: uuid.New(),
		People: []Person{
			{FirstName: firstName, LastName: lastName},
		},
		Day:     fullDay,
		Evening: true,
	}
}

// DisplayName returns the formatted display name for the invite
// Single person: "Sarah Johnson"
// Couple, same last name: "Sarah & Michael Johnson"
// Couple, different last names: "Sarah Johnson & Michael Smith"
func (i *Invite) DisplayName() string {
	if len(i.People) == 0 {
		return ""
	}

	if len(i.People) == 1 {
		return i.People[0].FullName()
	}

	// Two or more people
	if i.ShareLastName() {
		return i.People[0].FirstName + " & " + i.People[1].FirstName + " " + i.People[0].LastName
	}

	return i.People[0].FullName() + " & " + i.People[1].FullName()
}

// ShareLastName returns true if the first two people share the same last name
func (i *Invite) ShareLastName() bool {
	if len(i.People) < 2 {
		return false
	}
	return i.People[0].LastName == i.People[1].LastName
}

// FirstPersonFullName returns the full name of the first person
func (i *Invite) FirstPersonFullName() string {
	if len(i.People) > 0 {
		return i.People[0].FullName()
	}
	return ""
}

// SecondPersonFullName returns the full name of the second person, or empty string if no second person
func (i *Invite) SecondPersonFullName() string {
	if len(i.People) > 1 {
		return i.People[1].FullName()
	}
	return ""
}

// PeopleCount returns the number of people on this invite
func (i *Invite) PeopleCount() int {
	return len(i.People)
}

// InviteWithRSVPs combines an invite with its RSVP history
type InviteWithRSVPs struct {
	Invite *Invite
	RSVPs  []*RSVP
}

// DayAttendees returns the count of people confirming attendance for the day event
// Sums across all PersonAttendances in all RSVPs
func (i *InviteWithRSVPs) DayAttendees() int {
	count := 0
	for _, rsvp := range i.RSVPs {
		for _, attendance := range rsvp.Attendances {
			if attendance.AttendingDay {
				count++
			}
		}
	}
	return count
}

// EveningAttendees returns the count of people confirming attendance for the evening event
// Sums across all PersonAttendances in all RSVPs
func (i *InviteWithRSVPs) EveningAttendees() int {
	count := 0
	for _, rsvp := range i.RSVPs {
		for _, attendance := range rsvp.Attendances {
			if attendance.AttendingEvening {
				count++
			}
		}
	}
	return count
}

// LatestRSVP returns the most recent RSVP for this invite, or nil if no RSVPs exist
func (i *InviteWithRSVPs) LatestRSVP() *RSVP {
	if len(i.RSVPs) == 0 {
		return nil
	}
	latest := i.RSVPs[0]
	for _, rsvp := range i.RSVPs[1:] {
		if rsvp.Timestamp.After(latest.Timestamp) {
			latest = rsvp
		}
	}
	return latest
}

func (i *InviteWithRSVPs) String() string {
	var b strings.Builder
	b.WriteString("Invite:\n")
	fmt.Fprintf(&b, "  Id:                %s\n", i.Invite.Id)
	fmt.Fprintf(&b, "  DisplayName:       %s\n", i.Invite.DisplayName())
	fmt.Fprintf(&b, "  Day:               %t\n", i.Invite.Day)
	fmt.Fprintf(&b, "  Evening:           %t\n", i.Invite.Evening)
	b.WriteString("RSVPs:\n")
	for j, r := range i.RSVPs {
		fmt.Fprintf(&b, "\tRSVP #%d\n", j+1)
		lines := strings.Split(r.String(), "\n")
		for _, line := range lines {
			b.WriteString("\t " + line + "\n")
		}
	}
	return b.String()
}
