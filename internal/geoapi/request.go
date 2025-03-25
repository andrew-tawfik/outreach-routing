package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type GuestCoordinates struct {
	Long float64
	Lat  float64
}

// Two General Items to request

// First we need the coordinates of each address
func (g *Guest) GeocodeAddress() error {
	address := g.Address
	polishAddress(&address)
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=geojson", address)

	client := &http.Client{Timeout: 20 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "outreach-routing/1.0")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// fmt.Println("Response from nominatim", string(body))

	var nr NominatimResponse

	// unmashall the response body into a Features struct
	if err := json.Unmarshal(body, &nr); err != nil {
		log.Fatalf("could not deserialize response body")
	}
	city := "Ottawa"
	// Search within the the struct to extract for Ottawa address Coordinates
	coordinates, err := nr.locateCoordinatesByKeyword(city)
	if err != nil {
		return fmt.Errorf("could not extract coordinates: %v", err)
	}
	fmt.Println(coordinates)
	g.Coordinates = GuestCoordinates{Long: coordinates[0], Lat: coordinates[1]}

	// if found update GuestCoordinates
	return nil
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
