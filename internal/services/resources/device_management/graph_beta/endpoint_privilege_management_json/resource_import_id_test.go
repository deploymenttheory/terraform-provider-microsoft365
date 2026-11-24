package graphBetaEndpointPrivilegeManagement_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaEndpointPrivilegeManagement "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/endpoint_privilege_management_json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// TestUnitResourceEndpointPrivilegeManagementJson_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceEndpointPrivilegeManagementJson_ImportRejectsInvalidID(t *testing.T) {
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
						Config:        `resource "` + graphBetaEndpointPrivilegeManagement.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaEndpointPrivilegeManagement.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}
