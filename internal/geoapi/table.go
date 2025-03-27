package geoapi

type LocationRegistry struct {
	DistanceMatrix    [][]float64
	GuestCountByCoord map[GuestCoordinates]int
	CoordListString   string
}
