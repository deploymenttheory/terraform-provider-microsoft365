package graphBetaAppControlForBusinessBuiltInControls_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAppControlForBusinessBuiltInControls "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_built_in_controls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// TestUnitResourceAppControlForBusinessBuiltInControls_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceAppControlForBusinessBuiltInControls_ImportRejectsInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()

	tests := []struct {
		name     string
		importID string
	}{
		{name: "query string", importID: "00000000-0000-0000-0000-000000000000?$select=id"},
		{name: "path traversal", importID: "../00000000-0000-0000-0000-000000000000"},
		{name: "not a GUID", importID: "not-a-guid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:        `resource "` + graphBetaAppControlForBusinessBuiltInControls.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaAppControlForBusinessBuiltInControls.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}
