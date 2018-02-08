package sim

import (
	"fmt"
	"math/rand"
	"time"
)

type Sim struct {
	config        Config
	simStep       int
	medium        Medium
	microtubules  []Microtubule
	transmitter   NanoMachine
	receiver      NanoMachine
	molController MolController
	isFinish      bool
}

func (sim *Sim) createMedium() {
	m := createMedium()
	m.x = sim.config.mediumDimensionX
	m.y = sim.config.mediumDimensionY
	m.z = sim.config.mediumDimensionZ
	sim.medium = m
}

func (sim *Sim) createMicrotubules() {
	for _, v := range sim.config.microtubuleParams {
		m := Microtubule{startPosition: v.startPosition, endPosition: v.endPosition, name: "Microtubules"}
		m.setGrid(sim)
		sim.microtubules = append(sim.microtubules, m)
	}
}

func (sim *Sim) createNanoMachines() {
	var infoParam MoleculeParam
	var ackParam MoleculeParam
	for _, v := range sim.config.moleculeParams {
		if v.moleculeType == INFO {
			infoParam = v
		} else if v.moleculeType == ACK {
			ackParam = v
		}
	}

	// transmitter
	tx := createNanoMachine(sim.medium.grid, sim.config.transmitter, infoParam, "transmitter")
	tx.createInfoMolecule(&sim.molController, &sim.medium.grid)
	sim.transmitter = tx

	// receiver
	rx := createNanoMachine(sim.medium.grid, sim.config.receiver, ackParam, "receiver")
	sim.receiver = rx

	// intermediateNode
}

func (sim *Sim) createMolController() {
	rand.Seed(time.Now().UnixNano())
	sim.molController = MolController{
		molecules: make([]Molecule, 0),
		stepLength: FloatPosition{
			x: sim.config.stepLengthX,
			y: sim.config.stepLengthY,
			z: sim.config.stepLengthZ,
		},
		mediumLength: Position{
			x: sim.config.mediumDimensionX,
			y: sim.config.mediumDimensionY,
			z: sim.config.mediumDimensionZ,
		},
	}
}

func (sim *Sim) doNextStep() {
	// NanoMachine
	// sim.transmitter.doNextStep()
	// sim.receiver.doNextStep()

	// Molecule
	sim.molController.doNextStep(&sim.medium.grid)

	// 衝突
	// sim.molController.checkCollision(&sim.medium.grid, sim.config)
	sim.molController.checkCollision(sim)
}

func initSim(filename string, ptime bool) Sim {
	sim := Sim{config: createConfig(filename, ptime), simStep: 0}
	sim.createMolController()
	sim.createMedium()
	sim.createMicrotubules()
	sim.createNanoMachines()
	return sim
}

func Run(filename string, ptime bool) {
	sim := initSim(filename, ptime)
	for ; ; sim.simStep++ {
		sim.doNextStep()
		if sim.isFinish {
			fmt.Println(sim.simStep)
			break
		}
	}
}
