package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

// VehicleCard represents a single vehicle with its guest tiles
type VehicleCard struct {
	// Identity
	index   int
	vehicle *app.Vehicle
	grid    *VehicleGrid // Parent grid for drag coordination

	// UI Components
	tiles        []*GuestTile
	tileGrid     *fyne.Container
	card         *widget.Card
	capacityInfo *widget.Label

	// Layout info
	tileSize fyne.Size
}

// NewVehicleCard creates a new vehicle card
func NewVehicleCard(index int, vehicle *app.Vehicle, grid *VehicleGrid) *VehicleCard {
	vc := &VehicleCard{
		index:    index,
		vehicle:  vehicle,
		grid:     grid,
		tileSize: fyne.NewSize(220, 60), // Consistent tile size
	}

	vc.createTiles()
	return vc
}

// createTiles creates all tiles for this vehicle
func (vc *VehicleCard) createTiles() {
	guestCount := len(vc.vehicle.Guests)
	totalTiles := guestCount + 2 // Always g + 2 tiles

	vc.tiles = make([]*GuestTile, totalTiles)

	for i := 0; i < totalTiles; i++ {
		tile := NewGuestTile(vc.index, i, vc.grid, vc)

		if i < guestCount {
			// This tile should contain a guest
			tile.SetGuest(&vc.vehicle.Guests[i])
		}
		// Otherwise it remains empty

		vc.tiles[i] = tile
	}
}

// CreateCard builds the visual card widget
func (vc *VehicleCard) CreateCard() *widget.Card {
	title := fmt.Sprintf("Vehicle %d", vc.index+1)

	// Create capacity info label
	vc.capacityInfo = widget.NewLabel(vc.getCapacityText())
	vc.capacityInfo.TextStyle = fyne.TextStyle{Bold: true}

	// Create the tile grid
	vc.tileGrid = vc.createTileGrid()

	vc.card = widget.NewCard(title, "", vc.tileGrid)
	return vc.card
}

// createTileGrid builds the grid of guest tiles
func (vc *VehicleCard) createTileGrid() *fyne.Container {
	tiles := make([]fyne.CanvasObject, len(vc.tiles))

	for i, tile := range vc.tiles {
		tiles[i] = tile.CreateTile()
	}

	// Create vertical layout with consistent spacing
	vbox := container.NewVBox(tiles...)
	return vbox
}

// RefreshTiles rebuilds the tiles based on current vehicle state
func (vc *VehicleCard) RefreshTiles() {
	// Recreate tiles with current state
	vc.createTiles()

	// Update capacity info
	if vc.capacityInfo != nil {
		vc.capacityInfo.SetText(vc.getCapacityText())
	}

	// Rebuild the tile grid if it exists
	if vc.tileGrid != nil && vc.card != nil {
		// Clear existing tiles
		vc.tileGrid.Objects = nil

		// Add new tiles
		for _, tile := range vc.tiles {
			vc.tileGrid.Objects = append(vc.tileGrid.Objects, tile.CreateTile())
		}

		vc.tileGrid.Refresh()
		vc.card.Refresh()
	}
}

// getCapacityText returns the capacity display text
func (vc *VehicleCard) getCapacityText() string {
	used := 4 - vc.vehicle.SeatsRemaining // assuming max 4 seats
	return fmt.Sprintf("Capacity: %d/4 seats", used)
}

// IsTileEmpty checks if a specific tile is empty
func (vc *VehicleCard) IsTileEmpty(tileIndex int) bool {
	if tileIndex >= len(vc.tiles) {
		return false
	}
	return vc.tiles[tileIndex].IsEmpty()
}

// HasCapacityForGuest checks if the vehicle can accommodate the guest
func (vc *VehicleCard) HasCapacityForGuest(guest *app.Guest) bool {
	return true
}

// HideGuest hides a guest from a specific tile (during drag)
func (vc *VehicleCard) HideGuest(tileIndex int) {
	if tileIndex < len(vc.tiles) {
		vc.tiles[tileIndex].HideGuest()
	}
}

// ShowGuest shows a guest in a specific tile (after failed drag)
func (vc *VehicleCard) ShowGuest(tileIndex int) {
	if tileIndex < len(vc.tiles) {
		vc.tiles[tileIndex].ShowGuest()
	}
}

// RemoveAllHighlights removes highlights from all tiles
func (vc *VehicleCard) RemoveAllHighlights() {
	for _, tile := range vc.tiles {
		tile.RemoveHighlight()
	}
}

// GetTile returns a specific tile
func (vc *VehicleCard) GetTile(index int) *GuestTile {
	if index >= 0 && index < len(vc.tiles) {
		return vc.tiles[index]
	}
	return nil
}

// GetTilePosition returns the screen position of a specific tile
func (vc *VehicleCard) GetTilePosition(tileIndex int) fyne.Position {
	if vc.card == nil || tileIndex >= len(vc.tiles) {
		return fyne.NewPos(0, 0)
	}

	// Get the card's position
	cardPos := vc.card.Position()

	// Account for card header and capacity label
	headerHeight := float32(70) // Card title + capacity info + separator

	// Since tiles are stacked vertically, calculate Y position
	tileY := cardPos.Y + headerHeight + float32(tileIndex)*65 // 60px height + 5px spacing
	tileX := cardPos.X + 10                                   // 10px left padding from card edge

	return fyne.NewPos(tileX, tileY)
}

// GetTileSize returns the size of tiles in this vehicle card
func (vc *VehicleCard) GetTileSize() fyne.Size {
	return vc.tileSize // This is already set to fyne.NewSize(220, 60) in the struct
}
