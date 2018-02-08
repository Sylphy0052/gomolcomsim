package sim

import (
	"math"
	"math/rand"
)

type MolController struct {
	molecules    []Molecule
	stepLength   FloatPosition
	mediumLength Position
}

func (mc *MolController) addMol(m Molecule) {
	mc.molecules = append(mc.molecules, m)
}

func (mc *MolController) removeMol(m Molecule) {
	mols := []Molecule{}
	flg := true
	for _, v := range mc.molecules {
		if v.getName() == m.getName() && v.position == v.position && flg {
			flg = false
			continue
		} else {
			mols = append(mols, v)
		}
	}
	mc.molecules = mols
}

func (mc *MolController) doNextStep(g *Grid) {
	mols := []Molecule{}
	for _, m := range mc.molecules {
		g.removeObject(m, m.position)
		m.move(mc.stepLength, mc.mediumLength)
		mols = append(mols, m)
		g.addObject(m, m.position)
	}
	mc.molecules = mols
}

func checkReceive(m Molecule, objects []Object, g *Grid, mc *MolController, sim *Sim) {
	switch m.moleculeType {
	case INFO:
		for _, v := range objects {
			if v.getName() == "receiver" {
				g.removeObject(m, m.position)
				mc.removeMol(m)
				v.receiveMol(g, mc)
			}
		}
	case ACK:
		for _, v := range objects {
			if v.getName() == "transmitter" {
				g.removeObject(m, m.position)
				mc.removeMol(m)
				v.receiveMol(g, mc)
				sim.isFinish = true
			}
		}
	}
}

func checkCollision(m Molecule) {

}

func (mc *MolController) checkCollision(sim *Sim) {
	for _, m := range mc.molecules {
		currentGrid := sim.medium.grid
		currentObjects := currentGrid[m.position]
		checkReceive(m, currentObjects, &sim.medium.grid, mc, sim)
		if sim.config.useCollisions {
			checkCollision(m)
		}
	}
}

// func (mc *MolController) checkCollision(g *Grid, config Config) {
// 	for _, m := range mc.molecules {
// 		currentGrid := *g
// 		currentObjects := currentGrid[m.position]
// 		checkReceive(m, currentObjects, g, mc)
// 		if config.useCollisions {
// 			checkCollision(m)
// 		}
// 	}
// }

type Molecule struct {
	position       Position
	prevPosition   Position
	moleculeType   MoleculeType
	movementType   MovementType
	startTime      int
	endTime        int
	msgId          int
	adaptiveChange int
	size           float64
	name           string
	volume         float64
}

func (m Molecule) getName() string {
	return m.name
}

func (m Molecule) receiveMol(g *Grid, mc *MolController) {

}

func getNextPosition(m Molecule, stepLength FloatPosition) Position {
	currentPosition := m.position
	nextX := currentPosition.x + int(round(rand.Float64()*stepLength.x*2-stepLength.x))
	nextY := currentPosition.y + int(round(rand.Float64()*stepLength.y*2-stepLength.y))
	nextZ := currentPosition.z + int(round(rand.Float64()*stepLength.z*2-stepLength.z))
	return Position{
		x: nextX,
		y: nextY,
		z: nextZ,
	}
}

func (m *Molecule) move(stepLength FloatPosition, mediumLength Position) {
	m.prevPosition = m.position
	nextPosition := getNextPosition(*m, stepLength)
	// 範囲の外にいないか確認
	if nextPosition.x < -mediumLength.x/2 {
		nextPosition.x = -mediumLength.x / 2
	} else if nextPosition.x > mediumLength.x/2 {
		nextPosition.x = mediumLength.x / 2
	}
	if nextPosition.y < -mediumLength.y/2 {
		nextPosition.y = -mediumLength.y / 2
	} else if nextPosition.y > mediumLength.y/2 {
		nextPosition.y = mediumLength.y / 2
	}
	if nextPosition.z < -mediumLength.z/2 {
		nextPosition.z = -mediumLength.z / 2
	} else if nextPosition.z > mediumLength.z/2 {
		nextPosition.z = mediumLength.z / 2
	}
	m.position = nextPosition
}

func createMolecule(n NanoMachine, param MoleculeParam, msgId int, name string) Molecule {
	position := n.molReleasePosition
	moleculeType := param.moleculeType
	movementType := param.moleculeMovementType
	startTime := 0
	adaptiveChange := param.adaptiveChange
	size := param.size
	volume := math.Pow(param.size, 3.0)
	return Molecule{
		position:       position,
		moleculeType:   moleculeType,
		movementType:   movementType,
		startTime:      startTime,
		msgId:          msgId,
		adaptiveChange: adaptiveChange,
		size:           size,
		name:           name,
		volume:         volume,
	}
}
