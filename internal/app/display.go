package app

import (
	"container/list"
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
		return fmt.Sprintf("Vehicle %d: No guestsss", index+1)
	}

	return fmt.Sprintf("Vehicle %d: %s", index+1, v.displayGuests())
}

func (v *Vehicle) displayGuests() string {
	if len(v.Guests) == 0 {
		return "No guests"
	}
	var sb strings.Builder
	for i, g := range v.Guests {
		sb.WriteString(fmt.Sprintf("%s (%s)", g.Name, g.Address))
		if i < len(v.Guests)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

// UpdateRouteFromGuests rebuilds the Route.List based on current guests
func (v *Vehicle) UpdateRouteFromGuests(lr *LocationRegistry) {
	if len(v.Guests) == 0 {
		v.Route.List = nil
		v.Route.DestinationCount = 0
		return
	}

	// Create new linked list
	v.Route.List = list.New()
	v.Route.DestinationCount = 0

	// Find the index for each guest's address in the addressOrder
	for _, guest := range v.Guests {
		for idx, addr := range lr.CoordianteMap.AddressOrder {
			if addr == guest.Address {
				v.Route.List.PushBack(idx)
				v.Route.DestinationCount++
				break
			}
		}
	}
}
