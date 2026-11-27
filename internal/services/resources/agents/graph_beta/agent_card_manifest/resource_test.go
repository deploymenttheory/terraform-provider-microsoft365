package graphBetaAgentCardManifest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentCardManifest "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_card_manifest"
)

const manifestType = "microsoft365_graph_beta_agents_agent_card_manifest"

func TestUnitResourceAgentCardManifest_01_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	auth := mocks.NewMocks()
	auth.AuthMocks.RegisterMocks()
	var mu sync.Mutex
	remote := map[string]any{
		"id":                   "employee-assistant",
		"displayName":          "existing",
		"description":          "untouched",
		"createdBy":            "application",
		"createdDateTime":      "2026-11-27T08:00:00Z",
		"lastModifiedDateTime": "2026-11-27T08:00:00Z",
	}
	writes, forbidden := 0, 0
	endpoint := "https://graph.microsoft.com/beta/agentRegistry/agentCardManifests/employee-assistant"
	httpmock.RegisterResponder("GET", endpoint, func(_ *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		return httpmock.NewJsonResponse(200, remote)
	})
	httpmock.RegisterResponder("PATCH", endpoint, func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		for key, value := range body {
			remote[key] = value
		}
		writes++
		remote["lastModifiedDateTime"] = fmt.Sprintf("2026-11-27T08:00:%02dZ", writes)
		return httpmock.NewJsonResponse(200, remote)
	})
	for _, method := range []string{"POST", "DELETE"} {
		httpmock.RegisterResponder(method, endpoint, func(_ *http.Request) (*http.Response, error) {
			forbidden++
			return httpmock.NewStringResponse(405, ""), nil
		})
	}
	minimal := manifestConfig("employee-assistant", "managed", "")
	maximal := manifestConfig("employee-assistant", "maximal", maximalManifestProperties)
	cleared := manifestConfig("employee-assistant", "cleared", `
 owner_ids=[]
 default_input_modes=[]
 default_output_modes=[]
 skills=[]
 security_json="[]"
 security_schemes_json="{}"
`)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			mu.Lock()
			defer mu.Unlock()
			if remote["id"] != "employee-assistant" || remote["displayName"] != "cleared" || forbidden != 0 {
				return fmt.Errorf("destroy must leave the remote manifest intact")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{Config: minimal, Check: resource.TestCheckResourceAttr(manifestType+".test", "description", "untouched")},
			{
				Config: maximal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(manifestType+".test", "version", "2.0"),
					resource.TestCheckResourceAttr(manifestType+".test", "capabilities.streaming", "true"),
				),
			},
			{Config: maximal, PlanOnly: true},
			{ResourceName: manifestType + ".test", ImportState: true, ImportStateId: "employee-assistant", ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
			{
				Config: cleared,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(manifestType+".test", "owner_ids.#", "0"),
					resource.TestCheckResourceAttr(manifestType+".test", "skills.#", "0"),
				),
			},
			{Config: cleared, PlanOnly: true},
		},
	})
	require.Equal(t, 3, writes)
	require.Zero(t, forbidden)
}

func manifestConfig(id, name, properties string) string {
	return fmt.Sprintf(`
resource "%s" "test" {
 agent_card_manifest_id=%q
 display_name=%q
 %s
}
`, manifestType, id, name, properties)
}

const maximalManifestProperties = `
 description="Managed manifest"
 icon_url="https://example.test/icon.png"
 protocol_version="1.0"
 version="2.0"
 documentation_url="https://example.test/docs"
 managed_by="5b64cedc-ccc7-4896-8b1c-24a1bf84b101"
 originating_store="Terraform"
 owner_ids=["11111111-1111-1111-1111-111111111111"]
 default_input_modes=["application/json"]
 default_output_modes=["text/plain"]
 supports_authenticated_extended_card=false
 agent_provider={
  organization="Example"
  url="https://example.test"
 }
 capabilities={
  streaming=true
  push_notifications=false
  state_transition_history=true
  extensions=[{
   uri="https://example.test/extension"
   description="Extension"
   required=false
   params_json=jsonencode({enabled=true,retries=3,nested={values=["a","b"]}})
  }]
 }
 skills=[{
  id="lookup"
  display_name="Lookup"
  description="Find information"
  tags=["lookup"]
  examples=["Find a document"]
  input_modes=["text/plain"]
  output_modes=["application/json"]
 }]
 security_schemes_json=jsonencode({apiKey={type="apiKey",name="X-Key",in="header"}})
 security_json=jsonencode([{apiKey=[]}])
`

func TestUnitResourceAgentCardManifest_02_Schema(t *testing.T) {
	r := graphBetaAgentCardManifest.NewAgentCardManifestResource()
	var response frameworkresource.SchemaResponse
	r.Schema(context.Background(), frameworkresource.SchemaRequest{}, &response)
	require.False(t, response.Schema.ValidateImplementation(context.Background()).HasError())
}
