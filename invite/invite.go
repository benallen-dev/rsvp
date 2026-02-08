package invite

import "github.com/google/uuid"

type Invite struct {
	id      uuid.UUID
	name    string
	day     bool
	evening bool
}

func NewInvite(name string, fullDay bool) *Invite {
	return &Invite{
		id:      uuid.New(),
		name:    name,
		day:     fullDay,
		evening: true,
	}
}


