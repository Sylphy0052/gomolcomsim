package sim

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type MoleculeType int

const (
	INFO MoleculeType = iota
	ACK
	NOISE
)

func getMoleculeType(s string) MoleculeType {
	var val MoleculeType
	switch s {
	case "INFO":
		val = INFO
	case "ACK":
		val = ACK
	case "NOISE":
		val = NOISE
	}
	return val
}

type MovementType int

const (
	ACTIVE MovementType = iota
	PASSIVE
	NONE
)

func getMovementType(s string) MovementType {
	var val MovementType
	switch s {
	case "ACTIVE":
		val = ACTIVE
	case "PASSIVE":
		val = PASSIVE
	case "NONE":
		val = NONE
	}
	return val
}

type FloatPosition struct {
	x, y, z float64
}

func (fp FloatPosition) getPosition() (float64, float64, float64) {
	return fp.x, fp.y, fp.z
}

func (fp *FloatPosition) add(p FloatPosition) {
	fp.x, fp.y, fp.z = FloatPosition{x: fp.x + p.x, y: fp.y + p.y, z: fp.z + p.z}.getPosition()
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}

func (fp FloatPosition) toPosition() Position {
	x := int(round(fp.x))
	y := int(round(fp.y))
	z := int(round(fp.z))
	return Position{x: x, y: y, z: z}
}

type Position struct {
	x, y, z int
}

func (p Position) getPosition() (int, int, int) {
	return p.x, p.y, p.z
}

func createPosition(s string) Position {
	s = strings.Trim(s, "(")
	s = strings.Trim(s, ")")
	sp := strings.Split(s, ",")
	x, _ := strconv.Atoi(sp[0])
	y, _ := strconv.Atoi(sp[1])
	z, _ := strconv.Atoi(sp[2])
	p := Position{x, y, z}
	return p
}

type NanoMachineParam struct {
	center             Position
	radius             int
	molReleasePosition Position
}

func createNanoMachineParam(args []string) NanoMachineParam {
	n := NanoMachineParam{}
	n.center = createPosition(args[0] + args[1] + args[2])
	n.radius, _ = strconv.Atoi(args[3])
	n.molReleasePosition = createPosition(args[4] + args[5] + args[6])
	return n
}

type IntermediateNodeParam struct {
	center                 Position
	radius                 int
	infoMolReleasePosition Position
	ackMolReleasePosition  Position
}

func createIntermediateNodeParam(args []string) IntermediateNodeParam {
	inp := IntermediateNodeParam{}
	inp.center = createPosition(args[0] + args[1] + args[2])
	inp.radius, _ = strconv.Atoi(args[3])
	inp.infoMolReleasePosition = createPosition(args[4] + args[5] + args[6])
	inp.ackMolReleasePosition = createPosition(args[7] + args[8] + args[9])
	return inp
}

type MoleculeParam struct {
	num                  int
	moleculeType         MoleculeType
	moleculeMovementType MovementType
	adaptiveChange       int
	size                 float64
}

func createMoleculeParam(args []string) MoleculeParam {
	mp := MoleculeParam{}
	mp.num, _ = strconv.Atoi(args[0])
	mp.moleculeType = getMoleculeType(args[1])
	if len(args) != 3 {
		mp.moleculeMovementType = getMovementType(args[2])
		mp.adaptiveChange, _ = strconv.Atoi(args[3])
		mp.size, _ = strconv.ParseFloat(args[4], 32)
	} else {
		mp.size, _ = strconv.ParseFloat(args[2], 32)
	}
	return mp
}

type MicrotubuleParam struct {
	startPosition Position
	endPosition   Position
	velRail       int
}

func createMicrotubuleParam(args []string) MicrotubuleParam {
	mp := MicrotubuleParam{}
	mp.startPosition = createPosition(args[0] + args[1] + args[2])
	mp.endPosition = createPosition(args[3] + args[4] + args[5])
	return mp
}

func getBool(s string) bool {
	i, _ := strconv.Atoi(s)
	if i == 0 {
		return false
	}
	return true
}

type Config struct {
	configFilename      string
	mediumDimensionX    int
	mediumDimensionY    int
	mediumDimensionZ    int
	maxSimulationStep   int
	transmitter         NanoMachineParam
	receiver            NanoMachineParam
	intermediateNode    []IntermediateNodeParam
	numMessages         int
	numRetransmissions  int
	retransmitWaitTime  int
	stepLengthX         float64
	stepLengthY         float64
	stepLengthZ         float64
	moleculeParams      []MoleculeParam
	velRail             int
	probDRail           float64
	useCollisions       bool
	useAcknowledgements bool
	decomposing         bool
	microtubuleParams   []MicrotubuleParam
	outputFile          string
	batchRun            bool
	pwait               bool
}

func createConfig(configFilename string, pwait bool) Config {
	config := Config{configFilename: "./dat/" + configFilename, pwait: pwait}
	config.readDat()
	return config
}

func (c *Config) readDat() {
	fp, err := os.Open(c.configFilename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '*' {
			continue
		}
		params := strings.Split(line, " ")
		switch params[0] {
		case "mediumDimensionX":
			c.mediumDimensionX, _ = strconv.Atoi(params[1])
		case "mediumDimensionY":
			c.mediumDimensionY, _ = strconv.Atoi(params[1])
		case "mediumDimensionZ":
			c.mediumDimensionZ, _ = strconv.Atoi(params[1])
		case "maxSimulationStep":
			c.maxSimulationStep, _ = strconv.Atoi(params[1])
		case "transmitter":
			c.transmitter = createNanoMachineParam(params[1:])
		case "receiver":
			c.receiver = createNanoMachineParam(params[1:])
		case "intermediateNode":
			c.intermediateNode = append(c.intermediateNode, createIntermediateNodeParam(params[1:]))
		case "numMessages":
			c.numMessages, _ = strconv.Atoi(params[1])
		case "numRetransmissions":
			c.numRetransmissions, _ = strconv.Atoi(params[1])
		case "retransmitWaitTime":
			c.retransmitWaitTime, _ = strconv.Atoi(params[1])
		case "stepLengthX":
			c.stepLengthX, _ = strconv.ParseFloat(params[1], 32)
		case "stepLengthY":
			c.stepLengthY, _ = strconv.ParseFloat(params[1], 32)
		case "stepLengthZ":
			c.stepLengthZ, _ = strconv.ParseFloat(params[1], 32)
		case "moleculeParams":
			c.moleculeParams = append(c.moleculeParams, createMoleculeParam(params[1:]))
		case "velRail":
			c.velRail, _ = strconv.Atoi(params[1])
		case "probDRail":
			c.probDRail, _ = strconv.ParseFloat(params[1], 32)
		case "useCollisions":
			c.useCollisions = getBool(params[1])
		case "useAcknowledgements":
			c.useAcknowledgements = getBool(params[1])
		case "decomposing":
			c.decomposing = getBool(params[1])
		case "microtubuleParams":
			c.microtubuleParams = append(c.microtubuleParams, createMicrotubuleParam(params[1:]))
		case "outputFile":
			c.outputFile = params[1]
		}
	}
}
