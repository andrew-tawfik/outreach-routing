package ui

import (
	"fmt"
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
	dragContainer *fyne.Container // Container for the dragged guest visual
	draggedGuest  *app.Guest      // Currently dragged guest
	dragOrigin    VehiclePosition // Where drag started
	isDragging    bool
	dropHighlight *canvas.Rectangle // Shows valid drop zones

	// UI containers
	gridContainer   *fyne.Container // The actual grid of vehicle cards
	mainContainer   *fyne.Container // Max container with overlay
	scrollContainer *container.Scroll
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

	// Create drag container that will show the dragged guest
	vg.dragContainer = container.NewWithoutLayout()
	vg.dragContainer.Hide()

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

	// Create the scrollable content
	vg.scrollContainer = container.NewScroll(vg.gridContainer)
	vg.scrollContainer.SetMinSize(fyne.NewSize(900, 500))

	// Create the main container with overlay layer on top
	vg.mainContainer = container.NewMax(
		vg.scrollContainer,
		container.NewWithoutLayout(vg.dropHighlight, vg.dragContainer),
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

	// Use 3 columns for better spacing
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

	// Create visual representation of dragged guest
	vg.createDragVisual(guest)

	// Position it at the start position
	vg.dragContainer.Move(startPos)
	vg.dragContainer.Show()

	// Hide the original guest from its tile
	vg.vehicles[origin.VehicleIndex].HideGuest(origin.TileIndex)

	vg.config.InfoLog.Printf("Started dragging guest: %s from Vehicle %d, Tile %d",
		guest.Name, origin.VehicleIndex, origin.TileIndex)
}

// createDragVisual creates the visual representation of the dragged guest
func (vg *VehicleGrid) createDragVisual(guest *app.Guest) {
	// Create a semi-transparent version of the guest widget
	background := canvas.NewRectangle(color.NRGBA{70, 130, 180, 200}) // Semi-transparent blue
	background.CornerRadius = 3
	background.Resize(fyne.NewSize(190, 40))

	nameLabel := widget.NewLabel(guest.Name + " (" + fmt.Sprintf("%d", guest.GroupSize) + ")")
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	addressLabel := widget.NewLabel(guest.Address)
	if len(guest.Address) > 20 {
		addressLabel.SetText(guest.Address[:17] + "...")
	}

	content := container.NewVBox(nameLabel, addressLabel)
	visual := container.NewMax(background, container.NewPadded(content))
	visual.Resize(fyne.NewSize(190, 40))

	// Clear and add the new visual
	vg.dragContainer.Objects = []fyne.CanvasObject{visual}
	vg.dragContainer.Refresh()
}

// UpdateDrag updates the drag position and highlights valid drop zones
func (vg *VehicleGrid) UpdateDrag(newPos fyne.Position) {
	if !vg.isDragging {
		return
	}

	// Move the drag visual to follow the mouse
	vg.dragContainer.Move(newPos)

	// Calculate which tile is under the mouse
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
	vg.dragContainer.Hide()
	vg.dropHighlight.Hide()
	vg.dragContainer.Objects = nil

	// Refresh the entire grid
	vg.refreshAfterMove()
}

// positionToTile converts screen coordinates to vehicle/tile position
func (vg *VehicleGrid) positionToTile(pos fyne.Position) VehiclePosition {
	// Account for scroll offset
	scrollOffset := vg.scrollContainer.Offset

	// Adjust position by scroll offset
	adjustedPos := fyne.NewPos(pos.X+scrollOffset.X, pos.Y+scrollOffset.Y)

	// Check each vehicle's tiles
	for vIndex, vehicle := range vg.vehicles {
		if vehicle.card == nil {
			continue
		}

		// Get the card's position in the grid
		cardPos := vehicle.card.Position()
		cardSize := vehicle.card.Size()

		// Check if we're within this card
		if adjustedPos.X >= cardPos.X && adjustedPos.X <= cardPos.X+cardSize.Width &&
			adjustedPos.Y >= cardPos.Y && adjustedPos.Y <= cardPos.Y+cardSize.Height {

			// Now check which tile within the card
			for tIndex, tile := range vehicle.tiles {
				if tile.container == nil {
					continue
				}

				// Get tile position relative to card
				tilePos := tile.container.Position()
				tileSize := tile.container.Size()

				// Calculate absolute tile position
				absTileX := cardPos.X + tilePos.X
				absTileY := cardPos.Y + tilePos.Y

				if adjustedPos.X >= absTileX && adjustedPos.X <= absTileX+tileSize.Width &&
					adjustedPos.Y >= absTileY && adjustedPos.Y <= absTileY+tileSize.Height {
					return VehiclePosition{VehicleIndex: vIndex, TileIndex: tIndex}
				}
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

	if target.TileIndex < 0 || target.TileIndex >= len(vg.vehicles[target.VehicleIndex].tiles) {
		return false
	}

	vehicle := vg.vehicles[target.VehicleIndex]

	// Check if tile is empty and vehicle has capacity
	return vehicle.IsTileEmpty(target.TileIndex) &&
		vehicle.HasCapacityForGuest(vg.draggedGuest)
}

// highlightDropTarget shows visual feedback for valid drop zones
func (vg *VehicleGrid) highlightDropTarget(target VehiclePosition) {
	if target.VehicleIndex < 0 || target.VehicleIndex >= len(vg.vehicles) {
		vg.dropHighlight.Hide()
		return
	}

	vehicle := vg.vehicles[target.VehicleIndex]
	if vehicle.card == nil || target.TileIndex >= len(vehicle.tiles) {
		vg.dropHighlight.Hide()
		return
	}

	tile := vehicle.tiles[target.TileIndex]
	if tile.container == nil {
		vg.dropHighlight.Hide()
		return
	}

	// Get absolute position of the tile
	cardPos := vehicle.card.Position()
	tilePos := tile.container.Position()
	tileSize := tile.container.Size()

	// Account for scroll offset
	scrollOffset := vg.scrollContainer.Offset

	absolutePos := fyne.NewPos(
		cardPos.X+tilePos.X-scrollOffset.X,
		cardPos.Y+tilePos.Y-scrollOffset.Y,
	)

	vg.dropHighlight.Move(absolutePos)
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
		sourceVehicle.SeatsRemaining += vg.draggedGuest.GroupSize
	}

	// Add guest to target vehicle at the specific tile position
	targetVehicle := &vg.routeManager.Vehicles[to.VehicleIndex]

	// Insert at the tile position (but don't exceed current guest count)
	insertPos := to.TileIndex
	if insertPos > len(targetVehicle.Guests) {
		insertPos = len(targetVehicle.Guests)
	}

	// Insert guest at the specified position
	targetVehicle.Guests = append(targetVehicle.Guests[:insertPos],
		append([]app.Guest{*vg.draggedGuest}, targetVehicle.Guests[insertPos:]...)...)

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
