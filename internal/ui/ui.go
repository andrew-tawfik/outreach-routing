package ui

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func (cfg *Config) MakeUI() {
	// ---- Input & Run ----------------------------------------

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://docs.google.com/spreadsheets/d/...")

	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetText("…your output here…")
	outputEntry.Wrapping = fyne.TextWrapWord

	var currentGrid *VehicleGrid // Store reference to current grid

	runButton := widget.NewButton("Run", func() {
		// rp, err := ProcessEvent(urlEntry.Text)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }

		rp, err := ProcessJsonEvent(0)
		if err != nil {
			fmt.Println(" will post notification")
		} else {
			cfg.Rp = rp
			outputEntry.SetText(cfg.Rp.String())

			currentGrid = NewVehicleGrid(cfg.Rp.rm, cfg)
			cfg.VehicleSection.Objects = []fyne.CanvasObject{currentGrid}
			cfg.VehicleSection.Refresh()
		}
	})
	runButton.Importance = widget.HighImportance
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(500, 1)) // 1px height to keep it invisible

	rButton := container.NewHBox(
		runButton,
		spacer,
	)

	urlCard := widget.NewCard("Insert Google Sheet URL", "", container.NewVBox(
		container.NewBorder(nil, nil, nil, rButton, urlEntry),
	))

	// ---- Output & Actions ----------------------------------

	outputScroll := container.NewScroll(outputEntry)
	outputScroll.SetMinSize(fyne.NewSize(350, 300))

	outputSection := widget.NewCard("Results", "", container.NewBorder(nil, nil, nil, nil, outputScroll))

	spacer1 := canvas.NewRectangle(color.Transparent)
	spacer1.SetMinSize(fyne.NewSize(0, 100))

	spacer2 := canvas.NewRectangle(color.Transparent)
	spacer2.SetMinSize(fyne.NewSize(0, 20))

	homeTab := container.NewVBox(
		spacer2,
		urlCard, // Your input section
		spacer1,
		outputSection, // Results text
	)

	gradient := canvas.NewLinearGradient(
		color.NRGBA{55, 48, 163, 255}, // Dark purple
		color.NRGBA{75, 61, 96, 255},  // Dark violet-gray
		math.Pi/4,
	)

	buttonBar := CreateButtonBar(
		func() {
			if currentGrid != nil {
				//currentGrid.SubmitChanges()
				if cfg.Rp != nil {
					outputEntry.SetText(cfg.Rp.String())
					outputEntry.Refresh()
					cfg.InfoLog.Println("Changes submitted")
				}

			}
		},
		func() {
			if currentGrid != nil {
				currentGrid.ResetVehicles()
				cfg.InfoLog.Println("Reset to initial state")
			}
		},
	)

	buttonBarBg := canvas.NewRectangle(color.NRGBA{25, 25, 30, 255})
	buttonBarContainer := container.NewMax(
		buttonBarBg,
		container.NewPadded(buttonBar),
	)

	// Create Route Planning tab content
	routePlanningContent := container.NewBorder(
		buttonBarContainer, // Action buttons at top
		nil,
		nil,
		nil,
		cfg.VehicleSection, // Vehicle grid
	)

	routePlanningTab := container.NewMax(gradient, routePlanningContent)

	// Create the tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Home", homeTab),
		container.NewTabItem("Route Planning", routePlanningTab),
	)

	// Create a wrapper that handles mouse events
	wrapper := &mainContentWrapper{
		content: tabs,
		grid:    nil, // Will be updated when grid is created
	}
	wrapper.ExtendBaseWidget(wrapper)

	// Store reference to wrapper so we can update the grid reference
	cfg.MainWindow.SetContent(wrapper)

	// Set up keyboard handling for ESC key
	cfg.MainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape && currentGrid != nil && currentGrid.IsDragging() {
			currentGrid.CancelDrag()
		}
	})

	// Update wrapper's grid reference when grid changes
	originalRunFunc := runButton.OnTapped
	runButton.OnTapped = func() {
		originalRunFunc()
		wrapper.grid = currentGrid
	}
}

// mainContentWrapper wraps the main content and handles mouse events
type mainContentWrapper struct {
	widget.BaseWidget
	content fyne.CanvasObject
	grid    *VehicleGrid
}

func (w *mainContentWrapper) CreateRenderer() fyne.WidgetRenderer {
	return &mainContentRenderer{
		wrapper: w,
		objects: []fyne.CanvasObject{w.content},
	}
}

// Implement desktop.Mouseable interface
func (w *mainContentWrapper) MouseIn(*desktop.MouseEvent) {}
func (w *mainContentWrapper) MouseOut()                   {}

func (w *mainContentWrapper) MouseMoved(event *desktop.MouseEvent) {
	if w.grid != nil && w.grid.IsDragging() {
		w.grid.UpdateDrag(event.AbsolutePosition)
	}
}

func (w *mainContentWrapper) MouseDown(*desktop.MouseEvent) {}

func (w *mainContentWrapper) MouseUp(event *desktop.MouseEvent) {
	if w.grid != nil && w.grid.IsDragging() {
		w.grid.EndDrag(event.AbsolutePosition)
	}
}

// mainContentRenderer is the renderer for the main content wrapper
type mainContentRenderer struct {
	wrapper *mainContentWrapper
	objects []fyne.CanvasObject
}

func (r *mainContentRenderer) Layout(size fyne.Size) {
	r.wrapper.content.Resize(size)
	r.wrapper.content.Move(fyne.NewPos(0, 0))
}

func (r *mainContentRenderer) MinSize() fyne.Size {
	return r.wrapper.content.MinSize()
}

func (r *mainContentRenderer) Refresh() {
	r.wrapper.content.Refresh()
}

func (r *mainContentRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *mainContentRenderer) Destroy() {}
