package app

import "fmt"

// isExternal checks if a location is at the front or back of the current route.
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

// extendRoute attempts to extend the route by adding a new location
// to either the front or back, depending on where the existing location is.
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

// getRouteExtensionEndpoints returns the (existing, new) location order
// based on which one is already assigned to the given vehicle.
func (rm *RouteManager) getRouteExtensionEndpoints(vehicleIndex, location1, location2 int) (existingLocation, newLocation int) {
	if rm.ServedDestinations[location1] == vehicleIndex {
		return location1, location2
	}
	return location2, location1
}

// determineVehicle selects an appropriate vehicle index for a given action:
//
// action = 0: initiate a new route (needs a fresh vehicle with enough seats)
// action = 1: extend an existing route (choose the already-assigned vehicle)
// action = 2: create a solo route (like action 0 but for one guest only)
func (rm *RouteManager) determineVehicle(action, location1, location2 int) (int, error) {
	switch action {

	case 0:
		// Find an unused vehicle with enough seats for both guests
		for i, v := range rm.Vehicles {
			if rm.enoughSeatsToInitialize(i, location1, location2) && v.Route.List == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("no new vehicles available to initiate")

	case 1:
		// Determine which location is already assigned and return that vehicle
		assignedI := rm.ServedDestinations[location1]
		assignedJ := rm.ServedDestinations[location2]
		if assignedI != -1 {
			return assignedI, nil
		}
		return assignedJ, nil

	case 2:
		// Find an unused vehicle with exactly the number of seats needed for this guest
		for i, v := range rm.Vehicles {
			if rm.enoughSeatsToInitializeSolo(i, location1) && v.Route.List == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("no new vehicles available to initiate")
	}

	return -1, fmt.Errorf("unfinished function — invalid action code")
}

// enoughSeatsToInitialize checks if the vehicle has enough seats
// for both guest locations to be added in one new route.
func (rm *RouteManager) enoughSeatsToInitialize(vehicleIndex, locationI, locationJ int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtI := rm.DestinationGuestCount[locationI]
	guestsAtJ := rm.DestinationGuestCount[locationJ]

	underThreeStops := v.Route.DestinationCount < 3

	return v.SeatsRemaining >= guestsAtI+guestsAtJ && underThreeStops
}

// enoughSeatsToInitializeSolo checks if a vehicle has exactly the right
// number of seats remaining to accommodate a single location.
func (rm *RouteManager) enoughSeatsToInitializeSolo(vehicleIndex, locationI int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtI := rm.DestinationGuestCount[locationI]

	underThreeStops := v.Route.DestinationCount < 3
	return v.SeatsRemaining == guestsAtI && underThreeStops
}

// enoughSeatsToExtend checks if the vehicle has enough remaining seats
// to add the new location to an existing route.
func (rm *RouteManager) enoughSeatsToExtend(vehicleIndex, newLocation int) bool {
	v := &rm.Vehicles[vehicleIndex]
	guestsAtLocation := rm.DestinationGuestCount[newLocation]

	underThreeStops := v.Route.DestinationCount < 3

	return v.SeatsRemaining >= guestsAtLocation && underThreeStops
}

func (rm *RouteManager) determineGuestsInvolved(e *Event, lr *LocationRegistry) {

	for i := range rm.Vehicles {
		v := &rm.Vehicles[i]

		var nodeVisited []int
		for elem := v.Route.List.Front(); elem != nil; elem = elem.Next() {
			nodeVisited = append(nodeVisited, elem.Value.(int))
		}

		addresses := determineAddressesVisited(nodeVisited)
		v.findGuests(addresses, e, lr)
	}

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

func determineAddressesVisited(nodeVisited []int) []string {
	result := make([]string, 0, len(nodeVisited))
	for _, idx := range nodeVisited {
		result = append(result, addressOrder[idx])
	}
	return result
}
