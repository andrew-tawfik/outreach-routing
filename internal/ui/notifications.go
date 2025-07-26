package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)


func ShowErrorNotification(window fyne.Window, title, message string) {
	
	bgColor := color.NRGBA{R: 200, G: 50, B: 50, A: 230} 

	background := canvas.NewRectangle(bgColor)
	background.CornerRadius = 4

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(titleLabel, messageLabel)
	notification := container.NewStack(background, container.NewPadded(content))

	
	popup := widget.NewPopUp(notification, window.Canvas())
	popup.Resize(fyne.NewSize(300, 100))

	
	windowSize := window.Canvas().Size()
	popup.Move(fyne.NewPos(windowSize.Width-320, 20))

	popup.Show()

	
	go func() {
		time.Sleep(10 * time.Second)
		
		fyne.Do(func() {
			popup.Hide()
		})
	}()
}

func ShowMessage(window fyne.Window) *widget.PopUp {
	title := "Processing, please wait..."

	bgColor := color.NRGBA{R: 80, G: 80, B: 90, A: 255}
	background := canvas.NewRectangle(bgColor)
	background.CornerRadius = 4
	background.StrokeColor = color.NRGBA{R: 120, G: 120, B: 130, A: 255}
	background.StrokeWidth = 2

	messageLabel := widget.NewLabel(title)
	messageLabel.Alignment = fyne.TextAlignCenter
	messageLabel.Wrapping = fyne.TextWrapWord

	notification := container.NewStack(background, container.NewPadded(messageLabel))
	popup := widget.NewPopUp(notification, window.Canvas())
	popup.Resize(fyne.NewSize(220, 60))

	
	windowSize := window.Canvas().Size()
	popup.Move(fyne.NewPos(windowSize.Width-240, 50))

	return popup
}

func ShowSuccess(window fyne.Window) {
	title := "Operation Successful: Route Generated"

	bgColor := color.NRGBA{R: 76, G: 175, B: 80, A: 255}
	background := canvas.NewRectangle(bgColor)
	background.CornerRadius = 4
	background.StrokeColor = color.NRGBA{R: 56, G: 142, B: 60, A: 255}
	background.StrokeWidth = 1

	messageLabel := widget.NewLabel(title)
	messageLabel.Alignment = fyne.TextAlignCenter
	messageLabel.Wrapping = fyne.TextWrapWord

	notification := container.NewStack(background, container.NewPadded(messageLabel))
	popup := widget.NewPopUp(notification, window.Canvas())
	popup.Resize(fyne.NewSize(200, 60))

	
	windowSize := window.Canvas().Size()
	popup.Move(fyne.NewPos(windowSize.Width-220, 50))
	popup.Show()

	
	go func() {
		time.Sleep(3 * time.Second)
		
		fyne.Do(func() {
			popup.Hide()
		})
	}()
}
