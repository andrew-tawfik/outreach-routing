package app

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
	"github.com/mroth/weightedrand"
)

type Kmeans struct {
	Clusters []Cluster
}

type Cluster struct {
	centroid coordinates.GuestCoordinates
	vehicle  *Vehicle
}

func (km *Kmeans) GetName() string {
	return "Kmeans++"
}

func (km *Kmeans) StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error {

	// Create k means object
	km.init(lr, rm)
	err := km.determineCentroids(lr)
	if err != nil {
		return err
	}

	// perform standard k means

	// traveling sales person algorithm per cluster

	// vehicles will serve one cluster, update vehicle information

	return nil
}

func (km *Kmeans) init(lr *LocationRegistry, rm *RouteManager) {
	totalDestinationCount := 0

	for range lr.CoordianteMap.AddressOrder {
		totalDestinationCount += 1
	}

	vehicleCount := (totalDestinationCount + 2) / 3

	for i := 0; i < vehicleCount; i++ {
		newVehicle := Vehicle{}
		newVehicle.Route.List = list.New()
		rm.Vehicles = append(rm.Vehicles, newVehicle)

		newCluster := Cluster{vehicle: &rm.Vehicles[i]}
		km.Clusters = append(km.Clusters, newCluster)
	}
}

func (km *Kmeans) determineCentroids(lr *LocationRegistry) error {
	allGuestCoordinates := retreiveUniqueGuestCoordinates(lr) // this is my data
	centroids := make([]*coordinates.GuestCoordinates, 0)

	randomIndex := rand.Intn(len(allGuestCoordinates))
	km.Clusters[0].centroid = allGuestCoordinates[randomIndex]
	centroids = append(centroids, &km.Clusters[0].centroid) // append random element to centroids list

	for i := 1; i < len(km.Clusters); i++ {

		minSqDistances := make([]float64, 0, len(allGuestCoordinates))
		for _, x := range allGuestCoordinates {
			min := 9999.99
			for _, c := range centroids {
				sq_dist := sqDist(x, *c)
				if sq_dist < min {
					min = sq_dist
				}
			}
			minSqDistances = append(minSqDistances, min)
		}

		sumSqDistances := sumSqDist(minSqDistances)

		probabilities := getProbabilities(minSqDistances, sumSqDistances)

		randomGCIndex, err := km.getRandomCoordinateIndex(probabilities)
		if err != nil {
			return fmt.Errorf("unexpected error in centroid determination: %v", err)
		}
		km.Clusters[i].centroid = allGuestCoordinates[randomGCIndex]
		centroids = append(centroids, &km.Clusters[i].centroid)
	}
	return nil

}

func getProbabilities(distances []float64, divisor float64) []uint {
	probabilities := make([]uint, len(distances))
	sum := uint(0)

	// Convert all but the last element
	for i := 0; i < len(distances)-1; i++ {
		probabilities[i] = uint((distances[i] / divisor) * 100)
		sum += probabilities[i]
	}

	// Make the last element ensure sum equals 100
	probabilities[len(probabilities)-1] = 100 - sum
	sum += probabilities[len(probabilities)-1]

	if sum != 100 {
		panic(fmt.Sprintf("probabilities sum to %d, expected 100", sum))
	}

	return probabilities
}

func (km *Kmeans) getRandomCoordinateIndex(probabilities []uint) (int, error) {
	choices := make([]weightedrand.Choice, 0, len(probabilities))
	for i := range probabilities {
		choice := weightedrand.Choice{Item: i, Weight: probabilities[i]}
		choices = append(choices, choice)
	}

	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		return -1, err
	}

	selected := chooser.Pick()
	index, ok := selected.(int)
	if !ok {
		return -1, fmt.Errorf("failed to convert selected item to GuestCoordinates")
	}
	return index, nil
}

func sqDist(point, centroid coordinates.GuestCoordinates) float64 {
	x1, y1 := point.Long, point.Lat
	x2, y2 := centroid.Long, centroid.Lat

	dist := (math.Pow(x2-x1, 2)) + (math.Pow(y2-y1, 2))

	return dist
}

func sumSqDist(numbers []float64) float64 {
	var sum float64 = 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func retreiveUniqueGuestCoordinates(lr *LocationRegistry) []coordinates.GuestCoordinates {
	allGuestCoordinates := make([]coordinates.GuestCoordinates, 0)

	for gc := range lr.CoordianteMap.DestinationOccupancy {
		allGuestCoordinates = append(allGuestCoordinates, gc)
	}
	return allGuestCoordinates
}
