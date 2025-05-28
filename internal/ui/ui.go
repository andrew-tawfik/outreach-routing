package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (cfg *Config) MakeUI() {
	// ---- Input & Run ----------------------------------------
	urlTitle := widget.NewLabelWithStyle(
		"Insert a Google Sheet URL",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://docs.google.com/spreadsheets/d/...")

	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetText("…your output here…")
	outputEntry.Wrapping = fyne.TextWrapWord

	runButton := widget.NewButton("Run", func() {
		cfg.Rp = ProcessEvent(urlEntry.Text)
		outputEntry.SetText(cfg.Rp.String())

		grid := cfg.createVehicleGrid()
		cfg.VehicleSection.Objects = []fyne.CanvasObject{grid}
		cfg.VehicleSection.Refresh()
	})
	runButton.Importance = widget.HighImportance

	urlCard := widget.NewCard("Input", "", container.NewVBox(
		urlTitle,
		container.NewBorder(nil, nil, nil, runButton, urlEntry),
	))

	// ---- Output & Actions ----------------------------------
	submitButton := widget.NewButton("Submit", func() { /* TODO */ })
	submitButton.Importance = widget.HighImportance
	resetButton := widget.NewButton("Reset", func() { /* TODO */ })

	buttonBar := container.NewHBox(submitButton, resetButton)
	outputScroll := container.NewScroll(outputEntry)
	outputScroll.SetMinSize(fyne.NewSize(350, 300))
	outputCard := widget.NewCard("Results", "", container.NewVBox(buttonBar, outputScroll))

	topSection := container.New(layout.NewGridLayout(2), urlCard, outputCard)

	// ---- Assemble everything -------------------------------
	mainContent := container.NewBorder(
		topSection,         // top
		nil,                // bottom
		nil,                // left
		nil,                // right
		cfg.VehicleSection, // center
	)
	cfg.MainWindow.SetContent(mainContent)
}
