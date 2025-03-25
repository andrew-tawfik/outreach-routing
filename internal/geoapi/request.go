package geoapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Coordinates struct {
	Long float64
	Lat  float64
}

// Two General Items to request

// First we need the coordinates of each address
func GeocodeAddress(address string) (Coordinates, error) {
	polishAddress(&address)
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=geojson", address)

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Coordinates{}, err
	}
	req.Header.Set("User-Agent", "outreach-routing/1.0")

	resp, err := client.Do(req)

	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Coordinates{}, err
	}
	fmt.Println("Response from nominatim", string(body))

	return Coordinates{}, nil
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
