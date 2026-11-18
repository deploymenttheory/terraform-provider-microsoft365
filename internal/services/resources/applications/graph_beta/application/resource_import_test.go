package graphBetaApplication_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestUnitApplicationGeneratedImportEmptyCollections(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	fixture, err := helpers.ParseJSONFile("tests/responses/get_empty_application.json")
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(fixture), &body); err != nil {
		t.Fatal(err)
	}
	responder, err := httpmock.NewJsonResponder(http.StatusOK, body)
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/applications/11111111-1111-4111-8111-111111111111", responder)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/applications/11111111-1111-4111-8111-111111111111/owners", httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]any{"value": []any{}}))
	config, err := helpers.ParseHCLFile("tests/terraform/unit/import_empty_collections.tf")
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: config, ResourceName: resourceType + ".test",
			ImportStateId: "11111111-1111-4111-8111-111111111111",
			ImportState:   true, GenerateConfig: true, ImportStateKind: resource.ImportBlockWithID,
			ImportStateConfigExact: true,
			ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop),
			}},
		}},
	})
}
