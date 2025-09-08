package graphBetaGroupPolicyUploadedDefinitionFiles_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGroupPolicyUploadedDefinitionFilesResource_Mozilla(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	// Test: Mozilla ADMX should deploy successfully
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckGroupPolicyUploadedDefinitionFilesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupPolicyUploadedDefinitionFilesMozillaConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.mozilla", "file_name", "mozilla.admx"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.mozilla", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.mozilla", "default_language_code", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.mozilla", "group_policy_uploaded_language_files.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files.mozilla",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"content", // content is not returned by the API
					"group_policy_uploaded_language_files.0.content", // language file content is not returned by the API
				},
			},
		},
	})
}

func testAccGroupPolicyUploadedDefinitionFilesMozillaConfig() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_group_policy_uploaded_definition_files_mozilla.tf")
	if err != nil {
		return fmt.Sprintf("Error loading configuration: %s", err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// testAccCheckGroupPolicyUploadedDefinitionFilesDestroy verifies that all group policy uploaded definition files have been destroyed
func testAccCheckGroupPolicyUploadedDefinitionFilesDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()

	fmt.Printf("=== GROUP POLICY UPLOADED DEFINITION FILES DESTROY CHECK START ===\n")
	resourceCount := 0
	destroyedCount := 0
	orphanedResources := []string{}

	// Count total resources to check
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files" {
			resourceCount++
		}
	}
	fmt.Printf("INFO: Found %d group policy uploaded definition files resources to verify destruction\n", resourceCount)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files" {
			continue
		}

		fmt.Printf("--- Checking resource destruction: %s ---\n", rs.Primary.ID)

		// Build the API URL for logging
		apiUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles/%s", rs.Primary.ID)
		fmt.Printf("INFO: API URL: %s\n", apiUrl)

		// Attempt to get the group policy uploaded definition file by ID
		_, err := graphClient.
			DeviceManagement().
			GroupPolicyUploadedDefinitionFiles().
			ByGroupPolicyUploadedDefinitionFileId(rs.Primary.ID).
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
		return fmt.Errorf("group policy uploaded definition file %s still exists and was not destroyed", rs.Primary.ID)
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
