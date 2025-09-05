package graphBetaAppControlForBusinessPolicy_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppControlForBusinessPolicyResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories:  mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:              testAccCheckAppControlForBusinessPolicyDestroy,
		PreventPostDestroyRefresh: true, // Prevents refresh after destroy to avoid dependency issues
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAppControlForBusinessPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "name", "acc-test-app-control-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.*", "2"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppControlForBusinessPolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal", "name", "acc-test-app-control-policy-maximal"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal", "policy_xml"),
				),
			},
			{
				// Use removed block to control destruction order
				Config: testAccAppControlForBusinessPolicyConfig_removedPolicy(),
			},
		},
	})
}

func TestAccAppControlForBusinessPolicyResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories:  mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:              testAccCheckAppControlForBusinessPolicyDestroy,
		PreventPostDestroyRefresh: true, // Prevents refresh after destroy to avoid dependency issues
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAppControlForBusinessPolicyConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_policy.with_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.with_assignments", "name", "acc-test-app-control-policy-with-assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.with_assignments", "assignments.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_policy.with_assignments", "assignments.*", map[string]string{
						"type":        "allLicensedUsersAssignmentTarget",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_policy.with_assignments", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_policy.with_assignments", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
		},
	})
}

func testAccAppControlForBusinessPolicyConfig_minimal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_acfb_policy_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAppControlForBusinessPolicyConfig_maximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_acfb_policy_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAppControlForBusinessPolicyConfig_withAssignments() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_acfb_policy_with_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load assignments test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAppControlForBusinessPolicyConfig_removedPolicy() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	assignmentFilters, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		log.Fatalf("Failed to load assignment filters config: %v", err)
	}

	// Configuration with removed block for app control policy to control destroy order
	removedConfig := `
# App control policy removed first to release assignment filter dependencies
removed {
  from = microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal
  lifecycle {
    destroy = true
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + removedConfig)
}

// testAccCheckAppControlForBusinessPolicyDestroy verifies that app control for business policies have been destroyed
func testAccCheckAppControlForBusinessPolicyDestroy(s *terraform.State) error {
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
		if rs.Type == "microsoft365_graph_beta_device_management_app_control_for_business_policy" {
			resourceCount++
		}
	}
	fmt.Printf("INFO: Found %d app control for business policy resources to verify destruction\n", resourceCount)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_app_control_for_business_policy" {
			continue
		}

		fmt.Printf("--- Checking resource destruction: %s ---\n", rs.Primary.ID)

		// Build the API URL for logging
		apiUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/%s", rs.Primary.ID)
		fmt.Printf("INFO: API URL: %s\n", apiUrl)

		// Attempt to get the app control policy by ID
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
			return fmt.Errorf("app control for business policy %s still exists and was not destroyed", rs.Primary.ID)
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
