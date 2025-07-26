package database

import (
	"strconv"

	"gopkg.in/Iwark/spreadsheet.v2"
)


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
	Status      GuestStatus 
	Name        string
	GroupSize   int 
	PhoneNumber string
	Address     string
}


func processGuest(row *[]spreadsheet.Cell) (Guest, bool) {
	status := determineGuestStatus((*row)[0].Value)
	name := (*row)[1].Value
	count := (*row)[2].Value
	phone := (*row)[3].Value
	address := (*row)[4].Value
	validGuest := true

	if count == "" {
		count = "0"
	}
	
	iCount, err := strconv.Atoi(count)
	if err != nil {
		validGuest = false
	}

	
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
