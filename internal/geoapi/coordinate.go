package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// retreiveAddressCoordinate takes a raw address string, sends a geocode request to Nominatim,
// and returns the best-matched geographic coordinates (longitude, latitude).
func retreiveAddressCoordinate(address string) (coordinates.GuestCoordinates, error) {
	url := buildGeocodeURL(address)

	body, err := fetchGeocodeData(url)
	if err != nil {
		return coordinates.GuestCoordinates{}, err
	}

	coor, err := parseGeocodeResponse(body, "Ottawa") // Match based on city keyword
	if err != nil {
		return coordinates.GuestCoordinates{}, err
	}

	return coordinates.GuestCoordinates{Long: coor[0], Lat: coor[1]}, nil
}

// geocodeGuestAddress assigns GPS coordinates to the guest's address
func (g *Guest) geocodeGuestAddress() error {
	gc, err := retreiveAddressCoordinate(g.Address)
	if err != nil {
		return err
	}
	g.Coordinates = gc
	return nil
}

// buildGeocodeURL creates a Nominatim search URL from a sanitized address.
func buildGeocodeURL(address string) string {
	polishAddress(&address)
	return fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=geojson", address)
}

// fetchGeocodeData performs an HTTP GET request to the given URL with retry logic.
func fetchGeocodeData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "outreach-routing/1.0")

	resp, err := sendWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}
	return body, nil
}

// parseGeocodeResponse extracts coordinates from the raw response,
// filtered by a city keyword (e.g., "Ottawa") to retreive accurate location.
func parseGeocodeResponse(body []byte, city string) ([]float64, error) {
	var nr NominatimResponse
	if err := json.Unmarshal(body, &nr); err != nil {
		return nil, fmt.Errorf("could not deserialize response body: %v", err)
	}

	coordinates, err := nr.locateCoordinatesByKeyword(city)
	if err != nil {
		return nil, fmt.Errorf("could not extract coordinates, %v", err)
	}
	return coordinates, nil
}

// polishAddress cleans and formats an address string to make it suitable for OSRM request.
// - Removes apartment/unit suffixes.
// - Trims spaces and converts to lowercase.
// - Replaces spaces with '+' for URL encoding.
func polishAddress(rawAddress *string) {
	address := strings.ToLower(*rawAddress)

	if i := strings.Index(address, "apt"); i != -1 {
		address = address[:i]
	} else if i := strings.Index(address, "unit"); i != -1 {
		address = address[:i]
	}

	address = strings.TrimSpace(address)
	address = strings.ReplaceAll(address, " ", "+")

	*rawAddress = address
}

// sendWithRetry executes an HTTP request with retry logic and returns the response.
func sendWithRetry(req *http.Request) (*http.Response, error) {
	// httpClient is a shared HTTP client with timeout, reused for all geocoding requests.
	const maxAttempts = 3

	for attempts := 0; attempts < maxAttempts; attempts++ {
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("HTTP request failed for %s: %v\n", req.URL.String(), err)
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil
		} else {
			// Retry after short wait
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		// Don't retry on unexpected status
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil, fmt.Errorf("request to %s failed after %d attempts", req.URL.String(), maxAttempts)
}
