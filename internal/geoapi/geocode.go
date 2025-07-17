package geoapi

// GoogleGeocodeResponse represents the response from Google Maps Geocoding API
type GoogleGeocodeResponse struct {
	Results []GeocodeResult `json:"results"`
	Status  string          `json:"status"`
}

// GeocodeResult represents a single geocoding result
type GeocodeResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          GeocodeGeometry    `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
}

// AddressComponent represents a component of the address
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// GeocodeGeometry contains the location information
type GeocodeGeometry struct {
	Location     LatLng    `json:"location"`
	LocationType string    `json:"location_type"`
	Viewport     Viewport  `json:"viewport"`
}

// LatLng represents latitude and longitude coordinates
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Viewport represents the recommended viewport for displaying the result
type Viewport struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

