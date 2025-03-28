package geoapi

import "fmt"

func (e *Event) DisplayMatrix() {

	matrix := &e.GuestLocations.DistanceMatrix

	fmt.Println()
	fmt.Println()

	// Find max name length for padding
	maxNameLen := 0
	for _, name := range addressOrder {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
	}

	// Column width for numbers (float with 2 decimal places)
	cellWidth := 10

	// Print column headers
	fmt.Printf("%-*s", maxNameLen+2, "") // empty top-left corner
	for _, name := range addressOrder {
		fmt.Printf("%-*s", cellWidth, truncate(name, cellWidth-1))
	}
	fmt.Println()

	// Print each row
	for i, row := range *matrix {
		// Row header (name)
		fmt.Printf("%-*s", maxNameLen+2, addressOrder[i])

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
