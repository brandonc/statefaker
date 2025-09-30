package statefaker

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	fakeroptions "github.com/go-faker/faker/v4/pkg/options"
)

type StateV4 struct {
	Version          int                        `json:"version"`
	TerraformVersion string                     `json:"terraform_version"`
	Serial           int                        `json:"serial"`
	Lineage          string                     `json:"lineage"`
	Outputs          map[string]json.RawMessage `json:"outputs"`
	Resources        []ResourceV4               `json:"resources"`
	Source           string                     `json:"source,omitempty"`
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
	SensitiveAttributes   []string        `json:"sensitive_attributes" faker:"tfemptystringslice"`
	IdentitySchemaVersion int             `json:"identity_schema_version" faker:"oneof: 0,1"`
	Identity              json.RawMessage `json:"identity,omitempty" faker:"tfidentity"`
	Private               string          `json:"private,omitempty" faker:"tfprivate"`
	Dependencies          []string        `json:"dependencies,omitempty" faker:"tfdependencies"`
}

type OutputV4 struct {
	Value json.RawMessage `json:"value"`
	Type  json.RawMessage `json:"type"`
}

type ExampleAttributes struct {
	Name string `json:"name"`
	ARN  string `json:"arn"`
}

func NewFakeStateV4(opts ...Option) (*StateV4, error) {
	// Apply options with defaults
	options := ApplyOptions(opts...)

	// Generate multiple realistic resources
	var resourcesCollection []ResourceV4

	for i := 0; i < options.NumResources; i++ {
		mode := "managed"
		// 1 in 5 chance to be a data resource
		if rand.IntN(5) == 0 {
			mode = "data"
		}

		resourceType := generateResourceType()

		// Configurable chance to have a module address
		var moduleAddress string
		if rand.IntN(100) < options.ModuleChance {
			moduleAddress = generateModuleAddress()
		}

		// Generate instances - configurable chance to have multiple instances
		var instances []InstanceV4
		numInstances := 1
		if rand.IntN(100) < options.MultiInstanceChance {
			// Generate configurable range of instances
			instanceRange := options.MultiInstanceMax - options.MultiInstanceMin + 1
			numInstances = rand.IntN(instanceRange) + options.MultiInstanceMin
		}

		for j := 0; j < numInstances; j++ {
			var instance InstanceV4
			err := faker.FakeData(&instance)
			if err != nil {
				return nil, fmt.Errorf("failed to fake data for managed resource instance: %w", err)
			}

			// Set unique IndexKey for multiple instances
			if numInstances > 1 {
				instance.IndexKey = fmt.Sprintf("%s-%s-%d", faker.Word(fakeroptions.WithGenerateUniqueValues(true)), faker.Word(), j)
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

	for range options.NumOutputs {
		b, err := generateOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random output: %w", err)
		}
		outputsMap[fmt.Sprintf("%s_%s_%d", faker.Word(), faker.Word(), faker.UnixTime())] = b
	}

	state := &StateV4{
		Version:          4,
		TerraformVersion: "1.13.2",
		Serial:           1,
		Lineage:          faker.UUIDHyphenated(),
		Outputs:          outputsMap,
		Resources:        resourcesCollection,
		Source:           "statefaker",
	}

	return state, nil
}
