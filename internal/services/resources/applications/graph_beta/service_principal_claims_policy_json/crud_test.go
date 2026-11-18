package graphBetaServicePrincipalClaimsPolicyJSON

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

const (
	testParent    = "11111111-1111-4111-8111-111111111111"
	testPolicyURL = "https://graph.microsoft.com/beta/servicePrincipals/" + testParent + "/claimsPolicy"
)

func policyTestState(t *testing.T, r *ServicePrincipalClaimsPolicyJSONResource) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	data := ServicePrincipalClaimsPolicyJSONResourceModel{
		ID:                 types.StringValue(testParent),
		ServicePrincipalID: types.StringValue(testParent),
		PolicyID:           types.StringNull(),
		PolicyJSON:         NewPolicyJSONValue(emptyPolicy),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}),
		},
	}
	require.False(t, state.Set(context.Background(), &data).HasError())
	return state
}

func TestUnitClaimsPolicyReadPreservesUncertainState(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		statuses                 []int
		postWrite, missing, fail bool
	}{
		{name: "bad request", statuses: []int{400}, fail: true}, {name: "forbidden", statuses: []int{403}, fail: true}, {name: "confirmed missing", statuses: []int{404}, missing: true}, {name: "replica 404 recovers", statuses: []int{404, 200}}, {name: "post-write 404 retains ID", statuses: []int{404}, postWrite: true, fail: true}, {name: "empty response", statuses: []int{204}, fail: true}, {name: "invalid policy", statuses: []int{201}, fail: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				calls := 0
				transport.RegisterResponder("GET", testPolicyURL, func(_ *http.Request) (*http.Response, error) {
					status := tc.statuses[min(calls, len(tc.statuses)-1)]
					calls++
					if status == 200 {
						return httpmock.NewStringResponse(200, emptyPolicy), nil
					}
					if status == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					if status == 201 {
						return httpmock.NewStringResponse(201, `{}`), nil
					}
					return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "probe error"}})
				})
				r := NewServicePrincipalClaimsPolicyJSONResource().(*ServicePrincipalClaimsPolicyJSONResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := policyTestState(t, r)
				resp := resource.ReadResponse{State: state}
				ctx := context.Background()
				if tc.postWrite {
					ctx = context.WithValue(ctx, "retry_operation", "Create")
				}
				r.Read(ctx, resource.ReadRequest{State: state}, &resp)
				require.Equal(t, tc.fail, resp.Diagnostics.HasError(), resp.Diagnostics)
				require.Equal(t, tc.missing, resp.State.Raw.IsNull())
				if tc.missing {
					require.Equal(t, 3, calls)
				}
			})
		})
	}
}

func TestUnitClaimsPolicyUnknownValuesAndStateOnlyDestroy(t *testing.T) {
	_, err := constructResource(PolicyJSONValue{StringValue: basetypes.NewStringUnknown()})
	require.ErrorIs(t, err, errPolicyKnown)
	r := NewServicePrincipalClaimsPolicyJSONResource().(*ServicePrincipalClaimsPolicyJSONResource)
	var imported resource.ImportStateResponse
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "not-a-parent-uuid"}, &imported)
	require.True(t, imported.Diagnostics.HasError())
	state := policyTestState(t, r)
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	require.True(t, resp.State.Raw.IsNull(), "destroy must not need a Graph client or mutate policy settings")
}
