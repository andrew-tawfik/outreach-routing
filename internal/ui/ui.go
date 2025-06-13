package ui

import (
	"fyne.io/fyne/v2"
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
		cfg.Rp = ProcessJsonEvent()
		outputEntry.SetText(cfg.Rp.String())

		currentGrid = NewVehicleGrid(cfg.Rp.rm, cfg)
		cfg.VehicleSection.Objects = []fyne.CanvasObject{currentGrid}
		cfg.VehicleSection.Refresh()
	})
	runButton.Importance = widget.HighImportance

	urlCard := widget.NewCard("Insert Google Sheet URL", "", container.NewVBox(
		container.NewBorder(nil, nil, nil, runButton, urlEntry),
	))

	// ---- Output & Actions ----------------------------------
	submitButton := widget.NewButton("Submit", func() {
		if currentGrid != nil {
			currentGrid.SubmitChanges()
			cfg.InfoLog.Println("Changes submitted")
		}
	})
	submitButton.Importance = widget.HighImportance

	resetButton := widget.NewButton("Reset", func() {
		if currentGrid != nil {
			currentGrid.ResetVehicles()
			cfg.InfoLog.Println("Reset to initial state")
		}
	})

	buttonBar := container.NewHBox(submitButton, resetButton)
	outputScroll := container.NewScroll(outputEntry)
	outputScroll.SetMinSize(fyne.NewSize(350, 300))

	outputSection := container.NewBorder(
		widget.NewLabel("Results"), // top
		nil,                        // bottom
		nil,                        // left
		nil,                        // right
		outputScroll,               // center
	)

	homeTab := container.NewVBox(
		urlCard,       // Your input section
		outputSection, // Results text
	)

	// Create Route Planning tab content
	routePlanningTab := container.NewBorder(
		buttonBar, // Action buttons at top
		nil,
		nil,
		nil,
		cfg.VehicleSection, // Vehicle grid
	)

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
