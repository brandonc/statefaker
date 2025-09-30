package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/brandonc/go-statefaker.git/pkg/statefaker"
)

var numOutputs int
var numResources int
var percentMultiInstance int
var multiMaxInstances int
var multiMinInstances int
var percentModule int

func init() {
	defaults := statefaker.DefaultOptions()

	flag.IntVar(&numOutputs, "outputs", defaults.NumOutputs, "the number of outputs to generate")
	flag.IntVar(&numResources, "resources", defaults.NumResources, "the number of resources to generate")
	flag.IntVar(&percentMultiInstance, "pctmulti", defaults.MultiInstanceChance, "the percentage chance a resource is multi-instance")
	flag.IntVar(&multiMaxInstances, "multimax", defaults.MultiInstanceMax, "the maximum number of instances for multi-instance resources")
	flag.IntVar(&multiMinInstances, "multimin", defaults.MultiInstanceMin, "the minimum number of instances for multi-instance resources")
	flag.IntVar(&percentModule, "pctmodule", defaults.ModuleChance, "the percentage chance a resource appears within a module")
}

func main() {
	flag.Parse()

	sf, err := statefaker.NewFakeStateV4(
		statefaker.WithOutputs(numOutputs),
		statefaker.WithResources(numResources),
		statefaker.WithMultiInstanceChance(percentMultiInstance),
		statefaker.WithMultiInstanceMax(multiMaxInstances),
		statefaker.WithMultiInstanceMin(multiMinInstances),
		statefaker.WithModuleChance(percentModule),
	)
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
