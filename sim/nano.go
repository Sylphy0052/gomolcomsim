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
	currentId          int
}

func (n NanoMachine) getName() string { return n.name }

func (n NanoMachine) getVolume() float64 { return float64(0.0) }

func printGrid(g *Grid) {
	for k, v := range *g {
		if len(v) == 0 {
			continue
		}
		fmt.Println(k, ":", v)
	}
}

func (n *NanoMachine) receiveMol(m Molecule, mc *MolController, sim *Sim) {
	if m.msgId == n.currentId {
		switch n.getName() {
		case "receiver":
			fmt.Println("receive INFO mol")
			n.createAckMolecule(mc, &sim.medium.grid)
			n.currentId += 1

		case "transmitter":
			fmt.Println("receive ACK mol")
			if n.currentId == sim.config.numMessages {
				sim.isFinish = true
			} else {
				n.currentId += 1
				n.createInfoMolecule(mc, &sim.medium.grid)
			}
		}
	}
}

func (n NanoMachine) createInfoMolecule(mc *MolController, g *Grid) {
	num := n.molParam.num
	for i := 0; i < num; i++ {
		m := createMolecule(n, n.molParam, n.currentId, "INFO")
		g.addObject(m, n.molReleasePosition)
		mc.addMol(m)
	}
}

func (n NanoMachine) createAckMolecule(molecules *MolController, g *Grid) {
	num := n.molParam.num
	for i := 0; i < num; i++ {
		m := createMolecule(n, n.molParam, n.currentId, "ACK")
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
		currentId:          1,
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
				g.addObject(&tx, Position{x: i, y: j, z: k})
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

func (m Microtubule) getName() string { return m.name }

func (m Microtubule) receiveMol(mol Molecule, mc *MolController, sim *Sim) {}

func (m Microtubule) getVolume() float64 { return float64(0.0) }

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
