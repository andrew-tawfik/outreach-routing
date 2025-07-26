package app

import (
	"container/heap"
	"container/list"
	"fmt"
)

type ClarkeWright struct {
	savingList savingsList 
}

func (cw *ClarkeWright) GetName() string {
	return "Clarke-Wright Savings"
}



func (cw *ClarkeWright) determineSavingList(lr *LocationRegistry) {
	var value float64
	for i := range lr.DistanceMatrix {
		for j := range lr.DistanceMatrix[i] {

			if i == 0 || j == 0 || i == j { 
				continue
			}
			value = lr.retreiveValueFromPair(i, j)
			cw.addToSavingsList(i, j, value)
		}
	}
}



func (lr *LocationRegistry) retreiveValueFromPair(i, j int) float64 {
	depotToI := (lr.DistanceMatrix)[0][i]
	depotToJ := (lr.DistanceMatrix)[0][j]
	iToJ := (lr.DistanceMatrix)[i][j]

	result := depotToI + depotToJ - iToJ

	return result
}


func (cw *ClarkeWright) addToSavingsList(first, second int, value float64) {
	newSaving := saving{
		i:     first,
		j:     second,
		value: value,
	}
	heap.Push(&cw.savingList, newSaving)
}



func (cw *ClarkeWright) StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error {
	cw.InitSavings(lr)

	for cw.savingList.Len() > 0 {
		saving := heap.Pop(&cw.savingList).(saving)

		assignedI := rm.ServedDestinations[saving.i]
		assignedJ := rm.ServedDestinations[saving.j]

		switch {
		
		case assignedI == -1 && assignedJ == -1:
			rm.initiateNewRoute(saving.i, saving.j)

		
		case (assignedI == -1 && rm.DestinationGuestCount[saving.i] == maxVehicleSeats) ||
			(assignedJ == -1 && rm.DestinationGuestCount[saving.j] == maxVehicleSeats):

			if assignedI == -1 {
				rm.initializeSoloRoute(saving.i)
			} else {
				rm.initializeSoloRoute(saving.j)
			}

		
		case (assignedI == -1 && assignedJ != -1) || (assignedI != -1 && assignedJ == -1):
			vehicleIndex, err := rm.determineVehicle(1, saving.i, saving.j)
			if err != nil {
				return fmt.Errorf("failed to process: %w", err)
			}
			rm.canAttachToRoute(vehicleIndex, saving.i, saving.j)
		}
	}
	return nil 
}


func (rm *RouteManager) initializeSoloRoute(location1 int) {
	
	rm.AddNewVehicle()
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


func (rm *RouteManager) initiateNewRoute(location1, location2 int) {
	rm.AddNewVehicle()
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
