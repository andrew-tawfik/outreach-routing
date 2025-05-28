package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

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
	for i, c := range d.cfg.GuestContainers {
		if i == d.vehicleIndex {
			continue
		}
		if isOverlapping(d, c) {
			// Remove guest from current vehicle
			d.cfg.Rp.rm.Vehicles[d.vehicleIndex].Guests =
				removeGuest(d.cfg.Rp.rm.Vehicles[d.vehicleIndex].Guests, d.index)

			// Add guest to new vehicle
			d.cfg.Rp.rm.Vehicles[i].Guests = append(d.cfg.Rp.rm.Vehicles[i].Guests, d.Guest)

			// Rebuild the vehicle grid
			d.cfg.VehicleSection.Objects = []fyne.CanvasObject{d.cfg.createVehicleGrid()}
			d.cfg.VehicleSection.Refresh()
			return
		}
	}
}
