package database

import (
	"strconv"

	"gopkg.in/Iwark/spreadsheet.v2"
)

// GuestStatus represents the RSVP or participation status of a guest.
type GuestStatus int

// Defined guest statuses based on spreadsheet cell content.
const (
	Pending GuestStatus = iota
	Confirmed
	No
	GroceryOnly
	NotStarted
	Undecided
)

// Guest holds the structured information extracted from a spreadsheet row.
type Guest struct {
	Status      GuestStatus // Enum of guest participation
	Name        string
	GroupSize   int // Number of people in this guest's group
	PhoneNumber string
	Address     string
}

// processGuest converts a single row into a Guest struct.
func processGuest(row *[]spreadsheet.Cell) (Guest, bool) {
	status := determineGuestStatus((*row)[0].Value)
	name := (*row)[1].Value
	count := (*row)[2].Value
	phone := (*row)[3].Value
	address := (*row)[4].Value
	validGuest := true

	// Convert group size to integer and ensure it's > 0
	iCount, err := strconv.Atoi(count)
	if err != nil {
		validGuest = false
	}

	// Guests must have both a name and address
	if name == "" || address == "" {
		validGuest = false
	}

	if status == GroceryOnly {
		iCount = 0
	}

	return Guest{
		Status:      status,
		Name:        name,
		GroupSize:   iCount,
		PhoneNumber: phone,
		Address:     address,
	}, validGuest
}

// determineGuestStatus maps a spreadsheet status string to a GuestStatus enum.
func determineGuestStatus(cellContent string) GuestStatus {
	switch cellContent {
	case "Confirmed":
		return Confirmed
	case "Grocery Only":
		return GroceryOnly
	case "Pending":
		return Pending
	case "Not started":
		return NotStarted
	case "NO", "Not eligiable":
		return No
	default:
		return Undecided
	}
}
