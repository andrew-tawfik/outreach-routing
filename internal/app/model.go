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
	ServedDestinations    map[int]bool
	DestinationGuestCount []int
	SavingList            savingsList
}

func CreateRouteManager(lr *LocationRegistry, numVehicles int) *RouteManager {

	ao := &lr.CoordianteMap.AddressOrder
	destinationCount := &lr.CoordianteMap.DestinationOccupancy
	addrMap := &lr.CoordianteMap.CoordinateToAddress

	destinationGuestCount := make([]int, len(*ao))
	servedDestinations := make(map[int]bool)

	i := 0

	for _, addr := range *ao {
		guestCount := (*destinationCount)[(*addrMap)[addr]]
		destinationGuestCount[i] = guestCount
		servedDestinations[i] = false
		i++
	}
	vehicles := createVehicles(numVehicles)
	savings := &savingsList{}
	heap.Init(savings)

	return &RouteManager{
		Vehicles:              vehicles,
		ServedDestinations:    servedDestinations,
		DestinationGuestCount: destinationGuestCount,
		SavingList:            *savings,
	}
}

func createVehicles(numVehicles int) []Vehicle {
	v := Vehicle{SeatsRemaining: 5}
	vehicles := make([]Vehicle, 0, numVehicles)
	i := 0
	for i < numVehicles {
		vehicles = append(vehicles, v)
		i++
	}
	return vehicles
}
