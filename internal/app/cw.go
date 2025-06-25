package app

import (
	"container/heap"
	"container/list"
	"fmt"
)

type ClarkeWright struct {
	savingList savingsList // Heap of savings values for Clarke-Wright algorithm
}

func (cw *ClarkeWright) GetName() string {
	return "Clarke-Wright Savings"
}

// DetermineSavingList computes the "savings" between all possible guest pairs
// and populates the priority queue (heap) used to decide which routes to build first.
func (cw *ClarkeWright) determineSavingList(lr *LocationRegistry) {
	var value float64
	for i := range lr.DistanceMatrix {
		for j := range lr.DistanceMatrix[i] {

			if i == 0 || j == 0 || i == j { // Do not calculate with depot location
				continue
			}
			value = lr.retreiveValueFromPair(i, j)
			cw.addToSavingsList(i, j, value)
		}
	}
}

// retreiveValueFromPair calculates the savings value between two locations i and j
// using the Clarke-Wright formula: d(depot,i) + d(depot,j) - d(i,j)
func (lr *LocationRegistry) retreiveValueFromPair(i, j int) float64 {
	depotToI := (lr.DistanceMatrix)[0][i]
	depotToJ := (lr.DistanceMatrix)[0][j]
	iToJ := (lr.DistanceMatrix)[i][j]

	result := depotToI + depotToJ - iToJ

	return result
}

// addToSavingsList pushes a new savings record onto the heap (priority queue)
func (cw *ClarkeWright) addToSavingsList(first, second int, value float64) {
	newSaving := saving{
		i:     first,
		j:     second,
		value: value,
	}
	heap.Push(&cw.savingList, newSaving)
}

// StartRouteDispatch performs the Clarke-Wright dispatch loop,
// consuming the savings list in descending order and assigning locations to vehicles.
func (cw *ClarkeWright) StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error {
	cw.InitSavings(lr)

	for cw.savingList.Len() > 0 {
		saving := heap.Pop(&cw.savingList).(saving)

		assignedI := rm.ServedDestinations[saving.i]
		assignedJ := rm.ServedDestinations[saving.j]

		switch {
		// Case 1: neither location is assigned to a route yet
		case assignedI == -1 && assignedJ == -1:
			rm.initiateNewRoute(saving.i, saving.j)

		// Case 2: one is unassigned, but has too many passengers to be added to shared ride
		case (assignedI == -1 && rm.DestinationGuestCount[saving.i] == maxVehicleSeats) ||
			(assignedJ == -1 && rm.DestinationGuestCount[saving.j] == maxVehicleSeats):

			if assignedI == -1 {
				rm.initializeSoloRoute(saving.i)
			} else {
				rm.initializeSoloRoute(saving.j)
			}

		// Case 3: one location is already assigned; try extending the route
		case (assignedI == -1 && assignedJ != -1) || (assignedI != -1 && assignedJ == -1):
			vehicleIndex, err := rm.determineVehicle(1, saving.i, saving.j)
			if err != nil {
				return fmt.Errorf("failed to process: %w", err)
			}
			rm.canAttachToRoute(vehicleIndex, saving.i, saving.j)
		}
	}
	return nil // TODO error handling
}

// initializeSoloRoute assigns a single location to a new vehicle route
func (rm *RouteManager) initializeSoloRoute(location1 int) {
	// append new vehicle to slice
	rm.addNewVehicle()
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

	rm.update(vehicleToStart, location1)
}

// canAttachToRoute attempts to extend an existing route by adding a new location
// at either the beginning or end, depending on feasibility and seat availability.
func (rm *RouteManager) canAttachToRoute(vehicleIndex, location1, location2 int) {

	route := &rm.Vehicles[vehicleIndex].Route

	existingLocation, newLocation := rm.getRouteExtensionEndpoints(vehicleIndex, location1, location2)
	enoughSeats := rm.enoughSeatsToExtend(vehicleIndex, newLocation)

	if enoughSeats {
		if route.extendRoute(existingLocation, newLocation) {
			rm.update(vehicleIndex, newLocation)
		}

	}

}

// initiateNewRoute creates a new vehicle route between two unassigned locations
func (rm *RouteManager) initiateNewRoute(location1, location2 int) {
	rm.addNewVehicle()
	vehicleToStart, err := rm.determineVehicle(0, location1, location2)
	if err != nil || vehicleToStart == -1 {
		return
	}

	v := &rm.Vehicles[vehicleToStart]
	v.Route.List = list.New()
	v.Route.List.PushBack(location1)
	v.Route.List.PushBack(location2)

	rm.update(vehicleToStart, location1)
	rm.update(vehicleToStart, location2)
}

// update adjusts bookkeeping after assigning a location to a vehicle
func (rm *RouteManager) update(vehicleIndex, location int) {
	v := &rm.Vehicles[vehicleIndex]

	v.Route.DestinationCount++
	v.SeatsRemaining -= rm.DestinationGuestCount[location]
	rm.ServedDestinations[location] = vehicleIndex
}

func (cw *ClarkeWright) InitSavings(lr *LocationRegistry) {
	savings := &savingsList{}
	heap.Init(savings)
	cw.savingList = *savings
	cw.determineSavingList(lr)
}
