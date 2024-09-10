package models

type DrivingRoute struct {
	LoadNumber      int
	PickUp, DropOff Point
}

func (dr DrivingRoute) DurationFromOriginToPickUp() float64 {
	return dr.PickUp.DistanceTo(Point{0, 0})
}

func (dr DrivingRoute) DurationForHaul() float64 {
	return dr.PickUp.DistanceTo(dr.DropOff)
}

func (dr DrivingRoute) DurationFromDropOffToOrigin() float64 {
	return dr.DropOff.DistanceTo(Point{0, 0})
}
