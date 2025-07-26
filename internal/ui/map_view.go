package ui

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/config"
	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

const (
	staticMapsBaseURL = "https://maps.googleapis.com/maps/api/staticmap"
	mapWidth  = 640
	mapHeight = 640
)

var depotCoor = coordinates.GuestCoordinates{Long: -75.726118, Lat: 45.396826}



type MapView struct {
	widget.BaseWidget

	routingProcess *RoutingProcess
	config         *Config

	colorMap map[string]color.Color
	colors   []string
	
	mapImage      *canvas.Image
	legend        *fyne.Container
	errorLabel    *widget.Label
	mainContainer *fyne.Container
	apiKey        string

	autoRefresh bool
}

type MapsConfig struct {
	MapsAPIKey string `json:"maps_api_key"`
}

func LoadMapsConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var config MapsConfig
	err = json.Unmarshal(data, &config)
	return config.MapsAPIKey, err
}

func NewMapView(rp *RoutingProcess, cfg *Config) *MapView {
	mv := &MapView{
		routingProcess: rp,
		config:         cfg,
		errorLabel:     widget.NewLabel(""),
		autoRefresh:    false,
	}

	mv.CreateColorMapping()
	err := mv.getApiKey()
	if err != nil {
		fmt.Println("could not get api key because: ", err)
	}
	mv.ExtendBaseWidget(mv)
	return mv
}

func (mv *MapView) getApiKey() error {
	apiKey, err := config.GetEmbeddedMapsAPIKey()
	if err == nil && apiKey != "" {
		mv.apiKey = apiKey
		return nil
	}

	projectRoot, err := filepath.Abs(filepath.Join(".", ".."))
	if err != nil {
		return fmt.Errorf("failed to resolve project root:", err)
	}
	credentialsPath := filepath.Join(projectRoot, "maps_config.json")

	apiKeyFromFile, jsonErr := LoadMapsConfig(credentialsPath)
	if jsonErr != nil {
		return fmt.Errorf("failed to load api key:", err)
	}
	mv.apiKey = apiKeyFromFile
	return nil
}


func (mv *MapView) CreateRenderer() fyne.WidgetRenderer {
	
	mv.mapImage = &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	mv.mapImage.SetMinSize(fyne.NewSize(mapWidth, mapHeight))

	
	mapBorder := canvas.NewRectangle(color.Transparent)
	mapBorder.StrokeColor = color.Transparent
	mapBorder.StrokeWidth = 2
	mapBorder.CornerRadius = 8

	mapContainer := container.NewMax(
		mapBorder,
		container.NewPadded(mv.mapImage),
	)
	

	if mv.routingProcess != nil && mv.routingProcess.rm != nil {
		mv.legend = mv.createLegend()
	} else {
		mv.legend = container.NewVBox(widget.NewLabel("No route data available"))
	}

	
	content := container.NewHSplit(
		mapContainer, 
		mv.legend,    
	)
	content.SetOffset(0.6)

	mv.mainContainer = container.NewBorder(
		nil,           
		mv.errorLabel, 
		nil,           
		nil,           
		content,       
	)
	mv.mainContainer = container.NewPadded(mv.mainContainer)

	if mv.routingProcess != nil && mv.routingProcess.rm != nil {
		go mv.loadMap()
	} else {
		mv.showError("No routing data available. Please run the routing process first.")
	}

	return &mapViewRenderer{
		mapView: mv,
		objects: []fyne.CanvasObject{mv.mainContainer},
	}
}


