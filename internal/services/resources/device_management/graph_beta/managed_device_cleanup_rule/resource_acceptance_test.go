package graphBetaManagedDeviceCleanupRule_test

import (
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaManagedDeviceCleanupRule "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/managed_device_cleanup_rule"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaManagedDeviceCleanupRule.ResourceName

var testResource = graphBetaManagedDeviceCleanupRule.ManagedDeviceCleanupRuleTestResource{}

func TestAccResourceManagedDeviceCleanupRule_01_Platforms(t *testing.T) {
	platforms := []struct {
		name   string
		path   string
		resRef string
	}{
		{"all", "tests/terraform/acceptance/platform_all.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.all"},
		{"androidAOSP", "tests/terraform/acceptance/platform_androidAOSP.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.androidAOSP"},
		{"androidDeviceAdministrator", "tests/terraform/acceptance/platform_androidDeviceAdministrator.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.androidDeviceAdministrator"},
		{"androidDedicatedAndFullyManagedCorporateOwnedWorkProfile", "tests/terraform/acceptance/platform_androidDedicatedAndFullyManagedCorporateOwnedWorkProfile.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.androidDedicatedAndFullyManagedCorporateOwnedWorkProfile"},
		{"chromeOS", "tests/terraform/acceptance/platform_chromeOS.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.chromeOS"},
		{"androidPersonallyOwnedWorkProfile", "tests/terraform/acceptance/platform_androidPersonallyOwnedWorkProfile.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.androidPersonallyOwnedWorkProfile"},
		{"ios", "tests/terraform/acceptance/platform_ios.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.ios"},
		{"macOS", "tests/terraform/acceptance/platform_macOS.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.macOS"},
		{"windows", "tests/terraform/acceptance/platform_windows.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.windows"},
		{"windowsHolographic", "tests/terraform/acceptance/platform_windowsHolographic.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.windowsHolographic"},
	}

	for _, tc := range platforms {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				ExternalProviders: map[string]resource.ExternalProvider{
					"random": {
						Source:            "hashicorp/random",
						VersionConstraint: ">= 3.7.2",
					},
				},
				CheckDestroy: destroy.CheckDestroyedAllFunc(
					testResource,
					resourceType,
					30*time.Second,
				),
				Steps: []resource.TestStep{
					{
						Config: testAccConfigFromFile(tc.path),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(tc.resRef, "id"),
						),
					},
					{
						ResourceName:      tc.resRef,
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		})
	}
}

func TestAccResourceManagedDeviceCleanupRule_02_RequiredAndInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccManagedDeviceCleanupRule_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccManagedDeviceCleanupRule_invalidPlatform(),
				ExpectError: regexp.MustCompile("Attribute device_cleanup_rule_platform_type value must be one of"),
			},
		},
	})
}

func testAccConfigFromFile(path string) string {
	accTestConfig, err := helpers.ParseHCLFile(path)
	if err != nil {
		log.Fatalf("Failed to load test config %s: %v", path, err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccManagedDeviceCleanupRule_missingDisplayName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_missing_display_name.tf")
	if err != nil {
		log.Fatalf("Failed to load missing display name test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccManagedDeviceCleanupRule_invalidPlatform() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_invalid_platform.tf")
	if err != nil {
		log.Fatalf("Failed to load invalid platform test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
