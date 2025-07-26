package ui

import (
	"container/list"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

type VehicleManager struct {
	routeManager *app.RouteManager
	grid         *VehicleGrid
	config       *Config

	initialGuestState    map[int][]app.Guest
	initialRouteState    map[int]app.Route
	initialLocationState map[int][]coordinates.GuestCoordinates
	hasChanges           bool
}

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

func (vm *VehicleManager) captureInitialState() {
	for i, vehicle := range vm.routeManager.Vehicles {

		guestsCopy := make([]app.Guest, len(vehicle.Guests))
		copy(guestsCopy, vehicle.Guests)
		vm.initialGuestState[i] = guestsCopy

		locationsCopy := make([]coordinates.GuestCoordinates, len(vehicle.Locations))
		copy(locationsCopy, vehicle.Locations)
		vm.initialLocationState[i] = locationsCopy

		routeCopy := app.Route{
			DestinationCount: vehicle.Route.DestinationCount,
			List:             nil,
		}

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

func (vm *VehicleManager) MoveGuest(guest *app.Guest, fromVehicle, toVehicle int) error {
	rm := vm.config.Rp.rm

	if fromVehicle < 0 || fromVehicle >= len(rm.Vehicles) ||
		toVehicle < 0 || toVehicle >= len(rm.Vehicles) {
		return NewVehicleError("invalid vehicle index")
	}

	sourceVehicle := &rm.Vehicles[fromVehicle]
	targetVehicle := &rm.Vehicles[toVehicle]

	if targetVehicle.SeatsRemaining < guest.GroupSize {
		return NewVehicleError("insufficient capacity in target vehicle")
	}

	guestIndex := vm.findGuestIndex(sourceVehicle, guest)
	if guestIndex < 0 {
		return NewVehicleError("guest not found in source vehicle")
	}

	sourceVehicle.Guests = append(
		sourceVehicle.Guests[:guestIndex],
		sourceVehicle.Guests[guestIndex+1:]...,
	)
	sourceVehicle.SeatsRemaining += guest.GroupSize

	targetVehicle.Guests = append(targetVehicle.Guests, *guest)
	targetVehicle.SeatsRemaining -= guest.GroupSize

	vm.updateVehicleRoute(fromVehicle)
	vm.updateVehicleRoute(toVehicle)

	vm.hasChanges = true
	return nil
}

func (vm *VehicleManager) ResetToInitialState() {

	for i, originalGuests := range vm.initialGuestState {
		if i < len(vm.routeManager.Vehicles) {
			vehicle := &vm.routeManager.Vehicles[i]

			vehicle.Guests = make([]app.Guest, len(originalGuests))
			copy(vehicle.Guests, originalGuests)

			vehicle.SeatsRemaining = 4
			for _, guest := range vehicle.Guests {
				vehicle.SeatsRemaining -= guest.GroupSize
			}
		}
	}

	for i, originalRoute := range vm.initialRouteState {
		if i < len(vm.routeManager.Vehicles) {
			vehicle := &vm.routeManager.Vehicles[i]

			vehicle.Route.DestinationCount = originalRoute.DestinationCount

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

			vehicle.Locations = make([]coordinates.GuestCoordinates, len(originalLocations))
			copy(vehicle.Locations, originalLocations)
		}
	}

	vm.hasChanges = false

	vm.grid.refreshAfterMove()

}

func (vm *VehicleManager) updateVehicleRoute(vehicleIndex int) {
	if vehicleIndex < 0 || vehicleIndex >= len(vm.routeManager.Vehicles) {
		return
	}

	vehicle := &vm.routeManager.Vehicles[vehicleIndex]
	lr := vm.config.Rp.lr

	eventType := vm.config.Rp.ae.EventType

	vehicle.UpdateRouteFromGuests(lr, eventType)
}

func (vm *VehicleManager) updateAllVehicleRoutes() {
	lr := vm.config.Rp.lr
	eventType := vm.config.Rp.ae.EventType

	for i := range vm.routeManager.Vehicles {
		vm.routeManager.Vehicles[i].UpdateRouteFromGuests(lr, eventType)
	}
}

func (vm *VehicleManager) SubmitChanges() {
	if vm.hasChanges {

		vm.updateAllVehicleRoutes()

		vm.captureInitialState()
	}
}

func (vm *VehicleManager) HasChanges() bool {
	return vm.hasChanges
}

func (vm *VehicleManager) GetVehicleCapacityInfo(vehicleIndex int) VehicleCapacity {
	if vehicleIndex < 0 || vehicleIndex >= len(vm.routeManager.Vehicles) {
		return VehicleCapacity{}
	}

	vehicle := &vm.routeManager.Vehicles[vehicleIndex]
	used := 4 - vehicle.SeatsRemaining

	return VehicleCapacity{
		MaxSeats:       4,
		UsedSeats:      used,
		RemainingSeats: vehicle.SeatsRemaining,
		GuestCount:     len(vehicle.Guests),
	}
}

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

	if vm.findGuestIndex(sourceVehicle, guest) < 0 {
		return NewVehicleError("guest not found in source vehicle")
	}

	if targetVehicle.SeatsRemaining < guest.GroupSize {
		return NewVehicleError("insufficient capacity in target vehicle")
	}

	return nil
}

func (vm *VehicleManager) findGuestIndex(vehicle *app.Vehicle, guest *app.Guest) int {
	for i, g := range vehicle.Guests {
		if g.Name == guest.Name && g.Address == guest.Address {
			return i
		}
	}
	return -1
}

type VehicleCapacity struct {
	MaxSeats       int
	UsedSeats      int
	RemainingSeats int
	GuestCount     int
}

type VehicleInfo struct {
	Index    int
	Capacity VehicleCapacity
	Guests   []app.Guest
}

type VehicleError struct {
	message string
}

func NewVehicleError(message string) *VehicleError {
	return &VehicleError{message: message}
}

func (e *VehicleError) Error() string {
	return e.message
}
