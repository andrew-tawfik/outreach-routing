package app

import (
	"container/heap"
	"container/list"
)

type Route struct {
	List *list.List
}

type Vehicle struct {
	SeatsRemaining int
	Route          Route
}

type RouteManager struct {
	Vehicles              []Vehicle
	ServedDestinations    map[int]int
	DestinationGuestCount []int
	SavingList            savingsList
}

var maxVehicleSeats int = 4

// Delete later for debugging only
var addressOrder []string

func CreateRouteManager(lr *LocationRegistry, numVehicles int) *RouteManager {

	ao := &lr.CoordianteMap.AddressOrder
	destinationCount := &lr.CoordianteMap.DestinationOccupancy
	addrMap := &lr.CoordianteMap.CoordinateToAddress

	destinationGuestCount := make([]int, len(*ao))
	servedDestinations := make(map[int]int)

	i := 0

	for _, addr := range *ao {
		guestCount := (*destinationCount)[(*addrMap)[addr]]
		destinationGuestCount[i] = guestCount
		servedDestinations[i] = -1
		i++
	}
	vehicles := createVehicles(numVehicles)
	savings := &savingsList{}
	heap.Init(savings)

	//Delete Later
	addressOrder = lr.CoordianteMap.AddressOrder

	return &RouteManager{
		Vehicles:              vehicles,
		ServedDestinations:    servedDestinations,
		DestinationGuestCount: destinationGuestCount,
		SavingList:            *savings,
	}
}

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
