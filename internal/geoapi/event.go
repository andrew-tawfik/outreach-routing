package geoapi

import (
	"fmt"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
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
	Address     string
	Coordinates coordinates.GuestCoordinates
	PhoneNumber string
}


type Event struct {
	Guests         []Guest
	EventType      string
	GuestLocations LocationRegistry
	ApiErrors      ApiErrors
}


type LocationRegistry struct {
	DistanceMatrix [][]float64
	CoordianteMap  CoordinateMapping
}


type CoordinateMapping struct {
	DestinationOccupancy map[coordinates.GuestCoordinates]int
	CoordinateToAddress  map[string]coordinates.GuestCoordinates
	AddressOrder         []string
}


var coordListURL string



func (e *Event) filterGuestForService() {
	filteredGuests := make([]Guest, 0)

	for _, g := range e.Guests {
		if e.EventType == "Grocery" {
			g.GroupSize = 0
		}
		if g.Status == Confirmed || g.Status == GroceryOnly {
			filteredGuests = append(filteredGuests, g)
		}
	}
	e.Guests = filteredGuests
}


func (e *Event) RequestGuestCoordiantes() error {
	ResetGlobalState()

	e.filterGuestForService()
	e.initCoordinateMap()

	
	depotAddr := "555 Parkdale Ave"
	depotCoor, err := retreiveAddressCoordinate(depotAddr)
	if err != nil {
		return fmt.Errorf("failed to geocode SMSM address: %w", err)
	}
	e.GuestLocations.CoordianteMap.AddressOrder = append(e.GuestLocations.CoordianteMap.AddressOrder, depotAddr)

	depotCoorString := depotCoor.ToString()
	addToCoordListString(&depotCoorString)

	e.geocodeEvent()
	return nil
}

func (e *Event) initCoordinateMap() {
	
	if e.GuestLocations.CoordianteMap.DestinationOccupancy == nil &&
		e.GuestLocations.CoordianteMap.CoordinateToAddress == nil {

		e.GuestLocations.CoordianteMap.DestinationOccupancy = make(map[coordinates.GuestCoordinates]int)
		e.GuestLocations.CoordianteMap.CoordinateToAddress = make(map[string]coordinates.GuestCoordinates)
		e.GuestLocations.CoordianteMap.AddressOrder = make([]string, 0)
	}
}



func addToCoordListString(uniqueCoordinate *string) {
	coordListURL += *uniqueCoordinate

}


func (e *Event) isUnique(guestIndex int) (string, bool) {
	g := e.Guests[guestIndex]
	val, ok := e.GuestLocations.CoordianteMap.DestinationOccupancy[g.Coordinates]

	if ok {
		
		e.GuestLocations.CoordianteMap.DestinationOccupancy[g.Coordinates] = val + g.GroupSize
		return "", false
	}

	
	e.GuestLocations.CoordianteMap.DestinationOccupancy[g.Coordinates] = g.GroupSize
	e.GuestLocations.CoordianteMap.CoordinateToAddress[g.Address] = g.Coordinates
	e.GuestLocations.CoordianteMap.AddressOrder = append(e.GuestLocations.CoordianteMap.AddressOrder, g.Address)
	return g.Coordinates.ToString(), true
}

func ResetGlobalState() {
	coordListURL = ""
}
