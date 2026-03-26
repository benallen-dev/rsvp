package domain

type Person struct {
	FirstName string
	LastName  string
}

func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

type PersonAttendance struct {
	PersonIndex      int // 0 or 1 (index into Invite.People)
	AttendingDay     bool
	AttendingEvening bool
}
