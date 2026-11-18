package graphBetaOAuth2PermissionGrant

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

const testGrantID = "OpaqueGrant-ID_001"

func grantTestModel() OAuth2PermissionGrantResourceModel {
	return OAuth2PermissionGrantResourceModel{
		ID:          types.StringValue(testGrantID),
		ClientID:    types.StringValue("11111111-1111-4111-8111-111111111111"),
		ResourceID:  types.StringValue("22222222-2222-4222-8222-222222222222"),
		ConsentType: types.StringValue("AllPrincipals"),
		PrincipalID: types.StringNull(),
		Scopes:      types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Probe.Write"), types.StringValue("Probe.Read")}),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}),
		},
	}
}

func grantTestState(t *testing.T, r *OAuth2PermissionGrantResource) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	model := grantTestModel()
	require.False(t, state.Set(context.Background(), &model).HasError())
	return state
}

func grantResponse() map[string]any {
	return map[string]any{
		"id":          testGrantID,
		"clientId":    "11111111-1111-4111-8111-111111111111",
		"resourceId":  "22222222-2222-4222-8222-222222222222",
		"consentType": "AllPrincipals",
		"principalId": nil,
		"scope":       "Probe.Read Probe.Write",
	}
}

func grantErrorReply(status int, message string) (*http.Response, error) {
	return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": message}})
}

func TestUnitOAuth2PermissionGrantRequestPayload(t *testing.T) {
	for _, create := range []bool{true, false} {
		for _, principal := range []bool{false, true} {
			t.Run(map[bool]string{true: "create", false: "update"}[create]+"/"+map[bool]string{true: "principal", false: "all"}[principal], func(t *testing.T) {
				ctx := context.Background()
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				data := grantTestModel()
				if principal {
					data.ConsentType = types.StringValue("Principal")
					data.PrincipalID = types.StringValue("33333333-3333-4333-8333-333333333333")
				}
				body, err := constructResource(ctx, &data, create)
				require.NoError(t, err)
				request, err := client.Oauth2PermissionGrants().ToPostRequestInformation(ctx, body, nil)
				if !create {
					request, err = client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(testGrantID).ToPatchRequestInformation(ctx, body, nil)
				}
				require.NoError(t, err)
				var payload map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(request.Content, &payload))
				require.JSONEq(t, `"Probe.Read Probe.Write"`, string(payload["scope"]))
				if create {
					require.Contains(t, payload, "clientId")
					require.Contains(t, payload, "resourceId")
					require.Contains(t, payload, "consentType")
					require.NotEqual(t, "null", string(payload["startTime"]))
					require.NotEqual(t, "null", string(payload["expiryTime"]))
					if principal {
						require.Contains(t, payload, "principalId")
					} else {
						require.NotContains(t, payload, "principalId")
					}
				} else {
					require.Len(t, payload, 1, "PATCH must serialize scope only, including after import")
				}
			})
		}
	}
	t.Run("unknown scopes fail construction", func(t *testing.T) {
		data := grantTestModel()
		data.Scopes = types.SetUnknown(types.StringType)
		_, err := constructResource(context.Background(), &data, true)
		require.Error(t, err)
	})
}

func TestUnitOAuth2PermissionGrantStateScopeNormalization(t *testing.T) {
	remote := graphmodels.NewOAuth2PermissionGrant()
	remote.SetId(new(testGrantID))
	remote.SetScope(new(" Probe.Write\tProbe.Read  Probe.Read\n"))
	var data OAuth2PermissionGrantResourceModel
	MapRemoteStateToTerraform(context.Background(), &data, remote)
	require.Equal(t, testGrantID, data.ID.ValueString())
	require.Len(t, data.Scopes.Elements(), 2, "scope strings must have Terraform set semantics")
	remote.SetScope(new(""))
	MapRemoteStateToTerraform(context.Background(), &data, remote)
	require.False(t, data.Scopes.IsNull())
	require.Empty(t, data.Scopes.Elements())
}

func TestUnitOAuth2PermissionGrantConsistencyRequiresStableReads(t *testing.T) {
	ctx := context.Background()
	r := NewOAuth2PermissionGrantResource().(*OAuth2PermissionGrantResource)
	state := grantTestState(t, r)
	expected := grantTestModel()
	predicate := grantConsistencyPredicate(&expected)
	for range 4 {
		require.False(t, predicate(ctx, state))
	}
	require.True(t, predicate(ctx, state))
	stale := expected
	stale.Scopes = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Probe.Read")})
	require.False(t, state.Set(ctx, &stale).HasError())
	require.False(t, predicate(ctx, state), "a stale replica must reset the matching-read count")
	require.False(t, state.Set(ctx, &expected).HasError())
	for range 4 {
		require.False(t, predicate(ctx, state))
	}
	require.True(t, predicate(ctx, state))
}

