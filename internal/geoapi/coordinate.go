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

func (g *Guest) geocodeGuestAddress() error {

	apiKey, err := getApiKey()
	if err != nil {
		return err
	}
	gc, newAddr, err := retreiveGuestLocation(g.Address, apiKey)
	// TODO  add the formatted address as this will be a more complete address string
	if err != nil {
		fmt.Println("there was an error: ", err)
		return err
	}
	g.Coordinates = gc
	g.Address = newAddr
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

func sendWithRetry(req *http.Request) (*http.Response, error) {
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
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil, fmt.Errorf("request to %s failed after %d attempts", req.URL.String(), maxAttempts)
}
