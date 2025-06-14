package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// CreateButtonBar creates submit and reset buttons aligned to the right using standard Fyne buttons
func CreateButtonBar(submitFunc, resetFunc func()) fyne.CanvasObject {
	// Create standard Fyne buttons
	submitButton := widget.NewButton("Submit", submitFunc)
	submitButton.Importance = widget.MediumImportance

	resetButton := widget.NewButton("Reset", resetFunc)
	resetButton.Importance = widget.MediumImportance

	// Create spacers for layout
	middleSpacer := canvas.NewRectangle(color.Transparent)
	middleSpacer.SetMinSize(fyne.NewSize(10, 0)) // Space between buttons

	rightSpacer := canvas.NewRectangle(color.Transparent)
	rightSpacer.SetMinSize(fyne.NewSize(20, 0)) // Space from right edge

	return container.NewHBox(
		layout.NewSpacer(), // Pushes everything to the right
		submitButton,
		resetButton,
	)
}

// CreateButtonBarWithStyle creates a button bar with custom button styles
func CreateButtonBarWithStyle(submitFunc, resetFunc func(), submitStyle, resetStyle widget.Importance) fyne.CanvasObject {
	// Create standard Fyne buttons with custom importance
	submitButton := widget.NewButton("Submit", submitFunc)
	submitButton.Importance = submitStyle

	resetButton := widget.NewButton("Reset", resetFunc)
	resetButton.Importance = resetStyle

	// Create spacers for layout
	middleSpacer := canvas.NewRectangle(color.Transparent)
	middleSpacer.SetMinSize(fyne.NewSize(10, 0))

	rightSpacer := canvas.NewRectangle(color.Transparent)
	rightSpacer.SetMinSize(fyne.NewSize(20, 0))

	return container.NewHBox(
		layout.NewSpacer(),
		submitButton,
		middleSpacer,
		resetButton,
		rightSpacer,
	)
}
