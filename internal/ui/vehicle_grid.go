package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// VehicleGrid is the top-level container managing all vehicles and global drag state
type VehicleGrid struct {
	widget.BaseWidget

	// Core data
	vehicles       []*VehicleCard
	routeManager   *app.RouteManager
	config         *Config
	vehicleManager *VehicleManager

	// Drag state management
	dragOverlay   *canvas.Image   // Global overlay for dragging guests
	draggedGuest  *app.Guest      // Currently dragged guest
	dragOrigin    VehiclePosition // Where drag started
	isDragging    bool
	dropHighlight *canvas.Rectangle // Shows valid drop zones

	// UI containers
	gridContainer *fyne.Container // The actual grid of vehicle cards
	mainContainer *fyne.Container // Max container with overlay
}

// VehiclePosition identifies a specific location in the grid
type VehiclePosition struct {
	VehicleIndex int
	TileIndex    int
}

// NewVehicleGrid creates the main vehicle grid widget
func NewVehicleGrid(rm *app.RouteManager, cfg *Config) *VehicleGrid {
	vg := &VehicleGrid{
		routeManager:  rm,
		config:        cfg,
		vehicles:      make([]*VehicleCard, 0),
		dropHighlight: canvas.NewRectangle(color.NRGBA{0, 255, 0, 64}),
	}

	vg.dragOverlay = canvas.NewImageFromResource(nil)
	vg.dragOverlay.FillMode = canvas.ImageFillContain
	vg.dragOverlay.Hide()

	vg.dropHighlight.Hide()
	vg.dropHighlight.StrokeWidth = 3
	vg.dropHighlight.StrokeColor = color.NRGBA{0, 255, 0, 255}

	// Initialize vehicle manager
	vg.vehicleManager = NewVehicleManager(rm, vg)

	vg.ExtendBaseWidget(vg)

	// Initialize vehicles from route manager
	vg.refreshVehicles()

	return vg
}

// CreateRenderer creates the widget renderer for VehicleGrid
func (vg *VehicleGrid) CreateRenderer() fyne.WidgetRenderer {
	vg.gridContainer = vg.createVehicleCards()

	// Create the main container with overlay (similar to chess)
	vg.mainContainer = container.NewMax(
		vg.gridContainer,
		container.NewWithoutLayout(vg.dropHighlight, vg.dragOverlay),
	)

	scrollContainer := container.NewScroll(vg.mainContainer)

	return &vehicleGridRenderer{
		grid:    vg,
		scroll:  scrollContainer,
		objects: []fyne.CanvasObject{scrollContainer},
	}
}

// vehicleGridRenderer implements fyne.WidgetRenderer for VehicleGrid
type vehicleGridRenderer struct {
	grid    *VehicleGrid
	scroll  *container.Scroll
	objects []fyne.CanvasObject
}

func (r *vehicleGridRenderer) Layout(size fyne.Size) {
	r.scroll.Resize(size)
}

func (r *vehicleGridRenderer) MinSize() fyne.Size {
	return fyne.NewSize(900, 500) // Minimum size for the vehicle grid
}

func (r *vehicleGridRenderer) Refresh() {
	r.scroll.Refresh()
}

func (r *vehicleGridRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *vehicleGridRenderer) Destroy() {
	// Clean up any resources if needed
}

// createVehicleCards builds the grid of vehicle cards
func (vg *VehicleGrid) createVehicleCards() *fyne.Container {
	cards := make([]fyne.CanvasObject, 0, len(vg.vehicles))

	for _, vehicleCard := range vg.vehicles {
		cards = append(cards, vehicleCard.CreateCard())
	}

	// Use 3 columns instead of 4 for better spacing
	return container.NewGridWithColumns(3, cards...)
}

// refreshVehicles rebuilds the vehicle cards from current state
func (vg *VehicleGrid) refreshVehicles() {
	vg.vehicles = make([]*VehicleCard, len(vg.routeManager.Vehicles))

	for i := range vg.routeManager.Vehicles {
		vg.vehicles[i] = NewVehicleCard(i, &vg.routeManager.Vehicles[i], vg)
	}
}

// StartDrag initiates a drag operation
func (vg *VehicleGrid) StartDrag(guest *app.Guest, origin VehiclePosition, startPos fyne.Position) {
	vg.draggedGuest = guest
	vg.dragOrigin = origin
	vg.isDragging = true

	// Set up the drag overlay
	vg.dragOverlay.Move(startPos)
	vg.dragOverlay.Show()

	// Hide the original guest from its tile
	vg.vehicles[origin.VehicleIndex].HideGuest(origin.TileIndex)

	vg.config.InfoLog.Printf("Started dragging guest: %s from Vehicle %d, Tile %d",
		guest.Name, origin.VehicleIndex, origin.TileIndex)
}

