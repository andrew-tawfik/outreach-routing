package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// VehicleCard represents a single vehicle with its guest tiles
type VehicleCard struct {
	widget.BaseWidget

	// Identity
	index   int
	vehicle *app.Vehicle
	grid    *VehicleGrid // Parent grid for drag coordination

	// UI Components
	tiles    []*GuestTile
	tileGrid *fyne.Container
	card     *widget.Card

	// Visual
	capacityLabel *widget.Label
}
