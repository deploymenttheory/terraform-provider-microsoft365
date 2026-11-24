package graphBetaSettingsCatalogInventoryPolicy_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	terraformResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// loadUnitTestTerraform reads a unit test configuration from tests/terraform/unit.
func loadUnitTestTerraform(t *testing.T, filename string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		t.Fatalf("load %s: %v", filename, err)
	}
	return config
}

// TestUnitResourceSettingsCatalogInventoryPolicy_06_EmptyChildrenRequests creates and updates an inventory policy whose
// choice, choice collection and group collection values have no children, and checks the POST and PUT bodies against a
// request fixture. Each such value must be sent as "children": [], because Graph rejects "children": null (#4956).
func TestUnitResourceSettingsCatalogInventoryPolicy_06_EmptyChildrenRequests(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, inventoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer inventoryMock.CleanupMockState()
	capture := mocks.CaptureGraphRequests(t)

	resourceType := "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.empty_children"
	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: loadUnitTestTerraform(t, "resource_empty_children.tf"),
				Check:  terraformResource.TestCheckResourceAttr(resourceType, "configuration_policy.settings.#", "2"),
			},
			{
				Config: loadUnitTestTerraform(t, "resource_empty_children_updated.tf"),
				Check:  terraformResource.TestCheckResourceAttr(resourceType, "description", "Test inventory policy for settings without children - updated"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/empty_children.json")
}

// TestUnitResourceSettingsCatalogInventoryPolicy_07_MaximalRequests pins the full create and assign requests for an
// existing configuration, so changes to the shared settings constructor show up in the fixture.
func TestUnitResourceSettingsCatalogInventoryPolicy_07_MaximalRequests(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, inventoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer inventoryMock.CleanupMockState()
	capture := mocks.CaptureGraphRequests(t)

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: loadUnitTestTerraform(t, "resource_maximal.tf"),
				Check:  terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "name", "Test Maximal Inventory Policy - Unit"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/maximal.json")
}
