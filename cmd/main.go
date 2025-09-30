package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/brandonc/go-statefaker.git/pkg/statefaker"
)

var numOutputs int
var numResources int

func init() {
	flag.IntVar(&numOutputs, "outputs", 3, "the number of outputs to generate")
	flag.IntVar(&numResources, "resources", 3, "the number of resources to generate")
}

func main() {
	flag.Parse()

	sf, err := statefaker.NewFakeStateV4(numOutputs, numResources)
	if err != nil {
		panic(err)
	}

	// Marshal as json and print to stdout
	b, err := json.Marshal(sf)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
