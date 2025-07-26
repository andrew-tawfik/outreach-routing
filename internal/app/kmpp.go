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
	address         string
	clusterIndex    int
}

func (km *Kmeans) GetName() string {
	return "Kmeans++"
}

func (km *Kmeans) StartRouteDispatch(rm *RouteManager, lr *LocationRegistry) error {

	
	km.init(rm)
	err := km.determineCentroids(rm, lr)
	if err != nil {
		return err
	}

	
	km.clusterData()
	km.determineVehicleRoutes()
	

	

	

	return nil
}

func (km *Kmeans) init(rm *RouteManager) {
	totalDestinationCount := 0

	for range rm.CoordinateList {
		totalDestinationCount += 1
	}

	vehicleCount := ((totalDestinationCount + 2) / 3) + 2

	for i := 0; i < vehicleCount; i++ {
		newVehicle := Vehicle{Route: Route{List: list.New()}}
		rm.Vehicles = append(rm.Vehicles, newVehicle)

		newCluster := Cluster{vehicle: &rm.Vehicles[i], index: i}
		km.Clusters = append(km.Clusters, newCluster)
	}
}

func (km *Kmeans) determineCentroids(rm *RouteManager, lr *LocationRegistry) error {
	km.retreiveUniqueGuestCoordinates(rm, lr) 
	centroids := make([]*coordinates.GuestCoordinates, 0)

	randomIndex := rand.Intn(len(km.points))
	km.Clusters[0].centroid = km.points[randomIndex].guestCoordinate
	centroids = append(centroids, &km.Clusters[0].centroid) 

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

	
	for i := 0; i < len(distances)-1; i++ {
		probabilities[i] = uint((distances[i] / divisor) * 100)
		sum += probabilities[i]
	}

	
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

func dist(point, centroid coordinates.GuestCoordinates) float64 {
	x1, y1 := point.Long, point.Lat
	x2, y2 := centroid.Long, centroid.Lat

	dist := math.Sqrt((math.Pow(x2-x1, 2)) + (math.Pow(y2-y1, 2)))

	return dist
}

func sumSqDist(numbers []float64) float64 {
	var sum float64 = 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func (km *Kmeans) retreiveUniqueGuestCoordinates(rm *RouteManager, lr *LocationRegistry) {
	allPoints := make([]Point, 0)

	i := 1
	for _, gc := range rm.CoordinateList {
		addr := lr.CoordianteMap.AddressOrder[i]
		allPoints = append(allPoints, Point{guestCoordinate: gc, address: addr})
		i++
	}

	km.points = allPoints
}

func (km *Kmeans) clusterData() {
	maxIterations := 100 

	
	for i := range km.points {
		km.points[i].clusterIndex = -1
	}

	for iteration := 0; iteration < maxIterations; iteration++ {
		hasChanged := false

		
		for i := range km.Clusters {
			km.Clusters[i].bestThree = make([]*Point, 0, 3)
		}

		
		for i := range km.points {
			point := &km.points[i]
			bestClusterIndex, replacePosition := km.findBestClusterForPoint(point)

			if bestClusterIndex != -1 {
				cluster := &km.Clusters[bestClusterIndex]

				if replacePosition != -1 {
					
					cluster.bestThree[replacePosition] = point
				} else {
					
					cluster.bestThree = append(cluster.bestThree, point)
				}

				if point.clusterIndex != bestClusterIndex {
					point.clusterIndex = bestClusterIndex
					hasChanged = true
				}
			}
		}

		
		km.recalculateCentroids()

		
		if !hasChanged {
			break
		}
	}
}

func (km *Kmeans) findBestClusterForPoint(point *Point) (int, int) {
	bestClusterIndex := -1
	replacePosition := -1
	minDistance := math.Inf(1)

	
	for i := range km.Clusters {
		cluster := &km.Clusters[i]
		distance := dist(point.guestCoordinate, cluster.centroid)

		
		if len(cluster.bestThree) < 3 && distance < minDistance {
			minDistance = distance
			bestClusterIndex = i
			replacePosition = -1 
		}
	}

	
	if bestClusterIndex == -1 {
		for i := range km.Clusters {
			cluster := &km.Clusters[i]
			distance := dist(point.guestCoordinate, cluster.centroid)

			if len(cluster.bestThree) == 3 {
				
				farthestIndex, farthestDistance := km.findFarthestPointInCluster(cluster)

				
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






func (km *Kmeans) findFarthestPointInCluster(cluster *Cluster) (int, float64) {
	farthestIndex := -1
	farthestDistance := 0.0

	for i, p := range cluster.bestThree {
		distance := dist(p.guestCoordinate, cluster.centroid)
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
			continue 
		}

		var sumLat, sumLong float64
		for _, point := range cluster.bestThree {
			sumLat += point.guestCoordinate.Lat
			sumLong += point.guestCoordinate.Long
		}

		
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

func (km *Kmeans) printDistances() {
	for _, c := range km.Clusters {

		fmt.Printf("The distance between Centroid of Cluster %d: ", c.index)
		for _, p := range c.bestThree {
			f := dist(p.guestCoordinate, c.centroid)
			fmt.Printf("\n\twith point at address %s %f", p.address, f)
		}
		fmt.Println()
	}
}
