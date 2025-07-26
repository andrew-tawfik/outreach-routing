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


type GuestWidget struct {
	widget.BaseWidget

	
	guest        *app.Guest
	vehicleIndex int
	tileIndex    int
	grid         *VehicleGrid
	tile         *GuestTile 

	
	background *canvas.Rectangle
	content    *fyne.Container
}


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


func (gw *GuestWidget) CreateRenderer() fyne.WidgetRenderer {
	
	gw.background = canvas.NewRectangle(color.NRGBA{60, 60, 70, 255})
	gw.background.CornerRadius = 3

	

	name := gw.guest.Name
	if len(name) > 23 {
		name = name[:20] + "..."
	}
	nameAndGroup := fmt.Sprintf("%s (%d)", name, gw.guest.GroupSize)
	nameLabel := widget.NewLabel(nameAndGroup)
	nameLabel.TextStyle = fyne.TextStyle{Bold: false}
	nameLabel.TextStyle.Monospace = false

	
	gw.content = container.NewMax(
		gw.background,
		container.NewPadded(nameLabel),
	)

	return &guestWidgetRenderer{
		widget:  gw,
		objects: []fyne.CanvasObject{gw.content},
	}
}


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
	
	if r.widget.grid.IsDragging() &&
		r.widget.grid.GetDraggedGuest() == r.widget.guest {
		
	} else {
		
	}
	r.widget.background.Refresh()
}

func (r *guestWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *guestWidgetRenderer) Destroy() {}




func (gw *GuestWidget) Dragged(ev *fyne.DragEvent) {
	if !gw.grid.IsDragging() {
		
		origin := VehiclePosition{
			VehicleIndex: gw.vehicleIndex,
			TileIndex:    gw.tileIndex,
		}

		offset := ev.Position
		globalMousePos := ev.AbsolutePosition

		
		gw.grid.StartDrag(gw.guest, origin, globalMousePos, offset)
	}
}


func (gw *GuestWidget) DragEnd() {
	
}


func (gw *GuestWidget) Tapped(_ *fyne.PointEvent) {
}


func (gw *GuestWidget) TappedSecondary(_ *fyne.PointEvent) {
	
}


func (gw *GuestWidget) GetGuest() *app.Guest {
	return gw.guest
}


func (gw *GuestWidget) GetPosition() VehiclePosition {
	return VehiclePosition{
		VehicleIndex: gw.vehicleIndex,
		TileIndex:    gw.tileIndex,
	}
}
