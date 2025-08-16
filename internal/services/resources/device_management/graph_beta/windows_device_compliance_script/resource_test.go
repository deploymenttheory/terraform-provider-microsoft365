package graphBetaWindowsDeviceComplianceScript

import (
	"context"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_device_compliance_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// TestDeviceComplianceScriptResource_Schema tests the schema validation
func TestDeviceComplianceScriptResource_Schema(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: helpers.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadTestData(t, "tests/terraform/unit/resource_schema.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "display_name", "Test Device Compliance Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "description", "Test description for device compliance script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "publisher", "Test Publisher"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "detection_script_content", "Get-Process"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "id"),
				),
			},
		},
	})
}

// TestDeviceComplianceScriptResource_RunAsAccount tests different run as account types
func TestDeviceComplianceScriptResource_RunAsAccount(t *testing.T) {
	t.Parallel()

	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: helpers.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadTestData(t, "tests/terraform/unit/resource_run_as_user.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "display_name", "Test Device Compliance Script - User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "detection_script_content", "Get-ComputerInfo"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "id"),
				),
			},
		},
	})
}

// TestDeviceComplianceScriptResource_ErrorHandling tests error handling
func TestDeviceComplianceScriptResource_ErrorHandling(t *testing.T) {
	t.Parallel()

	setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: helpers.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadTestData(t, "tests/terraform/unit/resource_error.tf"),
				ExpectError: resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_windows_device_compliance_script.test", "id", func(value string) error { return nil }),
			},
		},
	})
}

func loadTestData(t *testing.T, path string) string {
	t.Helper()
	data, err := helpers.LoadTerraformTestData(path)
	if err != nil {
		t.Fatalf("Failed to load test data from %s: %v", path, err)
	}
	return data
}

func setupMockEnvironment() {
	httpmock.Activate()
	mocks.GlobalRegistry.Reset()
	mocks.GlobalRegistry.RegisterResourceMocks("windows_device_compliance_script")
}

func setupErrorMockEnvironment() {
	httpmock.Activate()
	mocks.GlobalRegistry.Reset()
	mocks.GlobalRegistry.RegisterResourceErrorMocks("windows_device_compliance_script")
}

func testAccCheckDeviceComplianceScriptExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return terraform.NewErrorf("Not found: %s", resourceName)
		}
		return nil
	}
}