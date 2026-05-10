package invite

import "github.com/google/uuid"

// Changing this? Don't forget to update the store tables!
type Invite struct {
	Id      uuid.UUID
	Name    string
	Day     bool
	Evening bool
}

func NewInvite(name string, fullDay bool) *Invite {
	return &Invite{
		Id:      uuid.New(),
		Name:    name,
		Day:     fullDay,
		Evening: true,
	}
}
