package graphBetaGroup_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestUnitResourceGroup_GeneratedImportRemoteDefaults(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		visibility, role         any
		wantVisibility, wantRole knownvalue.Check
	}{
		{"null fields", nil, nil, knownvalue.Null(), knownvalue.Null()},
		{"explicit defaults", "Private", false, knownvalue.StringExact("Private"), knownvalue.Bool(false)},
		{"public visibility", "Public", nil, knownvalue.StringExact("Public"), knownvalue.Null()},
		{"role assignable", "Private", true, knownvalue.StringExact("Private"), knownvalue.Bool(true)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, groupMock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer groupMock.CleanupMockState()

			fixture, err := helpers.ParseJSONFile("tests/responses/validate_create/post_group_scenario_1.json")
			if err != nil {
				t.Fatal(err)
			}
			var body map[string]any
			if err := json.Unmarshal([]byte(fixture), &body); err != nil {
				t.Fatal(err)
			}
			body["id"] = "11111111-1111-4111-8111-111111111111"
			body["visibility"] = tc.visibility
			body["isAssignableToRole"] = tc.role
			responder, err := httpmock.NewJsonResponder(http.StatusOK, body)
			if err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/11111111-1111-4111-8111-111111111111", responder)
			config, err := helpers.ParseHCLFile("tests/terraform/unit/import_remote_defaults.tf")
			if err != nil {
				t.Fatal(err)
			}

			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:        config,
					ResourceName:  resourceType + ".test",
					ImportStateId: "11111111-1111-4111-8111-111111111111",
					ImportState:   true, GenerateConfig: true,
					ImportStateKind:        resource.ImportBlockWithID,
					ImportStateConfigExact: true,
					ImportPlanChecks: resource.ImportPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop),
							plancheck.ExpectKnownValue(resourceType+".test", tfjsonpath.New("visibility"), tc.wantVisibility),
							plancheck.ExpectKnownValue(resourceType+".test", tfjsonpath.New("is_assignable_to_role"), tc.wantRole),
						},
					},
				}},
			})
		})
	}
}

func TestUnitResourceGroup_UpdateNoContent(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()
	config := `resource "microsoft365_graph_beta_groups_group" "test" {
  display_name = "unit-test-group-no-content"
  mail_nickname = "unit-test-group-no-content"
  mail_enabled = false
  security_enabled = true
  description = "Original group description"
}`
	updated := strings.Replace(config, "Original group description", "Updated group description", 1)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{Config: updated, Check: resource.TestCheckResourceAttr(resourceType+".test", "description", "Updated group description")},
			{Config: updated, PlanOnly: true},
		},
	})
}
