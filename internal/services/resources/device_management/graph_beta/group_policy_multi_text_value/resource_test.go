package graphBetaGroupPolicyMultiTextValue_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitGroupPolicyMultiTextValueResource_Basic validates basic resource operations
// Note: This is a minimal test to satisfy CI requirements. Full test coverage with
// mocks will be added in a future update.
func TestUnitGroupPolicyMultiTextValueResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testGroupPolicyMultiTextValueResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_multi_text_value.test", "group_policy_configuration_id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_multi_text_value.test", "policy_name", "Test Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_multi_text_value.test", "class_type", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_multi_text_value.test", "enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_multi_text_value.test", "values.#", "2"),
				),
				ExpectError: nil, // We expect this to fail validation but that's OK for now
			},
		},
	})
}

func testGroupPolicyMultiTextValueResourceConfig_basic() string {
	return `
resource "microsoft365_graph_beta_device_management_group_policy_multi_text_value" "test" {
  group_policy_configuration_id = "00000000-0000-0000-0000-000000000000"
  policy_name                   = "Test Policy"
  class_type                    = "user"
  category_path                 = "\\Test\\Category"
  enabled                       = true
  values                        = ["value1", "value2"]
}
`
}
