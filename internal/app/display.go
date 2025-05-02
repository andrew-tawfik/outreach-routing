package app

import (
	"fmt"
	"strings"
)

func (rm *RouteManager) Display(e *Event, lr *LocationRegistry) string {
	var b strings.Builder

	b.WriteString("Guest Dropoff Summary\n")
	b.WriteString("============================\n\n")

	for i := range rm.Vehicles {
		v := &rm.Vehicles[i] // take pointer into slice
		vehicleInfo := v.GetVehicleRouteInfo(i, e, lr)
		b.WriteString(vehicleInfo)
		b.WriteString("\n")
	}

	return b.String()
}

func (v *Vehicle) GetVehicleRouteInfo(index int, e *Event, lr *LocationRegistry) string {
	if v.Route.List == nil {
		return fmt.Sprintf("Vehicle %d: No guests", index+1)
	}

	// build the list of node IDs this vehicle visits
	var nodeVisited []int
	for elem := v.Route.List.Front(); elem != nil; elem = elem.Next() {
		nodeVisited = append(nodeVisited, elem.Value.(int))
	}

	addresses := determineAddressesVisited(nodeVisited)

	// this mutates v.Guests on the real slice element
	v.determineGuestsInvolved(addresses, e, lr)

	return fmt.Sprintf("Vehicle %d: %s", index+1, displayGuests(v.Guests))
}

func (v *Vehicle) determineGuestsInvolved(addresses []string, e *Event, lr *LocationRegistry) {
	var guestsInvolved []Guest
	for _, addr := range addresses {
		coord := lr.CoordianteMap.CoordinateToAddress[addr] // fixed typo
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

func displayGuests(guests []Guest) string {
	if len(guests) == 0 {
		return "No guests"
	}
	var sb strings.Builder
	for i, g := range guests {
		sb.WriteString(fmt.Sprintf("%s (%s)", g.Name, g.Address))
		if i < len(guests)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}
