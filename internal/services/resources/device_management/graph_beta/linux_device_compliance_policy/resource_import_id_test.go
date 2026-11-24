package graphBetaLinuxDeviceCompliancePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaLinuxDeviceCompliancePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/linux_device_compliance_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// TestUnitResourceLinuxDeviceCompliancePolicy_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceLinuxDeviceCompliancePolicy_ImportRejectsInvalidID(t *testing.T) {
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
						Config:        `resource "` + graphBetaLinuxDeviceCompliancePolicy.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}
