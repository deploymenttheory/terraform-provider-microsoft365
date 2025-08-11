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

	// TF_ACC must be specifically set to "1" to enable acceptance tests
	if tfAcc := os.Getenv("TF_ACC"); tfAcc != "1" {
		t.Fatalf("TF_ACC must be set to '1' for acceptance tests (got: '%s')", tfAcc)
	}
}
