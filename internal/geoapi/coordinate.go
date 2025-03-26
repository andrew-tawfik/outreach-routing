package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type GuestCoordinates struct {
	Long float64
	Lat  float64
}

// Global Client variable to handle requests
var httpClient = &http.Client{Timeout: 10 * time.Second}

// Two General Items to request

// First we need the coordinates of each address
func (g *Guest) geocodeAddress() error {
	url := buildGeocodeURL(g.Address)

	body, err := fetchGeocodeData(url)
	if err != nil {
		return err
	}

	coordinates, err := parseGeocodeResponse(body, "Ottawa")
	if err != nil {
		return err
	}

	g.Coordinates = GuestCoordinates{Long: coordinates[0], Lat: coordinates[1]}
	return nil
}

func buildGeocodeURL(address string) string {
	polishAddress(&address)
	return fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=geojson", address)
}

func fetchGeocodeData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "outreach-routing/1.0")

	var resp *http.Response
	for attempts := 0; attempts < 3; attempts++ { // Retry up to 3 times
		resp, err = httpClient.Do(req)
		if err != nil {
			fmt.Printf("HTTP request failed for %s: %v\n", url, err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusServiceUnavailable {
			time.Sleep(2 * time.Second)
			continue
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}
	return body, nil
}

func parseGeocodeResponse(body []byte, city string) ([]float64, error) {
	var nr NominatimResponse
	if err := json.Unmarshal(body, &nr); err != nil {
		return nil, fmt.Errorf("could not deserialize response body: %v", err)
	}

	coordinates, err := nr.locateCoordinatesByKeyword(city)
	if err != nil {
		return nil, fmt.Errorf("could not extract coordinates: %v", err)
	}
	return coordinates, nil
}

func polishAddress(rawAddress *string) {
	address := strings.ToLower(*rawAddress)

	// Remove anything after "apt" or "unit"
	if i := strings.Index(address, "apt"); i != -1 {
		address = address[:i]
	} else if i := strings.Index(address, "unit"); i != -1 {
		address = address[:i]
	}

	// Remove leading/trailing whitespace
	address = strings.TrimSpace(address)

	// Replace spaces with '+'
	address = strings.ReplaceAll(address, " ", "+")

	*rawAddress = address
}
