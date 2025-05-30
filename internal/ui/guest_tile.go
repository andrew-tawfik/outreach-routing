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
func NewGuestTile(vehicleIndex, tileIndex int, grid *VehicleGrid) *GuestTile {
	tile := &GuestTile{
		vehicleIndex: vehicleIndex,
		tileIndex:    tileIndex,
		grid:         grid,
		isEmpty:      true,
	}

	tile.setupVisuals()
	return tile
}

// setupVisuals initializes the visual components
func (gt *GuestTile) setupVisuals() {
	// Tile background - distinct color so you can see tile boundaries
	gt.background = canvas.NewRectangle(color.NRGBA{230, 230, 230, 255}) // Light gray
	gt.background.CornerRadius = 5

	// Prominent border so tiles are clearly visible
	gt.border = canvas.NewRectangle(color.Transparent)
	gt.border.StrokeColor = color.NRGBA{150, 150, 150, 255} // Darker gray border
	gt.border.StrokeWidth = 2                               // Thicker border
	gt.border.FillColor = color.Transparent
	gt.border.CornerRadius = 5

	// Placeholder text for empty tiles
	gt.placeholder = widget.NewLabel("Drop guest here")
	gt.placeholder.Alignment = fyne.TextAlignCenter
	gt.placeholder.TextStyle = fyne.TextStyle{Italic: true}
}

// CreateTile builds the visual representation of the tile
func (gt *GuestTile) CreateTile() fyne.CanvasObject {
	// Always show the tile background and border
	gt.container = container.NewMax(
		gt.background,
		gt.border,
	)

	if gt.isEmpty || gt.guest == nil {
		// Empty tile - show placeholder
		gt.container.Add(container.NewCenter(gt.placeholder))
	} else {
		// Tile with guest - guest widget sits ON TOP of tile
		gt.guestWidget = NewGuestWidget(gt.guest, gt.vehicleIndex, gt.tileIndex, gt.grid)
		guestContainer := container.NewPadded(gt.guestWidget.CreateWidget())
		gt.container.Add(guestContainer)
	}

	// Larger tile size so you can see the tile boundaries
	gt.container.Resize(fyne.NewSize(220, 50))
	return gt.container
}

// SetGuest assigns a guest to this tile
func (gt *GuestTile) SetGuest(guest *app.Guest) {
	gt.guest = guest
	gt.isEmpty = false
	gt.isHidden = false

	// Create the guest widget
	gt.guestWidget = NewGuestWidget(guest, gt.vehicleIndex, gt.tileIndex, gt.grid)

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

	gt.container.Objects = []fyne.CanvasObject{
		gt.background,
		gt.border,
	}

	if gt.isEmpty || gt.guest == nil || gt.isHidden {
		// Show placeholder
		gt.container.Objects = append(gt.container.Objects, gt.placeholder)
	} else if gt.guestWidget != nil {
		// Show guest widget
		gt.container.Objects = append(gt.container.Objects, gt.guestWidget.CreateWidget())
	}

	gt.container.Refresh()
}

// GetPosition returns the position of this tile
func (gt *GuestTile) GetPosition() VehiclePosition {
	return VehiclePosition{
		VehicleIndex: gt.vehicleIndex,
		TileIndex:    gt.tileIndex,
	}
}

// HighlightAsDropTarget highlights this tile as a valid drop target
func (gt *GuestTile) HighlightAsDropTarget() {
	gt.background.FillColor = color.NRGBA{200, 255, 200, 255} // Light green
	gt.border.StrokeColor = color.NRGBA{0, 200, 0, 255}       // Green border
	gt.border.StrokeWidth = 3                                 // Even thicker when highlighted
	if gt.container != nil {
		gt.container.Refresh()
	}
}

// RemoveHighlight removes drop target highlighting
func (gt *GuestTile) RemoveHighlight() {
	gt.background.FillColor = color.NRGBA{230, 230, 230, 255} // Back to light gray
	gt.border.StrokeColor = color.NRGBA{150, 150, 150, 255}   // Back to dark gray
	gt.border.StrokeWidth = 2                                 // Back to normal thickness
	if gt.container != nil {
		gt.container.Refresh()
	}
}
