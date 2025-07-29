package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/andrew-tawfik/outreach-routing/internal/ui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	cfg := &ui.Config{
		App:        a,
		InfoLog:    log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:   log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		MainWindow: a.NewWindow("Anba Abraam Service"),
	}

	cfg.VehicleSection = container.NewStack()
	cfg.MainWindow.Resize(fyne.NewSize(1200, 700))
	cfg.MainWindow.SetMaster()
	cfg.MakeUI()
	cfg.MainWindow.ShowAndRun()
}
