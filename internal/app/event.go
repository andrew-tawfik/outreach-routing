package app

type LocationRegistry struct {
	DistanceMatrix [][]float64
	CoordianteMap  CoordinateMapping
}

type CoordinateMapping struct {
	DestinationOccupancy map[GuestCoordinates]int
	CoordinateToAddress  map[GuestCoordinates]string
	AddressOrder         []string
}

type GuestCoordinates struct {
	Long float64
	Lat  float64
}

type Event struct {
	Guests    []Guest
	EventType string
}

type Guest struct {
	Name        string
	GroupSize   int
	Coordinates GuestCoordinates
}
