package graphBetaWindowsDeviceComplianceNotifications

import (
	"fmt"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
)

func TestConfigParsing() error {
	// Test parsing the minimal configuration
	_, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		return fmt.Errorf("failed to parse minimal config: %v", err)
	}

	// Test parsing the maximal configuration
	_, err = helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		return fmt.Errorf("failed to parse maximal config: %v", err)
	}

	// Test parsing the branding test configuration
	_, err = helpers.ParseHCLFile("tests/terraform/acceptance/resource_branding_test.tf")
	if err != nil {
		return fmt.Errorf("failed to parse branding test config: %v", err)
	}

	return nil
}