package graphBetaGroupPolicyDefinition_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	configuration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_configuration"
	definition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_definition"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Uses a built-in ListBox catalog definition and a freshly owned, unassigned parent.
// No existing policy or imported ADMX template is modified (issue #4905).
func TestAccResourceGroupPolicyDefinition_09_ListBox_ImportDriftUpdate(t *testing.T) {
	const address = resourceType + ".test_009"
	const initial = `[{"name":"https://example.test","value":"2"},{"name":"https://other.test","value":"1"}]`
	const updated = `[{"name":"https://updated.test","value":"3"}]`
	var configID, valueID, presentationID string
	capture := func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[address]
		if !ok {
			return fmt.Errorf("missing owned ListBox resource")
		}
		parts := strings.Split(r.Primary.ID, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid composite ID %q", r.Primary.ID)
		}
		configID, valueID = parts[0], parts[1]
		presentationID = r.Primary.Attributes["values.0.id"]
		if presentationID == "" {
			return fmt.Errorf("missing ListBox template ID")
		}
		return nil
	}
	checkValue := func(value string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(address, "values.#", "1"), resource.TestCheckResourceAttr(address, "values.0.value", value))
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		CheckDestroy: resource.ComposeTestCheckFunc(
			destroy.CheckDestroyedAllFunc(testResource, definition.ResourceName, 30*time.Second),
			destroy.CheckDestroyedAllFunc(configuration.GroupPolicyConfigurationTestResource{}, configuration.ResourceName, 30*time.Second),
		),
		Steps: []resource.TestStep{
			{Config: loadAcceptanceTestTerraform("009_listbox.tf"), Check: resource.ComposeTestCheckFunc(capture, checkValue(initial))},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: loadAcceptanceTestTerraform("009_listbox.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					if configID == "" || valueID == "" {
						t.Fatal("refusing to mutate without owned fixture IDs")
					}
					client, err := acceptance.TestGraphClient()
					if err != nil {
						t.Fatal(err)
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
					defer cancel()
					values, err := client.DeviceManagement().GroupPolicyConfigurations().ByGroupPolicyConfigurationId(configID).DefinitionValues().ByGroupPolicyDefinitionValueId(valueID).PresentationValues().Get(ctx, nil)
					if err != nil {
						t.Fatal(err)
					}
					if values == nil || len(values.GetValue()) != 1 {
						t.Fatal("expected exactly one owned presentation value")
					}
					pair := models.NewKeyValuePair()
					name, value := "https://drift.test", "4"
					pair.SetName(&name)
					pair.SetValue(&value)
					list := models.NewGroupPolicyPresentationValueList()
					typ := "#microsoft.graph.groupPolicyPresentationValueList"
					list.SetOdataType(&typ)
					list.SetId(values.GetValue()[0].GetId())
					list.SetValues([]models.KeyValuePairable{pair})
					list.SetAdditionalData(map[string]any{"presentation@odata.bind": fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('eee38118-0408-432d-b81c-632cc4f62f1a')/presentations('%s')", presentationID)})
					dv := models.NewGroupPolicyDefinitionValue()
					dv.SetId(&valueID)
					dv.SetAdditionalData(map[string]any{"definition@odata.bind": "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('eee38118-0408-432d-b81c-632cc4f62f1a')"})
					enabled := true
					dv.SetEnabled(&enabled)
					dv.SetPresentationValues([]models.GroupPolicyPresentationValueable{list})
					body := devicemanagement.NewGroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody()
					body.SetAdded([]models.GroupPolicyDefinitionValueable{})
					body.SetUpdated([]models.GroupPolicyDefinitionValueable{dv})
					body.SetDeletedIds([]string{})
					if err := client.DeviceManagement().GroupPolicyConfigurations().ByGroupPolicyConfigurationId(configID).UpdateDefinitionValues().Post(ctx, body, nil); err != nil {
						t.Fatal(err)
					}
					t.Log("Changed only the newly owned ListBox fixture to exercise remote drift")
				},
				Config: loadAcceptanceTestTerraform("009_listbox.tf"), PlanOnly: true, ExpectNonEmptyPlan: true,
			},
			{Config: loadAcceptanceTestTerraform("009_listbox.tf"), Check: checkValue(initial)},
			{Config: loadAcceptanceTestTerraform("009_listbox_updated.tf"), Check: checkValue(updated)},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: loadAcceptanceTestTerraform("009_listbox_updated.tf"), PlanOnly: true},
		},
	})
}
