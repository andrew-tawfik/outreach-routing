package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

var adjustedY float32 = 125


type VehicleGrid struct {
	widget.BaseWidget

	
	vehicles       []*VehicleCard
	routeManager   *app.RouteManager
	config         *Config
	vehicleManager *VehicleManager

	
	dragOverlay  *fyne.Container 
	dragVisual   *fyne.Container 
	draggedGuest *app.Guest      
	dragOrigin   VehiclePosition 
	isDragging   bool

	dragPosition    fyne.Position 
	currentMousePos fyne.Position

	
	gridContainer   *fyne.Container 
	mainContainer   *fyne.Container 
	scrollContainer *container.Scroll
	eventType       string
}


type VehiclePosition struct {
	VehicleIndex int
	TileIndex    int
}


func NewVehicleGrid(rm *app.RouteManager, cfg *Config) *VehicleGrid {
	vg := &VehicleGrid{
		routeManager: rm,
		config:       cfg,
		vehicles:     make([]*VehicleCard, 0),
		eventType:    cfg.Rp.ae.EventType,
	}

	
	vg.dragOverlay = container.NewWithoutLayout()
	vg.dragOverlay.Hide()

	
	vg.vehicleManager = NewVehicleManager(rm, vg, cfg)

	vg.ExtendBaseWidget(vg)

	
	vg.refreshVehicles()

	return vg
}


func (vg *VehicleGrid) CreateRenderer() fyne.WidgetRenderer {
	vg.gridContainer = vg.createVehicleCards()

	
	vg.scrollContainer = container.NewScroll(vg.gridContainer)
	vg.scrollContainer.SetMinSize(fyne.NewSize(900, 500))

	
	vg.mainContainer = container.NewMax(
		vg.scrollContainer,
		vg.dragOverlay,
	)

	return &vehicleGridRenderer{
		grid:    vg,
		objects: []fyne.CanvasObject{vg.mainContainer},
	}
}


type vehicleGridRenderer struct {
	grid    *VehicleGrid
	objects []fyne.CanvasObject
}

func (r *vehicleGridRenderer) Layout(size fyne.Size) {
	r.grid.mainContainer.Resize(size)
}

func (r *vehicleGridRenderer) MinSize() fyne.Size {
	return fyne.NewSize(900, 500)
}

func (r *vehicleGridRenderer) Refresh() {
	r.grid.mainContainer.Refresh()
}

func (r *vehicleGridRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *vehicleGridRenderer) Destroy() {}


func (vg *VehicleGrid) createVehicleCards() *fyne.Container {
	cards := make([]fyne.CanvasObject, 0, len(vg.vehicles))

	for _, vehicleCard := range vg.vehicles {
		cards = append(cards, vehicleCard.CreateCard())
	}

	
	return container.NewGridWithColumns(4, cards...)
}


func (vg *VehicleGrid) refreshVehicles() {
	vg.vehicles = make([]*VehicleCard, len(vg.routeManager.Vehicles))

	for i := range vg.routeManager.Vehicles {
		vg.vehicles[i] = NewVehicleCard(i, &vg.routeManager.Vehicles[i], vg)
	}
}


func (vg *VehicleGrid) StartDrag(guest *app.Guest, origin VehiclePosition, startPos fyne.Position, offset fyne.Position) {
	vg.draggedGuest = guest
	vg.dragOrigin = origin
	vg.isDragging = true

	
	vg.dragPosition = startPos

	
	vg.createDragVisual(guest)

	
	vg.vehicles[origin.VehicleIndex].HideGuest(origin.TileIndex)

	
	vg.dragOverlay.Hide()
}


func (vg *VehicleGrid) createDragVisual(guest *app.Guest) {
	
	background := canvas.NewRectangle(color.NRGBA{60, 60, 70, 255}) 
	background.CornerRadius = 3
	background.Resize(fyne.NewSize(190, 40))

	nameLabel := widget.NewLabel(guest.Name + " (" + fmt.Sprintf("%d", guest.GroupSize) + ")")
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	vg.dragVisual = container.NewMax(background, container.NewPadded(nameLabel))
	vg.dragVisual.Resize(fyne.NewSize(190, 40))
}

func (vg *VehicleGrid) UpdateDrag(globalPos fyne.Position) {
	if !vg.isDragging || vg.dragVisual == nil {
		return
	}

	vg.currentMousePos = fyne.Position{X: globalPos.X, Y: globalPos.Y - 20}

	
	if !vg.dragOverlay.Visible() {
		vg.dragOverlay.Objects = []fyne.CanvasObject{vg.dragVisual}
		vg.dragOverlay.Show()
	}

	
	vg.dragPosition = fyne.NewPos(
		globalPos.X-95,        
		globalPos.Y-adjustedY, 
	)

	
	vg.dragVisual.Move(vg.dragPosition)

	
	vg.highlightValidDropTargets()
}


