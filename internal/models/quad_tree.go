package models

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

const (
	MAX_CAPACITY = 4
)

type Point struct {
	X, Y float64
}

// For this solution, also equal to time to drive
func (p Point) DistanceTo(other Point) float64 {
	return math.Sqrt(math.Pow(p.X-other.X, 2) + math.Pow(p.Y-other.Y, 2))
}

type Bounds struct {
	MinX, MinY, MaxX, MaxY float64
}

func (b Bounds) contains(p Point) bool {
	return p.X >= b.MinX && p.X < b.MaxX &&
		p.Y >= b.MinY && p.Y < b.MaxY
}

type QuadTree struct {
	boundary Bounds
	routes   []DrivingRoute
	divided  bool

	northWest, northEast, southWest, southEast *QuadTree
}

func NewQuadTree(boundary Bounds) *QuadTree {
	return &QuadTree{
		boundary: boundary,
		routes:   make([]DrivingRoute, 0),
		divided:  false,
	}
}

func (qt *QuadTree) Subdivide() {
	halfX := (qt.boundary.MaxX - qt.boundary.MinX) / 2
	halfY := (qt.boundary.MaxY - qt.boundary.MinY) / 2

	nw := Bounds{qt.boundary.MinX, halfY, halfX, qt.boundary.MaxY}
	ne := Bounds{halfX, halfY, qt.boundary.MaxX, qt.boundary.MaxY}
	sw := Bounds{qt.boundary.MinX, qt.boundary.MinY, halfX, halfY}
	se := Bounds{halfX, qt.boundary.MinY, qt.boundary.MaxX, halfY}

	qt.northWest = NewQuadTree(nw)
	qt.northEast = NewQuadTree(ne)
	qt.southWest = NewQuadTree(sw)
	qt.southEast = NewQuadTree(se)

	qt.divided = true
}

func (qt *QuadTree) Insert(dr DrivingRoute) bool {
	if !qt.boundary.contains(dr.PickUp) {
		return false
	}

	if len(qt.routes) < MAX_CAPACITY {
		qt.routes = append(qt.routes, dr)
		return true
	}

	if !qt.divided {
		qt.Subdivide()
	}

	if qt.northWest.Insert(dr) {
		return true
	} else if qt.northEast.Insert(dr) {
		return true
	} else if qt.southWest.Insert(dr) {
		return true
	} else if qt.southEast.Insert(dr) {
		return true
	}

	return false
}

func (qt *QuadTree) FindNearestValidPickUp(
	target Point,
	alreadyVisited []int,
	driverTimeRemaining float64,
) (
	DrivingRoute,
	bool,
) {
	dr, minDistance := qt.FindNearestValidPickUpHelper(target, &DrivingRoute{}, math.MaxFloat64, alreadyVisited, driverTimeRemaining)

	return dr, minDistance != math.MaxFloat64 // determines if valid path found
}

func (qt *QuadTree) FindNearestValidPickUpHelper(
	target Point,
	nearestSoFar *DrivingRoute,
	minDistSoFar float64,
	alreadyVisitedList []int,
	driverTimeRemaining float64,
) (
	DrivingRoute,
	float64,
) {
	for _, dr := range qt.routes {
		visitedAlready := false
		for _, val := range alreadyVisitedList {
			if val == dr.LoadNumber {
				visitedAlready = true
				break
			}
		}
		if visitedAlready {
			continue
		}
		if target.DistanceTo(dr.PickUp)+dr.DurationForHaul()+dr.DurationFromDropOffToOrigin() <= driverTimeRemaining {
			d := target.DistanceTo(dr.PickUp)
			if d < minDistSoFar {
				minDistSoFar = d
				*nearestSoFar = dr
			}
		}
	}

	if qt.divided {
		temp, minDist := qt.northWest.FindNearestValidPickUpHelper(target, nearestSoFar, minDistSoFar, alreadyVisitedList, driverTimeRemaining)
		if minDist < minDistSoFar {
			*nearestSoFar = temp
			minDistSoFar = minDist
		}
		temp, minDist = qt.northEast.FindNearestValidPickUpHelper(target, nearestSoFar, minDistSoFar, alreadyVisitedList, driverTimeRemaining)
		if minDist < minDistSoFar {
			*nearestSoFar = temp
			minDistSoFar = minDist
		}
		temp, minDist = qt.southWest.FindNearestValidPickUpHelper(target, nearestSoFar, minDistSoFar, alreadyVisitedList, driverTimeRemaining)
		if minDist < minDistSoFar {
			*nearestSoFar = temp
			minDistSoFar = minDist
		}
		temp, minDist = qt.southEast.FindNearestValidPickUpHelper(target, nearestSoFar, minDistSoFar, alreadyVisitedList, driverTimeRemaining)
		if minDist < minDistSoFar {
			*nearestSoFar = temp
			minDistSoFar = minDist
		}
	}

	return *nearestSoFar, minDistSoFar
}

// Expected format of (x,y)
func StringToCartesianCoordinate(coordStr string) (Point, error) {
	re := regexp.MustCompile(`^\(\s*([-+]?\d*\.?\d+)\s*,\s*([-+]?\d*\.?\d+)\s*\)$`)

	if !re.MatchString(coordStr) {
		return Point{}, errors.New("invalid coordinate format")
	}

	matches := re.FindStringSubmatch(coordStr)
	if len(matches) != 3 {
		return Point{}, errors.New("invalid coordinate format")
	}

	x, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return Point{}, fmt.Errorf("invalid value for x: %v", err)
	}

	y, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return Point{}, fmt.Errorf("invalid value for y: %v", err)
	}

	return Point{X: x, Y: y}, nil
}
