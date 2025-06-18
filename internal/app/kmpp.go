package app

import "github.com/andrew-tawfik/outreach-routing/internal/coordinates"

type Kmeans struct {
	Clusters []Cluster
}

type Cluster struct {
	centroid *coordinates.GuestCoordinates
	vehicle  *Vehicle
}

func (km *Kmeans) GetName() string {
	return "Kmeans++"
}

func (km *Kmeans) StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error {

	// Determine centroids (++)
	//centroids := km.determineCentroids(lr)

	// perform standard k means

	// traveling sales person algorithm per cluster

	// vehicles will serve one cluster, update vehicle information

	return nil
}

func (km *Kmeans) determineCentroids(lr *LocationRegistry) {
	// allGuestCoordinates := km.retreiveUniqueGuestCoordinates(lr)

}

func (km *Kmeans) retreiveUniqueGuestCoordinates(lr *LocationRegistry) []coordinates.GuestCoordinates {
	allGuestCoordinates := make([]coordinates.GuestCoordinates, 0)

	for gc := range lr.CoordianteMap.DestinationOccupancy {
		allGuestCoordinates = append(allGuestCoordinates, gc)
	}
	return allGuestCoordinates
}
