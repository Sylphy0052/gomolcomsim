package sim

type Grid map[Position][]string

type Medium struct {
	x, y, z     int
	grid        Grid
	garbageSpot Position
}
