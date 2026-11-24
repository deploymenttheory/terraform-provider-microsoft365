package graphBetaDeviceConfigurationTemplatesJson_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// TestUnitResourceDeviceConfigurationTemplatesJson_ImportRejectsInvalidID checks that ImportState rejects an ID that
// isn't a GUID, before the ID can reach a Graph request URL.
func TestUnitResourceDeviceConfigurationTemplatesJson_ImportRejectsInvalidID(t *testing.T) {
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
						Config:        `resource "` + resourceType + `" "test" {}`,
						ResourceName:  resourceType + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}
