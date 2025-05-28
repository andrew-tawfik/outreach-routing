package ui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

type GuestTile struct {
	widget.BaseWidget

	// Identity
	vehicleIndex int
	tileIndex    int
	grid         *VehicleGrid

	// Content
	guest       *app.Guest
	guestWidget *GuestWidget

	// Visual
	background  *canvas.Rectangle
	border      *canvas.Rectangle
	placeholder *widget.Label
}
