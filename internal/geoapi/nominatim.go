package geoapi

import (
	"fmt"
	"strings"
)

type NominatimResponse struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	DisplayName string `json:"display_name"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
}

func (nr *NominatimResponse) locateCoordinatesByKeyword(keyword string) (coordinates []float64, err error) {
	for _, f := range nr.Features {
		if strings.Contains(f.Properties.DisplayName, keyword) {
			return f.Geometry.Coordinates, nil
		}
	}
	return coordinates, fmt.Errorf("no feature associated with %s", keyword)
}
