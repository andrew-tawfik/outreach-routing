package app

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

func (rm *RouteManager) Display(e *Event, lr *LocationRegistry) string {
	var b strings.Builder

	for i := range rm.Vehicles {
		v := &rm.Vehicles[i] 
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

	
	for _, guest := range v.Guests {
		result.WriteString(v.formatGuestEntry(guest))
	}

	return result.String()
}

func (v *Vehicle) formatGuestEntry(guest Guest) string {
	var entry strings.Builder

	
	guestName := guest.Name
	if guest.GroupSize > 1 {
		guestName = fmt.Sprintf("%s (Group of %d)", guest.Name, guest.GroupSize)
	}

	entry.WriteString(fmt.Sprintf("• %s\n", guestName))
	entry.WriteString(fmt.Sprintf("    ‣ %s\n", guest.Address))

	
	if guest.PhoneNumber != "" {
		entry.WriteString(fmt.Sprintf("    ‣ %s\n", guest.PhoneNumber))
	} else {
		entry.WriteString("  No number\n")
	}

	return entry.String()
}


func (v *Vehicle) UpdateRouteFromGuests(lr *LocationRegistry, eventType string) {
	if len(v.Guests) == 0 {
		v.Route.List = nil
		v.Route.DestinationCount = 0
		v.Locations = make([]coordinates.GuestCoordinates, 0)
		return
	}

	
	v.Route.List = list.New()
	v.Route.DestinationCount = 0
	v.Locations = make([]coordinates.GuestCoordinates, 0)

	grocAdj := 0
	if eventType == "Grocery" {
		grocAdj = 1
	}

	
	seenAddresses := make(map[string]bool)

	
	for _, guest := range v.Guests {
		
		if seenAddresses[guest.Address] {
			continue
		}

		for idx, addr := range lr.CoordianteMap.AddressOrder {
			if addr == guest.Address {
				newIdx := idx - grocAdj
				v.Route.List.PushBack(newIdx)
				v.Route.DestinationCount++

				gc := lr.CoordianteMap.CoordinateToAddress[addr]
				v.Locations = append(v.Locations, gc)

				
				seenAddresses[guest.Address] = true
				break
			}
		}
	}
}