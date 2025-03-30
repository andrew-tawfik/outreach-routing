package app

import "github.com/andrew-tawfik/outreach-routing/internal/coordinates"

type LocationRegistry struct {
	DistanceMatrix [][]float64
	CoordianteMap  CoordinateMapping
}

type CoordinateMapping struct {
	DestinationOccupancy map[coordinates.GuestCoordinates]int
	CoordinateToAddress  map[string]coordinates.GuestCoordinates
	AddressOrder         []string
}

type GuestCoordinates struct {
	Long float64
	Lat  float64
}

type Event struct {
	Guests    []Guest
	EventType string
}

type Guest struct {
	Name        string
	GroupSize   int
	Coordinates coordinates.GuestCoordinates
}
