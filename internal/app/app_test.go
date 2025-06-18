package app

import (
	"testing"
)

// TestOrchestateDispatch_RoutingValidation tests the core routing constraints
func TestOrchestateDispatch_RoutingValidation(t *testing.T) {
	// Load real data from JSON
	event, lr, err := loadTestData(t)
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Execute the function
	rm := OrchestateDispatch(lr, event)

	if rm == nil {
		t.Fatal("OrchestateDispatch returned nil RouteManager")
	}

	// 1. Check that all guests were served
	t.Run("AllGuestsServed", func(t *testing.T) {
		// Count total guests in input
		totalInputGuests := 0
		for _, occupancy := range lr.CoordianteMap.DestinationOccupancy {
			totalInputGuests += occupancy
		}

		// Count total guests assigned to vehicles
		totalAssignedGuests := 0
		for _, vehicle := range rm.Vehicles {
			for _, guest := range vehicle.Guests {
				totalAssignedGuests += guest.GroupSize
			}
		}

		if totalAssignedGuests != totalInputGuests {
			t.Errorf("Not all guests were served: input had %d guests, but only %d were assigned to vehicles",
				totalInputGuests, totalAssignedGuests)
		}

		// Check that no destination is marked as unserved if it has guests
		for i, guestCount := range rm.DestinationGuestCount {
			if guestCount > 0 && rm.ServedDestinations[i] == -1 {
				addr := lr.CoordianteMap.AddressOrder[i]
				t.Errorf("Destination %d (%s) has %d guests but is marked as unserved",
					i, addr, guestCount)
			}
		}
	})

	// 2. Check that no two vehicles visit the same location
	t.Run("NoVehicleDuplication", func(t *testing.T) {
		locationToVehicle := make(map[int]int) // maps destination index to vehicle index

		for vehicleIndex, vehicle := range rm.Vehicles {
			if vehicle.Route.List == nil {
				continue // skip vehicles with no routes
			}

			// Walk through the route and check each destination
			for elem := vehicle.Route.List.Front(); elem != nil; elem = elem.Next() {
				destIndex := elem.Value.(int)

				if previousVehicle, exists := locationToVehicle[destIndex]; exists {
					addr := "unknown"
					if destIndex < len(lr.CoordianteMap.AddressOrder) {
						addr = lr.CoordianteMap.AddressOrder[destIndex]
					}
					t.Errorf("Both vehicle %d and vehicle %d visit the same location %d (%s)",
						previousVehicle, vehicleIndex, destIndex, addr)
				} else {
					locationToVehicle[destIndex] = vehicleIndex
				}
			}
		}
	})

	// 3. Check vehicle capacity constraints
	t.Run("VehicleCapacityConstraints", func(t *testing.T) {
		for vehicleIndex, vehicle := range rm.Vehicles {
			totalSeatsUsed := 0
			for _, guest := range vehicle.Guests {
				totalSeatsUsed += guest.GroupSize
			}

			if totalSeatsUsed > maxVehicleSeats {
				t.Errorf("Vehicle %d exceeds capacity: using %d seats but max is %d",
					vehicleIndex, totalSeatsUsed, maxVehicleSeats)
			}

			expectedRemainingSeats := maxVehicleSeats - totalSeatsUsed
			if vehicle.SeatsRemaining != expectedRemainingSeats {
				t.Errorf("Vehicle %d seat calculation error: expected %d remaining seats, got %d",
					vehicleIndex, expectedRemainingSeats, vehicle.SeatsRemaining)
			}
		}
	})

	// 4. Check that ServedDestinations mapping is consistent
	t.Run("ServedDestinationConsistency", func(t *testing.T) {
		// Build a map of which destinations each vehicle actually visits
		actualVehicleDestinations := make(map[int][]int) // vehicle -> list of destinations

		for vehicleIndex, vehicle := range rm.Vehicles {
			if vehicle.Route.List == nil {
				continue
			}

			destinations := []int{}
			for elem := vehicle.Route.List.Front(); elem != nil; elem = elem.Next() {
				destIndex := elem.Value.(int)
				destinations = append(destinations, destIndex)
			}
			actualVehicleDestinations[vehicleIndex] = destinations
		}

		// Check that ServedDestinations matches actual routes
		for destIndex, assignedVehicle := range rm.ServedDestinations {
			if assignedVehicle == -1 {
				continue // unserved destination
			}

			// Check if this vehicle actually visits this destination
			vehicleDestinations, exists := actualVehicleDestinations[assignedVehicle]
			if !exists {
				t.Errorf("Destination %d is assigned to vehicle %d, but vehicle %d has no route",
					destIndex, assignedVehicle, assignedVehicle)
				continue
			}

			found := false
			for _, vDest := range vehicleDestinations {
				if vDest == destIndex {
					found = true
					break
				}
			}

			if !found {
				addr := "unknown"
				if destIndex < len(lr.CoordianteMap.AddressOrder) {
					addr = lr.CoordianteMap.AddressOrder[destIndex]
				}
				t.Errorf("Destination %d (%s) is assigned to vehicle %d, but vehicle %d doesn't visit it",
					destIndex, addr, assignedVehicle, assignedVehicle)
			}
		}
	})

	// 5. Check that guests are assigned to correct vehicles based on their addresses
	t.Run("GuestVehicleAssignment", func(t *testing.T) {
		for vehicleIndex, vehicle := range rm.Vehicles {
			for _, guest := range vehicle.Guests {
				// Find the destination index for this guest's address
				destIndex := -1
				for i, addr := range lr.CoordianteMap.AddressOrder {
					if addr == guest.Address {
						destIndex = i
						break
					}
				}

				if destIndex == -1 {
					t.Errorf("Vehicle %d has guest %s with unknown address %s",
						vehicleIndex, guest.Name, guest.Address)
					continue
				}

				// Check that this destination is assigned to this vehicle
				assignedVehicle := rm.ServedDestinations[destIndex]
				if assignedVehicle != vehicleIndex {
					t.Errorf("Guest %s (address %s, dest %d) is in vehicle %d but destination is assigned to vehicle %d",
						guest.Name, guest.Address, destIndex, vehicleIndex, assignedVehicle)
				}
			}
		}
	})
}

// Helper function to load test data from JSON
func loadTestData(t *testing.T) (*Event, *LocationRegistry, error) {
	// Load from default location
	event, lr, err := LoadAppDataFromFile("test_data.json")
	return &event, &lr, err
}
