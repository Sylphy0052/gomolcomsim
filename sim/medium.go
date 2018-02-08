package sim

import (
	"math"
)

type Grid map[Position][]Object

func createGrid() Grid { return make(map[Position][]Object) }

func (g Grid) addObject(o Object, p Position) {
	if g[p] == nil {
		g[p] = make([]Object, 0)
	}
	g[p] = append(g[p], o)
}

func (g Grid) removeObject(o Object, p Position) {
	newG := []Object{}
	flg := true
	for _, v := range g[p] {
		if v.getName() == o.getName() && flg {
			flg = false
			continue
		} else {
			newG = append(newG, v)
		}
	}
	g[p] = newG
}

func (g *Grid) setFloatPosition(o Object, fp FloatPosition) {
	floorX := int(math.Floor(fp.x))
	floorY := int(math.Floor(fp.y))
	floorZ := int(math.Floor(fp.z))
	ceilX := int(math.Ceil(fp.x))
	ceilY := int(math.Ceil(fp.x))
	ceilZ := int(math.Ceil(fp.x))
	x := make([]int, 0)
	y := make([]int, 0)
	z := make([]int, 0)
	x = append(x, floorX)
	x = append(x, ceilX)
	y = append(y, floorY)
	y = append(y, ceilY)
	z = append(z, floorZ)
	z = append(z, ceilZ)

	for _, vx := range x {
		for _, vy := range y {
			for _, vz := range z {
				g.addObject(o, Position{x: vx, y: vy, z: vz})
			}
		}
	}
}

type Object interface {
	getName() string
	receiveMol(m Molecule, mc *MolController, sim *Sim)
	getVolume() float64
}

type Medium struct {
	x, y, z int
	grid    Grid
}

func createMedium() Medium { return Medium{grid: createGrid()} }
