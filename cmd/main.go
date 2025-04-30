package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type config struct {
	App        fyne.App
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	MainWindow fyne.Window
}

var myApp config

func main() {

	//create a fyne application
	a := app.New()
	myApp.App = a

	// create loggers
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	myApp.MainWindow = a.NewWindow("Anba Abraam Service")
	myApp.MainWindow.Resize(fyne.Size{Width: 1200, Height: 700})
	myApp.MainWindow.SetMaster()

	myApp.makeUI()

	myApp.MainWindow.ShowAndRun()
}
