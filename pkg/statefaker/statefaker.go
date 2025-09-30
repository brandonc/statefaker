package statefaker

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

type StateV4 struct {
	Version          int                        `json:"version"`
	TerraformVersion string                     `json:"terraform_version"`
	Serial           int                        `json:"serial"`
	Lineage          string                     `json:"lineage"`
	Outputs          map[string]json.RawMessage `json:"outputs"`
	Resources        []ResourceV4               `json:"resources"`
}

type ResourceV4 struct {
	Module    string       `json:"module,omitempty"`
	Mode      string       `json:"mode"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Instances []InstanceV4 `json:"instances"`
}

type InstanceV4 struct {
	IndexKey              string          `json:"index_key,omitempty"`
	SchemaVersion         int             `json:"schema_version"`
	Attributes            json.RawMessage `json:"attributes" faker:"tfattributes"`
	SensitiveAttributes   []string        `json:"sensitive_attributes"`
	IdentitySchemaVersion int             `json:"identity_schema_version"`
}

type OutputV4 struct {
	Value json.RawMessage `json:"value"`
	Type  json.RawMessage `json:"type"`
}

type ExampleAttributes struct {
	Name string `json:"name"`
	ARN  string `json:"arn"`
}

func NewFakeStateV4(outputs, resources int) (*StateV4, error) {
	// Generate multiple realistic resources
	var resourcesCollection []ResourceV4

	for i := 0; i < resources; i++ {
		mode := "managed"
		// 1 in 5 chance to be a data resource
		if rand.IntN(5) == 0 {
			mode = "data"
		}

		resourceType := generateResourceType()

		// 30% chance to have a module address
		var moduleAddress string
		if rand.IntN(10) < 3 {
			moduleAddress = generateModuleAddress()
		}

		// Generate instances - 20% chance to have multiple instances (dozens)
		var instances []InstanceV4
		numInstances := 1
		if rand.IntN(5) < 1 {
			// Generate dozens of instances (12-48)
			numInstances = rand.IntN(37) + 12
		}

		for j := 0; j < numInstances; j++ {
			var instance InstanceV4
			err := faker.FakeData(&instance)
			if err != nil {
				return nil, fmt.Errorf("failed to fake data for managed resource instance: %w", err)
			}

			instance.SchemaVersion = 0
			instance.IdentitySchemaVersion = 0
			instance.SensitiveAttributes = []string{}

			// Set unique IndexKey for multiple instances
			if numInstances > 1 {
				instance.IndexKey = fmt.Sprintf("%s-%s-%d", faker.Word(options.WithGenerateUniqueValues(true)), faker.Word(), j)
			} else {
				instance.IndexKey = ""
			}

			instances = append(instances, instance)
		}

		faker.ResetUnique()

		resource := ResourceV4{
			Mode:      mode,
			Type:      resourceType,
			Name:      generateResourceName(),
			Module:    moduleAddress,
			Provider:  generateProviderString(resourceType, moduleAddress),
			Instances: instances,
		}

		resourcesCollection = append(resourcesCollection, resource)
	}

	faker.ResetUnique()

	// Generate realistic outputs
	outputsMap := make(map[string]json.RawMessage)

	for range outputs {
		b, err := randomOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random output: %w", err)
		}
		outputsMap[fmt.Sprintf("%s_%s_%d", faker.Word(), faker.Word(), faker.UnixTime())] = b
	}

	state := &StateV4{
		Version:          4,
		TerraformVersion: "1.5.6",
		Serial:           1,
		Lineage:          faker.UUIDHyphenated(),
		Outputs:          outputsMap,
		Resources:        resourcesCollection,
	}

	return state, nil
}
