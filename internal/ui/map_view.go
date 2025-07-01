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
	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

const (
	// Google Maps Static API base URL
	staticMapsBaseURL = "https://maps.googleapis.com/maps/api/staticmap"

	// Map dimensions
	mapWidth  = 640
	mapHeight = 640
)

var depotCoor = coordinates.GuestCoordinates{Long: -75.726118, Lat: 45.396826}

// ColorMap stores 13 distinct colors for vehicles plus brown for depot

type MapView struct {
	widget.BaseWidget

	routingProcess *RoutingProcess
	config         *Config

	colorMap map[string]color.Color
	colors   []string
	// UI components
	mapImage      *canvas.Image
	legend        *fyne.Container
	errorLabel    *widget.Label
	mainContainer *fyne.Container
	apiKey        string
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
	}

	mv.CreateColorMapping()
	mv.getApiKey()
	mv.ExtendBaseWidget(mv)
	return mv
}

func (mv *MapView) getApiKey() error {
	projectRoot, err := filepath.Abs(filepath.Join(".", ".."))
	if err != nil {
		return fmt.Errorf("failed to resolve project root:", err)
	}
	credentialsPath := filepath.Join(projectRoot, "maps_config.json")

	apiKey, jsonErr := LoadMapsConfig(credentialsPath)
	if jsonErr != nil {
		return fmt.Errorf("failed to load api key:", err)
	}
	mv.apiKey = apiKey
	return nil
}

// CreateRenderer creates the renderer for the map view
func (mv *MapView) CreateRenderer() fyne.WidgetRenderer {
	// Create initial UI components
	mv.mapImage = &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	mv.mapImage.SetMinSize(fyne.NewSize(mapWidth, mapHeight))

	// Create map container with border
	mapBorder := canvas.NewRectangle(color.NRGBA{200, 200, 200, 255})
	mapBorder.StrokeColor = color.NRGBA{150, 150, 150, 255}
	mapBorder.StrokeWidth = 2
	mapBorder.CornerRadius = 8

	mapContainer := container.NewMax(
		mapBorder,
		container.NewPadded(mv.mapImage),
	)

	// Create legend
	mv.legend = mv.createLegend()

	// Create main layout
	content := container.NewBorder(
		nil,           // top
		mv.errorLabel, // bottom
		nil,           // left
		mv.legend,     // right
		mapContainer,  // center
	)

	mv.mainContainer = container.NewPadded(content)

	// Load the map
	go mv.loadMap()

	return &mapViewRenderer{
		mapView: mv,
		objects: []fyne.CanvasObject{mv.mainContainer},
	}
}

// createLegend creates the legend showing vehicle colors and addresses
func (mv *MapView) createLegend() *fyne.Container {
	if mv.routingProcess == nil || mv.routingProcess.rm == nil {
		return container.NewVBox(widget.NewLabel("No route data available"))
	}

	legendTitle := widget.NewLabelWithStyle("Route Legend", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	legendItems := container.NewVBox(legendTitle, widget.NewSeparator())

	// Add depot marker
	depotRow := mv.createLegendRow("black", "Depot", "555 Parkdale Ave")
	legendItems.Add(depotRow)
	legendItems.Add(widget.NewSeparator())

	// Add vehicle routes
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

			// Add each address for this vehicle
			for j, guest := range vehicle.Guests {
				markerLabel := fmt.Sprintf("%c", 'A'+j)
				row := mv.createLegendRow(vehicleColor, markerLabel, guest.Address)
				legendItems.Add(row)
			}

			legendItems.Add(widget.NewSeparator())
		}
	}

	// Wrap in a card
	legendCard := widget.NewCard("", "", container.NewScroll(legendItems))
	legendCard.Resize(fyne.NewSize(300, mapHeight))

	return container.NewMax(legendCard)
}

// createLegendRow creates a single row in the legend
func (mv *MapView) createLegendRow(colorName, label, address string) *fyne.Container {
	// Create color indicator
	markerColor := mv.colorMap[colorName]
	colorBox := canvas.NewRectangle(markerColor)
	colorBox.SetMinSize(fyne.NewSize(20, 20))
	colorBox.StrokeColor = color.Black
	colorBox.StrokeWidth = 1
	colorBox.CornerRadius = 10

	// Create label
	markerLabel := widget.NewLabel(label)
	markerLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Truncate address if too long
	displayAddress := address
	if len(displayAddress) > 30 {
		displayAddress = displayAddress[:27] + "..."
	}
	addressLabel := widget.NewLabel(displayAddress)

	return container.NewHBox(
		colorBox,
		markerLabel,
		addressLabel,
	)
}

func NRGBAToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("0x%02X%02X%02X", r, g, b)
}

func (mv *MapView) CreateColorMapping() {

	// ColorMap stores 13 distinct colors for vehicles plus brown for depot
	ColorMap := map[string]color.Color{
		"brown":   color.NRGBA{139, 69, 19, 255},   // Depot color
		"red":     color.NRGBA{255, 0, 0, 255},     // Vehicle 1
		"blue":    color.NRGBA{0, 0, 255, 255},     // Vehicle 2
		"green":   color.NRGBA{0, 128, 0, 255},     // Vehicle 3
		"orange":  color.NRGBA{255, 165, 0, 255},   // Vehicle 4
		"purple":  color.NRGBA{128, 0, 128, 255},   // Vehicle 5
		"cyan":    color.NRGBA{0, 255, 255, 255},   // Vehicle 6
		"magenta": color.NRGBA{255, 0, 255, 255},   // Vehicle 7
		"pink":    color.NRGBA{255, 192, 203, 255}, // Vehicle 8
		"teal":    color.NRGBA{0, 128, 128, 255},   // Vehicle 9
		"indigo":  color.NRGBA{75, 0, 130, 255},    // Vehicle 10
		"gold":    color.NRGBA{255, 215, 0, 255},   // Vehicle 11 // too similar to orange
		"lime":    color.NRGBA{50, 205, 50, 255},   // Vehicle 12

	}

	// VehicleColors is the ordered list of colors for vehicles (excluding depot)
	VehicleColors := []string{"red", "blue", "green", "orange", "purple", "cyan",
		"magenta", "pink", "teal", "indigo", "gold", "lime"}

	mv.colorMap = ColorMap
	mv.colors = VehicleColors

}

// buildMapURL constructs the Google Static Maps API URL
func (mv *MapView) buildMapURL() string {
	params := url.Values{}
	params.Set("size", fmt.Sprintf("%dx%d", mapWidth, mapHeight))
	params.Set("key", mv.apiKey)
	params.Set("maptype", "roadmap")

	// Add depot marker with brown color
	depotColor := mv.colorMap["brown"]
	params.Add("markers", fmt.Sprintf("color:%s|label:✞|%f,%f ", NRGBAToHex(depotColor), depotCoor.Lat, depotCoor.Long))

	// Add markers for each vehicle's route
	for i, vehicle := range mv.routingProcess.rm.Vehicles {
		if len(vehicle.Guests) > 0 {
			// Use vehicle color from the ordered list
			vehicleColor := mv.colors[i%len(mv.colors)]

			// Build markers string for this vehicle
			var markersList []string
			for j, loc := range vehicle.Locations {
				label := fmt.Sprintf("%c", 'A'+j)
				marker := fmt.Sprintf("label:%s|%f,%f", label, loc.Lat, loc.Long)
				markersList = append(markersList, marker)
			}

			if len(markersList) > 0 {
				markersParam := fmt.Sprintf("color:%s|%s", vehicleColor, strings.Join(markersList, "|"))
				params.Add("markers", markersParam)
			}
		}
	}

	// Auto-adjust zoom to show all markers
	params.Set("zoom", "11") // Default zoom, API will auto-adjust if needed

	return staticMapsBaseURL + "?" + params.Encode()
}

// loadMap fetches and displays the map from Google Static Maps API
func (mv *MapView) loadMap() {

	if mv.routingProcess == nil || mv.routingProcess.rm == nil {
		mv.showError("No route data available")
		return
	}

	// Build the map URL
	mapURL := mv.buildMapURL()

	// Fetch the image
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

	// Read image data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		mv.showError(fmt.Sprintf("Failed to read map data: %v", err))
		return
	}

	// Create resource and update image
	resource := fyne.NewStaticResource("map.png", data)

	// Update UI on main thread
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

// Refresh updates the map view with new data
func (mv *MapView) Refresh() {
	mv.legend = mv.createLegend()
	mv.mainContainer.Refresh()
	go mv.loadMap()
}

// mapViewRenderer implements the renderer for MapView
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
