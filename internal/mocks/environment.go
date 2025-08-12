package mocks

import (
	"os"
	"testing"
)

func SetupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// TestAccPreCheck verifies necessary test prerequisites
func TestAccPreCheck(t *testing.T) {
	// Skip if not running acceptance tests
	tfAcc := os.Getenv("TF_ACC")
	if tfAcc != "1" {
		t.Skipf("Acceptance tests skipped unless env 'TF_ACC' set to '1' (current value: %q)", tfAcc)
		return
	}

	// Check for required environment variables
	requiredEnvVars := []string{
		"M365_CLIENT_ID",
		"M365_CLIENT_SECRET",
		"M365_TENANT_ID",
		"M365_AUTH_METHOD",
		"M365_CLOUD",
	}

	for _, envVar := range requiredEnvVars {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}
}
