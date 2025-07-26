package app

import "fmt"


func (r *Route) isExternal(locationIndex int) (string, bool) {
	if r.List == nil || r.List.Len() == 0 {
		return "", false
	}

	back := r.List.Back()
	front := r.List.Front()

	if back != nil && locationIndex == back.Value.(int) {
		return "back", true
	} else if front != nil && locationIndex == front.Value.(int) {
		return "front", true
	}

	return "", false
}



func (r *Route) extendRoute(existingLocation, newLocation int) bool {
	orientation, ok := r.isExternal(existingLocation)
	if !ok {
		return false
	}

	if orientation == "back" {
		r.List.PushBack(newLocation)
		return true
	} else {
		r.List.PushFront(newLocation)
		return true
	}
}



func (rm *RouteManager) getRouteExtensionEndpoints(vehicleIndex, location1, location2 int) (existingLocation, newLocation int) {
	if rm.ServedDestinations[location1] == vehicleIndex {
		return location1, location2
	}
	return location2, location1
}






func (rm *RouteManager) determineVehicle(action, location1, location2 int) (int, error) {
	switch action {

	case 0:
		
		for i, v := range rm.Vehicles {
			if rm.enoughSeatsToInitialize(i, location1, location2) && v.Route.List == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("no new vehicles available to initiate")

	case 1:
		
		assignedI := rm.ServedDestinations[location1]
		assignedJ := rm.ServedDestinations[location2]
		if assignedI != -1 {
			return assignedI, nil
		}
		return assignedJ, nil

	case 2:
		
		for i, v := range rm.Vehicles {
			if rm.enoughSeatsToInitializeSolo(i, location1) && v.Route.List == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("no new vehicles available to initiate")
	}

	return -1, fmt.Errorf("unfinished function — invalid action code")
}



func (rm *RouteManager) enoughSeatsToInitialize(vehicleIndex, locationI, locationJ int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtI := rm.DestinationGuestCount[locationI]
	guestsAtJ := rm.DestinationGuestCount[locationJ]

	underThreeStops := v.Route.DestinationCount < 3

	return v.SeatsRemaining >= guestsAtI+guestsAtJ && underThreeStops
}



func (rm *RouteManager) enoughSeatsToInitializeSolo(vehicleIndex, locationI int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtI := rm.DestinationGuestCount[locationI]

	underThreeStops := v.Route.DestinationCount < 3
	return v.SeatsRemaining == guestsAtI && underThreeStops
}



func (rm *RouteManager) enoughSeatsToExtend(vehicleIndex, newLocation int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtLocation := rm.DestinationGuestCount[newLocation]

	underThreeStops := v.Route.DestinationCount < 3

	return v.SeatsRemaining >= guestsAtLocation && underThreeStops
}

func (v *Vehicle) findGuests(addresses []string, e *Event, lr *LocationRegistry) {
	var guestsInvolved []Guest
	for _, addr := range addresses {
		coord := lr.CoordianteMap.CoordinateToAddress[addr]
		for _, g := range e.Guests {
			if g.Coordinates.Long == coord.Long && g.Coordinates.Lat == coord.Lat {
				guestsInvolved = append(guestsInvolved, g)
			}
		}
	}
	v.Guests = guestsInvolved
}

func determineAddressesVisited(nodeVisited []int, e *Event) []string {
	adj := 0
	if e.EventType == "Grocery" {
		adj++
	}
	result := make([]string, 0, len(nodeVisited))
	for _, idx := range nodeVisited {
		result = append(result, addressOrder[idx+adj])
	}
	return result
}


func (rm *RouteManager) AddNewVehicle() {
	newVehicle := Vehicle{SeatsRemaining: maxVehicleSeats}
	rm.Vehicles = append(rm.Vehicles, newVehicle)
}
