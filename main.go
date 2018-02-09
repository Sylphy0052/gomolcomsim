package main

import (
	"io/ioutil"

	"./sim"
)

func searchDat() {
	datRoot := "./dat/"
	filename, err := ioutil.ReadDir(datRoot)
	if err != nil {
		panic(err)
	}

	ptime := true

	for _, v := range filename {
		sim.Run(datRoot+v.Name(), ptime)
	}
}

func main() {
	searchDat()
}
