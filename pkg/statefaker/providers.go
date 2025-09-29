package statefaker

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"reflect"

	"github.com/go-faker/faker/v4"
)

// Helper functions for generating realistic AWS data
func generateAWSAccountID() string {
	return fmt.Sprintf("%012d", rand.IntN(1000000000000))
}

func generateAWSRegion() string {
	regions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1", "ca-central-1"}
	return regions[rand.IntN(len(regions))]
}

func generateARN(service, resource string) string {
	return fmt.Sprintf("arn:aws:%s:%s:%s:%s", service, generateAWSRegion(), generateAWSAccountID(), resource)
}

func generateAccessKeyID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := make([]byte, 20)
	for i := range key {
		key[i] = charset[rand.IntN(len(charset))]
	}
	return "AKIA" + string(key[4:])
}

func generateS3BucketName() string {
	prefixes := []string{"turo", "company", "app", "data", "backup", "logs", "config"}
	suffixes := []string{"prod", "staging", "dev", "test", "ml", "analytics", "artifacts"}
	middle := []string{"xyz", "main", "core", "service", "data", "bucket"}

	return fmt.Sprintf("%s-%s-%s",
		prefixes[rand.IntN(len(prefixes))],
		middle[rand.IntN(len(middle))],
		suffixes[rand.IntN(len(suffixes))])
}

func generateUserName() string {
	roles := []string{"reader", "writer", "admin", "analyst", "developer", "operator"}
	teams := []string{"ml", "data", "api", "web", "mobile", "infra", "security"}

	if rand.IntN(2) == 0 {
		return fmt.Sprintf("%s-%s", teams[rand.IntN(len(teams))], roles[rand.IntN(len(roles))])
	}
	return roles[rand.IntN(len(roles))]
}

func generateResourceType() string {
	resourceTypes := []string{
		"aws_s3_bucket", "aws_iam_user", "aws_iam_role", "aws_lambda_function",
		"aws_ec2_instance", "aws_rds_instance", "aws_dynamodb_table", "aws_vpc",
		"aws_security_group", "aws_route53_zone", "aws_cloudfront_distribution",
		"aws_ecs_cluster", "aws_eks_cluster", "aws_api_gateway_rest_api",
	}
	return resourceTypes[rand.IntN(len(resourceTypes))]
}

func generateResourceName() string {
	names := []string{"main", "primary", "secondary", "backup", "test", "prod", "staging", "dev", "example"}
	return names[rand.IntN(len(names))]
}

func getProviderFromResourceType(resourceType string) string {
	if len(resourceType) >= 3 && resourceType[:3] == "aws" {
		return "aws"
	}
	if len(resourceType) >= 5 && resourceType[:5] == "azurerm" {
		return "azurerm"
	}
	if len(resourceType) >= 6 && resourceType[:6] == "google" {
		return "google"
	}
	if len(resourceType) >= 10 && resourceType[:10] == "kubernetes" {
		return "kubernetes"
	}
	// Default to aws for most cases
	return "aws"
}

// randomComplexOutput generates complex realistic output structures
func randomComplexOutput(output *OutputV4) {
	outputTypes := []func(*OutputV4){
		generateS3BucketPolicyOutput,
		generateUserMapOutput,
		generateDatabaseConfigOutput,
		generateNetworkConfigOutput,
		generateSecurityGroupOutput,
	}

	generator := outputTypes[rand.IntN(len(outputTypes))]
	generator(output)
}

func generateS3BucketPolicyOutput(output *OutputV4) {
	bucketName := generateS3BucketName()
	accountID := generateAWSAccountID()
	userName := generateUserName()

	policy := map[string]any{
		"Version": "2012-10-17",
		"Statement": []map[string]any{
			{
				"Effect": "Allow",
				"Action": "s3:ListBucket",
				"Resource": []string{
					generateARN("s3", bucketName+"/*"),
					generateARN("s3", bucketName),
				},
				"Principal": map[string]string{
					"AWS": generateARN("iam", fmt.Sprintf("user/%s", userName)),
				},
			},
			{
				"Effect": "Allow",
				"Action": "s3:GetObject",
				"Resource": []string{
					generateARN("s3", bucketName+"/*"),
				},
				"Principal": map[string]string{
					"AWS": fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, userName),
				},
			},
		},
	}

	policyJSON, _ := json.Marshal(policy)
	typeJSON, _ := json.Marshal("string")
	output.Type = json.RawMessage(typeJSON)
	output.Value = json.RawMessage(policyJSON)
}

