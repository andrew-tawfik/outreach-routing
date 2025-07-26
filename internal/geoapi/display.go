package geoapi

import (
	"fmt"
)

func (e *Event) DisplayEvent() {
	space(2)
	e.displayMatrix()

	space(3)
	
	e.displayCountPerAddress()

	space(3)
	
	for _, g := range e.Guests {
		g.displayGuestInformation()
	}

	space(3)
	
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

	
	maxNameLen := 0
	for _, name := range e.GuestLocations.CoordianteMap.AddressOrder {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
	}

	
	cellWidth := 10

	
	fmt.Printf("%-*s", maxNameLen+2, "") 
	for _, name := range e.GuestLocations.CoordianteMap.AddressOrder {
		fmt.Printf("%-*s", cellWidth, truncate(name, cellWidth-1))
	}
	fmt.Println()

	
	for i, row := range *matrix {
		
		fmt.Printf("%-*s", maxNameLen+2, e.GuestLocations.CoordianteMap.AddressOrder[i])

		
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
