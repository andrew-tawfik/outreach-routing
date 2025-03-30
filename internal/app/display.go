package app

import (
	"fmt"
)

func (e *Event) Display() {
	fmt.Println("=== Event ===")
	fmt.Printf("Type: %s\n", e.EventType)
	fmt.Println("Guests:")
	for i, guest := range e.Guests {
		fmt.Printf("  [%d] Name: %s | Group Size: %d | Coordinates: (Lat: %.6f, Long: %.6f)\n",
			i, guest.Name, guest.GroupSize, guest.Coordinates.Long, guest.Coordinates.Lat)
	}
	fmt.Println()
}

func (r *LocationRegistry) Display() {
	fmt.Println("=== Location Registry ===")
	fmt.Println("Destination Occupancy:")
	for coord, count := range r.CoordianteMap.DestinationOccupancy {
		fmt.Printf("  (%0.6f, %0.6f) → %d guests\n", coord.Long, coord.Lat, count)
	}

	fmt.Println("\nCoordinate to Address:")
	for coord, addr := range r.CoordianteMap.CoordinateToAddress {
		fmt.Printf("  (%0.6f, %0.6f) → %s\n", coord.Long, coord.Lat, addr)
	}

	fmt.Println("\nDistance Matrix:")
	r.displayMatrix()
	fmt.Println()
}

func (lr *LocationRegistry) displayMatrix() {

	matrix := &lr.DistanceMatrix

	fmt.Println("Guest Information")

	// Find max name length for padding
	maxNameLen := 0
	for _, name := range lr.CoordianteMap.AddressOrder {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
	}

	// Column width for numbers (float with 2 decimal places)
	cellWidth := 10

	// Print column headers
	fmt.Printf("%-*s", maxNameLen+2, "") // empty top-left corner
	for _, name := range lr.CoordianteMap.AddressOrder {
		fmt.Printf("%-*s", cellWidth, truncate(name, cellWidth-1))
	}
	fmt.Println()

	// Print each row
	for i, row := range *matrix {
		// Row header (name)
		fmt.Printf("%-*s", maxNameLen+2, lr.CoordianteMap.AddressOrder[i])

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
