package statefaker

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
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
	Mode      string       `json:"mode"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Instances []InstanceV4 `json:"instances"`
}

type InstanceV4 struct {
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

func NewStateV4(outputs, resources int) (*StateV4, error) {
	// Generate multiple realistic resources
	var resourcesCollection []ResourceV4

	for i := 0; i < resources; i++ {
		mode := "managed"
		// 1 in 5 chance to be a data resource
		if rand.IntN(5) == 0 {
			mode = "data"
		}

		var instance InstanceV4
		err := faker.FakeData(&instance)
		if err != nil {
			return nil, fmt.Errorf("failed to fake data for managed resource instance: %w", err)
		}

		instance.SchemaVersion = 0
		instance.IdentitySchemaVersion = 0

		resourceType := generateResourceType()
		resource := ResourceV4{
			Mode:      mode,
			Type:      resourceType,
			Name:      generateResourceName(),
			Provider:  fmt.Sprintf("provider[\"registry.terraform.io/hashicorp/%s\"]", getProviderFromResourceType(resourceType)),
			Instances: []InstanceV4{instance},
		}

		resourcesCollection = append(resourcesCollection, resource)
	}

	// Generate 1-2 data resources
	numDataResources := 1 + rand.IntN(2)
	dataResourceTypes := []string{
		"aws_ami", "aws_availability_zones", "aws_caller_identity",
		"aws_region", "aws_s3_bucket", "aws_iam_policy_document",
	}

	for i := 0; i < numDataResources; i++ {
		var instance InstanceV4
		err := faker.FakeData(&instance)
		if err != nil {
			return nil, fmt.Errorf("failed to fake data for data resource instance: %w", err)
		}

		instance.SchemaVersion = 0
		instance.IdentitySchemaVersion = 0

		resourceType := dataResourceTypes[rand.IntN(len(dataResourceTypes))]
		resource := ResourceV4{
			Mode:      "data",
			Type:      resourceType,
			Name:      generateResourceName(),
			Provider:  fmt.Sprintf("provider[\"registry.terraform.io/hashicorp/%s\"]", getProviderFromResourceType(resourceType)),
			Instances: []InstanceV4{instance},
		}
		resourcesCollection = append(resourcesCollection, resource)
	}

	// Generate realistic outputs
	outputsMap := make(map[string]json.RawMessage)

	for i := 0; i < outputs; i++ {
		b, err := randomOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random output: %w", err)
		}
		outputsMap[faker.Word()] = b
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
