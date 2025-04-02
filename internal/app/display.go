package app

import (
	"fmt"
)

func (rm *RouteManager) Display(e *Event, lr *LocationRegistry) {
	space(4)
	fmt.Println("Guest Dropoff Summary")
	fmt.Println("============================")
	space(1)
	for i, v := range rm.Vehicles {
		v.DisplayVehicleRoute(i, e, lr)
	}
}

func (v *Vehicle) DisplayVehicleRoute(index int, e *Event, lr *LocationRegistry) {
	nodeVisited := make([]int, 0)
	if v.Route.List != nil {
		for elem := v.Route.List.Front(); elem != nil; elem = elem.Next() {
			nodeVisited = append(nodeVisited, elem.Value.(int))
		}
		addressesVisited := determineAddressesVisited(nodeVisited)
		guests := determineGuestsInvolved(addressesVisited, e, lr)
		fmt.Printf("Vehicle %d: %s ) \n", index, displayGuests(guests))
	}
}

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

func space(lines int) {
	for range lines {
		fmt.Println()
	}
}
