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

func (cfg *Config) createVehicleGrid() fyne.CanvasObject {
	cards := make([]fyne.CanvasObject, 0, len(cfg.Rp.rm.Vehicles))
	for i, v := range cfg.Rp.rm.Vehicles {
		cards = append(cards, cfg.makeVehicleCard(v, i))
	}
	grid := container.NewGridWithColumns(4, cards...)
	return container.NewScroll(grid)
}

func (cfg *Config) makeVehicleCard(v app.Vehicle, idx int) fyne.CanvasObject {
	cfg.InfoLog.Printf("Vehicle %d has %d guests", idx+1, len(v.Guests))
	title := fmt.Sprintf("Vehicle %d", idx+1)
	guestsBox := container.NewWithoutLayout()
	cfg.GuestContainers = append(cfg.GuestContainers, guestsBox)
	for i, g := range v.Guests {
		dg := NewDraggableGuest(g, idx, i, cfg)
		dg.Move(fyne.NewPos(0, float32(i*50))) // manually space out guests vertically
		guestsBox.Add(dg)
	}

	scroll := container.NewVScroll(guestsBox)
	scroll.SetMinSize(fyne.NewSize(0, 250))

	return widget.NewCard(title, "", scroll)
}

type VehicleGrid struct {
	widget.BaseWidget

	// Core data
	vehicles     []*VehicleCard
	routeManager *app.RouteManager

	// Drag state
	dragOverlay  *fyne.Container // Global overlay for dragging
	draggedGuest *GuestWidget    // Currently dragged guest
	dragOrigin   VehiclePosition // Where drag started

	// Visual state
	dropHighlight *canvas.Rectangle // Shows valid drop zones
	gridContainer *fyne.Container   // The actual grid of vehicle cards
}

// VehiclePosition identifies a specific location in the grid
type VehiclePosition struct {
	VehicleIndex int
	TileIndex    int
}

// NewVehicleGrid creates the main vehicle grid widget
func NewVehicleGrid(rm *app.RouteManager) *VehicleGrid {
	vg := &VehicleGrid{
		routeManager:  rm,
		vehicles:      make([]*VehicleCard, 0),
		dropHighlight: canvas.NewRectangle(color.NRGBA{0, 255, 0, 64}),
	}

	vg.dropHighlight.Hide()
	//vg.ExtendBaseWidget(vg)

	// Create the drag overlay container
	vg.dragOverlay = container.NewWithoutLayout()

	// Initialize vehicles from route manager
	vg.refreshVehicles()

	return vg
}

// refreshVehicles rebuilds the vehicle cards from current state
func (vg *VehicleGrid) refreshVehicles() {
	vg.vehicles = make([]*VehicleCard, len(vg.routeManager.Vehicles))

	for i, vehicle := range vg.routeManager.Vehicles {
		card := NewVehicleCard(i, &vehicle, vg)
		vg.vehicles[i] = card
	}
}

func NewVehicleCard(index int, vehicle *app.Vehicle, grid *VehicleGrid) *VehicleCard {}
