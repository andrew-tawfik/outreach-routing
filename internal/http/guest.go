package http

type GuestStatus int

const (
	Pending GuestStatus = iota
	Confirmed
	No
	GroceryOnly
	NotStarted
	Undecided
)

type Guest struct {
	Status  GuestStatus
	Name    string
	Address string
}

type Event struct {
	Guests    []Guest
	EventType string
}
