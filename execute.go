package main

import (
	"os"

	"./sim"
)

func main() {
	file_name := os.Args[1]
	ptimeString := os.Args[2]
	ptime := false
	if ptimeString == "true" {
		ptime = true
	}
	sim.Run(file_name, ptime)
}
