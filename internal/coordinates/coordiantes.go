package coordinates

import "fmt"

type GuestCoordinates struct {
	Long float64
	Lat  float64
}

func (gc *GuestCoordinates) ToString() string {
	return fmt.Sprintf("%f,%f;", gc.Long, gc.Lat)
}
