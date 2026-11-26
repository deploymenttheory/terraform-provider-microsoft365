package graphBetaServicePrincipalTokenSigningCertificate_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitSigningCertificatePrincipalReplication(t *testing.T) {
	for _, transient := range []bool{true, false} {
		name := "permanent_validation_error"
		if transient {
			name = "principal_not_yet_replicated"
		}
		t.Run(name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, state := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer state.CleanupMockState()
			firstRequest := false
			httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+/addTokenSigningCertificate$`, func(req *http.Request) (*http.Response, error) {
				firstRequest = true
				state.RegisterMocks()
				status := http.StatusBadRequest
				body := `{"error":{"code":"Request_BadRequest","message":"A different validation failure must not be retried."}}`
				if transient {
					status = http.StatusNotFound
					body = `{"error":{"code":"Request_ResourceNotFound","message":"Resource '11111111-1111-1111-1111-111111111111' does not exist or one of its queried reference-property objects are not present."}}`
				}
				response := httpmock.NewStringResponse(status, body)
				response.Header.Set("Content-Type", "application/json")
				return response, nil
			})
			step := resource.TestStep{Config: loadUnitTestTerraform("resource_01.tf")}
			if transient {
				step.Check = resource.TestCheckResourceAttr(resourceType+".test", "service_principal_id", testServicePrincipalID)
			} else {
				step.ExpectError = regexp.MustCompile("A different validation failure must not be retried")
			}
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{step}})
			if !firstRequest {
				t.Fatal("test did not exercise the failing POST")
			}
		})
	}
}
