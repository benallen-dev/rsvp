package invite

import "github.com/google/uuid"

type Invite struct {
	id uuid.UUID
	name string
	day bool
	evening bool
}
	
func NewInvite(name string, day, evening bool) *Invite {
	return &Invite{
		id: uuid.New(),
		name: name,
		day: day,
		evening: evening,
	}
}
