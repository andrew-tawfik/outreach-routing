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

	iCount, err := strconv.Atoi(count)
	if err != nil || iCount <= 0 { // Ensure iCount > 0
		validGuest = false
	}

	if name == "" || address == "" {
		validGuest = false
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
	if cellContent == "Confirmed" {
		return Confirmed
	} else if cellContent == "Grocery Only" {
		return GroceryOnly
	} else if cellContent == "Pending" {
		return Pending
	} else if cellContent == "Not started" {
		return NotStarted
	} else if cellContent == "NO" || cellContent == "Not eligiable" {
		return No
	} else {
		return Undecided
	}
}
