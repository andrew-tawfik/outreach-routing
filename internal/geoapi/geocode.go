package geoapi


type GoogleGeocodeResponse struct {
	Results []GeocodeResult `json:"results"`
	Status  string          `json:"status"`
}


type GeocodeResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          GeocodeGeometry    `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
}


type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}


type GeocodeGeometry struct {
	Location     LatLng    `json:"location"`
	LocationType string    `json:"location_type"`
	Viewport     Viewport  `json:"viewport"`
}


type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}


type Viewport struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

