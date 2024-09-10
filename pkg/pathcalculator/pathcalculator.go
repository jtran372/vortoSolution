package pathcalculator

import (
	"vrpSolution/internal/models"
)

func CalculateOptimalPaths(qt *models.QuadTree, drs []models.DrivingRoute) [][]models.DrivingRoute {
	calculatedRoutes := [][]models.DrivingRoute{}
	// represents routes and whether they have been visited
	visitedList := []int{}

	for len(visitedList) < len(drs) {
		driverTimeRemaining := float64(12 * 60)
		currentPath := []models.DrivingRoute{}
		// From Origin
		ln, _ := qt.FindNearestValidPickUp(models.Point{X: 0, Y: 0}, visitedList, driverTimeRemaining)

		currentPath = append(currentPath, ln)
		visitedList = append(visitedList, ln.LoadNumber)
		driverTimeRemaining = driverTimeRemaining - (ln.DurationFromOriginToPickUp() + ln.DurationForHaul())

		// From DropOff
		for ln, ok := qt.FindNearestValidPickUp(ln.DropOff, visitedList, driverTimeRemaining); ok; ln, ok = qt.FindNearestValidPickUp(ln.DropOff, visitedList, driverTimeRemaining) {
			driverTimeRemaining = driverTimeRemaining - (currentPath[len(currentPath)-1].DropOff.DistanceTo(ln.PickUp) + ln.DurationForHaul())
			currentPath = append(currentPath, ln)
			visitedList = append(visitedList, ln.LoadNumber)
		}

		calculatedRoutes = append(calculatedRoutes, currentPath)
	}

	return calculatedRoutes
}
