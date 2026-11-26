package graphBetaApplicationsTokenLifetimePolicy_test

import (
	"net/http"
	"path"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitResourceTokenLifetimePolicy_TransientRefresh404(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	reads := 0
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		reads++
		// Create's read succeeds; the subsequent refresh reaches a stale replica.
		if reads == 2 {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Replica has not observed this policy"}}`), nil
		}
		return httpmock.NewJsonResponse(200, map[string]any{
			"id": path.Base(req.URL.Path), "displayName": "unit-test-token-lifetime-policy-min",
			"definition":            []string{`{"TokenLifetimePolicy":{"Version":1,"AccessTokenLifetime":"01:00:00"}}`},
			"isOrganizationDefault": false,
		})
	})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: testConfigMinimal(),
			Check:  resource.TestCheckResourceAttrSet(resourceType+".minimal", "id"),
		}},
	})
}
