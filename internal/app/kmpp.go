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
	points   []Point
}

type Cluster struct {
	centroid coordinates.GuestCoordinates
	vehicle  *Vehicle
	index    int
}

type Point struct {
	guestCoordinate coordinates.GuestCoordinates
	cluster         *Cluster
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
	km.clusterData()
	km.determineVehicleRoutes()

	// traveling sales person algorithm per cluster

	// vehicles will serve one cluster, update vehicle information

	return nil
}

func (km *Kmeans) init(lr *LocationRegistry, rm *RouteManager) {
	totalDestinationCount := 0

	for range lr.CoordianteMap.DestinationOccupancy {
		totalDestinationCount += 1
	}

	vehicleCount := (totalDestinationCount + 2) / 3

	for i := 0; i < vehicleCount; i++ {
		newVehicle := Vehicle{Route: Route{List: list.New()}}
		rm.Vehicles = append(rm.Vehicles, newVehicle)

		newCluster := Cluster{vehicle: &rm.Vehicles[i], index: i}
		km.Clusters = append(km.Clusters, newCluster)
	}
}

func (km *Kmeans) determineCentroids(lr *LocationRegistry) error {
	km.retreiveUniqueGuestCoordinates(lr) // this is my data
	centroids := make([]*coordinates.GuestCoordinates, 0)

	randomIndex := rand.Intn(len(km.points))
	km.Clusters[0].centroid = km.points[randomIndex].guestCoordinate
	centroids = append(centroids, &km.Clusters[0].centroid) // append random element to centroids list

	for i := 1; i < len(km.Clusters); i++ {

		minSqDistances := make([]float64, 0, len(km.points))
		for _, x := range km.points {
			min := 9999.99
			for _, c := range centroids {
				sq_dist := sqDist(x.guestCoordinate, *c)
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
		km.Clusters[i].centroid = km.points[randomGCIndex].guestCoordinate
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

func (km *Kmeans) retreiveUniqueGuestCoordinates(lr *LocationRegistry) {
	allPoints := make([]Point, 0)

	for gc := range lr.CoordianteMap.DestinationOccupancy {
		allPoints = append(allPoints, Point{guestCoordinate: gc})
	}

	km.points = allPoints
}

func (km *Kmeans) clusterData() {
	hasChanged := true

	// populate the clusters
	for hasChanged {
		for i, p := range km.points {

			min := math.Inf(1)
			var closestCluster int
			for _, cluster := range km.Clusters {
				distance := sqDist(p.guestCoordinate, cluster.centroid)
				if distance < min {
					min = distance
					closestCluster = cluster.index
				}
			}

			if km.points[i].cluster == nil || km.points[i].cluster.index != closestCluster {
				km.points[i].cluster = &km.Clusters[closestCluster]
				hasChanged = true
			} else {
				hasChanged = false
			}
		}
	}
}

func (km *Kmeans) determineVehicleRoutes() {
	for i, p := range km.points {
		p.cluster.vehicle.Route.List.PushBack(i)
		fmt.Printf("\nAdded Point %d to Vehicle %d", i, p.cluster.index)
	}
}
