package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andrew-tawfik/outreach-routing/internal/coordinates"
)

type AppData struct {
	Event            Event
	LocationRegistry LocationRegistry
}

func CoordinateKey(c coordinates.GuestCoordinates) string {
	return fmt.Sprintf("%.6f,%.6f", c.Long, c.Lat)
}

type SerializableCoordinateMapping struct {
	DestinationOccupancy map[string]int
	CoordinateToAddress  map[string]string
	AddressOrder         []string
}

type SerializableLocationRegistry struct {
	DistanceMatrix [][]float64
	CoordinateMap  SerializableCoordinateMapping
}

type SerializableAppData struct {
	Event            Event
	LocationRegistry SerializableLocationRegistry
}

func ConvertToSerializable(lr LocationRegistry) SerializableLocationRegistry {
	occupancy := make(map[string]int)
	for coord, count := range lr.CoordianteMap.DestinationOccupancy {
		key := CoordinateKey(coord)
		occupancy[key] = count
	}

	address := make(map[string]string)
	for addr, coord := range lr.CoordianteMap.CoordinateToAddress {
		key := addr
		address[key] = CoordinateKey(coord)
	}

	return SerializableLocationRegistry{
		DistanceMatrix: lr.DistanceMatrix,
		CoordinateMap: SerializableCoordinateMapping{
			DestinationOccupancy: occupancy,
			CoordinateToAddress:  address,
			AddressOrder:         lr.CoordianteMap.AddressOrder,
		},
	}
}

func ConvertFromSerializable(slr SerializableLocationRegistry) LocationRegistry {
	reverseOccupancy := make(map[coordinates.GuestCoordinates]int)
	for key, count := range slr.CoordinateMap.DestinationOccupancy {
		var lat, long float64
		fmt.Sscanf(key, "%f,%f", &long, &lat)
		reverseOccupancy[coordinates.GuestCoordinates{Long: long, Lat: lat}] = count
	}

	reverseAddress := make(map[string]coordinates.GuestCoordinates)
	for addr, coord := range slr.CoordinateMap.CoordinateToAddress {
		var lat, long float64
		fmt.Sscanf(coord, "%f,%f", &long, &lat)
		reverseAddress[addr] = coordinates.GuestCoordinates{Long: long, Lat: lat}
	}

	return LocationRegistry{
		DistanceMatrix: slr.DistanceMatrix,
		CoordianteMap: CoordinateMapping{
			DestinationOccupancy: reverseOccupancy,
			CoordinateToAddress:  reverseAddress,
			AddressOrder:         slr.CoordinateMap.AddressOrder,
		},
	}
}

func SaveAppDataToFile(filename string, event Event, lr LocationRegistry) error {
	serializable := SerializableAppData{
		Event:            event,
		LocationRegistry: ConvertToSerializable(lr),
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(serializable)
}

func LoadAppDataFromFile(filename string) (Event, LocationRegistry, error) {
	var serializable SerializableAppData

	file, err := os.Open(filename)
	if err != nil {
		return Event{}, LocationRegistry{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&serializable)
	if err != nil {
		return Event{}, LocationRegistry{}, err
	}

	lr := ConvertFromSerializable(serializable.LocationRegistry)
	return serializable.Event, lr, nil
}
