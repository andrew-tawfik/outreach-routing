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

// GuestWidget represents a draggable guest within a tile
type GuestWidget struct {
	// Data
	guest        *app.Guest
	vehicleIndex int
	tileIndex    int
	grid         *VehicleGrid

	// Visual components
	nameLabel    *widget.Label
	groupLabel   *widget.Label
	addressLabel *widget.Label
	background   *canvas.Rectangle
	container    *fyne.Container

	// State
	isDragging bool
	isSelected bool
}

// NewGuestWidget creates a new guest widget
func NewGuestWidget(guest *app.Guest, vehicleIndex, tileIndex int, grid *VehicleGrid) *GuestWidget {
	gw := &GuestWidget{
		guest:        guest,
		vehicleIndex: vehicleIndex,
		tileIndex:    tileIndex,
		grid:         grid,
	}

	gw.setupVisuals()
	return gw
}

// setupVisuals initializes the visual components
func (gw *GuestWidget) setupVisuals() {
	// Background with rounded corners and nice color
	gw.background = canvas.NewRectangle(color.NRGBA{70, 130, 180, 255}) // Steel blue
	gw.background.CornerRadius = 3

	// These labels are used for state management, but actual display is handled in CreateWidget
	gw.nameLabel = widget.NewLabel(gw.guest.Name)
	gw.nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	gw.groupLabel = widget.NewLabel(fmt.Sprintf("%d", gw.guest.GroupSize))

	gw.addressLabel = widget.NewLabel(gw.guest.Address)
}

// CreateWidget builds the visual representation
func (gw *GuestWidget) CreateWidget() fyne.CanvasObject {
	// Compact guest widget that sits ON TOP of tile
	nameAndGroup := fmt.Sprintf("%s (%d)", gw.guest.Name, gw.guest.GroupSize)

	nameLabel := widget.NewLabel(nameAndGroup)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Truncate address to fit in smaller guest widget
	address := gw.guest.Address
	if len(address) > 20 {
		address = address[:17] + "..."
	}
	addressLabel := widget.NewLabel(address)

	// Vertical layout for guest info
	content := container.NewVBox(
		nameLabel,
		addressLabel,
	)

	gw.container = container.NewMax(
		gw.background,
		container.NewPadded(content),
	)

	// Make it interactive with drag overlay
	overlay := gw.createInteractiveOverlay()
	gw.container = container.NewWithoutLayout(
		gw.background,
		container.NewPadded(content),
		overlay,
	)

	// Guest widget is smaller than tile, so you can see tile edges
	gw.container.Resize(fyne.NewSize(190, 40))
	return gw.container
}

// createInteractiveOverlay creates an invisible overlay for handling interactions
func (gw *GuestWidget) createInteractiveOverlay() fyne.CanvasObject {
	overlay := canvas.NewRectangle(color.Transparent)

	// Create a custom draggable widget
	draggable := &DraggableGuestOverlay{
		guestWidget: gw,
	}

	return container.NewMax(overlay, draggable)
}

// DraggableGuestOverlay handles drag interactions
type DraggableGuestOverlay struct {
	widget.BaseWidget
	guestWidget *GuestWidget
}

// Tapped handles tap events
func (dgo *DraggableGuestOverlay) Tapped(_ *fyne.PointEvent) {
	gw := dgo.guestWidget
	gw.isSelected = !gw.isSelected
	gw.updateVisualState()
	gw.grid.config.InfoLog.Printf("Tapped guest: %s", gw.guest.Name)
}

// Dragged handles drag events
func (dgo *DraggableGuestOverlay) Dragged(ev *fyne.DragEvent) {
	gw := dgo.guestWidget
	if !gw.isDragging {
		// Start drag operation
		gw.isDragging = true
		startPos := gw.container.Position().Add(ev.Position)
		origin := VehiclePosition{
			VehicleIndex: gw.vehicleIndex,
			TileIndex:    gw.tileIndex,
		}
		gw.grid.StartDrag(gw.guest, origin, startPos)
	} else {
		// Update drag position
		newPos := gw.container.Position().Add(ev.Position)
		gw.grid.UpdateDrag(newPos)
	}
}

// DragEnd handles the end of drag operations
func (dgo *DraggableGuestOverlay) DragEnd() {
	gw := dgo.guestWidget
	if gw.isDragging {
		gw.isDragging = false
		endPos := gw.container.Position()
		gw.grid.EndDrag(endPos)
	}
}

// CreateRenderer creates the renderer for the overlay
func (dgo *DraggableGuestOverlay) CreateRenderer() fyne.WidgetRenderer {
	return &draggableGuestOverlayRenderer{
		overlay: dgo,
		objects: []fyne.CanvasObject{},
	}
}

type draggableGuestOverlayRenderer struct {
	overlay *DraggableGuestOverlay
	objects []fyne.CanvasObject
}

func (r *draggableGuestOverlayRenderer) Layout(size fyne.Size) {
	// The overlay fills the entire guest widget area
}

func (r *draggableGuestOverlayRenderer) MinSize() fyne.Size {
	return fyne.NewSize(190, 40) // Match the guest widget size (smaller than tile)
}

func (r *draggableGuestOverlayRenderer) Refresh() {
	// Nothing to refresh for the invisible overlay
}

func (r *draggableGuestOverlayRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *draggableGuestOverlayRenderer) Destroy() {}

// updateVisualState updates the guest's visual appearance based on state
func (gw *GuestWidget) updateVisualState() {
	if gw.isSelected {
		gw.background.FillColor = color.NRGBA{255, 165, 0, 255} // Orange when selected
	} else if gw.isDragging {
		gw.background.FillColor = color.NRGBA{255, 100, 100, 255} // Red when dragging
	} else {
		gw.background.FillColor = color.NRGBA{70, 130, 180, 255} // Default steel blue
	}

	if gw.container != nil {
		gw.container.Refresh()
	}
}

// Hide temporarily hides the guest widget (during drag)
func (gw *GuestWidget) Hide() {
	if gw.container != nil {
		gw.container.Hide()
	}
}

// Show shows the guest widget again
func (gw *GuestWidget) Show() {
	if gw.container != nil {
		gw.container.Show()
	}
}

// GetGuest returns the guest data
func (gw *GuestWidget) GetGuest() *app.Guest {
	return gw.guest
}

// GetPosition returns the position of this guest widget
func (gw *GuestWidget) GetPosition() VehiclePosition {
	return VehiclePosition{
		VehicleIndex: gw.vehicleIndex,
		TileIndex:    gw.tileIndex,
	}
}
