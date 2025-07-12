package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

var adjustedY float32 = 125

// VehicleGrid is the top-level container managing all vehicles and global drag state
type VehicleGrid struct {
	widget.BaseWidget

	// Core data
	vehicles       []*VehicleCard
	routeManager   *app.RouteManager
	config         *Config
	vehicleManager *VehicleManager

	// Drag state management
	dragOverlay  *fyne.Container // Overlay for drag visual
	dragVisual   *fyne.Container // The actual dragged guest visual
	draggedGuest *app.Guest      // Currently dragged guest
	dragOrigin   VehiclePosition // Where drag started
	isDragging   bool

	dragPosition    fyne.Position // Current position of dragged guest
	currentMousePos fyne.Position

	// UI containers
	gridContainer   *fyne.Container // The actual grid of vehicle cards
	mainContainer   *fyne.Container // Max container with overlay
	scrollContainer *container.Scroll
	eventType       string
}

// VehiclePosition identifies a specific location in the grid
type VehiclePosition struct {
	VehicleIndex int
	TileIndex    int
}

// NewVehicleGrid creates the main vehicle grid widget
func NewVehicleGrid(rm *app.RouteManager, cfg *Config) *VehicleGrid {
	vg := &VehicleGrid{
		routeManager: rm,
		config:       cfg,
		vehicles:     make([]*VehicleCard, 0),
		eventType:    cfg.Rp.ae.EventType,
	}

	// Create drag overlay that will show the dragged guest
	vg.dragOverlay = container.NewWithoutLayout()
	vg.dragOverlay.Hide()

	// Initialize vehicle manager
	vg.vehicleManager = NewVehicleManager(rm, vg, cfg)

	vg.ExtendBaseWidget(vg)

	// Initialize vehicles from route manager
	vg.refreshVehicles()

	return vg
}

// CreateRenderer creates the widget renderer for VehicleGrid
func (vg *VehicleGrid) CreateRenderer() fyne.WidgetRenderer {
	vg.gridContainer = vg.createVehicleCards()

	// Create the scrollable content
	vg.scrollContainer = container.NewScroll(vg.gridContainer)
	vg.scrollContainer.SetMinSize(fyne.NewSize(900, 500))

	// Create the main container with overlay layer on top
	vg.mainContainer = container.NewMax(
		vg.scrollContainer,
		vg.dragOverlay,
	)

	return &vehicleGridRenderer{
		grid:    vg,
		objects: []fyne.CanvasObject{vg.mainContainer},
	}
}

// vehicleGridRenderer implements fyne.WidgetRenderer for VehicleGrid
type vehicleGridRenderer struct {
	grid    *VehicleGrid
	objects []fyne.CanvasObject
}

func (r *vehicleGridRenderer) Layout(size fyne.Size) {
	r.grid.mainContainer.Resize(size)
}

func (r *vehicleGridRenderer) MinSize() fyne.Size {
	return fyne.NewSize(900, 500)
}

func (r *vehicleGridRenderer) Refresh() {
	r.grid.mainContainer.Refresh()
}

func (r *vehicleGridRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *vehicleGridRenderer) Destroy() {}

// createVehicleCards builds the grid of vehicle cards
func (vg *VehicleGrid) createVehicleCards() *fyne.Container {
	cards := make([]fyne.CanvasObject, 0, len(vg.vehicles))

	for _, vehicleCard := range vg.vehicles {
		cards = append(cards, vehicleCard.CreateCard())
	}

	// Use 4 columns for better spacing
	return container.NewGridWithColumns(4, cards...)
}

// refreshVehicles rebuilds the vehicle cards from current state
func (vg *VehicleGrid) refreshVehicles() {
	vg.vehicles = make([]*VehicleCard, len(vg.routeManager.Vehicles))

	for i := range vg.routeManager.Vehicles {
		vg.vehicles[i] = NewVehicleCard(i, &vg.routeManager.Vehicles[i], vg)
	}
}

// StartDrag initiates a drag operation
func (vg *VehicleGrid) StartDrag(guest *app.Guest, origin VehiclePosition, startPos fyne.Position, offset fyne.Position) {
	vg.draggedGuest = guest
	vg.dragOrigin = origin
	vg.isDragging = true

	// Initialize drag position (we'll update it in UpdateDrag)
	vg.dragPosition = startPos

	// Create visual representation of dragged guest
	vg.createDragVisual(guest)

	// Hide the original guest from its tile immediately
	vg.vehicles[origin.VehicleIndex].HideGuest(origin.TileIndex)

	// Don't show the drag visual yet - wait for first mouse move
	vg.dragOverlay.Hide()

	vg.config.InfoLog.Printf("Started dragging guest: %s from Vehicle %d, Tile %d",
		guest.Name, origin.VehicleIndex, origin.TileIndex)
}

