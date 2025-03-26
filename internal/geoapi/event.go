package geoapi

import "fmt"

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
	Coordinates GuestCoordinates
}

type Event struct {
	Guests            []Guest
	EventType         string
	CoordinatesString string
}

var eventCoordinateSet = make(map[GuestCoordinates]bool)

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
	for i := range e.Guests {
		err := e.Guests[i].geocodeAddress()
		if err != nil {
			return err
		}
		coor, unique := e.isUnique(i)
		if unique {
			e.AddToSet(&coor)
		}

	}
	fmt.Println("Retrived all coordinates successfully")
	return nil
}

func (e *Event) AddToSet(uniqueCoordinate *string) {
	e.CoordinatesString += *uniqueCoordinate
}

func (e *Event) isUnique(guestIndex int) (string, bool) {
	g := e.Guests[guestIndex]
	if eventCoordinateSet[g.Coordinates] {
		return "", false
	} else {
		eventCoordinateSet[g.Coordinates] = true
		return fmt.Sprintf("%f,%f;", g.Coordinates.Long, g.Coordinates.Lat), true
	}
}
