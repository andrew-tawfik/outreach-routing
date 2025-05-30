package ui

import (
	"fyne.io/fyne/v2"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// DragManager handles global drag state and operations
type DragManager struct {
	// State
	isDragging   bool
	draggedGuest *app.Guest
	dragOrigin   VehiclePosition
	currentHover VehiclePosition

	// Visual feedback
	validTargets    []*GuestTile
	highlightedTile *GuestTile

	// References
	grid *VehicleGrid
}

// NewDragManager creates a new drag manager
func NewDragManager(grid *VehicleGrid) *DragManager {
	return &DragManager{
		grid:         grid,
		validTargets: make([]*GuestTile, 0),
	}
}

// StartDrag initiates a drag operation with the specified guest
func (dm *DragManager) StartDrag(guest *app.Guest, origin VehiclePosition) {
	dm.isDragging = true
	dm.draggedGuest = guest
	dm.dragOrigin = origin
	dm.currentHover = VehiclePosition{VehicleIndex: -1, TileIndex: -1}

	// Calculate all valid drop targets
	dm.calculateValidTargets()

	dm.grid.config.InfoLog.Printf("Drag started for guest: %s from Vehicle %d, Tile %d",
		guest.Name, origin.VehicleIndex, origin.TileIndex)
}

// UpdateDrag updates the drag state based on current mouse position
func (dm *DragManager) UpdateDrag(mousePos fyne.Position) {
	if !dm.isDragging {
		return
	}

	// Determine which tile is under the mouse
	newHover := dm.positionToTile(mousePos)

	if newHover != dm.currentHover {
		// Remove highlight from previous tile
		if dm.highlightedTile != nil {
			dm.highlightedTile.RemoveHighlight()
			dm.highlightedTile = nil
		}

		// Highlight new tile if it's valid
		if dm.isValidTarget(newHover) {
			tile := dm.getTileAt(newHover)
			if tile != nil {
				tile.HighlightAsDropTarget()
				dm.highlightedTile = tile
			}
		}

		dm.currentHover = newHover
	}
}

// EndDrag completes the drag operation
func (dm *DragManager) EndDrag(endPos fyne.Position) bool {
	if !dm.isDragging {
		return false
	}

	success := false
	finalTarget := dm.positionToTile(endPos)

	if dm.isValidTarget(finalTarget) {
		// Perform the move
		success = dm.performMove(dm.dragOrigin, finalTarget)
		if success {
			dm.grid.config.InfoLog.Printf("Successfully moved guest %s from Vehicle %d to Vehicle %d",
				dm.draggedGuest.Name, dm.dragOrigin.VehicleIndex, finalTarget.VehicleIndex)
		}
	}

	// Clean up drag state
	dm.cleanupDrag()

	return success
}

// CancelDrag cancels the current drag operation
func (dm *DragManager) CancelDrag() {
	if dm.isDragging {
		dm.grid.config.InfoLog.Printf("Drag cancelled for guest: %s", dm.draggedGuest.Name)
		dm.cleanupDrag()
	}
}

// cleanupDrag resets the drag state and removes visual feedback
func (dm *DragManager) cleanupDrag() {
	// Remove any highlights
	if dm.highlightedTile != nil {
		dm.highlightedTile.RemoveHighlight()
		dm.highlightedTile = nil
	}

	// Clear valid targets
	for _, tile := range dm.validTargets {
		tile.RemoveHighlight()
	}
	dm.validTargets = dm.validTargets[:0]

	// Reset state
	dm.isDragging = false
	dm.draggedGuest = nil
	dm.dragOrigin = VehiclePosition{VehicleIndex: -1, TileIndex: -1}
	dm.currentHover = VehiclePosition{VehicleIndex: -1, TileIndex: -1}
}

// calculateValidTargets determines all tiles that can accept the dragged guest
func (dm *DragManager) calculateValidTargets() {
	dm.validTargets = dm.validTargets[:0]

	for vIndex, vehicle := range dm.grid.vehicles {
		// Skip the origin vehicle for now (could allow reordering later)
		if vIndex == dm.dragOrigin.VehicleIndex {
			continue
		}

		// Check if vehicle has capacity
		if !vehicle.HasCapacityForGuest(dm.draggedGuest) {
			continue
		}

		// Add all empty tiles in this vehicle to valid targets
		for _, tile := range vehicle.tiles {
			if tile.IsEmpty() {
				dm.validTargets = append(dm.validTargets, tile)
			}
		}
	}

	dm.grid.config.InfoLog.Printf("Found %d valid drop targets", len(dm.validTargets))
}

// isValidTarget checks if the given position is a valid drop target
func (dm *DragManager) isValidTarget(pos VehiclePosition) bool {
	if pos.VehicleIndex < 0 || pos.TileIndex < 0 {
		return false
	}

	for _, tile := range dm.validTargets {
		tilePos := tile.GetPosition()
		if tilePos.VehicleIndex == pos.VehicleIndex && tilePos.TileIndex == pos.TileIndex {
			return true
		}
	}
	return false
}

// getTileAt returns the tile at the specified position
func (dm *DragManager) getTileAt(pos VehiclePosition) *GuestTile {
	if pos.VehicleIndex < 0 || pos.VehicleIndex >= len(dm.grid.vehicles) {
		return nil
	}

	vehicle := dm.grid.vehicles[pos.VehicleIndex]
	if pos.TileIndex < 0 || pos.TileIndex >= len(vehicle.tiles) {
		return nil
	}

	return vehicle.tiles[pos.TileIndex]
}

// positionToTile converts screen coordinates to tile position
func (dm *DragManager) positionToTile(pos fyne.Position) VehiclePosition {
	// Check each vehicle's tiles
	for vIndex, vehicle := range dm.grid.vehicles {
		for tIndex := range vehicle.tiles {
			tilePos := vehicle.GetTilePosition(tIndex)
			tileSize := vehicle.GetTileSize()

			// Check if position is within this tile
			if pos.X >= tilePos.X && pos.X <= tilePos.X+tileSize.Width &&
				pos.Y >= tilePos.Y && pos.Y <= tilePos.Y+tileSize.Height {
				return VehiclePosition{VehicleIndex: vIndex, TileIndex: tIndex}
			}
		}
	}

	return VehiclePosition{VehicleIndex: -1, TileIndex: -1}
}

// performMove executes the actual guest movement between vehicles
func (dm *DragManager) performMove(from, to VehiclePosition) bool {
	// Validate the move
	if from.VehicleIndex < 0 || to.VehicleIndex < 0 {
		return false
	}

	// Get source and target vehicles
	sourceVehicle := &dm.grid.routeManager.Vehicles[from.VehicleIndex]
	targetVehicle := &dm.grid.routeManager.Vehicles[to.VehicleIndex]

	// Find and remove guest from source vehicle
	guestIndex := dm.findGuestIndex(sourceVehicle, dm.draggedGuest)
	if guestIndex < 0 {
		return false
	}

	// Remove guest from source
	sourceVehicle.Guests = append(
		sourceVehicle.Guests[:guestIndex],
		sourceVehicle.Guests[guestIndex+1:]...,
	)
	sourceVehicle.SeatsRemaining += dm.draggedGuest.GroupSize

	// Add guest to target
	targetVehicle.Guests = append(targetVehicle.Guests, *dm.draggedGuest)
	targetVehicle.SeatsRemaining -= dm.draggedGuest.GroupSize

	// Refresh both vehicle cards
	dm.grid.vehicles[from.VehicleIndex].RefreshTiles()
	dm.grid.vehicles[to.VehicleIndex].RefreshTiles()

	return true
}

// findGuestIndex finds the index of a guest in a vehicle's guest list
func (dm *DragManager) findGuestIndex(vehicle *app.Vehicle, guest *app.Guest) int {
	for i, g := range vehicle.Guests {
		if g.Name == guest.Name && g.Address == guest.Address {
			return i
		}
	}
	return -1
}

// IsDragging returns whether a drag operation is currently active
func (dm *DragManager) IsDragging() bool {
	return dm.isDragging
}

// GetDraggedGuest returns the currently dragged guest
func (dm *DragManager) GetDraggedGuest() *app.Guest {
	return dm.draggedGuest
}
