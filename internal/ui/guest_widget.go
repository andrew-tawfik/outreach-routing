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
	tile         *GuestTile // Parent tile reference

	// Visual components
	background *canvas.Rectangle
	content    *fyne.Container
}

// NewGuestWidget creates a new guest widget
func NewGuestWidget(guest *app.Guest, vehicleIndex, tileIndex int, grid *VehicleGrid, tile *GuestTile) *GuestWidget {
	gw := &GuestWidget{
		guest:        guest,
		vehicleIndex: vehicleIndex,
		tileIndex:    tileIndex,
		grid:         grid,
		tile:         tile,
	}

	gw.ExtendBaseWidget(gw)
	return gw
}

// CreateRenderer creates the renderer for the guest widget
func (gw *GuestWidget) CreateRenderer() fyne.WidgetRenderer {
	// Background with rounded corners
	gw.background = canvas.NewRectangle(color.NRGBA{60, 60, 70, 255})
	gw.background.CornerRadius = 3

	// Create labels

	name := gw.guest.Name
	if len(name) > 23 {
		name = name[:20] + "..."
	}
	nameAndGroup := fmt.Sprintf("%s (%d)", name, gw.guest.GroupSize)
	nameLabel := widget.NewLabel(nameAndGroup)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	nameLabel.TextStyle.Monospace = false

	// Content container
	gw.content = container.NewMax(
		gw.background,
		container.NewPadded(nameLabel),
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
	if r.widget.grid.IsDragging() &&
		r.widget.grid.GetDraggedGuest() == r.widget.guest {
		//r.widget.background.FillColor = color.NRGBA{255, 100, 100, 100} // Semi-transparent during drag
	} else {
		//r.widget.background.FillColor = color.NRGBA{70, 130, 180, 255} // Normal blue
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
	if !gw.grid.IsDragging() {
		// Start drag
		origin := VehiclePosition{
			VehicleIndex: gw.vehicleIndex,
			TileIndex:    gw.tileIndex,
		}

		offset := ev.Position
		globalMousePos := ev.AbsolutePosition

		// Start the drag with proper positions
		gw.grid.StartDrag(gw.guest, origin, globalMousePos, offset)
	}
}

// DragEnd handles the end of drag operations
func (gw *GuestWidget) DragEnd() {
	// The drag end is handled by mouse up event in the grid
}

// Tapped handles tap/click events
func (gw *GuestWidget) Tapped(_ *fyne.PointEvent) {
}

// TappedSecondary handles right-click events
func (gw *GuestWidget) TappedSecondary(_ *fyne.PointEvent) {
	// Could show context menu here
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
