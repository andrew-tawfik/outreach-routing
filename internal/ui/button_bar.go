package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)


func CreateButtonBar(submitFunc, resetFunc func()) fyne.CanvasObject {
	
	submitButton := widget.NewButton("Submit", submitFunc)
	submitButton.Importance = widget.MediumImportance

	resetButton := widget.NewButton("Reset", resetFunc)
	resetButton.Importance = widget.MediumImportance

	
	middleSpacer := canvas.NewRectangle(color.Transparent)
	middleSpacer.SetMinSize(fyne.NewSize(10, 0)) 

	rightSpacer := canvas.NewRectangle(color.Transparent)
	rightSpacer.SetMinSize(fyne.NewSize(20, 0)) 

	return container.NewHBox(
		layout.NewSpacer(), 
		submitButton,
		resetButton,
	)
}


func CreateButtonBarWithStyle(submitFunc, resetFunc func(), submitStyle, resetStyle widget.Importance) fyne.CanvasObject {
	
	submitButton := widget.NewButton("Submit", submitFunc)
	submitButton.Importance = submitStyle

	resetButton := widget.NewButton("Reset", resetFunc)
	resetButton.Importance = resetStyle

	
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