func (vg *VehicleGrid) EndDrag(globalPos fyne.Position) {
	if !vg.isDragging {
		return
	}

	
	vg.currentMousePos = fyne.Position{X: globalPos.X, Y: globalPos.Y - 20}

	
	targetPos := vg.positionToTile()

	if vg.isValidDropTarget(targetPos) {
		vg.performMove(vg.dragOrigin, targetPos)
	} else {
		
		vg.vehicles[vg.dragOrigin.VehicleIndex].ShowGuest(vg.dragOrigin.TileIndex)
	}

	
	vg.cleanupDrag()
}


func (vg *VehicleGrid) CancelDrag() {
	if vg.isDragging {
		
		vg.vehicles[vg.dragOrigin.VehicleIndex].ShowGuest(vg.dragOrigin.TileIndex)
		vg.cleanupDrag()
	}
}


func (vg *VehicleGrid) cleanupDrag() {
	vg.isDragging = false
	vg.draggedGuest = nil
	vg.dragOverlay.Hide()
	vg.dragOverlay.Objects = nil
	vg.dragVisual = nil

	
	for _, vehicle := range vg.vehicles {
		vehicle.RemoveAllHighlights()
	}
}


func (vg *VehicleGrid) positionToTile() VehiclePosition {
	return vg.findTileAtPosition(vg.currentMousePos)
}


func (vg *VehicleGrid) highlightValidDropTargets() {
	targetPos := vg.positionToTile()

	for vIndex, vehicle := range vg.vehicles {
		for tIndex, tile := range vehicle.tiles {
			if vIndex == targetPos.VehicleIndex && tIndex == targetPos.TileIndex &&
				vg.isValidDropTarget(targetPos) {
				tile.HighlightAsDropTarget()
			} else {
				tile.RemoveHighlight()
			}
		}
	}
}


func (vg *VehicleGrid) isValidDropTarget(target VehiclePosition) bool {
	if target.VehicleIndex < 0 || target.VehicleIndex >= len(vg.vehicles) {
		return false
	}

	if target.TileIndex < 0 || target.TileIndex >= len(vg.vehicles[target.VehicleIndex].tiles) {
		return false
	}

	
	if target.VehicleIndex == vg.dragOrigin.VehicleIndex &&
		target.TileIndex == vg.dragOrigin.TileIndex {
		return true
	}

	vehicle := vg.vehicles[target.VehicleIndex]

	
	if target.VehicleIndex == vg.dragOrigin.VehicleIndex {
		
		
		return target.TileIndex >= len(vg.routeManager.Vehicles[target.VehicleIndex].Guests)
	}

	
	return vehicle.IsTileEmpty(target.TileIndex)
}

func (vg *VehicleGrid) performMove(from, to VehiclePosition) {
	rm := vg.config.Rp.rm
	lr := vg.config.Rp.lr

	
	if from.VehicleIndex == to.VehicleIndex {
		vehicle := &rm.Vehicles[from.VehicleIndex]

		
		guestIndex := vg.findGuestIndex(vehicle, vg.draggedGuest)
		if guestIndex < 0 {
			return
		}

		
		guest := vehicle.Guests[guestIndex]

		
		insertPos := to.TileIndex
		if insertPos > len(vehicle.Guests) {
			insertPos = len(vehicle.Guests)
		}

		
		
		if insertPos > guestIndex {
			insertPos--
		}

		
		vehicle.Guests = append(
			vehicle.Guests[:guestIndex],
			vehicle.Guests[guestIndex+1:]...,
		)

		
		vehicle.Guests = append(
			vehicle.Guests[:insertPos],
			append([]app.Guest{guest}, vehicle.Guests[insertPos:]...)...,
		)

		
		vehicle.UpdateRouteFromGuests(lr, vg.eventType)

		
		vg.refreshAfterMove()
		return
	}

	
	sourceVehicle := &rm.Vehicles[from.VehicleIndex]
	guestIndex := vg.findGuestIndex(sourceVehicle, vg.draggedGuest)
	guest := sourceVehicle.Guests[guestIndex]

	if guestIndex >= 0 {
		
		sourceVehicle.Guests = append(
			sourceVehicle.Guests[:guestIndex],
			sourceVehicle.Guests[guestIndex+1:]...,
		)
		sourceVehicle.SeatsRemaining += vg.draggedGuest.GroupSize

		
		sourceVehicle.UpdateRouteFromGuests(lr, vg.eventType)
	}

	
	targetVehicle := &rm.Vehicles[to.VehicleIndex]

	
	insertPos := to.TileIndex
	if insertPos > len(targetVehicle.Guests) {
		insertPos = len(targetVehicle.Guests)
	}

	
	targetVehicle.Guests = append(targetVehicle.Guests[:insertPos],
		append([]app.Guest{guest}, targetVehicle.Guests[insertPos:]...)...)

	targetVehicle.SeatsRemaining -= guest.GroupSize

	
	targetVehicle.UpdateRouteFromGuests(lr, vg.eventType)

	
	vg.vehicleManager.hasChanges = true

	
	vg.refreshAfterMove()
}


