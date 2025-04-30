package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (app *config) makeUI() {

	// URL Section
	urlTitle := widget.NewLabelWithStyle(
		"Insert a Google Sheet URL",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://docs.google.com/spreadsheets/d/...")

	runButton := widget.NewButton("Run", func() {
		// TODO: Add functionality to fetch and parse Google Sheet
	})
	runButton.Importance = widget.HighImportance

	urlContent := container.NewVBox(
		urlTitle,
		container.NewBorder(nil, nil, nil, runButton, urlEntry),
	)

	app.MainWindow.SetContent(urlContent)

}
