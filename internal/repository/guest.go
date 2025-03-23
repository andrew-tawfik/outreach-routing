package repository

type GuestStatus int

const (
	Pending GuestStatus = iota
	Confirmed
	No
	GroceryOnly
	NotStarted
)

type Guest struct {
	Name        string
	PhoneNumber string
	Address     string
	Status      GuestStatus
}
