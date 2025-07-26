package app

import (
	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
	"github.com/andrew-tawfik/outreach-routing/internal/geoapi"
)


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
	ApiErrors geoapi.ApiErrors
}


type Guest struct {
	Name        string
	GroupSize   int
	Coordinates coordinates.GuestCoordinates
	Address     string
	PhoneNumber string
}
