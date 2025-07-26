package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)


type GuestTile struct {
	
	vehicleIndex int
	tileIndex    int
	grid         *VehicleGrid
	card         *VehicleCard 

	
	guest       *app.Guest
	guestWidget *GuestWidget

	
	background  *canvas.Rectangle
	border      *canvas.Rectangle
	placeholder *widget.Label
	container   *fyne.Container

	
	isEmpty  bool
	isHidden bool
}


func NewGuestTile(vehicleIndex, tileIndex int, grid *VehicleGrid, card *VehicleCard) *GuestTile {
	tile := &GuestTile{
		vehicleIndex: vehicleIndex,
		tileIndex:    tileIndex,
		grid:         grid,
		card:         card,
		isEmpty:      true,
	}

	tile.setupVisuals()
	return tile
}


func (gt *GuestTile) setupVisuals() {
	
	gt.background = canvas.NewRectangle(color.NRGBA{70, 70, 80, 255})
	gt.background.CornerRadius = 5

	
	gt.border = canvas.NewRectangle(color.NRGBA{80, 80, 90, 255})
	gt.border.StrokeColor = color.NRGBA{80, 80, 90, 255}
	gt.border.StrokeWidth = 2
	gt.border.FillColor = color.NRGBA{80, 80, 90, 255}
	gt.border.CornerRadius = 5

	
	gt.placeholder = widget.NewLabel("")
	gt.placeholder.Alignment = fyne.TextAlignCenter
	gt.placeholder.TextStyle = fyne.TextStyle{Italic: true}
}


func (gt *GuestTile) CreateTile() fyne.CanvasObject {
	
	baseObjects := []fyne.CanvasObject{gt.background, gt.border}

	if gt.isEmpty || gt.guest == nil || gt.isHidden {
		
		placeholderContainer := container.NewCenter(gt.placeholder)
		baseObjects = append(baseObjects, placeholderContainer)
	} else {
		
		gt.guestWidget = NewGuestWidget(gt.guest, gt.vehicleIndex, gt.tileIndex, gt.grid, gt)
		
		guestContainer := container.NewPadded(gt.guestWidget)
		baseObjects = append(baseObjects, guestContainer)
	}

	gt.container = container.NewMax(baseObjects...)
	gt.container.Resize(fyne.NewSize(220, 70))
	return gt.container
}


func (gt *GuestTile) SetGuest(guest *app.Guest) {
	gt.guest = guest
	gt.isEmpty = false
	gt.isHidden = false

	if gt.guestWidget != nil {
		gt.guestWidget.vehicleIndex = gt.vehicleIndex
		gt.guestWidget.tileIndex = gt.tileIndex
	}

	
	gt.refreshContainer()
}


func (gt *GuestTile) ClearGuest() {
	gt.guest = nil
	gt.guestWidget = nil
	gt.isEmpty = true
	gt.isHidden = false

	gt.refreshContainer()
}


func (gt *GuestTile) IsEmpty() bool {
	return gt.isEmpty || gt.guest == nil
}


func (gt *GuestTile) HideGuest() {
	gt.isHidden = true
	gt.refreshContainer()
}


func (gt *GuestTile) ShowGuest() {
	gt.isHidden = false
	gt.refreshContainer()
}


func (gt *GuestTile) refreshContainer() {
	if gt.container == nil {
		return
	}

	
	objects := []fyne.CanvasObject{gt.background, gt.border}

	if gt.isEmpty || gt.guest == nil || gt.isHidden {
		
		placeholderContainer := container.NewCenter(gt.placeholder)
		objects = append(objects, placeholderContainer)
	} else {
		
		if gt.guestWidget == nil {
			gt.guestWidget = NewGuestWidget(gt.guest, gt.vehicleIndex, gt.tileIndex, gt.grid, gt)
		}
		guestContainer := container.NewPadded(gt.guestWidget)
		objects = append(objects, guestContainer)
	}

	
	gt.container.Objects = objects
	gt.container.Refresh()
}


func (gt *GuestTile) GetPosition() VehiclePosition {
	return VehiclePosition{
		VehicleIndex: gt.vehicleIndex,
		TileIndex:    gt.tileIndex,
	}
}


func (gt *GuestTile) ContainsPosition(pos fyne.Position) bool {
	if gt.container == nil || gt.card == nil || gt.card.card == nil {
		return false
	}

	
	cardPos := gt.card.card.Position()

	
	tilePos := gt.container.Position()
	tileSize := gt.container.Size()

	
	absX := cardPos.X + tilePos.X
	absY := cardPos.Y + tilePos.Y

	
	return pos.X >= absX && pos.X <= absX+tileSize.Width &&
		pos.Y >= absY && pos.Y <= absY+tileSize.Height
}


func (gt *GuestTile) HighlightAsDropTarget() {
	gt.background.FillColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255} 
	gt.border.StrokeColor = color.NRGBA{R: 76, G: 175, B: 80, A: 255}   
	gt.border.StrokeWidth = 3
	gt.background.Refresh()
	gt.border.Refresh()
}


func (gt *GuestTile) RemoveHighlight() {
	gt.background.FillColor = color.NRGBA{70, 70, 80, 255} 
	gt.border.StrokeColor = color.NRGBA{80, 80, 90, 255}   
	gt.border.StrokeWidth = 2                              
	gt.background.Refresh()
	gt.border.Refresh()
}


func (gt *GuestTile) GetGuest() *app.Guest {
	return gt.guest
}
