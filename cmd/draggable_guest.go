package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// DraggableGuest is a composite widget: it wraps a widget.Card.
type DraggableGuest struct {
	widget.BaseWidget   // register as custom widget
	Guest               app.Guest
	vehicleIndex, index int
	card                *widget.Card
}

// NewDraggableGuest creates the widget but does not yet wire up drag logic.
func NewDraggableGuest(g app.Guest, vehicleIndex, index int) *DraggableGuest {
	title := fmt.Sprintf("%s (group of %d): %s", g.Name, g.GroupSize, g.Address)
	card := widget.NewCard("", "", widget.NewLabel(title))

	dg := &DraggableGuest{
		Guest:        g,
		vehicleIndex: vehicleIndex,
		index:        index,
		card:         card,
	}
	dg.ExtendBaseWidget(dg) // tell Fyne this is a custom widget
	return dg
}

func (dg *DraggableGuest) Dragged(e *fyne.DragEvent) {
	pos := dg.Position().Add(e.Dragged)
	dg.Move(pos)

}

func (dg *DraggableGuest) DragEnd() {

}

func isOverlapping(obj1, obj2 fyne.CanvasObject) bool {

	return false
}

// CreateRenderer hands off all drawing/layout to the card’s own renderer.
func (d *DraggableGuest) CreateRenderer() fyne.WidgetRenderer {
	// grab the card’s built-in renderer
	cardR := d.card.CreateRenderer()

	// wrap it so we satisfy fyne.WidgetRenderer
	return &guestRenderer{
		guest:        d,
		card:         d.card,
		cardRenderer: cardR,
	}
}

// guestRenderer simply proxies everything to the card’s renderer.
type guestRenderer struct {
	guest        *DraggableGuest
	card         *widget.Card
	cardRenderer fyne.WidgetRenderer
}

func (r *guestRenderer) Layout(size fyne.Size) {
	r.cardRenderer.Layout(size)
}

func (r *guestRenderer) MinSize() fyne.Size {
	return r.cardRenderer.MinSize()
}

func (r *guestRenderer) Refresh() {
	r.cardRenderer.Refresh()
}

func (r *guestRenderer) Objects() []fyne.CanvasObject {
	return r.cardRenderer.Objects()
}

func (r *guestRenderer) Destroy() {
	r.cardRenderer.Destroy()
}
