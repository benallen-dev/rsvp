package invite

import "github.com/google/uuid"

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
