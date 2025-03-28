package geoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type LocationRegistry struct {
	DistanceMatrix    [][]float64
	GuestCountByCoord map[GuestCoordinates]int
	CoordListString   string
}

type OSRMResponse struct {
	Sources   []Source    `json:"sources"`
	Distances [][]float64 `json:"distances"`
	Status    string      `json:"code"`
}

type Source struct {
	Location []float64 `json:"location"`
	Name     string    `json:"name"`
}

func (e *Event) RetreiveDistanceMatrix() {
	e.GuestLocations.CoordListString = strings.TrimSuffix(e.GuestLocations.CoordListString, ";")
	url := buildDistanceMatrixURL(&e.GuestLocations.CoordListString)
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

func buildDistanceMatrixURL(coordinatesList *string) string {
	url := fmt.Sprintf("http://router.project-osrm.org/table/v1/driving/%s?annotations=distance", *coordinatesList)
	return url
}

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
