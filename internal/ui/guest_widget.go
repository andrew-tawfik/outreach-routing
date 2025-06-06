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
	gw.background = canvas.NewRectangle(color.NRGBA{70, 130, 180, 255}) // Steel blue
	gw.background.CornerRadius = 3

	// Create labels
	nameAndGroup := fmt.Sprintf("%s (%d)", gw.guest.Name, gw.guest.GroupSize)
	nameLabel := widget.NewLabel(nameAndGroup)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	nameLabel.TextStyle.Monospace = false

	// Truncate address if too long
	address := gw.guest.Address
	if len(address) > 25 {
		address = address[:22] + "..."
	}
	addressLabel := widget.NewLabel(address)
	addressLabel.TextStyle.Monospace = false

	// Set text color to white for better contrast
	nameLabel.Importance = widget.HighImportance
	addressLabel.Importance = widget.MediumImportance

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
	if r.widget.grid.IsDragging() &&
		r.widget.grid.GetDraggedGuest() == r.widget.guest {
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
	if !gw.grid.IsDragging() {
		// Start drag
		origin := VehiclePosition{
			VehicleIndex: gw.vehicleIndex,
			TileIndex:    gw.tileIndex,
		}

		offset := ev.Position
		globalMousePos := ev.AbsolutePosition

		// Log positions
		widgetPos := gw.Position()
		widgetSize := gw.Size()

		gw.grid.config.InfoLog.Printf("=== DRAG START ===")
		gw.grid.config.InfoLog.Printf("Guest: %s", gw.guest.Name)
		gw.grid.config.InfoLog.Printf("Widget Position: X=%.2f, Y=%.2f", widgetPos.X, widgetPos.Y)
		gw.grid.config.InfoLog.Printf("Widget Size: W=%.2f, H=%.2f", widgetSize.Width, widgetSize.Height)
		gw.grid.config.InfoLog.Printf("Mouse Click Offset (relative to widget): X=%.2f, Y=%.2f", offset.X, offset.Y)
		gw.grid.config.InfoLog.Printf("Mouse Absolute Position: X=%.2f, Y=%.2f", globalMousePos.X, globalMousePos.Y)
		gw.grid.config.InfoLog.Printf("=================")

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
	gw.grid.config.InfoLog.Printf("Tapped guest: %s", gw.guest.Name)
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