func generateUserMapOutput(output *OutputV4) {
	users := make(map[string]map[string]any)

	for i := 0; i < rand.IntN(5)+2; i++ {
		userName := generateUserName()
		users[userName] = map[string]any{
			"access_key_id":               generateAccessKeyID(),
			"encrypted_secret_access_key": faker.Password(),
			"pgp_key_name": map[string]string{
				"name":              "aws-pgp-v0-2020-07-08.pgp.base64",
				"public_key_base64": faker.Password(), // Simplified for example
			},
		}
	}

	typeStructure := []any{
		"object",
		map[string][]any{},
	}
	typeJSON, _ := json.Marshal(typeStructure)
	valueJSON, _ := json.Marshal(users)
	output.Type = json.RawMessage(typeJSON)
	output.Value = json.RawMessage(valueJSON)
}

func generateDatabaseConfigOutput(output *OutputV4) {
	config := map[string]any{
		"endpoint":                fmt.Sprintf("%s.%s.rds.amazonaws.com", faker.Username(), generateAWSRegion()),
		"port":                    5432,
		"database":                faker.Username(),
		"username":                faker.Username(),
		"password":                faker.Password(),
		"ssl_mode":                "require",
		"max_connections":         rand.IntN(100) + 10,
		"backup_retention_period": rand.IntN(30) + 1,
	}

	typeStructure := []any{
		"object",
		map[string]string{
			"endpoint":                "string",
			"port":                    "number",
			"database":                "string",
			"username":                "string",
			"password":                "string",
			"ssl_mode":                "string",
			"max_connections":         "number",
			"backup_retention_period": "number",
		},
	}
	typeJSON, _ := json.Marshal(typeStructure)
	valueJSON, _ := json.Marshal(config)
	output.Type = json.RawMessage(typeJSON)
	output.Value = json.RawMessage(valueJSON)
}

func generateNetworkConfigOutput(output *OutputV4) {
	config := map[string]any{
		"vpc_id": fmt.Sprintf("vpc-%s", faker.UUIDDigit()),
		"subnet_ids": []string{
			fmt.Sprintf("subnet-%s", faker.UUIDDigit()),
			fmt.Sprintf("subnet-%s", faker.UUIDDigit()),
		},
		"security_group_ids": []string{
			fmt.Sprintf("sg-%s", faker.UUIDDigit()),
		},
		"availability_zones": []string{
			generateAWSRegion() + "a",
			generateAWSRegion() + "b",
		},
		"cidr_block": "10.0.0.0/16",
	}

	typeStructure := []any{
		"object",
		map[string]any{
			"vpc_id":             "string",
			"subnet_ids":         []string{"string"},
			"security_group_ids": []string{"string"},
			"availability_zones": []string{"string"},
			"cidr_block":         "string",
		},
	}
	typeJSON, _ := json.Marshal(typeStructure)
	valueJSON, _ := json.Marshal(config)
	output.Type = json.RawMessage(typeJSON)
	output.Value = json.RawMessage(valueJSON)
}

func generateSecurityGroupOutput(output *OutputV4) {
	rules := make([]map[string]any, rand.IntN(5)+1)
	for i := range rules {
		rules[i] = map[string]any{
			"type":        []string{"ingress", "egress"}[rand.IntN(2)],
			"protocol":    []string{"tcp", "udp", "icmp"}[rand.IntN(3)],
			"from_port":   rand.IntN(65535),
			"to_port":     rand.IntN(65535),
			"cidr_blocks": []string{"0.0.0.0/0"},
		}
	}

	config := map[string]any{
		"id":          fmt.Sprintf("sg-%s", faker.UUIDDigit()),
		"name":        fmt.Sprintf("%s-sg", generateResourceName()),
		"description": faker.Sentence(),
		"rules":       rules,
		"vpc_id":      fmt.Sprintf("vpc-%s", faker.UUIDDigit()),
	}

	typeStructure := []any{
		"object",
		map[string]any{
			"id":          "string",
			"name":        "string",
			"description": "string",
			"rules":       []string{"object"},
			"vpc_id":      "string",
		},
	}
	typeJSON, _ := json.Marshal(typeStructure)
	valueJSON, _ := json.Marshal(config)
	output.Type = json.RawMessage(typeJSON)
	output.Value = json.RawMessage(valueJSON)
}

