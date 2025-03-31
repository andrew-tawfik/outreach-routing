package app

import (
	"container/heap"
	"container/list"
	"log"
)

func (rm *RouteManager) DetermineSavingList(lr *LocationRegistry) {
	var value float64
	for i := range lr.DistanceMatrix {
		for j := range lr.DistanceMatrix[i] {
			if i == 0 || j == 0 || i == j {
				continue
			}
			value = lr.retreiveValueFromPair(i, j)
			rm.addToSavingsList(i, j, value)
		}
	}
}

func (lr *LocationRegistry) retreiveValueFromPair(i, j int) float64 {

	// Clarke-Wright Algorithm Formula: d(D, i) + d(D, j) - d(i, j)
	depotToI := (lr.DistanceMatrix)[0][i]
	depotToJ := (lr.DistanceMatrix)[0][j]
	iToJ := (lr.DistanceMatrix)[i][j]

	result := depotToI + depotToJ - iToJ

	return result
}

func (rm *RouteManager) addToSavingsList(first, second int, value float64) {
	newSaving := saving{
		i:     first,
		j:     second,
		value: value,
	}
	heap.Push(&rm.SavingList, newSaving)
}

func (rm *RouteManager) TestRemoveAll() {
	rm.SavingList.popAll()
}

func (rm *RouteManager) StartRouteDispatch() {

	for rm.SavingList.Len() > 0 {
		saving := heap.Pop(&rm.SavingList).(saving)

		assignedI := rm.ServedDestinations[saving.i]
		assignedJ := rm.ServedDestinations[saving.j]

		// a. neither i nor j have already been assigned to a route
		if assignedI == -1 && assignedJ == -1 {
			rm.initiateNewRoute(saving.i, saving.j)

		} else if (assignedI == -1 && rm.DestinationGuestCount[saving.i] == maxVehicleSeats) ||
			(assignedJ == -1 && rm.DestinationGuestCount[saving.j] == maxVehicleSeats) {

			if assignedI == -1 {
				rm.initializeSoloRoute(saving.i)
			} else {
				rm.initializeSoloRoute(saving.j)
			}

		} else if (assignedI == -1 && assignedJ != -1) || (assignedI != -1 && assignedJ == -1) {

			// Put a function that might add a new location to a vehicle route
			vehicleIndex, err := rm.determineVehicle(1, saving.i, saving.j)
			if err != nil {
				log.Fatalf("Cannot determine vehicle route: %v", err)
			}

			rm.canAttachToRoute(vehicleIndex, saving.i, saving.j)
		}

	}
}

func (rm *RouteManager) initializeSoloRoute(location1 int) {
	vehicleToStart, err := rm.determineVehicle(2, location1, -1)
	if err != nil {
		return
	}
	v := &rm.Vehicles[vehicleToStart]
	if vehicleToStart == -1 {
		return
	}
	v.Route.List = list.New()
	v.Route.List.PushBack(location1)

	// Add location 1 and 2
	rm.update(vehicleToStart, location1)
}

func (rm *RouteManager) canAttachToRoute(vehicleIndex, location1, location2 int) {

	//
	route := &rm.Vehicles[vehicleIndex].Route

	existingLocation, newLocation := rm.getRouteExtensionEndpoints(vehicleIndex, location1, location2)
	enoughSeats := rm.enoughSeatsToExtend(vehicleIndex, newLocation)

	if enoughSeats {
		if route.extendRoute(existingLocation, newLocation) {
			rm.update(vehicleIndex, newLocation)
		}

	}

}

func (rm *RouteManager) initiateNewRoute(location1, location2 int) {
	// Initilize the route
	vehicleToStart, err := rm.determineVehicle(0, location1, location2)
	if err != nil {
		return
	}
	v := &rm.Vehicles[vehicleToStart]
	if vehicleToStart == -1 {
		return
	}
	v.Route.List = list.New()
	v.Route.List.PushBack(location1)
	v.Route.List.PushBack(location2)

	// Add location 1 and 2
	rm.update(vehicleToStart, location1)
	rm.update(vehicleToStart, location2)
}

func (rm *RouteManager) update(vehicleIndex, location int) {
	v := &rm.Vehicles[vehicleIndex]

	v.SeatsRemaining -= rm.DestinationGuestCount[location]
	rm.ServedDestinations[location] = vehicleIndex
}
