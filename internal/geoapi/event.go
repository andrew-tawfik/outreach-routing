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
	Address     string
	Coordinates GuestCoordinates
}

type Event struct {
	Guests            []Guest
	EventType         string
	CoordinatesString string
}

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
		e.AddToCoordinatesList(i)
		if err != nil {
			return err
		}
	}
	fmt.Println("Retrived all coordinates successfully")
	return nil
}

func (e *Event) AddToCoordinatesList(guestIndex int) {
	g := e.Guests[guestIndex]
	str := fmt.Sprintf("%f,%f;", g.Coordinates.Long, g.Coordinates.Lat)
	e.CoordinatesString += str
}