// createDragVisual creates the visual representation of the dragged guest
func (vg *VehicleGrid) createDragVisual(guest *app.Guest) {
	// Create a semi-transparent version of the guest widget
	background := canvas.NewRectangle(color.NRGBA{60, 60, 70, 255}) // Semi-transparent blue
	background.CornerRadius = 3
	background.Resize(fyne.NewSize(190, 40))

	nameLabel := widget.NewLabel(guest.Name + " (" + fmt.Sprintf("%d", guest.GroupSize) + ")")
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	vg.dragVisual = container.NewMax(background, container.NewPadded(nameLabel))
	vg.dragVisual.Resize(fyne.NewSize(190, 40))
}

func (vg *VehicleGrid) UpdateDrag(globalPos fyne.Position) {
	if !vg.isDragging || vg.dragVisual == nil {
		return
	}

	vg.currentMousePos = fyne.Position{X: globalPos.X, Y: globalPos.Y - 20}

	// Show on first move if hidden
	if !vg.dragOverlay.Visible() {
		vg.dragOverlay.Objects = []fyne.CanvasObject{vg.dragVisual}
		vg.dragOverlay.Show()
	}

	// Update the drag position - center the widget on cursor
	vg.dragPosition = fyne.NewPos(
		globalPos.X-95,        // Center horizontally
		globalPos.Y-adjustedY, // Center vertically with offset correction
	)

	// Move the drag visual to the new position
	vg.dragVisual.Move(vg.dragPosition)

	// Highlight valid drop targets
	vg.highlightValidDropTargets()
}

// EndDrag completes the drag operation
func (vg *VehicleGrid) EndDrag(globalPos fyne.Position) {
	if !vg.isDragging {
		return
	}

	// Update mouse position one last time
	vg.currentMousePos = fyne.Position{X: globalPos.X, Y: globalPos.Y - 20}

	// Find the target position based on mouse position
	targetPos := vg.positionToTile()

	if vg.isValidDropTarget(targetPos) {
		vg.performMove(vg.dragOrigin, targetPos)
		vg.config.InfoLog.Printf("Moved guest %s from Vehicle %d to Vehicle %d",
			vg.draggedGuest.Name, vg.dragOrigin.VehicleIndex, targetPos.VehicleIndex)
	} else {
		// Return guest to original position
		vg.vehicles[vg.dragOrigin.VehicleIndex].ShowGuest(vg.dragOrigin.TileIndex)
		vg.config.InfoLog.Printf("Invalid drop for guest: %s", vg.draggedGuest.Name)
	}

	// Clean up drag state
	vg.cleanupDrag()
}

// CancelDrag cancels the current drag operation
func (vg *VehicleGrid) CancelDrag() {
	if vg.isDragging {
		// Show the guest back in original position
		vg.vehicles[vg.dragOrigin.VehicleIndex].ShowGuest(vg.dragOrigin.TileIndex)
		vg.cleanupDrag()
		vg.config.InfoLog.Printf("Drag cancelled")
	}
}

// cleanupDrag resets all drag state
func (vg *VehicleGrid) cleanupDrag() {
	vg.isDragging = false
	vg.draggedGuest = nil
	vg.dragOverlay.Hide()
	vg.dragOverlay.Objects = nil
	vg.dragVisual = nil

	// Remove all highlights
	for _, vehicle := range vg.vehicles {
		vehicle.RemoveAllHighlights()
	}
}

// positionToTile converts global screen coordinates to vehicle/tile position
func (vg *VehicleGrid) positionToTile() VehiclePosition {
	return vg.findTileAtPosition(vg.currentMousePos)
}

// highlightValidDropTargets highlights all valid drop targets
func (vg *VehicleGrid) highlightValidDropTargets() {
	targetPos := vg.positionToTile()

	for vIndex, vehicle := range vg.vehicles {
		for tIndex, tile := range vehicle.tiles {
			if vIndex == targetPos.VehicleIndex && tIndex == targetPos.TileIndex &&
				vg.isValidDropTarget(targetPos) {
				tile.HighlightAsDropTarget()
			} else {
				tile.RemoveHighlight()
			}
		}
	}
}

