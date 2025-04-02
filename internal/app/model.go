package app

import (
	"container/heap"
	"container/list"
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
}

// RouteManager manages the state of all vehicles, routing decisions, and guest assignments.
type RouteManager struct {
	Vehicles              []Vehicle   // List of available vehicles
	ServedDestinations    map[int]int // Maps destination index to assigned vehicle index
	DestinationGuestCount []int       // Number of guests at each destination index
	SavingList            savingsList // Heap of savings values for Clarke-Wright algorithm
}

// maxVehicleSeats defines the seat capacity of each vehicle
var maxVehicleSeats int = 4

// addressOrder is a temporary global variable used for displaying purposes
var addressOrder []string

// CreateRouteManager initializes a RouteManager instance from a LocationRegistry and
// number of available vehicles. It sets up guest counts, initializes vehicles, and prepares the savings heap.
func CreateRouteManager(lr *LocationRegistry, numVehicles int) *RouteManager {

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
	vehicles := createVehicles(numVehicles)
	savings := &savingsList{}
	heap.Init(savings)

	addressOrder = lr.CoordianteMap.AddressOrder

	return &RouteManager{
		Vehicles:              vehicles,
		ServedDestinations:    servedDestinations,
		DestinationGuestCount: destinationGuestCount,
		SavingList:            *savings,
	}
}

// createVehicles returns a slice of Vehicles, each initialized with max seats and an empty route.
func createVehicles(numVehicles int) []Vehicle {
	v := Vehicle{SeatsRemaining: maxVehicleSeats}
	vehicles := make([]Vehicle, 0, numVehicles)
	i := 0
	for i < numVehicles {
		vehicles = append(vehicles, v)
		i++
	}
	return vehicles
}
