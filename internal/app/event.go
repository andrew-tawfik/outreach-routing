package app

import "github.com/andrew-tawfik/outreach-routing/internal/coordinates"

// LocationRegistry holds all spatial data necessary for route planning.
type LocationRegistry struct {
	DistanceMatrix [][]float64       // Distance between all pairs of guest locations (including depot at index 0)
	CoordianteMap  CoordinateMapping // Metadata for interpreting coordinate
}

// CoordinateMapping organizes geospatial guest information.
type CoordinateMapping struct {
	DestinationOccupancy map[coordinates.GuestCoordinates]int    // Number of guests at each unique coordinate
	CoordinateToAddress  map[string]coordinates.GuestCoordinates // Maps address strings to coordinate objects
	AddressOrder         []string                                // Ordered list of addresses to preserve consistent indexing
}

// GuestCoordinates represents a pair of geographic coordinates (longitude and latitude).
type GuestCoordinates struct {
	Long float64
	Lat  float64
}

// Event represent a real-world event (e.g., dinner or grocery delivery) with its guests.
type Event struct {
	Guests    []Guest
	EventType string
}

// Guest represents a single person or group needing transportation.
type Guest struct {
	Name        string
	GroupSize   int
	Coordinates coordinates.GuestCoordinates
	Address     string
	PhoneNumber string
}
