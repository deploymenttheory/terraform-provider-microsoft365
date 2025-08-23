package graphBetaWindowsDeviceCompliancePolicy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsDeviceCompliancePolicyResource_CustomCompliance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceCompliancePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCustomCompliance(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "display_name", "acc-test-windows-device-compliance-policy-custom-compliance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "description", "acc-test-windows-device-compliance-policy-custom-compliance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "custom_compliance_required", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "device_compliance_policy_script.device_compliance_script_id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "device_compliance_policy_script.rules_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "scheduled_actions_for_rule.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccWindowsDeviceCompliancePolicyResource_DeviceHealth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceCompliancePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDeviceHealth(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_health", "display_name", "acc-test-windows-device-compliance-policy-device-health"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_health", "device_health.bit_locker_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_health", "device_health.secure_boot_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_health", "device_health.code_integrity_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_health", "id"),
				),
			},
		},
	})
}

func TestAccWindowsDeviceCompliancePolicyResource_DeviceProperties(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceCompliancePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDeviceProperties(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_properties", "display_name", "acc-test-windows-device-compliance-policy-device-properties"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_properties", "device_properties.os_minimum_version", "10.0.22631.5768"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_properties", "device_properties.os_maximum_version", "10.0.26100.9999"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_properties", "device_properties.valid_operating_system_build_ranges.#", "2"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.device_properties", "id"),
				),
			},
		},
	})
}

func TestAccWindowsDeviceCompliancePolicyResource_MicrosoftDefenderForEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceCompliancePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMicrosoftDefenderForEndpoint(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.microsoft_defender_for_endpoint", "display_name", "acc-test-windows-device-compliance-policy-microsoft-defender-for-endpoint"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.microsoft_defender_for_endpoint", "microsoft_defender_for_endpoint.device_threat_protection_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.microsoft_defender_for_endpoint", "microsoft_defender_for_endpoint.device_threat_protection_required_security_level", "medium"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.microsoft_defender_for_endpoint", "id"),
				),
			},
		},
	})
}

func TestAccWindowsDeviceCompliancePolicyResource_WSL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceCompliancePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWSL(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "display_name", "acc-test-windows-device-compliance-policy-wsl"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "wsl_distributions.#", "2"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "id"),
				),
			},
		},
	})
}

func TestAccWindowsDeviceCompliancePolicyResource_WSLAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceCompliancePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWSLAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl_assignments", "display_name", "acc-test-windows-device-compliance-policy-wsl-assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl_assignments", "assignments.#", "6"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl_assignments", "id"),
				),
			},
		},
	})
}

func testAccConfigCustomCompliance() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	scriptConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/windows_device_compliance_script.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load script dependency: %s", err.Error()))
	}

	notificationConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_compliance_notification_template.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load notification template dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/compliance_policy_custom_compliance.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + scriptConfig + "\n\n" + notificationConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigDeviceHealth() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	notificationConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_compliance_notification_template.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load notification template dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/compliance_policy_device_health.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + notificationConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigDeviceProperties() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	notificationConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_compliance_notification_template.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load notification template dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/compliance_policy_device_properties.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + notificationConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigMicrosoftDefenderForEndpoint() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	notificationConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_compliance_notification_template.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load notification template dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/compliance_policy_microsoft_defender_for_endpoint.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + notificationConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigWSL() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	notificationConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_compliance_notification_template.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load notification template dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/compliance_policy_wsl.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + notificationConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigWSLAssignments() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	assignmentFilterConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load assignment filter dependency: %s", err.Error()))
	}

	notificationConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_compliance_notification_template.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load notification template dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/compliance_policy_wsl_assignments.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + assignmentFilterConfig + "\n\n" + notificationConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccCheckWindowsDeviceCompliancePolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_device_compliance_policy" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			DeviceCompliancePolicies().
			ByDeviceCompliancePolicyId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows device compliance policy %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows device compliance policy %s still exists", rs.Primary.ID)
	}
	return nil
}