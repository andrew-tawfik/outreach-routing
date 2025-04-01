package geoapi

import (
	"fmt"
	"strings"
)

// NominatimResponse represents the response from the Nominatim API
type NominatimResponse struct {
	Features []Feature `json:"features"`
}

// Feature represents one individual match in the Nominatim results,
type Feature struct {
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

// Properties includes name of location
type Properties struct {
	DisplayName string `json:"display_name"`
}

// Geometry stores the geographic coordinates for a location in [longitude, latitude] order.
type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
}

// locateCoordinatesByKeyword searches the list of geocoded features for a match
// containing the given keyword (e.g., "Ottawa") in the display name.
func (nr *NominatimResponse) locateCoordinatesByKeyword(keyword string) (coordinates []float64, err error) {
	for _, f := range nr.Features {
		if strings.Contains(f.Properties.DisplayName, keyword) {
			return f.Geometry.Coordinates, nil
		}
	}
	return coordinates, fmt.Errorf("no feature associated with %s", keyword)
}
