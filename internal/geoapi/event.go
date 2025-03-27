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
	Guests         []Guest
	EventType      string
	GuestLocations LocationRegistry
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
	if e.GuestLocations.GuestCountByCoord == nil {
		e.GuestLocations.GuestCountByCoord = make(map[GuestCoordinates]int)
	}

	depotCoor, err := retreiveAddress("555 Parkdale Ave")
	if err != nil {
		return err
	}

	depotCoorString := depotCoor.toString()
	e.AddToCoordListString(&depotCoorString)

	for i := range e.Guests {
		err := e.Guests[i].geocodeGuestAddress()
		if err != nil {
			return err
		}
		coor, unique := e.isUnique(i)
		if unique {
			e.AddToCoordListString(&coor)
		}

	}
	fmt.Println("Retrived all coordinates successfully")
	return nil
}

func (e *Event) AddToCoordListString(uniqueCoordinate *string) {
	e.GuestLocations.CoordListString += *uniqueCoordinate
}

func (e *Event) isUnique(guestIndex int) (string, bool) {
	g := e.Guests[guestIndex]
	val, ok := e.GuestLocations.GuestCountByCoord[g.Coordinates]

	if ok {
		// Already exists, update and return false
		e.GuestLocations.GuestCountByCoord[g.Coordinates] = val + g.GroupSize
		return "", false
	}

	// First time seeing this coordinate
	e.GuestLocations.GuestCountByCoord[g.Coordinates] = g.GroupSize
	return g.Coordinates.toString(), true
}

func (gc *GuestCoordinates) toString() string {
	return fmt.Sprintf("%f,%f;", gc.Long, gc.Lat)

}