func (vg *VehicleGrid) findGuestIndex(vehicle *app.Vehicle, guest *app.Guest) int {
	for i, g := range vehicle.Guests {
		if g.Name == guest.Name && g.Address == guest.Address {
			return i
		}
	}
	return -1
}


func (vg *VehicleGrid) refreshAfterMove() {
	
	for _, vehicle := range vg.vehicles {
		vehicle.RemoveAllHighlights()
	}

	
	for i := range vg.vehicles {
		vg.vehicles[i].RefreshTiles()
	}

	
	vg.gridContainer.Objects = nil
	for _, vehicleCard := range vg.vehicles {
		vg.gridContainer.Objects = append(vg.gridContainer.Objects, vehicleCard.CreateCard())
	}
	vg.gridContainer.Refresh()
}


func (vg *VehicleGrid) ResetVehicles() {
	if vg.vehicleManager != nil {
		vg.vehicleManager.ResetToInitialState()
	}
}


func (vg *VehicleGrid) SubmitChanges() {
	if vg.vehicleManager != nil {
		vg.vehicleManager.SubmitChanges()
	}
}


func (vg *VehicleGrid) MouseIn(*desktop.MouseEvent) {}
func (vg *VehicleGrid) MouseOut()                   {}

func (vg *VehicleGrid) MouseMoved(event *desktop.MouseEvent) {
	if vg.isDragging {
		vg.UpdateDrag(event.AbsolutePosition)
	}
}


func (vg *VehicleGrid) IsDragging() bool {
	return vg.isDragging
}


func (vg *VehicleGrid) GetDraggedGuest() *app.Guest {
	return vg.draggedGuest
}




func (vg *VehicleGrid) findTileAtPosition(mousePos fyne.Position) VehiclePosition {
	
	mainPos := vg.mainContainer.Position()

	
	scrollOffset := vg.scrollContainer.Offset

	
	
	adjustedMouseY := mousePos.Y - adjustedY 

	contentPos := fyne.NewPos(
		mousePos.X-mainPos.X+scrollOffset.X,
		adjustedMouseY-mainPos.Y+scrollOffset.Y,
	)

	
	for vIndex, vehicle := range vg.vehicles {
		if vehicle.card == nil {
			continue
		}

		cardPos := vehicle.card.Position()
		cardSize := vehicle.card.Size()

		
		if contentPos.X >= cardPos.X && contentPos.X <= cardPos.X+cardSize.Width &&
			contentPos.Y >= cardPos.Y && contentPos.Y <= cardPos.Y+cardSize.Height {

			
			for tIndex, tile := range vehicle.tiles {
				if tile.container == nil {
					continue
				}

				
				tilePos := tile.container.Position()
				tileSize := tile.container.Size()

				
				tilePosInCard := fyne.NewPos(
					cardPos.X+tilePos.X,
					cardPos.Y+tilePos.Y,
				)

				
				if contentPos.X >= tilePosInCard.X && contentPos.X <= tilePosInCard.X+tileSize.Width &&
					contentPos.Y >= tilePosInCard.Y && contentPos.Y <= tilePosInCard.Y+tileSize.Height {

					return VehiclePosition{VehicleIndex: vIndex, TileIndex: tIndex}
				}
			}
		}
	}

	return VehiclePosition{VehicleIndex: -1, TileIndex: -1}
}


func (vg *VehicleGrid) GetDragPosition() fyne.Position {
	return vg.dragPosition
}


func (vg *VehicleGrid) updateVehicleRoutes() {
	rm := vg.config.Rp.rm
	lr := vg.config.Rp.lr

	for i := range rm.Vehicles {
		rm.Vehicles[i].UpdateRouteFromGuests(lr, vg.eventType)
	}
}
