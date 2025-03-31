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

	if action == 0 {
		for i, v := range rm.Vehicles {
			if rm.enoughSeatsToInitialize(i, location1, location2) && v.Route.List == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("no new vehicles available to initiate")

	} else if action == 1 {
		assignedI := rm.ServedDestinations[location1]
		assignedJ := rm.ServedDestinations[location2]
		if assignedI != -1 {
			return assignedI, nil
		} else {
			return assignedJ, nil
		}

	} else if action == 2 {
		for i, v := range rm.Vehicles {
			if rm.enoughSeatsToInitializeSolo(i, location1) && v.Route.List == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("no new vehicles available to initiate")

	}

	return -1, fmt.Errorf("unfinished function change later")
}

func (rm *RouteManager) enoughSeatsToInitialize(vehicleIndex, locationI, locationJ int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtI := rm.DestinationGuestCount[locationI]
	guestsAtJ := rm.DestinationGuestCount[locationJ]

	return v.SeatsRemaining >= guestsAtI+guestsAtJ
}

func (rm *RouteManager) enoughSeatsToInitializeSolo(vehicleIndex, locationI int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtI := rm.DestinationGuestCount[locationI]
	return v.SeatsRemaining == guestsAtI
}

func (rm *RouteManager) enoughSeatsToExtend(vehicleIndex, newLocation int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtLocation := rm.DestinationGuestCount[newLocation]

	return v.SeatsRemaining >= guestsAtLocation
}
