package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
		myApp.rp = ProcessEvent(urlEntry.Text)
	})
	runButton.Importance = widget.HighImportance

	urlContent := container.NewVBox(
		urlTitle,
		container.NewBorder(nil, nil, nil, runButton, urlEntry),
	)
	urlCard := widget.NewCard("", "", urlContent)
	//urlCard.Resize(fyne.Size{Width: 550, Height: 150})

	// Action buttons
	submitButton := widget.NewButton("Submit", func() {
		// TODO: Add submit functionality
	})
	submitButton.Importance = widget.HighImportance

	resetButton := widget.NewButton("Reset", func() {
		// TODO: Add reset functionality
	})

	buttonContainer := container.NewHBox(
		submitButton,
		resetButton,
	)

	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetText("…your output here…")
	outputEntry.Wrapping = fyne.TextWrapWord

	outputScrollContainer := container.NewScroll(outputEntry)
	outputScrollContainer.SetMinSize(fyne.Size{Width: 350, Height: 120})

	outputContent := container.NewVBox(
		buttonContainer,
		outputScrollContainer,
	)

	outputCard := widget.NewCard("", "", outputContent)
	//outputCard.Resize(fyne.Size{Width: 550, Height: 150})

	// Top section with URL and output side by side using HBox
	// a 2-column grid will give each cell half the width
	topSection := container.New(layout.NewGridLayout(2), urlCard, outputCard)

	app.MainWindow.SetContent(topSection)
}
