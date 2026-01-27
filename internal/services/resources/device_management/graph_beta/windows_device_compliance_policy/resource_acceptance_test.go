package graphBetaWindowsDeviceCompliancePolicy_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsDeviceCompliancePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_device_compliance_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaWindowsDeviceCompliancePolicy.ResourceName

	// testResource is the test resource implementation for windows device compliance policies
	testResource = graphBetaWindowsDeviceCompliancePolicy.WindowsDeviceCompliancePolicyTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Scenario 1 - custom compliance
func TestAccResourceWindowsDeviceCompliancePolicy_01_CustomCompliance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating custom compliance policy")
				},
				Config: loadAcceptanceTestTerraform("compliance_policy_custom_compliance.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows device compliance policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".custom_compliance").ExistsInGraph(testResource),
					check.That(resourceType+".custom_compliance").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_compliance").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-wdcp-custom-compliance-[a-z0-9]{8}$`)),
					check.That(resourceType+".custom_compliance").Key("description").MatchesRegex(regexp.MustCompile(`^acc-test-wdcp-custom-compliance-[a-z0-9]{8}$`)),
					check.That(resourceType+".custom_compliance").Key("custom_compliance_required").HasValue("true"),
					check.That(resourceType+".custom_compliance").Key("device_compliance_policy_script.device_compliance_script_id").Exists(),
					check.That(resourceType+".custom_compliance").Key("device_compliance_policy_script.rules_content").Exists(),
					check.That(resourceType+".custom_compliance").Key("scheduled_actions_for_rule.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing custom compliance policy")
				},
				ResourceName:            resourceType + ".custom_compliance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 02: Scenario 2 - device health
func TestAccResourceWindowsDeviceCompliancePolicy_02_DeviceHealth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating device health policy")
				},
				Config: loadAcceptanceTestTerraform("compliance_policy_device_health.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows device compliance policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".device_health").ExistsInGraph(testResource),
					check.That(resourceType+".device_health").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".device_health").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-wdcp-device-health-[a-z0-9]{8}$`)),
					check.That(resourceType+".device_health").Key("device_health.bit_locker_enabled").HasValue("true"),
					check.That(resourceType+".device_health").Key("device_health.secure_boot_enabled").HasValue("true"),
					check.That(resourceType+".device_health").Key("device_health.code_integrity_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing device health policy")
				},
				ResourceName:            resourceType + ".device_health",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 03: Scenario 3 - device properties
func TestAccResourceWindowsDeviceCompliancePolicy_03_DeviceProperties(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating device properties policy")
				},
				Config: loadAcceptanceTestTerraform("compliance_policy_device_properties.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows device compliance policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".device_properties").ExistsInGraph(testResource),
					check.That(resourceType+".device_properties").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".device_properties").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-dcnt-device-properties-[a-z0-9]{8}$`)),
					check.That(resourceType+".device_properties").Key("device_properties.os_minimum_version").HasValue("10.0.22631.5768"),
					check.That(resourceType+".device_properties").Key("device_properties.os_maximum_version").HasValue("10.0.26100.9999"),
					check.That(resourceType+".device_properties").Key("device_properties.valid_operating_system_build_ranges.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing device properties policy")
				},
				ResourceName:            resourceType + ".device_properties",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 04: Scenario 4 - Microsoft Defender for Endpoint
func TestAccResourceWindowsDeviceCompliancePolicy_04_MicrosoftDefenderForEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Microsoft Defender for Endpoint policy")
				},
				Config: loadAcceptanceTestTerraform("compliance_policy_microsoft_defender_for_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows device compliance policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".microsoft_defender_for_endpoint").ExistsInGraph(testResource),
					check.That(resourceType+".microsoft_defender_for_endpoint").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".microsoft_defender_for_endpoint").Key("display_name").HasValue("acc-test-windows-device-compliance-policy-microsoft-defender-for-endpoint"),
					check.That(resourceType+".microsoft_defender_for_endpoint").Key("microsoft_defender_for_endpoint.device_threat_protection_enabled").HasValue("true"),
					check.That(resourceType+".microsoft_defender_for_endpoint").Key("microsoft_defender_for_endpoint.device_threat_protection_required_security_level").HasValue("medium"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Microsoft Defender for Endpoint policy")
				},
				ResourceName:            resourceType + ".microsoft_defender_for_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 05: Scenario 5 - WSL
func TestAccResourceWindowsDeviceCompliancePolicy_05_WSL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating WSL policy")
				},
				Config: loadAcceptanceTestTerraform("compliance_policy_wsl.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows device compliance policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".wsl").ExistsInGraph(testResource),
					check.That(resourceType+".wsl").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wsl").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-wdcp-wsl-[a-z0-9]{8}$`)),
					check.That(resourceType+".wsl").Key("wsl_distributions.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing WSL policy")
				},
				ResourceName:            resourceType + ".wsl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 06: Scenario 6 - WSL assignments
func TestAccResourceWindowsDeviceCompliancePolicy_06_WSL_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating WSL policy with assignments")
				},
				Config: loadAcceptanceTestTerraform("compliance_policy_wsl_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows device compliance policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".wsl_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".wsl_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wsl_assignments").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-wdcp-wsl-assignments-[a-z0-9]{8}$`)),
					check.That(resourceType+".wsl_assignments").Key("assignments.#").HasValue("6"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing WSL policy with assignments")
				},
				ResourceName:            resourceType + ".wsl_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
