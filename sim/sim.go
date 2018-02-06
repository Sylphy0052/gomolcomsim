package sim

type Sim struct {
	config  *Config
	simStep int
	medium  Medium
}

type Molecule struct {
	position     Position
	moleculeType MoleculeType
	movementType MovementType
	startTime    int
	endTime      int
	msgId        int
	medium       Medium
}

func (sim *Sim) createMedium() {
	m := Medium{}
	m.x = sim.config.mediumDimensionX
	m.y = sim.config.mediumDimensionY
	m.z = sim.config.mediumDimensionZ
	m.garbageSpot = Position{m.x * 2, m.y * 2, m.z * 2}
	sim.medium = m
}

func (sim *Sim) createMicrotubules() {
	for _, v := range sim.config.microtubuleParams {
		start, end := v.startPosition, v.endPosition

	}
}

func (sim *Sim) createNanoMachines() {

}

func Run(filename string, ptime bool) {
	sim := Sim{config: createConfig(filename, ptime), simStep: 0}
	sim.createMedium()
	sim.createMicrotubules()
	sim.createNanoMachines()
}
