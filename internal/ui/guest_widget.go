package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

type GuestWidget struct {
	widget.BaseWidget

	guestData   app.Guest
	isDragging  bool
	dragManager *dragManager
}

type DraggableGuest struct {
	widget.BaseWidget
	Guest        app.Guest
	vehicleIndex int
	index        int
	cfg          *Config

	label    *widget.Label
	bg       *canvas.Rectangle
	selected bool
}

// Factory
func NewDraggableGuest(g app.Guest, vehicleIndex, index int, cfg *Config) *DraggableGuest {
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

func (d *DraggableGuest) CreateRenderer() fyne.WidgetRenderer {
	d.bg.Hide()

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
