package geoapi

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