func TestUnitOAuth2PermissionGrantReadSemantics(t *testing.T) {
	for _, tc := range []struct {
		name          string
		statuses      []int
		missing, fail bool
		postWrite     bool
		principalLag  bool
		timeout       time.Duration
	}{
		{name: "transient 404", statuses: []int{404, 200}},
		{name: "principal replication error recovers", statuses: []int{400, 200}, principalLag: true},
		{name: "confirmed missing", statuses: []int{404}, missing: true},
		{name: "bad request preserves state", statuses: []int{400}, fail: true},
		{name: "permission failure preserves state", statuses: []int{403}, fail: true},
		{name: "empty success preserves state", statuses: []int{204}, fail: true},
		{name: "post-create 404 must retry", statuses: []int{404}, postWrite: true, fail: true},
		{name: "cancellation preserves state", statuses: []int{404}, fail: true, timeout: time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := context.Background()
				if tc.timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tc.timeout)
					defer cancel()
				}
				if tc.postWrite {
					ctx = context.WithValue(ctx, "retry_operation", "Create")
				}
				transport := httpmock.NewMockTransport()
				calls := 0
				transport.RegisterResponder("GET", "https://graph.microsoft.com/beta/oauth2PermissionGrants/"+testGrantID, func(_ *http.Request) (*http.Response, error) {
					status := tc.statuses[min(calls, len(tc.statuses)-1)]
					calls++
					if status == 200 {
						return httpmock.NewJsonResponse(200, grantResponse())
					}
					if status == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					message := "probe response"
					if tc.principalLag {
						message = "Filter is invalid. principalId"
					}
					return grantErrorReply(status, message)
				})
				r := NewOAuth2PermissionGrantResource().(*OAuth2PermissionGrantResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := grantTestState(t, r)
				resp := resource.ReadResponse{State: state}
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

func TestUnitOAuth2PermissionGrantDeleteSemantics(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		deleteStatus, getStatus int
		transient               bool
		fail                    bool
		timeout                 time.Duration
	}{
		{name: "successful delete", deleteStatus: 204, getStatus: 404},
		{name: "already absent", deleteStatus: 404, getStatus: 404},
		{name: "missing write replica does not prove deletion", deleteStatus: 400, getStatus: 404, transient: true},
		{name: "permission failure", deleteStatus: 403, fail: true},
		{name: "bad request", deleteStatus: 400, fail: true},
		{name: "read permission failure", deleteStatus: 204, getStatus: 403, fail: true},
		{name: "nil read model", deleteStatus: 204, getStatus: 204, fail: true},
		{name: "visible grant times out", deleteStatus: 204, getStatus: 200, fail: true, timeout: time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := context.Background()
				if tc.timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tc.timeout)
					defer cancel()
				}
				transport := httpmock.NewMockTransport()
				getCalls := 0
				url := "https://graph.microsoft.com/beta/oauth2PermissionGrants/" + testGrantID
				transport.RegisterResponder("DELETE", url, func(_ *http.Request) (*http.Response, error) {
					if tc.deleteStatus == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					message := "delete error"
					if tc.transient {
						message = "Permission being updated or deleted is not found."
					}
					return grantErrorReply(tc.deleteStatus, message)
				})
				transport.RegisterResponder("GET", url, func(_ *http.Request) (*http.Response, error) {
					getCalls++
					if tc.transient && getCalls == 1 {
						return httpmock.NewJsonResponse(200, grantResponse())
					}
					if tc.getStatus == 200 {
						return httpmock.NewJsonResponse(200, grantResponse())
					}
					if tc.getStatus == 204 {
						return httpmock.NewStringResponse(204, ""), nil
					}
					return grantErrorReply(tc.getStatus, "read error")
				})
				r := NewOAuth2PermissionGrantResource().(*OAuth2PermissionGrantResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := grantTestState(t, r)
				resp := resource.DeleteResponse{State: state}
				r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)
				require.Equal(t, tc.fail, resp.Diagnostics.HasError(), resp.Diagnostics)
				require.Equal(t, !tc.fail, resp.State.Raw.IsNull())
				if !tc.fail {
					require.GreaterOrEqual(t, getCalls, 3)
				}
			})
		})
	}
}
