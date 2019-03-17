package sim

import (
	"fmt"
	"log"
	"os"
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
	finishStep    int
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

func writeResult(filename string, str string) {
	writeFileName := "./result/" + filename
	file, err := os.OpenFile(writeFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, str)
}

func Run(filename string, ptime bool) {
	sim := initSim(filename, ptime)
	start := time.Now()
	for ; ; sim.simStep++ {
		sim.doNextStep()

		if sim.simStep%1000 == 999 {
			fmt.Println(sim.simStep+1, time.Now().Sub(start).Nanoseconds()/1000000, "ms")
			start = time.Now()
		}

		if sim.isFinish && !ptime {
			writeResult(sim.config.outputFile, fmt.Sprint(sim.simStep))
			break
		} else if ptime && sim.isFinish && len(sim.molController.molecules) == 0 {
			writeResult(sim.config.outputFile, fmt.Sprint(sim.simStep)+","+fmt.Sprint(sim.finishStep))
			break
		}
	}
}
