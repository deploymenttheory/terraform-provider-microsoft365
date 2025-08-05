package graphBetaTermsAndConditions_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccPreCheck verifies necessary test prerequisites
func testAccPreCheck(t *testing.T) {
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

// testAccCheckTermsAndConditionsDestroy verifies that terms and conditions have been destroyed
func testAccCheckTermsAndConditionsDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_terms_and_conditions" {
			continue
		}

		// Attempt to get the terms and conditions by ID
		_, err := graphClient.
			DeviceManagement().
			TermsAndConditions().
			ByTermsAndConditionsId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if terms and conditions %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("terms and conditions %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccTermsAndConditionsResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTermsAndConditionsDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccTermsAndConditionsConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_terms_and_conditions.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "display_name", "Test Acceptance Terms and Conditions"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "title", "Company Terms"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "body_text", "These are the basic terms and conditions."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "acceptance_statement", "I accept these terms"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "version", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_terms_and_conditions.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccTermsAndConditionsConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_terms_and_conditions.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "display_name", "Test Acceptance Terms and Conditions - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "title", "Complete Company Terms and Conditions"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "body_text", "These are the comprehensive terms and conditions that all users must read and accept before accessing company resources."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "acceptance_statement", "I have read and agree to abide by all terms and conditions outlined above"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "version", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "assignments.#", "3"),
				),
			},
		},
	})
}

func TestAccTermsAndConditionsResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTermsAndConditionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTermsAndConditionsConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_terms_and_conditions.description", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.description", "display_name", "Test Description Terms and Conditions"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.description", "description", "This is a test terms and conditions with description"),
				),
			},
		},
	})
}

func TestAccTermsAndConditionsResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTermsAndConditionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTermsAndConditionsConfig_assignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_terms_and_conditions.assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.assignments", "display_name", "Test Assignments Terms and Conditions"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.assignments", "description", "Terms and conditions policy with assignments for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.assignments", "assignments.#", "3"),
				),
			},
		},
	})
}

// Test configuration functions

func testAccTermsAndConditionsConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfigWithProvider(config)
}

func testAccTermsAndConditionsConfig_maximal() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfigWithProvider(dependencies + "\n" + config)
}

func testAccTermsAndConditionsConfig_description() string {
	config := mocks.LoadTerraformConfigFile("resource_description.tf")
	return acceptance.ConfigWithProvider(config)
}

func testAccTermsAndConditionsConfig_assignments() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_assignments.tf")
	return acceptance.ConfigWithProvider(dependencies + "\n" + config)
}