func (mv *MapView) createLegend() *fyne.Container {
	if mv.routingProcess == nil || mv.routingProcess.rm == nil {
		return container.NewVBox(widget.NewLabel("No route data available"))
	}

	legendTitle := widget.NewLabelWithStyle("Route Legend", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	legendItems := container.NewVBox(legendTitle, widget.NewSeparator())

	
	depotRow := mv.createLegendRow("brown", "Depot", "555 Parkdale Ave", nil)
	legendItems.Add(depotRow)
	legendItems.Add(widget.NewSeparator())

	
	for i, vehicle := range mv.routingProcess.rm.Vehicles {
		if len(vehicle.Guests) > 0 {
			colorIndex := i % len(mv.colors)
			vehicleColor := mv.colors[colorIndex]

			vehicleLabel := widget.NewLabelWithStyle(
				fmt.Sprintf("Vehicle %d", i+1),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)

			legendItems.Add(vehicleLabel)

			groceryAdj := 0
			if mv.routingProcess.ae.EventType == "Grocery" {
				groceryAdj = 1
			}
			
			for elem := vehicle.Route.List.Front(); elem != nil; elem = elem.Next() {
				addressIndex := elem.Value.(int) + groceryAdj
				addr := mv.routingProcess.lr.CoordianteMap.AddressOrder[addressIndex]
				coor := mv.routingProcess.lr.CoordianteMap.CoordinateToAddress[addr]
				markerLabel := mv.determineMarkerLabel(&vehicle, &coor)

				guestInfo := mv.getGuestInfoForAddress(addr, &vehicle)
				row := mv.createLegendRow(vehicleColor, markerLabel, addr, guestInfo)
				legendItems.Add(row)
			}

			legendItems.Add(widget.NewSeparator())
		}
	}

	
	legendCard := widget.NewCard("", "", container.NewScroll(legendItems))
	legendCard.Resize(fyne.NewSize(300, mapHeight))

	return container.NewMax(legendCard)
}

func (mv *MapView) determineMarkerLabel(vehicle *app.Vehicle, coor *coordinates.GuestCoordinates) string {
	colorIndex := -1
	for i := range vehicle.Locations {
		match := vehicle.Locations[i]
		if match.Long == coor.Long && match.Lat == coor.Lat {
			colorIndex = i
			break
		}
	}
	return fmt.Sprintf("%c", 'A'+colorIndex)

}

func (mv *MapView) getGuestInfoForAddress(address string, vehicle *app.Vehicle) []string {
	var names []string
	for _, guest := range vehicle.Guests {
		if guest.Address == address {
			guestName := guest.Name
			if len(guestName) > 25 {
				guestName = guestName[:22] + "..."
			}
			names = append(names, guestName)
		}
	}
	return names
}


func (mv *MapView) createLegendRow(colorName, label, address string, guestInfo []string) *fyne.Container {
	
	markerColor := mv.colorMap[colorName]
	colorBox := canvas.NewRectangle(markerColor)
	colorBox.SetMinSize(fyne.NewSize(20, 20))
	colorBox.StrokeColor = color.Black
	colorBox.StrokeWidth = 1
	colorBox.CornerRadius = 10

	
	markerLabel := widget.NewLabel(label)
	markerLabel.TextStyle = fyne.TextStyle{Bold: true}

	
	displayAddress := address
	if len(displayAddress) > 30 {
		displayAddress = displayAddress[:27] + "..."
	}

	
	var contentLabel *widget.Label
	if len(guestInfo) == 0 {
		
		contentLabel = widget.NewLabel(displayAddress)
	} else {
		var content strings.Builder
		content.WriteString(displayAddress)

		for _, guest := range guestInfo {
			content.WriteString("\n → ")
			content.WriteString(guest)
		}
		contentLabel = widget.NewLabel(content.String())
	}

	contentLabel.TextStyle.Monospace = true 

	return container.NewHBox(
		colorBox,
		markerLabel,
		contentLabel,
	)
}

func NRGBAToHex(c color.Color) string {
	rgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	return fmt.Sprintf("0x%02X%02X%02X", rgba.R, rgba.G, rgba.B)
}

func (mv *MapView) CreateColorMapping() {

	
	ColorMap := map[string]color.Color{
		"brown":   color.NRGBA{139, 69, 19, 255},   
		"red":     color.NRGBA{255, 0, 0, 255},     
		"blue":    color.NRGBA{0, 0, 255, 255},     
		"green":   color.NRGBA{0, 128, 0, 255},     
		"orange":  color.NRGBA{255, 165, 0, 255},   
		"purple":  color.NRGBA{128, 0, 128, 255},   
		"cyan":    color.NRGBA{0, 255, 255, 255},   
		"magenta": color.NRGBA{255, 0, 255, 255},   
		"pink":    color.NRGBA{255, 192, 203, 255}, 
		"teal":    color.NRGBA{0, 128, 128, 255},   
		"indigo":  color.NRGBA{75, 0, 130, 255},    
		"gold":    color.NRGBA{255, 215, 0, 255},   
		"lime":    color.NRGBA{50, 205, 50, 255},   

	}

	
	VehicleColors := []string{"red", "blue", "green", "orange", "purple", "cyan",
		"magenta", "pink", "teal", "indigo", "gold", "lime"}

	mv.colorMap = ColorMap
	mv.colors = VehicleColors

}


func (mv *MapView) buildMapURL() string {
	params := url.Values{}
	params.Set("size", fmt.Sprintf("%dx%d", mapWidth, mapHeight))
	params.Set("key", mv.apiKey)
	params.Set("maptype", "roadmap")
	params.Set("center", fmt.Sprintf("%f,%f", depotCoor.Lat, depotCoor.Long))

	
	depotColor := mv.colorMap["brown"]
	params.Add("markers", fmt.Sprintf("color:%s|label:M|%f,%f", NRGBAToHex(depotColor), depotCoor.Lat, depotCoor.Long))

	
	for i, vehicle := range mv.routingProcess.rm.Vehicles {
		if len(vehicle.Guests) > 0 {
			
			vehicleColorName := mv.colors[i%len(mv.colors)]
			vehicleColor := mv.colorMap[vehicleColorName]
			hexColor := NRGBAToHex(vehicleColor)

			
			for j, loc := range vehicle.Locations {
				label := fmt.Sprintf("%c", 'A'+j)
				marker := fmt.Sprintf("color:%s|label:%s|%f,%f", hexColor, label, loc.Lat, loc.Long)
				params.Add("markers", marker)
			}
		}
	}

	completeURL := staticMapsBaseURL + "?" + params.Encode()
	return completeURL
}


func (mv *MapView) loadMap() {

	if mv.routingProcess == nil || mv.routingProcess.rm == nil {
		mv.showError("No route data available")
		return
	}

	
	mapURL := mv.buildMapURL()

	
	resp, err := http.Get(mapURL)
	if err != nil {
		mv.showError(fmt.Sprintf("Failed to load map: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		mv.showError(fmt.Sprintf("Failed to load map: HTTP %d", resp.StatusCode))
		return
	}

	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		mv.showError(fmt.Sprintf("Failed to read map data: %v", err))
		return
	}

	
	resource := fyne.NewStaticResource("map.png", data)

	
	fyne.Do(func() {
		mv.mapImage.Resource = resource
		mv.mapImage.Refresh()
		mv.errorLabel.SetText("")
	})
}

func (mv *MapView) showError(message string) {
	fyne.Do(func() {
		mv.errorLabel.SetText(message)
		mv.errorLabel.TextStyle = fyne.TextStyle{Bold: true}
		mv.errorLabel.Refresh()
	})
}


func (mv *MapView) Refresh() {
	if mv.mainContainer != nil && mv.autoRefresh {
		mv.legend = mv.createLegend()
		mv.mainContainer.Refresh()
		go mv.loadMap()
	}
}


type mapViewRenderer struct {
	mapView *MapView
	objects []fyne.CanvasObject
}

func (r *mapViewRenderer) Layout(size fyne.Size) {
	r.mapView.mainContainer.Resize(size)
}

func (r *mapViewRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1000, 700)
}

func (r *mapViewRenderer) Refresh() {
	r.mapView.mainContainer.Refresh()
}

func (r *mapViewRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *mapViewRenderer) Destroy() {}


func (mv *MapView) ForceRefresh() {
	if mv.mainContainer != nil {
		mv.legend = mv.createLegend()
		mv.mainContainer.Refresh()
		go mv.loadMap()
	}
}
