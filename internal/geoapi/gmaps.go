package geoapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

var geocodeMapsBaseURL string = "https://maps.googleapis.com/maps/api/geocode/json"

func (e *Event) geocodeEvent() {

	apiErrors := ApiErrors{
		FailedGuests: make([]FailedGuest, 0),
	}

	// Geocode all guests and track unique coordinates
	for i := range e.Guests {
		err := e.Guests[i].geocodeGuestAddress()
		if err != nil {

			apiErrors.FailedGuests = append(apiErrors.FailedGuests, FailedGuest{
				Name:    e.Guests[i].Name,
				Address: e.Guests[i].Address,
				Reason:  err.Error(),
			})
			continue
		}
		coor, unique := e.isUnique(i)
		if unique {
			addToCoordListString(&coor)
		}
	}

	e.ApiErrors = apiErrors

}

func buildGeoMapURL(address, apiKey string) string {

	params := url.Values{}
	address += " Ottawa, ON, Canada"
	params.Set("address", address)
	params.Set("key", apiKey)

	completeURL := geocodeMapsBaseURL + "?" + params.Encode()
	return completeURL
}

func retreiveGuestLocation(gAddress, apiKey string) (coordinates.GuestCoordinates, string, error) {
	url := buildGeoMapURL(gAddress, apiKey)

	body, err := fetchGeocodeData(url)
	if err != nil {
		return coordinates.GuestCoordinates{}, "", err
	}

	coor, newAddr, err := parseGoogleGeocodeResponse(body) // Match based on city keyword
	if err != nil {
		return coordinates.GuestCoordinates{}, "", err
	}

	return coordinates.GuestCoordinates{Long: coor[0], Lat: coor[1]}, newAddr, nil
}

func getApiKey() (string, error) {
	projectRoot, err := filepath.Abs(filepath.Join(".", ".."))
	if err != nil {
		return "", fmt.Errorf("failed to resolve project root:", err)
	}
	credentialsPath := filepath.Join(projectRoot, "maps_config.json")

	apiKey, jsonErr := LoadMapsConfig(credentialsPath)
	if jsonErr != nil {
		return "", fmt.Errorf("failed to load api key:", err)
	}
	return apiKey, nil
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

func parseGoogleGeocodeResponse(body []byte) ([]float64, string, error) {
	var response GoogleGeocodeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, "", fmt.Errorf("could not deserialize response body: %v", err)
	}

	// Check if the request was successful
	if response.Status != "OK" {
		return nil, "", fmt.Errorf("geocoding failed with status: %s", response.Status)
	}

	// Check if we have results
	if len(response.Results) == 0 {
		return nil, "", fmt.Errorf("no geocoding results found")
	}

	// Get the first (best) result
	result := response.Results[0]

	// Extract coordinates - Google returns lat, lng but your system expects [lat, lng]
	coordinates := []float64{
		result.Geometry.Location.Lng,
		result.Geometry.Location.Lat,
	}

	formatted_address := result.FormattedAddress

	return coordinates, formatted_address, nil
}
