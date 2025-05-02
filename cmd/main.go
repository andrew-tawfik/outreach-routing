package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type config struct {
	App            fyne.App
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	MainWindow     fyne.Window
	rp             *RoutingProcess
	vehicleSection *fyne.Container
}

func main() {
	a := app.New()
	cfg := &config{
		App:        a,
		InfoLog:    log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:   log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		MainWindow: a.NewWindow("Anba Abraam Service"),
	}

	cfg.vehicleSection = container.NewStack()
	cfg.MainWindow.Resize(fyne.NewSize(1200, 700))
	cfg.MainWindow.SetMaster()
	cfg.makeUI()
	cfg.MainWindow.ShowAndRun()
}
