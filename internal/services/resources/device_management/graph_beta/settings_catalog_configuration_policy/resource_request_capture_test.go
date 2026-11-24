package graphBetaSettingsCatalogConfigurationPolicy_test

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

// TestUnitResourceSettingsCatalogConfigurationPolicy_04_EmptyChildrenRequests creates and updates a policy whose
// choice, choice collection and group collection values have no children, and checks the POST and PUT bodies against
// a request fixture. Each such value must be sent as "children": [], because Graph rejects "children": null (#4956).
func TestUnitResourceSettingsCatalogConfigurationPolicy_04_EmptyChildrenRequests(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()
	capture := mocks.CaptureGraphRequests(t)

	resourceType := "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.empty_children"
	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: loadUnitTestTerraform(t, "resource_empty_children.tf"),
				Check:  terraformResource.TestCheckResourceAttr(resourceType, "configuration_policy.settings.#", "4"),
			},
			{
				Config: loadUnitTestTerraform(t, "resource_empty_children_updated.tf"),
				Check:  terraformResource.TestCheckResourceAttr(resourceType, "description", "Test policy for settings without children - updated"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/empty_children.json")
}

// TestUnitResourceSettingsCatalogConfigurationPolicy_05_MaximalRequests pins the full create request for an existing
// configuration with nested group, choice and simple settings, so changes to the constructor show up in the fixture.
func TestUnitResourceSettingsCatalogConfigurationPolicy_05_MaximalRequests(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()
	capture := mocks.CaptureGraphRequests(t)

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: loadUnitTestTerraform(t, "resource_maximal.tf"),
				Check:  terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "name", "macos mdm filevault2 settings"),
			},
		},
	})
	capture.AssertMatchesFixture(t, "tests/requests/maximal.json")
}
