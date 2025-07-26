package ui

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func (cfg *Config) MakeUI() {
	

	var wrapper *mainContentWrapper

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://docs.google.com/spreadsheets/d/...")

	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetText("…your output here…")
	outputEntry.Wrapping = fyne.TextWrapWord

	var currentGrid *VehicleGrid 
	var tabs *container.AppTabs  
	var mapView *MapView

	runButton := widget.NewButton("Run", func() {
		var popup *widget.PopUp
		var result *RoutingProcess = nil
		var processErr error

		
		fyne.Do(func() {
			popup = ShowMessage(cfg.MainWindow)
			popup.Show()
		})

		
		go func() {
			
			
			result, processErr = ProcessEvent(urlEntry.Text)

			
			fyne.Do(func() {
				
				if popup != nil {
					popup.Hide()
				}

				
				if processErr != nil {
					fyne.Do(func() {
						ShowErrorNotification(cfg.MainWindow, "Processing Error", processErr.Error())
					})
					return
				} else {
					fyne.Do(func() {
						ShowSuccess(cfg.MainWindow)
					})
				}
				
				cfg.Rp = result
				outputEntry.SetText(result.String())
				currentGrid = NewVehicleGrid(result.rm, cfg)
				cfg.VehicleSection.Objects = []fyne.CanvasObject{currentGrid}
				cfg.VehicleSection.Refresh()

				if wrapper != nil {
					wrapper.grid = currentGrid
				}

				mapView = NewMapView(cfg.Rp, cfg)

				if tabs != nil {
					tabs.Items[2].Content = mapView
					tabs.Refresh()
				}
			})
		}()
	})

	runButton.Importance = widget.HighImportance
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(500, 1)) 

	rButton := container.NewHBox(
		runButton,
		spacer,
	)

	urlCard := widget.NewCard("Insert Google Sheet URL", "", container.NewVBox(
		container.NewBorder(nil, nil, nil, rButton, urlEntry),
	))

	

	outputScroll := container.NewScroll(outputEntry)
	
	

	outputSection := widget.NewCard("Guest Dropoff Summary", "", outputScroll)

	
	spacer2 := canvas.NewRectangle(color.Transparent)
	spacer2.SetMinSize(fyne.NewSize(0, 20))

	
	homeTab := container.NewBorder(
		container.NewVBox(spacer2, urlCard), 
		nil,                                 
		nil,                                 
		nil,                                 
		outputSection,                       
	)

	gradient := canvas.NewLinearGradient(
		color.NRGBA{55, 48, 163, 255}, 
		color.NRGBA{75, 61, 96, 255},  
		math.Pi/4,
	)

	buttonBar := CreateButtonBar(
		func() {
			if currentGrid != nil {
				if cfg.Rp != nil {
					outputEntry.SetText(cfg.Rp.String())
					outputEntry.Refresh()
					cfg.InfoLog.Println("Changes submitted")
				}

				mapView = NewMapView(cfg.Rp, cfg)

				if tabs != nil {
					tabs.Items[2].Content = mapView
					tabs.Refresh()
				}

			}
		},
		func() {

			if currentGrid != nil {
				currentGrid.ResetVehicles()
				cfg.InfoLog.Println("Reset to initial state")

				if cfg.Rp != nil {
					outputEntry.SetText(cfg.Rp.String())
					outputEntry.Refresh()
					cfg.InfoLog.Println("Changes submitted")
				}

				mapView = NewMapView(cfg.Rp, cfg)
				if tabs != nil {
					tabs.Items[2].Content = mapView
					tabs.Refresh()
				}
			}
		},
	)

	buttonBarBg := canvas.NewRectangle(color.NRGBA{25, 25, 30, 255})
	buttonBarContainer := container.NewMax(
		buttonBarBg,
		container.NewPadded(buttonBar),
	)

	
	routePlanningContent := container.NewBorder(
		buttonBarContainer, 
		nil,
		nil,
		nil,
		cfg.VehicleSection, 
	)

	routePlanningTab := container.NewMax(gradient, routePlanningContent)

	mapTabPlaceholder := container.NewCenter(
		widget.NewLabel("Run the routing process to see the map visualization"),
	)

	
	tabs = container.NewAppTabs(
		container.NewTabItem("Home", homeTab),
		container.NewTabItem("Route Planning", routePlanningTab),
		container.NewTabItem("Map", mapTabPlaceholder),
	)

	
	wrapper = &mainContentWrapper{
		content: tabs,
		grid:    nil, 
	}
	wrapper.ExtendBaseWidget(wrapper)

	
	cfg.MainWindow.SetContent(wrapper)

	
	cfg.MainWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape && currentGrid != nil && currentGrid.IsDragging() {
			currentGrid.CancelDrag()
		}
	})

	
	originalRunFunc := runButton.OnTapped
	runButton.OnTapped = func() {
		originalRunFunc()
		wrapper.grid = currentGrid
	}
}


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
