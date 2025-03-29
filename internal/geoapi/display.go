package geoapi

import (
	"fmt"
)

func (e *Event) DisplayEvent() {
	space(2)
	e.displayMatrix()

	space(3)
	// Display people per location
	e.displayCountPerAddress()

	space(3)
	// Display guest information
	for _, g := range e.Guests {
		g.displayGuestInformation()
	}

	space(3)
	// Display event type
	fmt.Println("Event type: ", e.EventType)
}

func (e *Event) displayCountPerAddress() {
	total := 0
	fmt.Println("(Coordinate): Number of People at this Address")
	for address, count := range e.GuestLocations.CoordianteMap.DestinationOccupancy {
		fmt.Printf("(%f, %f): %d \n", address.Long, address.Lat, count)
		total += count
	}

	fmt.Println("Total people that require service. ", total)
}

func (g *Guest) displayGuestInformation() {
	fmt.Printf("Guest: %s needs a ride to %s(%f, %f) \n", g.Name, g.Address, g.Coordinates.Long, g.Coordinates.Lat)
}

func (e *Event) displayMatrix() {

	matrix := &e.GuestLocations.DistanceMatrix

	fmt.Println("Guest Information")

	// Find max name length for padding
	maxNameLen := 0
	for _, name := range e.GuestLocations.CoordianteMap.AddressOrder {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
	}

	// Column width for numbers (float with 2 decimal places)
	cellWidth := 10

	// Print column headers
	fmt.Printf("%-*s", maxNameLen+2, "") // empty top-left corner
	for _, name := range e.GuestLocations.CoordianteMap.AddressOrder {
		fmt.Printf("%-*s", cellWidth, truncate(name, cellWidth-1))
	}
	fmt.Println()

	// Print each row
	for i, row := range *matrix {
		// Row header (name)
		fmt.Printf("%-*s", maxNameLen+2, e.GuestLocations.CoordianteMap.AddressOrder[i])

		// Row values
		for _, val := range row {
			fmt.Printf("%*.2f", cellWidth, val)
		}
		fmt.Println()
	}
}

func truncate(name string, max int) string {
	if len(name) <= max {
		return name
	}
	return name[:max-1] + "…"
}

func space(lines int) {
	for i := 0; i < lines; i++ {
		fmt.Println()
	}
}
