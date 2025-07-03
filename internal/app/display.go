package app

import (
	"container/list"
	"fmt"
	"strings"
)

func (rm *RouteManager) Display(e *Event, lr *LocationRegistry) string {
	var b strings.Builder

	for i := range rm.Vehicles {
		v := &rm.Vehicles[i] // take pointer into slice
		vehicleInfo := v.GetVehicleRouteInfo(i, e, lr)
		b.WriteString(vehicleInfo)
		b.WriteString("\n")
	}
	return b.String()
}

func (v *Vehicle) GetVehicleRouteInfo(index int, e *Event, lr *LocationRegistry) string {
	if v.Route.List == nil || len(v.Guests) == 0 {
		return fmt.Sprintf("Driver %d: No guests assigned", index+1)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Driver %d:\n", index+1))

	// Format each guest with bullet points
	for _, guest := range v.Guests {
		result.WriteString(v.formatGuestEntry(guest))
	}

	return result.String()
}

func (v *Vehicle) formatGuestEntry(guest Guest) string {
	var entry strings.Builder

	// Guest name with group size if > 1
	guestName := guest.Name
	if guest.GroupSize > 1 {
		guestName = fmt.Sprintf("%s (Group of %d)", guest.Name, guest.GroupSize)
	}

	entry.WriteString(fmt.Sprintf("• %s\n", guestName))
	entry.WriteString(fmt.Sprintf("    ‣ %s\n", guest.Address))

	// Phone number - handle empty numbers
	if guest.PhoneNumber != "" {
		entry.WriteString(fmt.Sprintf("    ‣ %s\n", guest.PhoneNumber))
	} else {
		entry.WriteString("  No number\n")
	}

	return entry.String()
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
