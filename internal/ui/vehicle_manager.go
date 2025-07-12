package ui

import (
	"container/list"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

// VehicleManager handles vehicle state management and operations
type VehicleManager struct {
	routeManager *app.RouteManager
	grid         *VehicleGrid
	config       *Config

	// State tracking - now includes routes
	initialGuestState    map[int][]app.Guest // Original guest assignments
	initialRouteState    map[int]app.Route   // Original route assignments
	initialLocationState map[int][]coordinates.GuestCoordinates
	hasChanges           bool
}

// NewVehicleManager creates a new vehicle manager
func NewVehicleManager(rm *app.RouteManager, grid *VehicleGrid, cfg *Config) *VehicleManager {
	vm := &VehicleManager{
		routeManager:         rm,
		grid:                 grid,
		initialGuestState:    make(map[int][]app.Guest),
		initialRouteState:    make(map[int]app.Route),
		initialLocationState: make(map[int][]coordinates.GuestCoordinates),
		config:               cfg,
	}

	vm.captureInitialState()
	return vm
}

// captureInitialState saves the current vehicle assignments AND routes for reset functionality
func (vm *VehicleManager) captureInitialState() {
	for i, vehicle := range vm.routeManager.Vehicles {
		// Deep copy the guests slice
		guestsCopy := make([]app.Guest, len(vehicle.Guests))
		copy(guestsCopy, vehicle.Guests)
		vm.initialGuestState[i] = guestsCopy

		locationsCopy := make([]coordinates.GuestCoordinates, len(vehicle.Locations))
		copy(locationsCopy, vehicle.Locations)
		vm.initialLocationState[i] = locationsCopy

		// Deep copy the route - need to copy the linked list
		routeCopy := app.Route{
			DestinationCount: vehicle.Route.DestinationCount,
			List:             nil, // Will be set below
		}

		// Copy the linked list if it exists
		if vehicle.Route.List != nil {
			routeCopy.List = list.New()
			for elem := vehicle.Route.List.Front(); elem != nil; elem = elem.Next() {
				routeCopy.List.PushBack(elem.Value)
			}
		}

		vm.initialRouteState[i] = routeCopy
	}
	vm.hasChanges = false
}

// MoveGuest moves a guest from one vehicle to another
func (vm *VehicleManager) MoveGuest(guest *app.Guest, fromVehicle, toVehicle int) error {
	rm := vm.config.Rp.rm

	if fromVehicle < 0 || fromVehicle >= len(rm.Vehicles) ||
		toVehicle < 0 || toVehicle >= len(rm.Vehicles) {
		return NewVehicleError("invalid vehicle index")
	}

	sourceVehicle := &rm.Vehicles[fromVehicle]
	targetVehicle := &rm.Vehicles[toVehicle]

	// Check capacity
	if targetVehicle.SeatsRemaining < guest.GroupSize {
		return NewVehicleError("insufficient capacity in target vehicle")
	}

	// Find and remove guest from source
	guestIndex := vm.findGuestIndex(sourceVehicle, guest)
	if guestIndex < 0 {
		return NewVehicleError("guest not found in source vehicle")
	}

	// Perform the move
	sourceVehicle.Guests = append(
		sourceVehicle.Guests[:guestIndex],
		sourceVehicle.Guests[guestIndex+1:]...,
	)
	sourceVehicle.SeatsRemaining += guest.GroupSize

	targetVehicle.Guests = append(targetVehicle.Guests, *guest)
	targetVehicle.SeatsRemaining -= guest.GroupSize

	// Update routes for both vehicles
	vm.updateVehicleRoute(fromVehicle)
	vm.updateVehicleRoute(toVehicle)

	vm.hasChanges = true
	return nil
}

// ResetToInitialState restores all vehicles to their original state INCLUDING routes
func (vm *VehicleManager) ResetToInitialState() {
	// First restore guest assignments
	for i, originalGuests := range vm.initialGuestState {
		if i < len(vm.routeManager.Vehicles) {
			vehicle := &vm.routeManager.Vehicles[i]

			// Restore guest list
			vehicle.Guests = make([]app.Guest, len(originalGuests))
			copy(vehicle.Guests, originalGuests)

			// Recalculate seat capacity
			vehicle.SeatsRemaining = 4 // Reset to max capacity
			for _, guest := range vehicle.Guests {
				vehicle.SeatsRemaining -= guest.GroupSize
			}
		}
	}

	// Then restore route assignments
	for i, originalRoute := range vm.initialRouteState {
		if i < len(vm.routeManager.Vehicles) {
			vehicle := &vm.routeManager.Vehicles[i]

			// Restore route destination count
			vehicle.Route.DestinationCount = originalRoute.DestinationCount

			// Restore the linked list
			if originalRoute.List == nil {
				vehicle.Route.List = nil
			} else {
				vehicle.Route.List = list.New()
				for elem := originalRoute.List.Front(); elem != nil; elem = elem.Next() {
					vehicle.Route.List.PushBack(elem.Value)

				}
			}
		}
	}

	for i, originalLocations := range vm.initialLocationState {
		if i < len(vm.routeManager.Vehicles) {
			vehicle := &vm.routeManager.Vehicles[i]

			// Restore locations list
			vehicle.Locations = make([]coordinates.GuestCoordinates, len(originalLocations))
			copy(vehicle.Locations, originalLocations)
		}
	}

	vm.hasChanges = false

	// Refresh the grid display
	vm.grid.refreshAfterMove()

	vm.config.InfoLog.Println("Reset to initial state completed - both guests and routes restored")
}

// updateVehicleRoute updates the route for a specific vehicle based on its current guests
func (vm *VehicleManager) updateVehicleRoute(vehicleIndex int) {
	if vehicleIndex < 0 || vehicleIndex >= len(vm.routeManager.Vehicles) {
		return
	}

	vehicle := &vm.routeManager.Vehicles[vehicleIndex]
	lr := vm.config.Rp.lr

	eventType := vm.config.Rp.ae.EventType
	// Use the existing UpdateRouteFromGuests method
	vehicle.UpdateRouteFromGuests(lr, eventType)
}

// updateAllVehicleRoutes updates routes for all vehicles
func (vm *VehicleManager) updateAllVehicleRoutes() {
	lr := vm.config.Rp.lr
	eventType := vm.config.Rp.ae.EventType

	for i := range vm.routeManager.Vehicles {
		vm.routeManager.Vehicles[i].UpdateRouteFromGuests(lr, eventType)
	}
}

// SubmitChanges applies the current state as the new baseline
func (vm *VehicleManager) SubmitChanges() {
	if vm.hasChanges {
		// Make sure all routes are up to date before capturing state
		vm.updateAllVehicleRoutes()

		vm.captureInitialState() // Make current state the new baseline
		vm.grid.config.InfoLog.Println("Changes submitted successfully")
	}
}

// HasChanges returns whether there are unsaved changes
func (vm *VehicleManager) HasChanges() bool {
	return vm.hasChanges
}

// GetVehicleCapacityInfo returns capacity information for a vehicle
func (vm *VehicleManager) GetVehicleCapacityInfo(vehicleIndex int) VehicleCapacity {
	if vehicleIndex < 0 || vehicleIndex >= len(vm.routeManager.Vehicles) {
		return VehicleCapacity{}
	}

	vehicle := &vm.routeManager.Vehicles[vehicleIndex]
	used := 4 - vehicle.SeatsRemaining // Assuming max 4 seats

	return VehicleCapacity{
		MaxSeats:       4,
		UsedSeats:      used,
		RemainingSeats: vehicle.SeatsRemaining,
		GuestCount:     len(vehicle.Guests),
	}
}

// GetAllVehicleInfo returns summary information for all vehicles
func (vm *VehicleManager) GetAllVehicleInfo() []VehicleInfo {
	info := make([]VehicleInfo, len(vm.routeManager.Vehicles))

	for i, vehicle := range vm.routeManager.Vehicles {
		info[i] = VehicleInfo{
			Index:    i,
			Capacity: vm.GetVehicleCapacityInfo(i),
			Guests:   make([]app.Guest, len(vehicle.Guests)),
		}
		copy(info[i].Guests, vehicle.Guests)
	}

	return info
}

// ValidateMove checks if a move is valid without executing it
func (vm *VehicleManager) ValidateMove(guest *app.Guest, fromVehicle, toVehicle int) error {
	if fromVehicle < 0 || fromVehicle >= len(vm.routeManager.Vehicles) ||
		toVehicle < 0 || toVehicle >= len(vm.routeManager.Vehicles) {
		return NewVehicleError("invalid vehicle index")
	}

	if fromVehicle == toVehicle {
		return NewVehicleError("source and target vehicles are the same")
	}

	sourceVehicle := &vm.routeManager.Vehicles[fromVehicle]
	targetVehicle := &vm.routeManager.Vehicles[toVehicle]

	// Check if guest exists in source vehicle
	if vm.findGuestIndex(sourceVehicle, guest) < 0 {
		return NewVehicleError("guest not found in source vehicle")
	}

	// Check target capacity
	if targetVehicle.SeatsRemaining < guest.GroupSize {
		return NewVehicleError("insufficient capacity in target vehicle")
	}

	return nil
}

// findGuestIndex finds the index of a guest in a vehicle's guest list
func (vm *VehicleManager) findGuestIndex(vehicle *app.Vehicle, guest *app.Guest) int {
	for i, g := range vehicle.Guests {
		if g.Name == guest.Name && g.Address == guest.Address {
			return i
		}
	}
	return -1
}

// Supporting types

// VehicleCapacity holds capacity information for a vehicle
type VehicleCapacity struct {
	MaxSeats       int
	UsedSeats      int
	RemainingSeats int
	GuestCount     int
}

// VehicleInfo holds summary information for a vehicle
type VehicleInfo struct {
	Index    int
	Capacity VehicleCapacity
	Guests   []app.Guest
}

// VehicleError represents errors related to vehicle operations
type VehicleError struct {
	message string
}

func NewVehicleError(message string) *VehicleError {
	return &VehicleError{message: message}
}

func (e *VehicleError) Error() string {
	return e.message
}
