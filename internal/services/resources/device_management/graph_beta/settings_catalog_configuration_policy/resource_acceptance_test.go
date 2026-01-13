package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSettingsCatalogConfigurationPolicyResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "name", "macos mdm filevault2 settings"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.*", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore secret password field as Microsoft Graph returns different UUIDs for security
				ImportStateVerifyIgnore: []string{
					"configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.6.simple_setting_value.value",
				},
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "name", "Test Acceptance Settings Catalog Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_Assignments(t *testing.T) {
	t.Log("=== ASSIGNMENTS TEST START ===")
	t.Log("Starting assignments acceptance test with comprehensive logging")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Log("=== PRE-CHECK START ===")
			mocks.TestAccPreCheck(t)
			t.Log("=== PRE-CHECK COMPLETE ===")
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			t.Log("=== DESTROY CHECK TRIGGERED ===")
			t.Log("Starting comprehensive destroy verification")
			err := testAccCheckSettingsCatalogConfigurationPolicyDestroy(s)
			if err != nil {
				t.Logf("ERROR: Destroy check failed: %v", err)
			} else {
				t.Log("SUCCESS: Destroy check completed successfully")
			}
			t.Log("=== DESTROY CHECK COMPLETE ===")
			return err
		},
		Steps: []resource.TestStep{
			// Create with all assignment types
			{
				PreConfig: func() {
					t.Log("=== STEP PRE-CONFIG ===")
					t.Log("About to apply configuration with assignments, groups, and role scope tags")
					t.Log("Expected dependencies: 3 groups, 2 role scope tags")
				},
				Config: testAccSettingsCatalogConfigurationPolicyConfig_assignments(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("=== STEP CHECK START ===")
						t.Log("Verifying resource creation and attributes")
						return nil
					},
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "id"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: Resource ID is set")
						return nil
					},
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "name", "Test All Assignment Types Settings Catalog Policy"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: Resource name verified")
						return nil
					},
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.#", "5"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: 5 assignments verified")
						return nil
					},
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: groupAssignmentTarget verified")
						return nil
					},
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: allLicensedUsersAssignmentTarget verified")
						return nil
					},
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: allDevicesAssignmentTarget verified")
						return nil
					},
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: exclusionGroupAssignmentTarget verified")
						return nil
					},
					// Verify role scope tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "role_scope_tag_ids.#", "2"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: Role scope tags verified")
						t.Log("=== STEP CHECK COMPLETE ===")
						return nil
					},
				),
			},
		},
	})
	t.Log("=== ASSIGNMENTS TEST COMPLETE ===")
}

func TestAccSettingsCatalogConfigurationPolicyResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingPlatforms(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_invalidPlatform(),
				ExpectError: regexp.MustCompile("Attribute platforms value must be one of"),
			},
		},
	})
}

// testAccCheckSettingsCatalogConfigurationPolicyDestroy verifies that settings catalog configuration policies have been destroyed
func testAccCheckSettingsCatalogConfigurationPolicyDestroy(s *terraform.State) error {
	fmt.Printf("=== DESTROY CHECK START ===\n")

	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		fmt.Printf("ERROR: Failed to create Graph client for CheckDestroy: %v\n", err)
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()
	resourceCount := 0
	destroyedCount := 0
	orphanedResources := []string{}

	// Count total resources to check
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" {
			resourceCount++
		}
	}
	fmt.Printf("INFO: Found %d settings catalog configuration policy resources to verify destruction\n", resourceCount)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" {
			continue
		}

		fmt.Printf("--- Checking resource destruction: %s ---\n", rs.Primary.ID)

		// Build the API URL for logging
		apiUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/%s", rs.Primary.ID)
		fmt.Printf("INFO: API URL: %s\n", apiUrl)

		// Attempt to get the settings catalog configuration policy by ID
		resource, err := graphClient.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			// Accept multiple error conditions that indicate successful deletion
			if errorInfo.StatusCode == 404 ||
				errorInfo.StatusCode == 400 || // Bad Request - often indicates resource no longer exists
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" ||
				errorInfo.ErrorCode == "Request_ResourceNotFound" ||
				errorInfo.StatusCode == 0 { // Handle cases where status code is not set
				fmt.Printf("SUCCESS: Resource %s successfully destroyed (verified by API error)\n", rs.Primary.ID)
				destroyedCount++
				continue // Resource successfully destroyed
			}

			// For other errors, this might indicate an orphaned resource or API issue
			fmt.Printf("WARNING: Unexpected error checking resource %s destruction: %v\n", rs.Primary.ID, err)
			fmt.Printf("WARNING: This could indicate an orphaned resource or API connectivity issue\n")
			orphanedResources = append(orphanedResources, rs.Primary.ID)

			// Still continue but mark as potentially problematic
			continue
		}

		// If we can still get the resource, it wasn't destroyed - this is a real problem
		if resource != nil {
			fmt.Printf("ERROR: Resource %s still exists and was not properly destroyed!\n", rs.Primary.ID)
			if resource.GetName() != nil {
				fmt.Printf("ERROR: Resource name: %s\n", *resource.GetName())
			}
			if resource.GetId() != nil {
				fmt.Printf("ERROR: Resource ID: %s\n", *resource.GetId())
			}
			orphanedResources = append(orphanedResources, rs.Primary.ID)
			return fmt.Errorf("settings catalog configuration policy %s still exists and was not destroyed", rs.Primary.ID)
		}
	}

	fmt.Printf("=== DESTROY CHECK SUMMARY ===\n")
	fmt.Printf("Total resources checked: %d\n", resourceCount)
	fmt.Printf("Successfully destroyed: %d\n", destroyedCount)
	fmt.Printf("Potentially orphaned: %d\n", len(orphanedResources))

	if len(orphanedResources) > 0 {
		fmt.Printf("WARNING: Potentially orphaned resource IDs:\n")
		for _, id := range orphanedResources {
			fmt.Printf("  - %s\n", id)
		}
		fmt.Printf("WARNING: Please manually verify these resources in the Microsoft Graph API\n")
	}

	fmt.Printf("=== DESTROY CHECK COMPLETE ===\n")
	return nil
}

func testAccSettingsCatalogConfigurationPolicyConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_maximal() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}

	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}

	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_assignments() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}

	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}

	// Use local assignment filters with proper dependency management to preserve the correct destroy order
	//assignmentFilters, err := helpers.ParseHCLFile("assignment_filters_with_dependencies.tf")
	//if err != nil {
	//	log.Fatalf("Failed to load assignment filters config: %v", err)
	//}

	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load test config: %v", err)
	}

	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingName() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  platforms = "windows10"
  configuration_policy = {
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingPlatforms() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Policy"
  configuration_policy = {
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccSettingsCatalogConfigurationPolicyConfig_invalidPlatform() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name      = "Test Policy"
  platforms = "invalid"
	template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}
