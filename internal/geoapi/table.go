package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// OSRMResponse represents the JSON structure returned by OSRM's /table API.
type OSRMResponse struct {
	Sources   []Source    `json:"sources"`   // List of input points (with coordinates and names)
	Distances [][]float64 `json:"distances"` // Matrix of distances between each pair of points
	Status    string      `json:"code"`      // Should be "Ok" if request was successful
}

// Source represents a single input point in the OSRM response.
type Source struct {
	Location []float64 `json:"location"`
	Name     string    `json:"name"`
}

// RetreiveDistanceMatrix constructs and sends a request to OSRM.
// then parses the resulting distance matrix and stores it in the Event.
func (e *Event) RetreiveDistanceMatrix() {
	coordListURL = strings.TrimSuffix(coordListURL, ";")
	url := buildDistanceMatrixURL(&coordListURL)
	jsonresp, err := fetchDistanceMatrix(&url)
	if err != nil {
		log.Fatalf("%v", err)
	}

	matrix, err := parseOsrmResponse(&jsonresp)

	if err != nil {
		log.Fatalf("%v", err)
	}

	e.GuestLocations.DistanceMatrix = matrix

}

// buildDistanceMatrixURL returns a fully constructed OSRM request URL
func buildDistanceMatrixURL(coordinatesList *string) string {
	url := fmt.Sprintf("http://router.project-osrm.org/table/v1/driving/%s?annotations=distance", *coordinatesList)
	return url
}

// fetchDistanceMatrix sends a GET request to OSRM and returns the response body.
func fetchDistanceMatrix(url *string) ([]byte, error) {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "outreach-routing/1.0")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return nil, fmt.Errorf("osrm server: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}

	return body, nil
}

// parseOsrmResponse parses the OSRM JSON response body into a 2D distance matrix.
func parseOsrmResponse(body *[]byte) ([][]float64, error) {
	var osrm OSRMResponse
	if err := json.Unmarshal(*body, &osrm); err != nil {
		return nil, fmt.Errorf("could not deserialize response body: %v", err)
	}

	if osrm.Status != "Ok" {
		return nil, fmt.Errorf("OSRM Status: %s", osrm.Status)
	}

	return osrm.Distances, nil

}
