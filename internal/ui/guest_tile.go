package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// GuestTile represents a single tile that can contain a guest
type GuestTile struct {
	// Identity
	vehicleIndex int
	tileIndex    int
	grid         *VehicleGrid
	card         *VehicleCard // Reference to parent card

	// Content
	guest       *app.Guest
	guestWidget *GuestWidget

	// Visual components
	background  *canvas.Rectangle
	border      *canvas.Rectangle
	placeholder *widget.Label
	container   *fyne.Container

	// State
	isEmpty  bool
	isHidden bool
}

// NewGuestTile creates a new guest tile
func NewGuestTile(vehicleIndex, tileIndex int, grid *VehicleGrid, card *VehicleCard) *GuestTile {
	tile := &GuestTile{
		vehicleIndex: vehicleIndex,
		tileIndex:    tileIndex,
		grid:         grid,
		card:         card,
		isEmpty:      true,
	}

	tile.setupVisuals()
	return tile
}

// setupVisuals initializes the visual components
func (gt *GuestTile) setupVisuals() {
	// Tile background - distinct color so you can see tile boundaries
	gt.background = canvas.NewRectangle(color.NRGBA{70, 70, 80, 255})
	gt.background.CornerRadius = 5

	// Prominent border so tiles are clearly visible
	gt.border = canvas.NewRectangle(color.NRGBA{80, 80, 90, 255})
	gt.border.StrokeColor = color.NRGBA{80, 80, 90, 255}
	gt.border.StrokeWidth = 2
	gt.border.FillColor = color.NRGBA{80, 80, 90, 255}
	gt.border.CornerRadius = 5

	// Placeholder text for empty tiles
	gt.placeholder = widget.NewLabel("")
	gt.placeholder.Alignment = fyne.TextAlignCenter
	gt.placeholder.TextStyle = fyne.TextStyle{Italic: true}
}

// CreateTile builds the visual representation of the tile
func (gt *GuestTile) CreateTile() fyne.CanvasObject {
	// Create base container with background and border
	baseObjects := []fyne.CanvasObject{gt.background, gt.border}

	if gt.isEmpty || gt.guest == nil || gt.isHidden {
		// Empty tile - show placeholder
		placeholderContainer := container.NewCenter(gt.placeholder)
		baseObjects = append(baseObjects, placeholderContainer)
	} else {
		// Tile with guest - create and add guest widget
		gt.guestWidget = NewGuestWidget(gt.guest, gt.vehicleIndex, gt.tileIndex, gt.grid, gt)
		// Center the guest widget in the tile with some padding
		guestContainer := container.NewPadded(gt.guestWidget)
		baseObjects = append(baseObjects, guestContainer)
	}

	gt.container = container.NewMax(baseObjects...)
	gt.container.Resize(fyne.NewSize(220, 70))
	return gt.container
}

// SetGuest assigns a guest to this tile
func (gt *GuestTile) SetGuest(guest *app.Guest) {
	gt.guest = guest
	gt.isEmpty = false
	gt.isHidden = false

	if gt.guestWidget != nil {
		gt.guestWidget.vehicleIndex = gt.vehicleIndex
		gt.guestWidget.tileIndex = gt.tileIndex
	}

	// Update the container
	gt.refreshContainer()
}

// ClearGuest removes the guest from this tile
func (gt *GuestTile) ClearGuest() {
	gt.guest = nil
	gt.guestWidget = nil
	gt.isEmpty = true
	gt.isHidden = false

	gt.refreshContainer()
}

// IsEmpty returns whether the tile is empty
func (gt *GuestTile) IsEmpty() bool {
	return gt.isEmpty || gt.guest == nil
}

// HideGuest temporarily hides the guest (during drag operations)
func (gt *GuestTile) HideGuest() {
	gt.isHidden = true
	gt.refreshContainer()
}

// ShowGuest shows the guest again (after failed drag)
func (gt *GuestTile) ShowGuest() {
	gt.isHidden = false
	gt.refreshContainer()
}

// refreshContainer updates the tile's visual content
func (gt *GuestTile) refreshContainer() {
	if gt.container == nil {
		return
	}

	// Rebuild objects array
	objects := []fyne.CanvasObject{gt.background, gt.border}

	if gt.isEmpty || gt.guest == nil || gt.isHidden {
		// Show placeholder for empty or hidden tiles
		placeholderContainer := container.NewCenter(gt.placeholder)
		objects = append(objects, placeholderContainer)
	} else {
		// Show guest widget
		if gt.guestWidget == nil {
			gt.guestWidget = NewGuestWidget(gt.guest, gt.vehicleIndex, gt.tileIndex, gt.grid, gt)
		}
		guestContainer := container.NewPadded(gt.guestWidget)
		objects = append(objects, guestContainer)
	}

	// Update container objects
	gt.container.Objects = objects
	gt.container.Refresh()
}

// GetPosition returns the position of this tile
func (gt *GuestTile) GetPosition() VehiclePosition {
	return VehiclePosition{
		VehicleIndex: gt.vehicleIndex,
		TileIndex:    gt.tileIndex,
	}
}

// ContainsPosition checks if a given position is within this tile
func (gt *GuestTile) ContainsPosition(pos fyne.Position) bool {
	if gt.container == nil || gt.card == nil || gt.card.card == nil {
		return false
	}

	// Get the card's position
	cardPos := gt.card.card.Position()

	// Get the tile's position relative to the card
	tilePos := gt.container.Position()
	tileSize := gt.container.Size()

	// Calculate absolute position
	absX := cardPos.X + tilePos.X
	absY := cardPos.Y + tilePos.Y

	// Check if position is within bounds
	return pos.X >= absX && pos.X <= absX+tileSize.Width &&
		pos.Y >= absY && pos.Y <= absY+tileSize.Height
}

// HighlightAsDropTarget highlights this tile as a valid drop target
func (gt *GuestTile) HighlightAsDropTarget() {
	gt.background.FillColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} // Almost black-green
	gt.border.StrokeColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255}   // Darker green border
	gt.border.StrokeWidth = 3
	gt.background.Refresh()
	gt.border.Refresh()
}

// RemoveHighlight removes drop target highlighting
func (gt *GuestTile) RemoveHighlight() {
	gt.background.FillColor = color.NRGBA{70, 70, 80, 255} // Back to light gray
	gt.border.StrokeColor = color.NRGBA{80, 80, 90, 255}   // Back to light border
	gt.border.StrokeWidth = 2                              // Back to normal thickness
	gt.background.Refresh()
	gt.border.Refresh()
}

// GetGuest returns the guest in this tile
func (gt *GuestTile) GetGuest() *app.Guest {
	return gt.guest
}
