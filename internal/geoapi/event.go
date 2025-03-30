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
}

type Event struct {
	Guests         []Guest
	EventType      string
	GuestLocations LocationRegistry
}

type LocationRegistry struct {
	DistanceMatrix [][]float64
	CoordianteMap  CoordinateMapping
}

type CoordinateMapping struct {
	DestinationOccupancy map[coordinates.GuestCoordinates]int
	CoordinateToAddress  map[coordinates.GuestCoordinates]string
	AddressOrder         []string
}

var coordListURL string

// Filter for only Confirmed or GroceryOnly
func (e *Event) FilterGuestForService() {
	filteredGuests := make([]Guest, 0)
	for _, g := range e.Guests {
		if g.Status == Confirmed || g.Status == GroceryOnly {
			filteredGuests = append(filteredGuests, g)
		}
	}
	e.Guests = filteredGuests
}

func (e *Event) RequestGuestCoordiantes() error {
	if e.GuestLocations.CoordianteMap.DestinationOccupancy == nil &&
		e.GuestLocations.CoordianteMap.CoordinateToAddress == nil {

		e.GuestLocations.CoordianteMap.DestinationOccupancy = make(map[coordinates.GuestCoordinates]int)
		e.GuestLocations.CoordianteMap.CoordinateToAddress = make(map[coordinates.GuestCoordinates]string)
		e.GuestLocations.CoordianteMap.AddressOrder = make([]string, 0)
	}

	depotAddr := "555 Parkdale Ave"
	depotCoor, err := retreiveAddressCoordinate(depotAddr)
	if err != nil {
		return err
	}
	e.GuestLocations.CoordianteMap.AddressOrder = append(e.GuestLocations.CoordianteMap.AddressOrder, depotAddr)

	depotCoorString := depotCoor.ToString()
	addToCoordListString(&depotCoorString)

	for i := range e.Guests {
		err := e.Guests[i].geocodeGuestAddress()
		if err != nil {
			return err
		}
		coor, unique := e.isUnique(i)
		if unique {
			addToCoordListString(&coor)
		}

	}
	fmt.Println("Retrived all coordinates successfully")
	return nil
}

func addToCoordListString(uniqueCoordinate *string) {
	coordListURL += *uniqueCoordinate

}

func (e *Event) isUnique(guestIndex int) (string, bool) {
	g := e.Guests[guestIndex]
	val, ok := e.GuestLocations.CoordianteMap.DestinationOccupancy[g.Coordinates]

	if ok {
		// Already exists, update and return false
		e.GuestLocations.CoordianteMap.DestinationOccupancy[g.Coordinates] = val + g.GroupSize
		return "", false
	}

	// First time seeing this coordinate
	e.GuestLocations.CoordianteMap.DestinationOccupancy[g.Coordinates] = g.GroupSize
	e.GuestLocations.CoordianteMap.CoordinateToAddress[g.Coordinates] = g.Address
	e.GuestLocations.CoordianteMap.AddressOrder = append(e.GuestLocations.CoordianteMap.AddressOrder, g.Address)
	return g.Coordinates.ToString(), true
}
