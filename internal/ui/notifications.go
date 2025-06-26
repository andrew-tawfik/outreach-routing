package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowInAppNotification shows a notification within the app window
func ShowErrorNotification(window fyne.Window, title, message string) {
	// Create notification card
	var bgColor color.Color
	bgColor = color.NRGBA{R: 200, G: 50, B: 50, A: 230} // Red

	background := canvas.NewRectangle(bgColor)
	background.CornerRadius = 4

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(titleLabel, messageLabel)
	notification := container.NewStack(background, container.NewPadded(content))

	// Create popup
	popup := widget.NewPopUp(notification, window.Canvas())
	popup.Resize(fyne.NewSize(300, 100))

	// Position at top-right
	windowSize := window.Canvas().Size()
	popup.Move(fyne.NewPos(windowSize.Width-320, 20))

	popup.Show()

	// Auto-hide after 3 seconds
	go func() {
		time.Sleep(10 * time.Second)
		popup.Hide()
	}()
}

func ShowMessage(window fyne.Window) {
	title := "Processing, will show results shortly."

	// Create notification card
	var bgColor color.Color
	bgColor = color.NRGBA{R: 240, G: 240, B: 240, A: 255} // soft gray
	background := canvas.NewRectangle(bgColor)
	background.CornerRadius = 4

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	messageLabel := widget.NewLabel(title)
	messageLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(titleLabel, messageLabel)
	notification := container.NewStack(background, container.NewPadded(content))

	// Create popup
	popup := widget.NewPopUp(notification, window.Canvas())
	popup.Resize(fyne.NewSize(300, 100))

	// Position at top-right
	windowSize := window.Canvas().Size()
	popup.Move(fyne.NewPos(windowSize.Width-320, 20))

	popup.Show()

	// Auto-hide after 3 seconds
	go func() {
		time.Sleep(10 * time.Second)
		popup.Hide()
	}()
}
