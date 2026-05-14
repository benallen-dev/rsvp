package invite

import (
	"github.com/google/uuid"
)

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

// func (i *Invite) Name() string {
// 	// If same last name, maybe reduce here
// 	return strings.Join(i.Names, ", ")
// }

// // The number of people invited
// func (i *Invite) Count() int {
// 	return len(i.Names)
// }

