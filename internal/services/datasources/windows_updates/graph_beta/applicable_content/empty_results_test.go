package graphBetaWindowsUpdatesApplicableContent_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitDatasourceApplicableContent_EmptyResults(t *testing.T) {
	for _, filtered := range []bool{false, true} {
		name := "empty_api_response"
		filter := ""
		if filtered {
			name = "all_entries_filtered_out"
			filter = `manufacturer = "nonexistent-test-manufacturer"`
		}
		t.Run(name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, mock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer mock.CleanupMockState()
			if !filtered {
				httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/[^/]+/applicableContent`, httpmock.NewStringResponder(200, `{"value":[]}`))
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: `data "microsoft365_graph_beta_windows_updates_applicable_content" "test" {
  audience_id = "f660d844-30b7-46e4-a6cf-47e36164d3cb"
  ` + filter + `
}
output "count" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.test.applicable_content)
}`,
					Check: resource.TestCheckOutput("count", "0"),
				}},
			})
		})
	}
}
