package graphBetaApplication_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func setupApplicationPlanningMock(t *testing.T, flags *bool, readOverride ...func(map[string]any) map[string]any) {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mocks.NewMocks().AuthMocks.RegisterMocks()
	fixture, err := helpers.ParseJSONFile("tests/responses/get_empty_application.json")
	require.NoError(t, err)
	var body map[string]any
	require.NoError(t, json.Unmarshal([]byte(fixture), &body))
	body["displayName"] = "application-role-origin"
	if flags != nil {
		body["isDeviceOnlyAuthSupported"] = *flags
		body["isFallbackPublicClient"] = *flags
	}
	deleted := false
	purged := false
	const deletedEndpoint = "https://graph.microsoft.com/beta/directory/deletedItems/11111111-1111-4111-8111-111111111111"
	httpmock.RegisterResponder("GET", deletedEndpoint, func(_ *http.Request) (*http.Response, error) {
		if deleted && !purged {
			return httpmock.NewJsonResponse(200, body)
		}
		return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "No deleted application"}})
	})
	httpmock.RegisterResponder("DELETE", deletedEndpoint, func(_ *http.Request) (*http.Response, error) {
		purged = true
		return httpmock.NewStringResponse(204, ""), nil
	})
	const endpoint = "https://graph.microsoft.com/beta/applications/11111111-1111-4111-8111-111111111111"
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/applications", func(_ *http.Request) (*http.Response, error) { return httpmock.NewJsonResponse(201, body) })
	httpmock.RegisterResponder("GET", endpoint, func(_ *http.Request) (*http.Response, error) {
		if deleted {
			return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Application deleted"}})
		}
		if len(readOverride) > 0 {
			return httpmock.NewJsonResponse(200, readOverride[0](body))
		}
		return httpmock.NewJsonResponse(200, body)
	})
	httpmock.RegisterResponder("GET", endpoint+"/owners", httpmock.NewJsonResponderOrPanic(200, map[string]any{"value": []any{}}))
	httpmock.RegisterResponder("PATCH", endpoint, func(req *http.Request) (*http.Response, error) {
		var changes map[string]any
		if err := json.NewDecoder(req.Body).Decode(&changes); err != nil {
			return nil, err
		}
		for key, value := range changes {
			if key == "appRoles" {
				for _, role := range value.([]any) {
					// The live API always supplies this read-only value.
					role.(map[string]any)["origin"] = "Application"
				}
			}
			body[key] = value
		}
		return httpmock.NewStringResponse(204, ""), nil
	})
	httpmock.RegisterResponder("DELETE", endpoint, func(_ *http.Request) (*http.Response, error) {
		deleted = true
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func TestUnitApplicationNewRoleOriginPlanning(t *testing.T) {
	setupApplicationPlanningMock(t, nil)
	const minimal = `resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "application-role-origin"
}`
	const withRole = `resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "application-role-origin"
  app_roles = [{
    id = "33333333-3333-4333-8333-333333333333"
    allowed_member_types = ["User"]
    description = "User role"
    display_name = "User role"
    is_enabled = true
    value = "User.Role"
  }]
}`
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: minimal, ResourceName: address, ImportState: true, ImportStateId: "11111111-1111-4111-8111-111111111111", ImportStatePersist: true},
			{
				Config: withRole,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectUnknownValue(address, tfjsonpath.New("app_roles").AtSliceIndex(0).AtMapKey("origin")),
				}},
				Check: resource.TestCheckTypeSetElemNestedAttrs(address, "app_roles.*", map[string]string{"origin": "Application"}),
			},
			{Config: withRole, PlanOnly: true},
		},
	})
}

func TestUnitApplicationAuthFlagsRemovalPlanning(t *testing.T) {
	setupApplicationPlanningMock(t, new(false))
	const configured = `resource "microsoft365_graph_beta_applications_application" "test" {
      display_name = "application-role-origin"
      is_device_only_auth_supported = false
      is_fallback_public_client = false
    }`
	const omitted = `resource "microsoft365_graph_beta_applications_application" "test" {
      display_name = "application-role-origin-updated"
    }`
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: configured},
			{Config: omitted, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate)}}},
			{Config: omitted, PlanOnly: true},
		},
	})
}

func TestUnitApplicationMetadataUpdateWaitsForStableReads(t *testing.T) {
	reads := 0
	setupApplicationPlanningMock(t, nil, func(body map[string]any) map[string]any {
		if body["description"] != "Updated metadata" {
			return body
		}
		reads++
		if reads >= 2 && reads <= 4 {
			stale := make(map[string]any, len(body))
			for key, value := range body {
				stale[key] = value
			}
			stale["description"] = "Empty collection import probe"
			return stale
		}
		return body
	})
	const before = `resource "microsoft365_graph_beta_applications_application" "test" {
      display_name = "application-role-origin"
    }`
	const after = `resource "microsoft365_graph_beta_applications_application" "test" {
      display_name = "application-role-origin"
      description = "Updated metadata"
    }`
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: before, ResourceName: address, ImportState: true, ImportStateId: "11111111-1111-4111-8111-111111111111", ImportStatePersist: true},
			{Config: after, Check: resource.TestCheckResourceAttr(address, "description", "Updated metadata")},
			{Config: after, PlanOnly: true},
		},
	})
	require.GreaterOrEqual(t, reads, 7)
}
