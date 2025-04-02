package geoapi

import (
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

// GuestStatus represents a guest's eligibility for routing services.
type GuestStatus int

// Enum of guest statuses
const (
	Pending GuestStatus = iota
	Confirmed
	No
	GroceryOnly
	NotStarted
	Undecided
)

// Guest represents a single person or group needing transport
type Guest struct {
	Status      GuestStatus
	Name        string
	GroupSize   int
	Address     string
	Coordinates coordinates.GuestCoordinates
}

// Event represents either Dinner or Grocery Run holding Guest and Address information
type Event struct {
	Guests         []Guest
	EventType      string
	GuestLocations LocationRegistry
}

// Location Registry holds the distance matrix and additional coordinate information
type LocationRegistry struct {
	DistanceMatrix [][]float64
	CoordianteMap  CoordinateMapping
}

// CoordinateMapping tracks coordinate-to-address associations and guest counts.
type CoordinateMapping struct {
	DestinationOccupancy map[coordinates.GuestCoordinates]int
	CoordinateToAddress  map[string]coordinates.GuestCoordinates
	AddressOrder         []string
}

// coordListURL is a semicolon-separated list of all coordinates (used for OSRM API request)
var coordListURL string

// FilterGuestForService removes guests who are not eligible for routing.
// Only guests marked Confirmed or GroceryOnly are kept.
func (e *Event) FilterGuestForService() {
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

// RequestGuestCoordiantes performs geocoding on all filtered guests
func (e *Event) RequestGuestCoordiantes() {
	// Initialize coordinate mapping if empty
	if e.GuestLocations.CoordianteMap.DestinationOccupancy == nil &&
		e.GuestLocations.CoordianteMap.CoordinateToAddress == nil {

		e.GuestLocations.CoordianteMap.DestinationOccupancy = make(map[coordinates.GuestCoordinates]int)
		e.GuestLocations.CoordianteMap.CoordinateToAddress = make(map[string]coordinates.GuestCoordinates)
		e.GuestLocations.CoordianteMap.AddressOrder = make([]string, 0)
	}

	// Always include the depot location as index 0
	depotAddr := "555 Parkdale Ave"
	depotCoor, err := retreiveAddressCoordinate(depotAddr)
	if err != nil {
		log.Println("Could not geocode SMSM address: ", err)
	}
	e.GuestLocations.CoordianteMap.AddressOrder = append(e.GuestLocations.CoordianteMap.AddressOrder, depotAddr)

	depotCoorString := depotCoor.ToString()
	addToCoordListString(&depotCoorString)

	// Geocode all guests and track unique coordinates
	for i := range e.Guests {
		err := e.Guests[i].geocodeGuestAddress()
		if err != nil {
			log.Printf("Warning: Failed to find address coordinates: %s. Reason: %v\n Program will proceed, please manually add Guest: %s",
				e.Guests[i].Address, err, e.Guests[i].Name)
			continue
		}
		coor, unique := e.isUnique(i)
		if unique {
			addToCoordListString(&coor)
		}

	}

}

// addToCoordListString appends a new semicolon-prefixed coordinate string
// to the global OSRM coordinate list.
func addToCoordListString(uniqueCoordinate *string) {
	coordListURL += *uniqueCoordinate

}

// isUnique checks if the given guest's coordinates are already known.
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
	e.GuestLocations.CoordianteMap.CoordinateToAddress[g.Address] = g.Coordinates
	e.GuestLocations.CoordianteMap.AddressOrder = append(e.GuestLocations.CoordianteMap.AddressOrder, g.Address)
	return g.Coordinates.ToString(), true
}
