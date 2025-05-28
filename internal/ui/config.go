package ui

import (
	"log"

	"fyne.io/fyne/v2"
)

type Config struct {
	App             fyne.App
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	MainWindow      fyne.Window
	Rp              *RoutingProcess
	VehicleSection  *fyne.Container
	GuestContainers []*fyne.Container
}
