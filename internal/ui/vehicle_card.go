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
	tilesPerRow int
	tileSize    fyne.Size
}

// NewVehicleCard creates a new vehicle card
func NewVehicleCard(index int, vehicle *app.Vehicle, grid *VehicleGrid) *VehicleCard {
	vc := &VehicleCard{
		index:       index,
		vehicle:     vehicle,
		grid:        grid,
		tilesPerRow: 1,                     // 1 tile per row for vertical alignment
		tileSize:    fyne.NewSize(220, 50), // Larger tiles
	}

	vc.RefreshTiles()
	return vc
}

// CreateCard builds the visual card widget
func (vc *VehicleCard) CreateCard() *widget.Card {
	title := fmt.Sprintf("Vehicle %d", vc.index+1)

	// Create capacity info label
	vc.capacityInfo = widget.NewLabel(vc.getCapacityText())

	// Create the tile grid
	vc.tileGrid = vc.createTileGrid()

	// Combine capacity info and tiles
	content := container.NewVBox(
		vc.capacityInfo,
		vc.tileGrid,
	)

	vc.card = widget.NewCard(title, "", content)
	return vc.card
}

// createTileGrid builds the grid of guest tiles
func (vc *VehicleCard) createTileGrid() *fyne.Container {
	guestCount := len(vc.vehicle.Guests)

	// Always g + 2 tiles (guests + 2 empty slots)
	numTiles := guestCount + 2
	tiles := make([]fyne.CanvasObject, numTiles)

	for i := 0; i < numTiles; i++ {
		if i < len(vc.tiles) {
			tiles[i] = vc.tiles[i].CreateTile()
		} else {
			// Create tile (guest or empty)
			tile := NewGuestTile(vc.index, i, vc.grid)
			if i < guestCount {
				// This should have a guest
				tile.SetGuest(&vc.vehicle.Guests[i])
			}
			vc.tiles = append(vc.tiles, tile)
			tiles[i] = tile.CreateTile()
		}
	}

	// Create vertical layout with minimal spacing
	vbox := container.NewVBox(tiles...)
	return vbox
}

// refreshTiles rebuilds the tiles based on current vehicle state
func (vc *VehicleCard) RefreshTiles() {
	vc.tiles = make([]*GuestTile, 0)
	guestCount := len(vc.vehicle.Guests)

	// Always create g + 2 tiles
	totalTiles := guestCount + 2

	for i := 0; i < totalTiles; i++ {
		tile := NewGuestTile(vc.index, i, vc.grid)

		if i < guestCount {
			// This tile should contain a guest
			tile.SetGuest(&vc.vehicle.Guests[i])
		}
		// Otherwise it remains empty

		vc.tiles = append(vc.tiles, tile)
	}

	// Update capacity info
	if vc.capacityInfo != nil {
		vc.capacityInfo.SetText(vc.getCapacityText())
	}
}

// getCapacityText returns the capacity display text
func (vc *VehicleCard) getCapacityText() string {
	used := 4 - vc.vehicle.SeatsRemaining // assuming max 4 seats
	return fmt.Sprintf("Capacity: %d/4 seats", used)
}

// GetTilePosition returns the screen position of a specific tile
func (vc *VehicleCard) GetTilePosition(tileIndex int) fyne.Position {
	if vc.card == nil || tileIndex >= len(vc.tiles) {
		return fyne.NewPos(0, 0)
	}

	// Calculate position within the card
	cardPos := vc.card.Position()

	// Account for card header and capacity label
	headerHeight := float32(50) // Header + capacity info height

	// Since tiles are stacked vertically, calculate Y position
	tileY := cardPos.Y + headerHeight + float32(tileIndex)*55 + 5 // 50px height + 5px spacing
	tileX := cardPos.X + 10                                       // 10px left padding

	return fyne.NewPos(tileX, tileY)
}

// GetTileSize returns the size of tiles in this vehicle card
func (vc *VehicleCard) GetTileSize() fyne.Size {
	return fyne.NewSize(220, 50) // Updated to match larger tile size
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
	return vc.vehicle.SeatsRemaining >= guest.GroupSize
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