// UpdateDrag updates the drag position and highlights valid drop zones
func (vg *VehicleGrid) UpdateDrag(newPos fyne.Position) {
	if !vg.isDragging {
		return
	}

	vg.dragOverlay.Move(newPos)

	// Calculate which tile is under the mouse and highlight it
	targetPos := vg.positionToTile(newPos)
	if vg.isValidDropTarget(targetPos) {
		vg.highlightDropTarget(targetPos)
	} else {
		vg.dropHighlight.Hide()
	}
}

// EndDrag completes the drag operation
func (vg *VehicleGrid) EndDrag(endPos fyne.Position) {
	if !vg.isDragging {
		return
	}

	targetPos := vg.positionToTile(endPos)

	if vg.isValidDropTarget(targetPos) {
		vg.performMove(vg.dragOrigin, targetPos)
		vg.config.InfoLog.Printf("Moved guest %s from Vehicle %d to Vehicle %d",
			vg.draggedGuest.Name, vg.dragOrigin.VehicleIndex, targetPos.VehicleIndex)
	} else {
		// Invalid drop - return to original position
		vg.vehicles[vg.dragOrigin.VehicleIndex].ShowGuest(vg.dragOrigin.TileIndex)
		vg.config.InfoLog.Printf("Invalid drop for guest: %s", vg.draggedGuest.Name)
	}

	// Clean up drag state
	vg.isDragging = false
	vg.draggedGuest = nil
	vg.dragOverlay.Hide()
	vg.dropHighlight.Hide()

	// Refresh the entire grid
	vg.refreshAfterMove()
}

// positionToTile converts screen coordinates to vehicle/tile position
func (vg *VehicleGrid) positionToTile(pos fyne.Position) VehiclePosition {
	// This is a simplified implementation - in practice you'd need to
	// calculate based on actual vehicle card positions and sizes
	for vIndex, vehicle := range vg.vehicles {
		for tIndex := range vehicle.tiles {
			tilePos := vehicle.GetTilePosition(tIndex)
			tileSize := vehicle.GetTileSize()

			if pos.X >= tilePos.X && pos.X <= tilePos.X+tileSize.Width &&
				pos.Y >= tilePos.Y && pos.Y <= tilePos.Y+tileSize.Height {
				return VehiclePosition{VehicleIndex: vIndex, TileIndex: tIndex}
			}
		}
	}

	return VehiclePosition{VehicleIndex: -1, TileIndex: -1}
}

// isValidDropTarget checks if the target position can accept the dragged guest
func (vg *VehicleGrid) isValidDropTarget(target VehiclePosition) bool {
	if target.VehicleIndex < 0 || target.VehicleIndex >= len(vg.vehicles) {
		return false
	}

	vehicle := vg.vehicles[target.VehicleIndex]

	// Check if tile is empty and vehicle has capacity
	return vehicle.IsTileEmpty(target.TileIndex) &&
		vehicle.HasCapacityForGuest(vg.draggedGuest)
}

// highlightDropTarget shows visual feedback for valid drop zones
func (vg *VehicleGrid) highlightDropTarget(target VehiclePosition) {
	if target.VehicleIndex < 0 {
		vg.dropHighlight.Hide()
		return
	}

	tilePos := vg.vehicles[target.VehicleIndex].GetTilePosition(target.TileIndex)
	tileSize := vg.vehicles[target.VehicleIndex].GetTileSize()

	vg.dropHighlight.Move(tilePos)
	vg.dropHighlight.Resize(tileSize)
	vg.dropHighlight.Show()
}

// performMove executes the actual guest movement between vehicles
func (vg *VehicleGrid) performMove(from, to VehiclePosition) {
	// Remove guest from source vehicle
	sourceVehicle := &vg.routeManager.Vehicles[from.VehicleIndex]
	guestIndex := vg.findGuestIndex(sourceVehicle, vg.draggedGuest)
	if guestIndex >= 0 {
		sourceVehicle.Guests = append(
			sourceVehicle.Guests[:guestIndex],
			sourceVehicle.Guests[guestIndex+1:]...,
		)
	}

	// Add guest to target vehicle
	targetVehicle := &vg.routeManager.Vehicles[to.VehicleIndex]
	targetVehicle.Guests = append(targetVehicle.Guests, *vg.draggedGuest)

	// Update seat counts
	sourceVehicle.SeatsRemaining += vg.draggedGuest.GroupSize
	targetVehicle.SeatsRemaining -= vg.draggedGuest.GroupSize
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
	vg.vehicleManager.ResetToInitialState()
}

// SubmitChanges applies current changes as the new baseline
func (vg *VehicleGrid) SubmitChanges() {
	vg.vehicleManager.SubmitChanges()
}
