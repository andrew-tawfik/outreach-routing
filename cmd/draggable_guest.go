package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

type DraggableGuest struct {
	widget.BaseWidget
	Guest        app.Guest
	vehicleIndex int
	index        int
	cfg          *config

	label    *widget.Label
	bg       *canvas.Rectangle
	selected bool
}

// Factory
func NewDraggableGuest(g app.Guest, vehicleIndex, index int, cfg *config) *DraggableGuest {
	d := &DraggableGuest{
		Guest:        g,
		vehicleIndex: vehicleIndex,
		index:        index,
		cfg:          cfg,
		label: widget.NewLabel(fmt.Sprintf(
			"%s (Group: %d)\n%s", g.Name, g.GroupSize, g.Address)),
		bg: canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 255, A: 255}),
	}
	d.ExtendBaseWidget(d)
	return d
}

// Tappable
func (d *DraggableGuest) Tapped(_ *fyne.PointEvent) {
	d.selected = !d.selected
	d.Refresh()
	fmt.Println("Tapped:", d.Guest.Name)
}

// Draggable
func (d *DraggableGuest) Dragged(ev *fyne.DragEvent) {
	newPos := d.Position().Add(ev.Dragged)
	d.Move(newPos)
	canvas.Refresh(d)
}

func (d *DraggableGuest) DragEnd() {
	for i, c := range d.cfg.guestContainers {
		if i == d.vehicleIndex {
			continue
		}
		if isOverlapping(d, c) {
			// Remove guest from current vehicle
			d.cfg.rp.rm.Vehicles[d.vehicleIndex].Guests =
				removeGuest(d.cfg.rp.rm.Vehicles[d.vehicleIndex].Guests, d.index)

			// Add guest to new vehicle
			d.cfg.rp.rm.Vehicles[i].Guests = append(d.cfg.rp.rm.Vehicles[i].Guests, d.Guest)

			// Rebuild the vehicle grid
			d.cfg.vehicleSection.Objects = []fyne.CanvasObject{d.cfg.createVehicleGrid()}
			d.cfg.vehicleSection.Refresh()
			return
		}
	}
}

func (d *DraggableGuest) CreateRenderer() fyne.WidgetRenderer {
	d.bg.Hide() // background shown only when selected

	objects := []fyne.CanvasObject{
		d.bg,
		d.label,
	}
	return &draggableGuestRenderer{
		guest:   d,
		bg:      d.bg,
		label:   d.label,
		objects: objects,
	}
}

type draggableGuestRenderer struct {
	guest   *DraggableGuest
	bg      *canvas.Rectangle
	label   *widget.Label
	objects []fyne.CanvasObject
}

func (r *draggableGuestRenderer) Layout(s fyne.Size) {
	r.bg.Resize(s)
	r.label.Resize(s)
}

func (r *draggableGuestRenderer) MinSize() fyne.Size {
	return r.label.MinSize()
}

func (r *draggableGuestRenderer) Refresh() {
	if r.guest.selected {
		r.bg.Show()
	} else {
		r.bg.Hide()
	}
	canvas.Refresh(r.guest)
}

func (r *draggableGuestRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *draggableGuestRenderer) Destroy() {}

func isOverlapping(a, b fyne.CanvasObject) bool {
	aPos, bPos := a.Position(), b.Position()
	aEnd, bEnd := aPos.Add(a.Size()), bPos.Add(b.Size())
	return aPos.X < bEnd.X && aEnd.X > bPos.X && aPos.Y < bEnd.Y && aEnd.Y > bPos.Y
}

func removeGuest(slice []app.Guest, index int) []app.Guest {
	return append(slice[:index], slice[index+1:]...)
}
