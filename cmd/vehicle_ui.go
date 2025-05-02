package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
)

func (cfg *config) createVehicleGrid() fyne.CanvasObject {
	cards := make([]fyne.CanvasObject, 0, len(cfg.rp.rm.Vehicles))
	for i, v := range cfg.rp.rm.Vehicles {
		cards = append(cards, cfg.makeVehicleCard(v, i))
	}
	grid := container.NewGridWithColumns(4, cards...)
	return container.NewScroll(grid)
}

func (cfg *config) makeVehicleCard(v app.Vehicle, idx int) fyne.CanvasObject {
	cfg.InfoLog.Printf("Vehicle %d has %d guests", idx+1, len(v.Guests))
	title := fmt.Sprintf("Vehicle %d", idx+1)
	guestsBox := container.NewVBox()
	for i, g := range v.Guests {
		dg := NewDraggableGuest(g, idx, i)
		guestsBox.Add(dg)
	}

	scroll := container.NewVScroll(guestsBox)
	scroll.SetMinSize(fyne.NewSize(0, 400))

	return widget.NewCard(title, "", scroll)
}
