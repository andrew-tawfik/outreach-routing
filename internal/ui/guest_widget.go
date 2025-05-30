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
	widget.BaseWidget

	// Data
	guest        *app.Guest
	vehicleIndex int
	tileIndex    int
	grid         *VehicleGrid

	// Visual components
	background *canvas.Rectangle
	content    *fyne.Container

	// Drag state
	dragStartPos fyne.Position
}

// NewGuestWidget creates a new guest widget
func NewGuestWidget(guest *app.Guest, vehicleIndex, tileIndex int, grid *VehicleGrid) *GuestWidget {
	gw := &GuestWidget{
		guest:        guest,
		vehicleIndex: vehicleIndex,
		tileIndex:    tileIndex,
		grid:         grid,
	}

	gw.ExtendBaseWidget(gw)
	return gw
}

// CreateRenderer creates the renderer for the guest widget
func (gw *GuestWidget) CreateRenderer() fyne.WidgetRenderer {
	// Background with rounded corners
	gw.background = canvas.NewRectangle(color.NRGBA{70, 130, 180, 255}) // Steel blue
	gw.background.CornerRadius = 3

	// Create labels
	nameAndGroup := fmt.Sprintf("%s (%d)", gw.guest.Name, gw.guest.GroupSize)
	nameLabel := widget.NewLabel(nameAndGroup)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Truncate address if too long
	address := gw.guest.Address
	if len(address) > 20 {
		address = address[:17] + "..."
	}
	addressLabel := widget.NewLabel(address)

	// Content container
	textContent := container.NewVBox(nameLabel, addressLabel)
	gw.content = container.NewMax(
		gw.background,
		container.NewPadded(textContent),
	)

	return &guestWidgetRenderer{
		widget:  gw,
		objects: []fyne.CanvasObject{gw.content},
	}
}

// guestWidgetRenderer implements the renderer for GuestWidget
type guestWidgetRenderer struct {
	widget  *GuestWidget
	objects []fyne.CanvasObject
}

func (r *guestWidgetRenderer) Layout(size fyne.Size) {
	r.widget.content.Resize(size)
}

func (r *guestWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(190, 40)
}

func (r *guestWidgetRenderer) Refresh() {
	// Update background color based on state
	if r.widget.grid.isDragging &&
		r.widget.grid.draggedGuest == r.widget.guest {
		r.widget.background.FillColor = color.NRGBA{255, 100, 100, 100} // Semi-transparent during drag
	} else {
		r.widget.background.FillColor = color.NRGBA{70, 130, 180, 255} // Normal blue
	}
	r.widget.background.Refresh()
}

func (r *guestWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *guestWidgetRenderer) Destroy() {}

// Interface implementations for drag handling

// Dragged handles drag events
func (gw *GuestWidget) Dragged(ev *fyne.DragEvent) {
	if !gw.grid.isDragging {
		// Start drag
		origin := VehiclePosition{
			VehicleIndex: gw.vehicleIndex,
			TileIndex:    gw.tileIndex,
		}

		// Calculate absolute position of the drag start
		absolutePos := gw.absolutePosition()
		startPos := absolutePos.Add(ev.Position)

		gw.grid.StartDrag(gw.guest, origin, startPos)
		gw.dragStartPos = ev.Position
	} else {
		// Continue drag - update position
		currentPos := gw.absolutePosition()
		dragPos := currentPos.Add(ev.Position)
		gw.grid.UpdateDrag(dragPos)
	}
}

// DragEnd handles the end of drag operations
func (gw *GuestWidget) DragEnd() {
	if gw.grid.isDragging && gw.grid.draggedGuest == gw.guest {
		// Get final position
		finalPos := gw.absolutePosition().Add(gw.dragStartPos)
		gw.grid.EndDrag(finalPos)
	}
}

// Tapped handles tap/click events
func (gw *GuestWidget) Tapped(_ *fyne.PointEvent) {
	gw.grid.config.InfoLog.Printf("Tapped guest: %s", gw.guest.Name)
}

// TappedSecondary handles right-click events
func (gw *GuestWidget) TappedSecondary(_ *fyne.PointEvent) {
	// Could show context menu here
}

// absolutePosition calculates the absolute screen position of this widget
func (gw *GuestWidget) absolutePosition() fyne.Position {
	// Get the position of this widget
	pos := gw.Position()

	// Walk up the widget tree to get absolute position
	parent := gw.Parent()
	for parent != nil {
		if parentWidget, ok := parent.(fyne.Widget); ok {
			pos = pos.Add(parentWidget.Position())
			parent = parentWidget.(*GuestWidget).Parent()
		} else {
			break
		}
	}

	return pos
}

// Parent returns the parent widget (needed for position calculation)
func (gw *GuestWidget) Parent() fyne.Widget {
	// This would need to be set by the tile when creating the widget
	// For now, return nil - the position calculation will still work
	// but might be slightly off in complex layouts
	return nil
}

// CreateWidget creates a simple visual representation (for compatibility)
func (gw *GuestWidget) CreateWidget() fyne.CanvasObject {
	return gw
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
