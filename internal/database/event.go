package database

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/Iwark/spreadsheet.v2"
)

type Event struct {
	Guests    []Guest
	EventType string
}

func (db *Database) ProcessEvent() (*Event, error) {

	sheet := &db.sheet.Sheets[0]
	et, err := determineEventType(sheet)
	if err != nil {
		log.Fatalf("Please fix spreadsheet title. %v", err)
	}

	firstRow := &db.sheet.Sheets[0].Rows[0]
	err = verifyColumnTitles(firstRow)

	if err != nil {
		log.Fatalf("Column title verification failed: %v", err)
	}

	guests := make([]Guest, 0, 30)

	for i := 1; i < len(db.sheet.Sheets[0].Rows); i++ {
		g, ok := processGuest(&db.sheet.Sheets[0].Rows[i])
		if ok {
			guests = append(guests, g)
		}
	}
	return &Event{Guests: guests, EventType: et}, nil
}

func determineEventType(sheet *spreadsheet.Sheet) (string, error) {
	title := sheet.Properties.Title

	if strings.Contains(title, "Dinner") {
		return "Dinner", nil
	} else if strings.Contains(title, "Grocery") {
		return "Grocery", nil
	} else {
		return "", fmt.Errorf("title must include either 'Dinner' or 'Grocery'")
	}
}

func verifyColumnTitles(row *[]spreadsheet.Cell) error {
	correctOrder := []string{"Status", "Name", "Group Size", "Number", "Address"}

	if len(*row) < len(correctOrder) {
		return fmt.Errorf("row too short: expected at least %d columns, got %d", len(correctOrder), len(*row))
	}

	for i, expected := range correctOrder {
		if (*row)[i].Value != expected {
			return fmt.Errorf("column %d mismatch: expected %q, got %q", i, expected, (*row)[i].Value)
		}
	}
	fmt.Println("Columns Verified !")

	return nil // all good
}
