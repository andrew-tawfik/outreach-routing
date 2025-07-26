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

var extraTiles = 1


type VehicleCard struct {
	
	index   int
	vehicle *app.Vehicle
	grid    *VehicleGrid 

	
	tiles        []*GuestTile
	tileGrid     *fyne.Container
	card         fyne.CanvasObject
	capacityInfo *widget.Label

	
	tileSize fyne.Size
}


func NewVehicleCard(index int, vehicle *app.Vehicle, grid *VehicleGrid) *VehicleCard {
	vc := &VehicleCard{
		index:    index,
		vehicle:  vehicle,
		grid:     grid,
		tileSize: fyne.NewSize(220, 60), 
	}

	vc.createTiles()
	return vc
}


func (vc *VehicleCard) createTiles() {
	guestCount := len(vc.vehicle.Guests)
	totalTiles := guestCount + extraTiles

	vc.tiles = make([]*GuestTile, totalTiles)

	for i := 0; i < totalTiles; i++ {
		tile := NewGuestTile(vc.index, i, vc.grid, vc)

		if i < guestCount {
			
			tile.SetGuest(&vc.vehicle.Guests[i])
		}
		

		vc.tiles[i] = tile
	}
}


func (vc *VehicleCard) CreateCard() fyne.CanvasObject {
	title := fmt.Sprintf("Vehicle %d", vc.index+1)

	
	background := canvas.NewRectangle(color.NRGBA{40, 40, 45, 255})
	background.CornerRadius = 12

	
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	
	vc.tileGrid = vc.createTileGrid()

	
	content := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		vc.tileGrid,
	)

	
	paddedContent := container.NewPadded(content)

	
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(5, 5)) 

	cardWithSpacing := container.NewBorder(
		spacer, 
		spacer, 
		spacer, 
		spacer, 
		container.NewMax(background, paddedContent), 
	)

	vc.card = cardWithSpacing
	return vc.card
}


func (vc *VehicleCard) createTileGrid() *fyne.Container {
	tiles := make([]fyne.CanvasObject, len(vc.tiles))

	for i, tile := range vc.tiles {
		tiles[i] = tile.CreateTile()
	}

	
	vbox := container.NewVBox(tiles...)
	return vbox
}


func (vc *VehicleCard) RefreshTiles() {
	
	vc.createTiles()

	
	if vc.capacityInfo != nil {
		vc.capacityInfo.SetText(vc.getCapacityText())
	}

	
	if vc.tileGrid != nil && vc.card != nil {
		
		vc.tileGrid.Objects = nil

		
		for _, tile := range vc.tiles {
			vc.tileGrid.Objects = append(vc.tileGrid.Objects, tile.CreateTile())
		}

		vc.tileGrid.Refresh()
		vc.card.Refresh()
	}
}


func (vc *VehicleCard) getCapacityText() string {
	used := 4 - vc.vehicle.SeatsRemaining 
	return fmt.Sprintf("Capacity: %d/4 seats", used)
}


func (vc *VehicleCard) IsTileEmpty(tileIndex int) bool {
	if tileIndex >= len(vc.tiles) {
		return false
	}
	return vc.tiles[tileIndex].IsEmpty()
}


func (vc *VehicleCard) HasCapacityForGuest(guest *app.Guest) bool {
	return true
}


func (vc *VehicleCard) HideGuest(tileIndex int) {
	if tileIndex < len(vc.tiles) {
		vc.tiles[tileIndex].HideGuest()
	}
}


func (vc *VehicleCard) ShowGuest(tileIndex int) {
	if tileIndex < len(vc.tiles) {
		vc.tiles[tileIndex].ShowGuest()
	}
}


func (vc *VehicleCard) RemoveAllHighlights() {
	for _, tile := range vc.tiles {
		tile.RemoveHighlight()
	}
}


func (vc *VehicleCard) GetTile(index int) *GuestTile {
	if index >= 0 && index < len(vc.tiles) {
		return vc.tiles[index]
	}
	return nil
}


func (vc *VehicleCard) GetTilePosition(tileIndex int) fyne.Position {
	if vc.card == nil || tileIndex >= len(vc.tiles) {
		return fyne.NewPos(0, 0)
	}

	
	cardPos := vc.card.Position()

	
	headerHeight := float32(70) 

	
	tileY := cardPos.Y + headerHeight + float32(tileIndex)*65 
	tileX := cardPos.X + 10                                   

	return fyne.NewPos(tileX, tileY)
}


func (vc *VehicleCard) GetTileSize() fyne.Size {
	return vc.tileSize 
}
