package statefaker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"reflect"
)

func tfattributesProvider(v reflect.Value) (any, error) {
	// Generate realistic resource attributes based on common AWS resource types
	attributeGenerators := []func() map[string]any{
		generateS3BucketAttributes,
		generateIAMUserAttributes,
		generateEC2InstanceAttributes,
		generateLambdaFunctionAttributes,
		generateRDSInstanceAttributes,
	}

	generator := attributeGenerators[rand.IntN(len(attributeGenerators))]
	resourceAttributes := generator()

	b, err := json.Marshal(resourceAttributes)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(b), nil
}

func tfidentityProvider(v reflect.Value) (any, error) {
	// Most of the time, generate an empty identity
	if rand.IntN(5) > 2 {
		return json.RawMessage(""), nil
	}

	// Generate a simple identity structure
	identity := map[string]any{
		"arn":        generateARN("iam", fmt.Sprintf("user/%s", generateUserName())),
		"account_id": generateAWSAccountID(),
		"region":     generateAWSRegion(),
	}

	b, err := json.Marshal(identity)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

func tfprivateProvider(v reflect.Value) (any, error) {
	// Occasionally generate a non-empty private field
	if rand.IntN(5) == 0 {
		// Generate some random bytes and encode as a base64 string
		bytes := make([]byte, 16)
		for i := range bytes {
			bytes[i] = byte(rand.IntN(256))
		}
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
	return "", nil
}

func tfdependenciesProvider(v reflect.Value) (any, error) {
	// Generate a list of dependencies (0-3) for the resource
	numDeps := rand.IntN(4)
	dependencies := make([]string, numDeps)

	for i := 0; i < numDeps; i++ {
		resourceType := generateResourceType()
		resourceName := generateResourceName()
		var moduleAddress string
		if rand.IntN(10) < 3 {
			moduleAddress = generateModuleAddress()
		}
		dep := fmt.Sprintf("%s.%s", resourceType, resourceName)
		if moduleAddress != "" {
			dep = fmt.Sprintf("%s.%s", moduleAddress, dep)
		}
		dependencies[i] = dep
	}

	return dependencies, nil
}

func tfemptystringsliceProvider(v reflect.Value) (any, error) {
	// Always return an empty string slice
	return []string{}, nil
}
