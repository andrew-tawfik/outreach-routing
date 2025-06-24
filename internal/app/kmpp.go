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
	centroid  coordinates.GuestCoordinates
	vehicle   *Vehicle
	index     int
	bestThree []*Point
}

type Point struct {
	guestCoordinate coordinates.GuestCoordinates
	clusterIndex    int
}

func (km *Kmeans) GetName() string {
	return "Kmeans++"
}

func (km *Kmeans) StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error {

	// Create k means object
	km.init(rm)
	err := km.determineCentroids(rm)
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

func (km *Kmeans) init(rm *RouteManager) {
	totalDestinationCount := 0

	for range rm.CoordinateList {
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

func (km *Kmeans) determineCentroids(rm *RouteManager) error {
	km.retreiveUniqueGuestCoordinates(rm) // this is my data
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

func (km *Kmeans) retreiveUniqueGuestCoordinates(rm *RouteManager) {
	allPoints := make([]Point, 0)

	for _, gc := range rm.CoordinateList {
		allPoints = append(allPoints, Point{guestCoordinate: gc})
	}

	km.points = allPoints
}

func (km *Kmeans) clusterData() {
	maxIterations := 100 // Prevent infinite loops

	// Initialize all points to cluster -1 (unassigned)
	for i := range km.points {
		km.points[i].clusterIndex = -1
	}

	for iteration := 0; iteration < maxIterations; iteration++ {
		hasChanged := false

		// Clear existing assignments
		for i := range km.Clusters {
			km.Clusters[i].bestThree = make([]*Point, 0, 3)
		}

		// Assign each point to the closest cluster (with capacity constraint)
		for i := range km.points {
			point := &km.points[i]
			bestClusterIndex, replacePosition := km.findBestClusterForPoint(point)

			if bestClusterIndex != -1 {
				cluster := &km.Clusters[bestClusterIndex]

				if replacePosition != -1 {
					// Replace the point at the specified position
					cluster.bestThree[replacePosition] = point
				} else {
					// Add to cluster (has space)
					cluster.bestThree = append(cluster.bestThree, point)
				}

				if point.clusterIndex != bestClusterIndex {
					point.clusterIndex = bestClusterIndex
					hasChanged = true
				}
			}
		}

		// Recalculate centroids
		km.recalculateCentroids()

		// If no assignments changed, we've converged
		if !hasChanged {
			break
		}
	}
}

func (km *Kmeans) findBestClusterForPoint(point *Point) (int, int) {
	bestClusterIndex := -1
	replacePosition := -1
	minDistance := math.Inf(1)

	// Try to find the best available cluster
	for i := range km.Clusters {
		cluster := &km.Clusters[i]
		distance := sqDist(point.guestCoordinate, cluster.centroid)

		// If cluster has space and this is the closest so far
		if len(cluster.bestThree) < 3 && distance < minDistance {
			minDistance = distance
			bestClusterIndex = i
			replacePosition = -1 // No replacement needed, just append
		}
	}

	// If no cluster with space was found, try to replace in existing clusters
	if bestClusterIndex == -1 {
		for i := range km.Clusters {
			cluster := &km.Clusters[i]
			distance := sqDist(point.guestCoordinate, cluster.centroid)

			if len(cluster.bestThree) == 3 {
				// Find the farthest point in this cluster
				farthestIndex, farthestDistance := km.findFarthestPointInCluster(cluster)

				// If current point is closer than the farthest point in cluster
				if distance < farthestDistance && distance < minDistance {
					minDistance = distance
					bestClusterIndex = i
					replacePosition = farthestIndex
				}
			}
		}
	}

	return bestClusterIndex, replacePosition
}

// Case 1: Cluster has less than 3. Add the points
// Case 2: Cluster already has 3.
// 		2a. If farther than the 3: Skip this Cluster look for the next smallest and available Cluster
// 		2b. If closer than one of the three, replace the farthest one of the three

func (km *Kmeans) findFarthestPointInCluster(cluster *Cluster) (int, float64) {
	farthestIndex := -1
	farthestDistance := 0.0

	for i, p := range cluster.bestThree {
		distance := sqDist(p.guestCoordinate, cluster.centroid)
		if distance > farthestDistance {
			farthestDistance = distance
			farthestIndex = i
		}
	}

	return farthestIndex, farthestDistance
}

func (km *Kmeans) recalculateCentroids() {
	for i := range km.Clusters {
		cluster := &km.Clusters[i]

		if len(cluster.bestThree) == 0 {
			continue // Keep existing centroid if no points assigned
		}

		var sumLat, sumLong float64
		for _, point := range cluster.bestThree {
			sumLat += point.guestCoordinate.Lat
			sumLong += point.guestCoordinate.Long
		}

		// Update centroid to average of assigned points
		cluster.centroid.Lat = sumLat / float64(len(cluster.bestThree))
		cluster.centroid.Long = sumLong / float64(len(cluster.bestThree))
	}
}

func (km *Kmeans) determineVehicleRoutes() {
	for i, point := range km.points {
		if point.clusterIndex >= 0 && point.clusterIndex < len(km.Clusters) {
			cluster := &km.Clusters[point.clusterIndex]
			cluster.vehicle.Route.List.PushBack(i)
		}
	}
}
