package statefaker

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestStateValid(t *testing.T) {
	// Generate a fake state with some outputs and resources
	state, err := NewFakeStateV4(
		WithOutputs(3),
		WithResources(5),
	)
	if err != nil {
		t.Fatalf("failed to generate fake state: %v", err)
	}

	// Marshal the state to JSON
	stateJSON, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal state to JSON: %v", err)
	}

	// Create a temporary file to write the state
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "test.tfstate")

	err = os.WriteFile(stateFile, stateJSON, 0644)
	if err != nil {
		t.Fatalf("failed to write state file: %v", err)
	}

	// Run terraform state list to validate the state file
	cmd := exec.Command("terraform", "state", "list", "-state="+stateFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("terraform state list failed: %v\nOutput: %s", err, string(output))
	}

	// Check that we got some output (indicating resources were found)
	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		t.Error("terraform state list returned no resources, but we expected some")
	}

	// Verify that the number of resources listed matches what we expect
	// Note: Some resources may have multiple instances, so we count unique resource names
	resourceLines := strings.Split(outputStr, "\n")
	uniqueResources := make(map[string]bool)

	for _, line := range resourceLines {
		if line == "" {
			continue
		}
		// Extract the resource name (everything before the first '[' for indexed resources)
		resourceName := strings.Split(line, "[")[0]
		uniqueResources[resourceName] = true
	}

	expectedResources := 5 // Based on our WithResources(5) option
	if len(uniqueResources) != expectedResources {
		t.Errorf("expected %d unique resources, but got %d. Unique resources: %v",
			expectedResources, len(uniqueResources), uniqueResources)
	}

	t.Logf("terraform state list succeeded with %d total instances from %d unique resources:\n%s",
		len(resourceLines), len(uniqueResources), outputStr)
}
