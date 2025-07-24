package graphBetaMacosCustomAttributeScriptAssignment_test

import (
	"net/http"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const unitTestProviderConfig = `
provider "microsoft365" {
  tenant_id = "00000000-0000-0000-0000-000000000001"
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "11111111-1111-1111-1111-111111111111"
    client_secret = "mock-secret-value"
  }
  cloud = "public"
}
`

func TestUnitMacosCustomAttributeScriptAssignmentResourceModel_Basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Mock POST for creating assignment
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/00000000-0000-0000-0000-000000000002/assignments",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(201, map[string]interface{}{
				"id": "00000000-0000-0000-0000-000000000002_00000000-0000-0000-0000-000000000003",
				"target": map[string]interface{}{
					"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
				},
			})
		})

	// Mock GET for reading assignment
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/00000000-0000-0000-0000-000000000002/assignments/00000000-0000-0000-0000-000000000002_00000000-0000-0000-0000-000000000003",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"id": "00000000-0000-0000-0000-000000000002_00000000-0000-0000-0000-000000000003",
				"target": map[string]interface{}{
					"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
				},
			})
		})

	// Mock DELETE for removing assignment
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/00000000-0000-0000-0000-000000000002/assignments/00000000-0000-0000-0000-000000000002_00000000-0000-0000-0000-000000000003",
		httpmock.NewBytesResponder(204, nil))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: unitTestProviderConfig + `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script_assignment" "test" {
  macos_custom_attribute_script_id = "00000000-0000-0000-0000-000000000002"
  target = {
    target_type = "allDevices"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script_assignment.test", "macos_custom_attribute_script_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script_assignment.test", "target.target_type", "allDevices"),
				),
			},
		},
	})
}