// isValidDropTarget checks if the target position can accept the dragged guest
func (vg *VehicleGrid) isValidDropTarget(target VehiclePosition) bool {
	if target.VehicleIndex < 0 || target.VehicleIndex >= len(vg.vehicles) {
		return false
	}

	if target.TileIndex < 0 || target.TileIndex >= len(vg.vehicles[target.VehicleIndex].tiles) {
		return false
	}

	// Allow same-position drops (will be ignored in performMove)
	if target.VehicleIndex == vg.dragOrigin.VehicleIndex &&
		target.TileIndex == vg.dragOrigin.TileIndex {
		return true
	}

	vehicle := vg.vehicles[target.VehicleIndex]

	// For same-vehicle moves, we need to check differently
	if target.VehicleIndex == vg.dragOrigin.VehicleIndex {
		// The tile might appear occupied but it's actually the dragged guest
		// Check if the tile index is beyond current guests (empty tiles)
		return target.TileIndex >= len(vg.routeManager.Vehicles[target.VehicleIndex].Guests)
	}

	// For different vehicles, check if tile is empty
	return vehicle.IsTileEmpty(target.TileIndex)
}

func (vg *VehicleGrid) performMove(from, to VehiclePosition) {
	rm := vg.config.Rp.rm
	lr := vg.config.Rp.lr

	// Special handling for same-vehicle moves
	if from.VehicleIndex == to.VehicleIndex {
		vehicle := &rm.Vehicles[from.VehicleIndex]

		// Find the guest's current position
		guestIndex := vg.findGuestIndex(vehicle, vg.draggedGuest)
		if guestIndex < 0 {
			return
		}

		// Save the guest
		guest := vehicle.Guests[guestIndex]

		// Calculate the actual insert position
		insertPos := to.TileIndex
		if insertPos > len(vehicle.Guests) {
			insertPos = len(vehicle.Guests)
		}

		// If we're moving to a position after the current position,
		// we need to adjust because removing will shift indices
		if insertPos > guestIndex {
			insertPos--
		}

		// Remove guest from current position
		vehicle.Guests = append(
			vehicle.Guests[:guestIndex],
			vehicle.Guests[guestIndex+1:]...,
		)

		// Insert guest at new position
		vehicle.Guests = append(
			vehicle.Guests[:insertPos],
			append([]app.Guest{guest}, vehicle.Guests[insertPos:]...)...,
		)

		// Update the route for this vehicle
		vehicle.UpdateRouteFromGuests(lr, vg.eventType)

		// Refresh just this vehicle
		vg.refreshAfterMove()
		return
	}

	// Different vehicle move
	sourceVehicle := &rm.Vehicles[from.VehicleIndex]
	guestIndex := vg.findGuestIndex(sourceVehicle, vg.draggedGuest)
	guest := sourceVehicle.Guests[guestIndex]

	if guestIndex >= 0 {
		// Remove from source vehicle
		sourceVehicle.Guests = append(
			sourceVehicle.Guests[:guestIndex],
			sourceVehicle.Guests[guestIndex+1:]...,
		)
		sourceVehicle.SeatsRemaining += vg.draggedGuest.GroupSize

		// Update source vehicle route
		sourceVehicle.UpdateRouteFromGuests(lr, vg.eventType)
	}

	// Add guest to target vehicle at the specific tile position
	targetVehicle := &rm.Vehicles[to.VehicleIndex]

	// Insert at the tile position (but don't exceed current guest count)
	insertPos := to.TileIndex
	if insertPos > len(targetVehicle.Guests) {
		insertPos = len(targetVehicle.Guests)
	}

	// Insert guest at the specified position
	targetVehicle.Guests = append(targetVehicle.Guests[:insertPos],
		append([]app.Guest{guest}, targetVehicle.Guests[insertPos:]...)...)

	targetVehicle.SeatsRemaining -= guest.GroupSize

	// Update target vehicle route
	targetVehicle.UpdateRouteFromGuests(lr, vg.eventType)

	// Mark that we have changes
	vg.vehicleManager.hasChanges = true

	// Refresh both vehicles
	vg.refreshAfterMove()
}

// findGuestIndex finds the index of a guest in a vehicle's guest list
func (vg *VehicleGrid) findGuestIndex(vehicle *app.Vehicle, guest *app.Guest) int {
	for i, g := range vehicle.Guests {
		if g.Name == guest.Name && g.Address == guest.Address {
			return i
		}
	}
	return -1
}

