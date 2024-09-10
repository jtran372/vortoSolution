package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"vrpSolution/internal/models"
	"vrpSolution/pkg/pathcalculator"
	"vrpSolution/pkg/txtparser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Error: You must provide the file path as a command-line argument.")
	}
	filePath := os.Args[1]

	drivingRoutes, err := txtparser.ParseTextFile(filePath)
	if err != nil {
		fmt.Println("Error parsing Txt:", err)
		return
	}

	// find max width, height routes
	minX, minY, maxX, maxY := 0.0, 0.0, 0.0, 0.0
	for _, dr := range drivingRoutes {
		if math.Max(dr.PickUp.X, dr.DropOff.X) > maxX {
			maxX = math.Max(dr.PickUp.X, dr.DropOff.X) + 1
		}
		if math.Max(dr.PickUp.Y, dr.DropOff.Y) > maxY {
			maxY = math.Max(dr.PickUp.Y, dr.DropOff.Y) + 1
		}
		if math.Min(dr.PickUp.X, dr.DropOff.X) < minX {
			minX = math.Min(dr.PickUp.X, dr.DropOff.X) - 1
		}
		if math.Min(dr.PickUp.Y, dr.DropOff.Y) < minY {
			minY = math.Min(dr.PickUp.Y, dr.DropOff.Y) - 1
		}
	}

	qt := models.NewQuadTree(models.Bounds{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY})

	for _, r := range drivingRoutes {
		qt.Insert(r)
	}

	// Simulate average target work hours to determine least cost
	determinedDrivingRoutes := pathcalculator.CalculateOptimalPaths(qt, drivingRoutes)

	for _, routes := range determinedDrivingRoutes {
		currentRoute := "["
		for i, route := range routes {
			if i == len(routes)-1 {
				currentRoute += strconv.Itoa(route.LoadNumber) + "]"
			} else {
				currentRoute += strconv.Itoa(route.LoadNumber) + ","
			}
		}
		fmt.Println(currentRoute)
	}

	// cost := float64(500*len(determinedDrivingRoutes)) + getDrivenMinutes(determinedDrivingRoutes)
	// fmt.Println("cost is ", cost)
}

func getDrivenMinutes(ddr [][]models.DrivingRoute) float64 {
	runningMinutes := 0.0
	for _, route := range ddr {
		currentRoute := 0.0
		currentRoute += route[0].DurationFromOriginToPickUp()
		for i, stop := range route {
			currentRoute += stop.DurationForHaul()
			if i != len(route)-1 {
				currentRoute += stop.DropOff.DistanceTo(route[i+1].PickUp)
			}
		}
		currentRoute += route[len(route)-1].DurationFromDropOffToOrigin()
		runningMinutes += currentRoute
	}

	return runningMinutes
}
