package invite

import (
	"fmt"
	"strings"
)

type InviteWithRSVPs struct {
	Invite *Invite
	RSVPs  []*RSVP
}

func (i *InviteWithRSVPs) String() string {
	var b strings.Builder
	b.WriteString("Invite:\n")
	fmt.Fprintf(&b, "  Id:                %s\n", i.Invite.Id)
	fmt.Fprintf(&b, "  Name:              %s\n", i.Invite.Name)
	fmt.Fprintf(&b, "  Day:               %t\n", i.Invite.Day)
	fmt.Fprintf(&b, "  Evening:           %t\n", i.Invite.Evening)
	b.WriteString("RSVPs:\n")
	for j, r := range i.RSVPs {
		fmt.Fprintf(&b, "\tRSVP #%d\n", j+1)
		lines := strings.SplitSeq(r.String(), "\n")
		for line := range lines {
			b.WriteString("\t " + line + "\n")
		}
	}
	return b.String()
}
