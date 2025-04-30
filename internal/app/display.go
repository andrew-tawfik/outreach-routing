package app

import (
	"fmt"
	"strings"
)

func (rm *RouteManager) Display(e *Event, lr *LocationRegistry) string {
	var b strings.Builder

	b.WriteString("Guest Dropoff Summary\n")
	b.WriteString("============================\n\n")

	for i, v := range rm.Vehicles {
		vehicleInfo := v.GetVehicleRouteInfo(i, e, lr)
		b.WriteString(vehicleInfo)
		b.WriteString("\n")
	}

	return b.String()
}

func (v *Vehicle) GetVehicleRouteInfo(index int, e *Event, lr *LocationRegistry) string {
	nodeVisited := make([]int, 0)
	if v.Route.List != nil {
		for elem := v.Route.List.Front(); elem != nil; elem = elem.Next() {
			nodeVisited = append(nodeVisited, elem.Value.(int))
		}
		addressesVisited := determineAddressesVisited(nodeVisited)
		guests := determineGuestsInvolved(addressesVisited, e, lr)
		return fmt.Sprintf("Vehicle %d: %s", index+1, displayGuests(guests))
	}
	return fmt.Sprintf("Vehicle %d: No guests", index+1)
}

// Keeping the existing helper functions unchanged
func determineGuestsInvolved(addressesVisited []string, e *Event, lr *LocationRegistry) []Guest {
	guestsInvolved := make([]Guest, 0)
	for _, v := range addressesVisited {
		coor := lr.CoordianteMap.CoordinateToAddress[v]
		for _, g := range e.Guests {
			if g.Coordinates.Long == coor.Long && g.Coordinates.Lat == coor.Lat {
				guestsInvolved = append(guestsInvolved, g)
			}
		}
	}
	return guestsInvolved
}

func determineAddressesVisited(nodeVisited []int) []string {
	addressesVisited := make([]string, 0)
	for _, val := range nodeVisited {
		addressesVisited = append(addressesVisited, addressOrder[val])
	}
	return addressesVisited
}

func displayGuests(guests []Guest) string {
	if len(guests) == 0 {
		return "No guests"
	}
	str := ""
	for i, g := range guests {
		str += fmt.Sprintf("%s (%s)", g.Name, g.Address)
		if i != len(guests)-1 {
			str += ", "
		}
	}
	return str
}
