package graphBetaDirectoryRoleAssignment

import (
	"context"
	"net/http"
	"testing"
	"testing/synctest"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

const testAssignmentID = "opaque-assignment-1"

func assignmentTestState(t *testing.T, r *DirectoryRoleAssignmentResource) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	data := DirectoryRoleAssignmentResourceModel{
		ID:               types.StringValue(testAssignmentID),
		RoleDefinitionID: types.StringValue("88d8e3e3-8f55-4a1e-953a-9b9898b8876b"),
		PrincipalID:      types.StringValue("11111111-1111-4111-8111-111111111111"),
		DirectoryScopeID: types.StringValue("/"),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}),
		},
	}
	require.False(t, state.Set(context.Background(), &data).HasError())
	return state
}

func TestUnitDirectoryRoleAssignmentReadPreservesInvalidResponses(t *testing.T) {
	for _, tc := range []struct {
		name          string
		status        int
		scope         any
		missing, fail bool
	}{
		{name: "permission denied", status: 403, fail: true},
		{name: "bad request", status: 400, fail: true},
		{name: "missing", status: 404, missing: true},
		{name: "nil model", status: 204, fail: true},
		{name: "missing scope", status: 200, scope: nil, fail: true},
		{name: "administrative unit", status: 200, scope: "/administrativeUnits/11111111-1111-4111-8111-111111111111", fail: true},
		{name: "tenant", status: 200, scope: "/"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := httpmock.NewMockTransport()
			transport.RegisterResponder(
				"GET",
				"https://graph.microsoft.com/beta/roleManagement/directory/roleAssignments/"+testAssignmentID,
				func(_ *http.Request) (*http.Response, error) {
					if tc.status == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					if tc.status != 200 {
						return httpmock.NewJsonResponse(tc.status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "probe response"}})
					}
					return httpmock.NewJsonResponse(
						200,
						map[string]any{
							"id":               testAssignmentID,
							"principalId":      "11111111-1111-4111-8111-111111111111",
							"roleDefinitionId": "88d8e3e3-8f55-4a1e-953a-9b9898b8876b",
							"directoryScopeId": tc.scope,
						},
					)
				},
			)
			r := NewDirectoryRoleAssignmentResource().(*DirectoryRoleAssignmentResource)
			r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
			state := assignmentTestState(t, r)
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.Equal(t, tc.fail, resp.Diagnostics.HasError(), resp.Diagnostics)
			require.Equal(t, tc.missing, resp.State.Raw.IsNull())
		})
	}
}

func TestUnitDirectoryRoleAssignmentDeletionConfirmsAbsence(t *testing.T) {
	for _, status := range []int{204, 403} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				gets := 0
				url := "https://graph.microsoft.com/beta/roleManagement/directory/roleAssignments/" + testAssignmentID
				transport.RegisterResponder("DELETE", url, func(_ *http.Request) (*http.Response, error) {
					if status == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": "Authorization_RequestDenied", "message": "denied"}})
				})
				transport.RegisterResponder("GET", url, func(_ *http.Request) (*http.Response, error) {
					gets++
					return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "missing"}})
				})
				r := NewDirectoryRoleAssignmentResource().(*DirectoryRoleAssignmentResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := assignmentTestState(t, r)
				resp := resource.DeleteResponse{State: state}
				r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
				require.Equal(t, status == 403, resp.Diagnostics.HasError(), resp.Diagnostics)
				if status == 204 {
					require.Equal(t, 3, gets)
					require.True(t, resp.State.Raw.IsNull())
				} else {
					require.Zero(t, gets)
					require.False(t, resp.State.Raw.IsNull())
				}
			})
		})
	}
}
