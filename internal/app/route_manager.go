package app

import (
	"container/list"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

// Route holds the doubly-linked list of destination indices assigned to a vehicle.
type Route struct {
	List             *list.List
	DestinationCount int
}

// Vehicle represents a single transport unit with a route and remaining seat capacity.
type Vehicle struct {
	SeatsRemaining int
	Route          Route
	Guests         []Guest
	Locations      []coordinates.GuestCoordinates
}

// RouteManager manages the state of all vehicles, routing decisions, and guest assignments.
type RouteManager struct {
	Vehicles              []Vehicle   // List of available vehicles
	ServedDestinations    map[int]int // Maps destination index to assigned vehicle index
	DestinationGuestCount []int       // Number of guests at each destination index
	CoordinateList        []coordinates.GuestCoordinates
}

// maxVehicleSeats defines the seat capacity of each vehicle
var maxVehicleSeats int = 4

// addressOrder is a temporary global variable used for displaying purposes
var addressOrder []string

type VRPAlgorithm interface {
	StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error
	GetName() string
}

// CreateRouteManager initializes a RouteManager instance from a LocationRegistry and
// number of available vehicles. It sets up guest counts, initializes vehicles, and prepares the savings heap.
func OrchestateDispatch(lr *LocationRegistry, e *Event) *RouteManager {

	ao := &lr.CoordianteMap.AddressOrder
	destinationCount := &lr.CoordianteMap.DestinationOccupancy
	addrMap := &lr.CoordianteMap.CoordinateToAddress

	// Convert address-based guest data to index-based arrays
	destinationGuestCount := make([]int, len(*ao))
	servedDestinations := make(map[int]int)

	for i, addr := range *ao {
		guestCount := (*destinationCount)[(*addrMap)[addr]]
		destinationGuestCount[i] = guestCount
		servedDestinations[i] = -1
	}

	// Create vehicles and initialize the savings heap
	vehicles := make([]Vehicle, 0, 10)

	addressOrder = lr.CoordianteMap.AddressOrder

	rm := &RouteManager{
		Vehicles:              vehicles,
		ServedDestinations:    servedDestinations,
		DestinationGuestCount: destinationGuestCount,
	}
	rm.createCoordinateList(lr)

	var strategy VRPAlgorithm
	if e.EventType == "Dinner" {
		strategy = &ClarkeWright{}
	} else {
		strategy = &Kmeans{}
	}
	strategy.StartRouteDispatch(rm, lr)

	rm.determineGuestsInvolved(e, lr)
	return rm
}

func (rm *RouteManager) determineGuestsInvolved(e *Event, lr *LocationRegistry) {

	for i := range rm.Vehicles {
		v := &rm.Vehicles[i]

		var nodeVisited []int
		for elem := v.Route.List.Front(); elem != nil; elem = elem.Next() {
			nodeVisited = append(nodeVisited, elem.Value.(int))
		}

		addresses := determineAddressesVisited(nodeVisited, e)
		v.determineCoordinates(addresses, lr)
		v.findGuests(addresses, e, lr)
	}

}

func (rm *RouteManager) createCoordinateList(lr *LocationRegistry) {
	ao := lr.CoordianteMap.AddressOrder
	coorList := make([]coordinates.GuestCoordinates, 0, len(ao))

	for i := 1; i < len(ao); i++ {
		s := ao[i]
		coor := lr.CoordianteMap.CoordinateToAddress[s]
		coorList = append(coorList, coor)
	}
	rm.CoordinateList = coorList

}

func (v *Vehicle) determineCoordinates(addresses []string, lr *LocationRegistry) {
	vehicleCoordinates := make([]coordinates.GuestCoordinates, 0)
	for _, a := range addresses {
		gc := lr.CoordianteMap.CoordinateToAddress[a]
		vehicleCoordinates = append(vehicleCoordinates, gc)
	}
	v.Locations = vehicleCoordinates
}
