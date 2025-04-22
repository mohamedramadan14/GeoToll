package main

import (
	"fmt"
	"math"

	"github.com/mohamedramadan14/roads-fees-system/types"
)

type DistanceServicer interface {
	CalculateDistance(types.OBUData, types.OBUData) (float64, error)
}

type calculatorService struct {
}

func NewDistanceService() DistanceServicer {
	return &calculatorService{}
}

func (s *calculatorService) CalculateDistance(lastOBUData, latestOBUData types.OBUData) (float64, error) {
	// Simulate distance calculation
	if lastOBUData.Lat == 0 && lastOBUData.Long == 0 {
		return 0, fmt.Errorf("invalid coordinates: lat and long cannot be zero")
	}
	distance := CalculateHaversineDistance(lastOBUData, latestOBUData)
	return distance, nil
}

func CalculateHaversineDistance(lastOBUData, latestOBUData types.OBUData) float64 {
	radius_km := 6371.0
	phi1 := lastOBUData.Lat * math.Pi / 180
	phi2 := latestOBUData.Lat * math.Pi / 180
	deltaPhi := (latestOBUData.Lat - lastOBUData.Lat) * math.Pi / 180
	deltaLambda := (latestOBUData.Long - lastOBUData.Long) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return radius_km * c
}
