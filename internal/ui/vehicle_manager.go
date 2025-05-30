package ui

import (
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// VehicleManager handles vehicle state management and operations
type VehicleManager struct {
	routeManager *app.RouteManager
	grid         *VehicleGrid

	// State tracking
	initialState map[int][]app.Guest // Original state for reset functionality
	hasChanges   bool
}

// NewVehicleManager creates a new vehicle manager
func NewVehicleManager(rm *app.RouteManager, grid *VehicleGrid) *VehicleManager {
	vm := &VehicleManager{
		routeManager: rm,
		grid:         grid,
		initialState: make(map[int][]app.Guest),
	}

	vm.captureInitialState()
	return vm
}

// captureInitialState saves the current vehicle assignments for reset functionality
func (vm *VehicleManager) captureInitialState() {
	for i, vehicle := range vm.routeManager.Vehicles {
		// Deep copy the guests slice
		guestsCopy := make([]app.Guest, len(vehicle.Guests))
		copy(guestsCopy, vehicle.Guests)
		vm.initialState[i] = guestsCopy
	}
	vm.hasChanges = false
}

// MoveGuest moves a guest from one vehicle to another
func (vm *VehicleManager) MoveGuest(guest *app.Guest, fromVehicle, toVehicle int) error {
	if fromVehicle < 0 || fromVehicle >= len(vm.routeManager.Vehicles) ||
		toVehicle < 0 || toVehicle >= len(vm.routeManager.Vehicles) {
		return NewVehicleError("invalid vehicle index")
	}

	sourceVehicle := &vm.routeManager.Vehicles[fromVehicle]
	targetVehicle := &vm.routeManager.Vehicles[toVehicle]

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

	vm.hasChanges = true
	return nil
}

// ResetToInitialState restores all vehicles to their original state
func (vm *VehicleManager) ResetToInitialState() {
	for i, originalGuests := range vm.initialState {
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

	vm.hasChanges = false

	// Refresh the grid display
	vm.grid.refreshAfterMove()
}

// SubmitChanges applies the current state as the new baseline
func (vm *VehicleManager) SubmitChanges() {
	if vm.hasChanges {
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
