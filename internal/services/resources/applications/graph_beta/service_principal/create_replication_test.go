package graphBetaServicePrincipal_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	servicePrincipalMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitServicePrincipalApplicationReplication(t *testing.T) {
	for _, transient := range []bool{true, false} {
		name := "permanent_validation_error"
		if transient {
			name = "application_not_yet_replicated"
		}
		t.Run(name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			state := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer state.CleanupMockState()
			firstRequest := false
			httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals$`, func(req *http.Request) (*http.Response, error) {
				firstRequest = true
				// Any retry uses the ordinary stateful successful responder.
				servicePrincipalMocks.RegisterServicePrincipalMockResponders()
				if transient {
					response := httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"The appId does not reference a valid application object.","details":[{"code":"NoBackingApplicationObject","message":"The backing application is not visible yet."}]}}`)
					response.Header.Set("Content-Type", "application/json")
					return response, nil
				}
				response := httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"A different validation failure must not be retried."}}`)
				response.Header.Set("Content-Type", "application/json")
				return response, nil
			})
			step := resource.TestStep{Config: loadUnitTestTerraform("resource_01_minimal.tf")}
			if transient {
				step.Check = resource.TestCheckResourceAttrSet(resourceType+".test_minimal", "id")
			} else {
				step.ExpectError = regexp.MustCompile("A different validation failure must not be retried")
			}
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{step}})
			if !firstRequest {
				t.Fatal("test did not exercise the failing create request")
			}
		})
	}
}