// Attribute generators for different resource types
func generateS3BucketAttributes() map[string]any {
	bucketName := generateS3BucketName()
	return map[string]any{
		"id":                          bucketName,
		"arn":                         generateARN("s3", bucketName),
		"bucket":                      bucketName,
		"bucket_domain_name":          fmt.Sprintf("%s.s3.amazonaws.com", bucketName),
		"bucket_regional_domain_name": fmt.Sprintf("%s.s3.%s.amazonaws.com", bucketName, generateAWSRegion()),
		"region":                      generateAWSRegion(),
		"versioning": []map[string]any{
			{
				"enabled":    rand.IntN(2) == 1,
				"mfa_delete": false,
			},
		},
		"server_side_encryption_configuration": []map[string]any{
			{
				"rule": []map[string]any{
					{
						"apply_server_side_encryption_by_default": []map[string]any{
							{
								"sse_algorithm": "AES256",
							},
						},
					},
				},
			},
		},
		"tags": map[string]string{
			"Environment": []string{"prod", "staging", "dev"}[rand.IntN(3)],
			"Team":        []string{"data", "ml", "web", "mobile"}[rand.IntN(4)],
		},
	}
}

func generateIAMUserAttributes() map[string]any {
	userName := generateUserName()
	return map[string]any{
		"id":                   userName,
		"arn":                  generateARN("iam", fmt.Sprintf("user/%s", userName)),
		"name":                 userName,
		"path":                 "/",
		"permissions_boundary": nil,
		"unique_id":            fmt.Sprintf("AIDA%s", faker.UUIDDigit()[:16]),
		"tags": map[string]string{
			"Role": []string{"reader", "writer", "admin"}[rand.IntN(3)],
			"Team": []string{"data", "ml", "security"}[rand.IntN(3)],
		},
	}
}

func generateEC2InstanceAttributes() map[string]any {
	instanceID := fmt.Sprintf("i-%s", faker.UUIDDigit()[:17])
	return map[string]any{
		"id":                     instanceID,
		"arn":                    generateARN("ec2", fmt.Sprintf("instance/%s", instanceID)),
		"instance_id":            instanceID,
		"instance_type":          []string{"t3.micro", "t3.small", "m5.large", "c5.xlarge"}[rand.IntN(4)],
		"ami":                    fmt.Sprintf("ami-%s", faker.UUIDDigit()[:17]),
		"availability_zone":      generateAWSRegion() + []string{"a", "b", "c"}[rand.IntN(3)],
		"private_ip":             fmt.Sprintf("10.0.%d.%d", rand.IntN(255), rand.IntN(255)),
		"public_ip":              fmt.Sprintf("%d.%d.%d.%d", rand.IntN(255), rand.IntN(255), rand.IntN(255), rand.IntN(255)),
		"subnet_id":              fmt.Sprintf("subnet-%s", faker.UUIDDigit()[:17]),
		"vpc_security_group_ids": []string{fmt.Sprintf("sg-%s", faker.UUIDDigit()[:17])},
		"key_name":               faker.Username(),
		"monitoring":             rand.IntN(2) == 1,
		"state":                  "running",
		"tags": map[string]string{
			"Name":        fmt.Sprintf("%s-instance", generateResourceName()),
			"Environment": []string{"prod", "staging", "dev"}[rand.IntN(3)],
		},
	}
}

func generateLambdaFunctionAttributes() map[string]any {
	functionName := fmt.Sprintf("%s-lambda", generateResourceName())
	return map[string]any{
		"id":               functionName,
		"arn":              generateARN("lambda", fmt.Sprintf("function:%s", functionName)),
		"function_name":    functionName,
		"role":             generateARN("iam", fmt.Sprintf("role/%s-lambda-role", generateResourceName())),
		"handler":          "index.handler",
		"runtime":          []string{"nodejs18.x", "python3.9", "java11", "go1.x"}[rand.IntN(4)],
		"memory_size":      []int{128, 256, 512, 1024}[rand.IntN(4)],
		"timeout":          rand.IntN(900) + 3,
		"last_modified":    faker.Date(),
		"source_code_hash": faker.UUIDDigit(),
		"version":          "$LATEST",
		"environment": []map[string]any{
			{
				"variables": map[string]string{
					"ENV":       []string{"prod", "staging", "dev"}[rand.IntN(3)],
					"LOG_LEVEL": []string{"DEBUG", "INFO", "WARN", "ERROR"}[rand.IntN(4)],
				},
			},
		},
		"tags": map[string]string{
			"Environment": []string{"prod", "staging", "dev"}[rand.IntN(3)],
			"Team":        []string{"backend", "data", "ml"}[rand.IntN(3)],
		},
	}
}