// refreshAfterMove rebuilds the UI after a successful move
func (vg *VehicleGrid) refreshAfterMove() {
	// Clean up any remaining highlights BEFORE refreshing
	for _, vehicle := range vg.vehicles {
		vehicle.RemoveAllHighlights()
	}

	// Refresh all vehicle cards to reflect new state
	for i := range vg.vehicles {
		vg.vehicles[i].RefreshTiles()
	}

	// Rebuild the grid container
	vg.gridContainer.Objects = nil
	for _, vehicleCard := range vg.vehicles {
		vg.gridContainer.Objects = append(vg.gridContainer.Objects, vehicleCard.CreateCard())
	}
	vg.gridContainer.Refresh()
}

// ResetVehicles resets all vehicles to their initial state
func (vg *VehicleGrid) ResetVehicles() {
	if vg.vehicleManager != nil {
		vg.vehicleManager.ResetToInitialState()
		vg.config.InfoLog.Println("Vehicles reset to initial state")
	}
}

// Updated SubmitChanges method
func (vg *VehicleGrid) SubmitChanges() {
	if vg.vehicleManager != nil {
		vg.vehicleManager.SubmitChanges()
		vg.config.InfoLog.Println("Changes submitted successfully")
	}
}

// Implement Mouseable interface to capture mouse events
func (vg *VehicleGrid) MouseIn(*desktop.MouseEvent) {}
func (vg *VehicleGrid) MouseOut()                   {}

func (vg *VehicleGrid) MouseMoved(event *desktop.MouseEvent) {
	if vg.isDragging {
		vg.UpdateDrag(event.AbsolutePosition)
	}
}

// IsDragging returns whether a drag is in progress
func (vg *VehicleGrid) IsDragging() bool {
	return vg.isDragging
}

// GetDraggedGuest returns the currently dragged guest
func (vg *VehicleGrid) GetDraggedGuest() *app.Guest {
	return vg.draggedGuest
}

// Add this method to vehicle_grid.go to fix position calculations

// findTileAtPosition finds which tile is at the given window position
func (vg *VehicleGrid) findTileAtPosition(mousePos fyne.Position) VehiclePosition {
	// Get the position of our main container in the window
	mainPos := vg.mainContainer.Position()

	// Adjust for scroll offset if content is scrolled
	scrollOffset := vg.scrollContainer.Offset

	// Calculate the position relative to the scrolled content
	// Add back the Y offset that seems to be present in your setup
	adjustedMouseY := mousePos.Y - adjustedY // Compensate for the same offset used in visual positioning

	contentPos := fyne.NewPos(
		mousePos.X-mainPos.X+scrollOffset.X,
		adjustedMouseY-mainPos.Y+scrollOffset.Y,
	)

	// Now check each vehicle card
	for vIndex, vehicle := range vg.vehicles {
		if vehicle.card == nil {
			continue
		}

		cardPos := vehicle.card.Position()
		cardSize := vehicle.card.Size()

		// Check if we're within this card's bounds
		if contentPos.X >= cardPos.X && contentPos.X <= cardPos.X+cardSize.Width &&
			contentPos.Y >= cardPos.Y && contentPos.Y <= cardPos.Y+cardSize.Height {

			// We're in this card, now find which tile
			for tIndex, tile := range vehicle.tiles {
				if tile.container == nil {
					continue
				}

				// Get tile position relative to the card
				tilePos := tile.container.Position()
				tileSize := tile.container.Size()

				// Calculate tile's absolute position within the card
				tilePosInCard := fyne.NewPos(
					cardPos.X+tilePos.X,
					cardPos.Y+tilePos.Y,
				)

				// Check if the position is within this tile
				if contentPos.X >= tilePosInCard.X && contentPos.X <= tilePosInCard.X+tileSize.Width &&
					contentPos.Y >= tilePosInCard.Y && contentPos.Y <= tilePosInCard.Y+tileSize.Height {

					return VehiclePosition{VehicleIndex: vIndex, TileIndex: tIndex}
				}
			}
		}
	}

	return VehiclePosition{VehicleIndex: -1, TileIndex: -1}
}

// GetDragPosition returns the current position of the dragged guest
func (vg *VehicleGrid) GetDragPosition() fyne.Position {
	return vg.dragPosition
}

// Call this method after any guest move operation in VehicleGrid
func (vg *VehicleGrid) updateVehicleRoutes() {
	rm := vg.config.Rp.rm
	lr := vg.config.Rp.lr

	for i := range rm.Vehicles {
		rm.Vehicles[i].UpdateRouteFromGuests(lr, vg.eventType)
	}
}
