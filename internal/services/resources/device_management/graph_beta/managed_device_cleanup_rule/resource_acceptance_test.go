package graphBetaManagedDeviceCleanupRule_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccManagedDeviceCleanupRuleResource_Platforms(t *testing.T) {
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
				CheckDestroy: testAccCheckManagedDeviceCleanupRuleDestroy,
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

func TestAccManagedDeviceCleanupRuleResource_RequiredAndInvalid(t *testing.T) {
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

func testAccCheckManagedDeviceCleanupRuleDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			ManagedDeviceCleanupRules().
			ByManagedDeviceCleanupRuleId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "404") || strings.Contains(strings.ToLower(errStr), "not found") || strings.Contains(strings.ToLower(errStr), "does not exist") {
				continue
			}
			continue
		}
		return fmt.Errorf("Managed Device Cleanup Rule %s still exists", rs.Primary.ID)
	}
	return nil
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
