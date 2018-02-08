package sim

import (
	"fmt"
	"math"
)

type NanoMachine struct {
	center             Position
	radius             int
	molReleasePosition Position
	molParam           MoleculeParam
	name               string
}

func (n NanoMachine) getName() string {
	return n.name
}

// func (n NanoMachine) doNextStep() {
// 	switch n.getName() {
// 	case "transmitter":
//
// 	case "receiver":
//
// 	}
// }

func printGrid(g *Grid) {
	for k, v := range *g {
		if len(v) == 0 {
			continue
		}
		fmt.Println(k, ":", v)
	}
}

func (n NanoMachine) receiveMol(g *Grid, mc *MolController) {
	switch n.getName() {
	case "receiver":
		fmt.Println("receive INFO mol")
		n.createAckMolecule(mc, g, 1)
	case "transmitter":
		fmt.Println("receive ACK mol")
	}
}

func (n NanoMachine) createInfoMolecule(mc *MolController, g *Grid, msgId int) {
	num := n.molParam.num
	for i := 0; i < num; i++ {
		m := createMolecule(n, n.molParam, msgId, "INFO")
		g.addObject(m, n.molReleasePosition)
		mc.addMol(m)
	}
}

func (n NanoMachine) createAckMolecule(molecules *MolController, g *Grid, msgId int) {
	num := n.molParam.num
	for i := 0; i < num; i++ {
		m := createMolecule(n, n.molParam, msgId, "ACK")
		g.addObject(m, n.molReleasePosition)
		molecules.addMol(m)
	}
}

func createNanoMachine(g Grid, tParam NanoMachineParam, molParam MoleculeParam, name string) NanoMachine {
	tx := NanoMachine{
		center:             tParam.center,
		radius:             tParam.radius,
		molReleasePosition: tParam.molReleasePosition,
		molParam:           molParam,
		name:               name,
	}
	startX := tParam.center.x - tParam.radius + 1
	endX := tParam.center.x + tParam.radius - 1
	startY := tParam.center.y - tParam.radius + 1
	endY := tParam.center.y + tParam.radius - 1
	startZ := tParam.center.z - tParam.radius + 1
	endZ := tParam.center.z + tParam.radius - 1

	for i := startX; i <= endX; i++ {
		for j := startY; j <= endY; j++ {
			for k := startZ; k <= endZ; k++ {
				g.addObject(tx, Position{x: i, y: j, z: k})
			}
		}
	}

	return tx
}

type Microtubule struct {
	startPosition Position
	endPosition   Position
	name          string
}

func (m Microtubule) getName() string {
	return m.name
}

func (m Microtubule) receiveMol(g *Grid, mc *MolController) {

}

func (m Microtubule) setGrid(sim *Sim) {
	vel := float64(sim.config.velRail)
	startX, startY, startZ := m.startPosition.getPosition()
	endX, endY, endZ := m.endPosition.getPosition()

	diffX := float64(endX - startX)
	diffY := float64(endY - startY)
	diffZ := float64(endZ - startZ)

	length := math.Sqrt(diffX*diffX + diffY*diffY + diffZ*diffZ)
	unitX := diffX * vel / length
	unitY := diffY * vel / length
	unitZ := diffZ * vel / length
	unitFloatPosition := FloatPosition{unitX, unitY, unitZ}

	grid := sim.medium.grid
	currentPosition := FloatPosition{float64(startX), float64(startY), float64(startZ)}
	for {
		grid.setFloatPosition(m, currentPosition)
		currentPosition.add(unitFloatPosition)
		if currentPosition.toPosition() == m.endPosition {
			break
		}
	}
}