func generateRDSInstanceAttributes() map[string]any {
	instanceID := fmt.Sprintf("%s-db", generateResourceName())
	return map[string]any{
		"id":                      instanceID,
		"arn":                     generateARN("rds", fmt.Sprintf("db:%s", instanceID)),
		"identifier":              instanceID,
		"engine":                  []string{"postgres", "mysql", "mariadb"}[rand.IntN(3)],
		"engine_version":          []string{"13.7", "14.2", "8.0.28"}[rand.IntN(3)],
		"instance_class":          []string{"db.t3.micro", "db.t3.small", "db.r5.large"}[rand.IntN(3)],
		"allocated_storage":       []int{20, 50, 100, 200}[rand.IntN(4)],
		"storage_type":            "gp2",
		"db_name":                 faker.Username(),
		"username":                faker.Username(),
		"port":                    []int{3306, 5432}[rand.IntN(2)],
		"endpoint":                fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", instanceID, faker.UUIDDigit()[:10], generateAWSRegion()),
		"hosted_zone_id":          fmt.Sprintf("Z%s", faker.UUIDDigit()[:13]),
		"status":                  "available",
		"multi_az":                rand.IntN(2) == 1,
		"backup_retention_period": rand.IntN(35) + 1,
		"backup_window":           "03:00-04:00",
		"maintenance_window":      "sun:04:00-sun:05:00",
		"storage_encrypted":       rand.IntN(2) == 1,
		"tags": map[string]string{
			"Environment": []string{"prod", "staging", "dev"}[rand.IntN(3)],
			"Team":        []string{"data", "backend", "analytics"}[rand.IntN(3)],
		},
	}
}

func randomOutput() (json.RawMessage, error) {
	var output OutputV4

	// Half the time, generate a simple output
	if rand.IntN(2) == 0 {
		switch rand.IntN(3) {
		case 0:
			typeJSON, _ := json.Marshal("string")
			valueJSON, _ := json.Marshal(faker.Sentence())
			output.Type = json.RawMessage(typeJSON)
			output.Value = json.RawMessage(valueJSON)
		case 1:
			typeJSON, _ := json.Marshal("number")
			valueJSON, _ := json.Marshal(rand.IntN(1000))
			output.Type = json.RawMessage(typeJSON)
			output.Value = json.RawMessage(valueJSON)
		case 2:
			typeJSON, _ := json.Marshal("bool")
			valueJSON, _ := json.Marshal(rand.IntN(2) == 0)
			output.Type = json.RawMessage(typeJSON)
			output.Value = json.RawMessage(valueJSON)
		}
	} else {
		// Generate a complex output
		randomComplexOutput(&output)
	}

	b, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(b), nil
}

func init() {
	_ = faker.AddProvider("tfresource_managed", func(v reflect.Value) (interface{}, error) {
		var instance1 InstanceV4
		err := faker.FakeData(&instance1)
		if err != nil {
			return nil, fmt.Errorf("failed to fake data for instance1: %w", err)
		}

		instance1.SchemaVersion = 0
		instance1.IdentitySchemaVersion = 0

		resourceType := generateResourceType()
		resourceName := generateResourceName()

		resource := ResourceV4{
			Mode: "managed",
			Type: resourceType,
			Name: resourceName,
			Provider: fmt.Sprintf("provider[\"registry.terraform.io/hashicorp/%s\"]",
				getProviderFromResourceType(resourceType)),
			Instances: []InstanceV4{instance1},
		}

		b, err := json.Marshal(resource)
		if err != nil {
			return nil, err
		}

		return json.RawMessage(b), nil
	})

	_ = faker.AddProvider("tfresource_data", func(v reflect.Value) (interface{}, error) {
		var instance1 InstanceV4
		err := faker.FakeData(&instance1)
		if err != nil {
			return nil, fmt.Errorf("failed to fake data for data instance: %w", err)
		}

		instance1.SchemaVersion = 0
		instance1.IdentitySchemaVersion = 0

		dataResourceTypes := []string{
			"aws_ami", "aws_availability_zones", "aws_caller_identity",
			"aws_region", "aws_s3_bucket", "aws_iam_policy_document",
			"aws_vpc", "aws_subnet", "aws_security_group",
		}

		resourceType := dataResourceTypes[rand.IntN(len(dataResourceTypes))]
		resourceName := generateResourceName()

		resource := ResourceV4{
			Mode: "data",
			Type: resourceType,
			Name: resourceName,
			Provider: fmt.Sprintf("provider[\"registry.terraform.io/hashicorp/%s\"]",
				getProviderFromResourceType(resourceType)),
			Instances: []InstanceV4{instance1},
		}

		b, err := json.Marshal(resource)
		if err != nil {
			return nil, err
		}

		return json.RawMessage(b), nil
	})

	_ = faker.AddProvider("tfattributes", func(v reflect.Value) (interface{}, error) {
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
	})
}
