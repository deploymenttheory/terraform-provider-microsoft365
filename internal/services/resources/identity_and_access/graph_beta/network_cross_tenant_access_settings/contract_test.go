package graphBetaNetworkCrossTenantAccessSettings

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func testModel() NetworkCrossTenantAccessSettingsResourceModel {
	return NetworkCrossTenantAccessSettingsResourceModel{
		ID: types.StringValue(singletonID), NetworkPacketTaggingStatus: types.StringValue("disabled"),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})},
	}
}

func testState(t *testing.T, r *NetworkCrossTenantAccessSettingsResource, m NetworkCrossTenantAccessSettingsResourceModel) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), m).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *NetworkCrossTenantAccessSettingsResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewNetworkCrossTenantAccessSettingsResource().(*NetworkCrossTenantAccessSettingsResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}

func reply(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if body != "" {
		_, _ = io.WriteString(w, body)
	}
}

func TestUnitResourceNetworkCrossTenantAccessSettings_10_ReadErrorsPreserveState(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404, 500} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				reply(w, code, `{"error":{"code":"TestError","message":"Synthetic read failure"}}`)
			})
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkCrossTenantAccessSettings_11_MalformedReadPreservesState(t *testing.T) {
	for _, body := range []string{`{}`, `{"networkPacketTaggingStatus":null}`, `{"networkPacketTaggingStatus":"futureValue"}`, `{"networkPacketTaggingStatus":"unknownFutureValue"}`} {
		t.Run(body, func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) { reply(w, 200, body) })
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkCrossTenantAccessSettings_12_CreateReadbackFailureRetainsID(t *testing.T) {
	patches := 0
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, "/networkAccess/settings/crossTenantAccess", q.URL.Path)
		if q.Method == http.MethodPatch {
			patches++
			reply(w, 204, "")
			return
		}
		reply(w, 403, `{"error":{"code":"Forbidden","message":"Read failed"}}`)
	})
	m := testModel()
	m.ID = types.StringUnknown()
	state := testState(t, r, m)
	resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	var id string
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, singletonID, id)
	require.Equal(t, 1, patches)
}

func TestUnitResourceNetworkCrossTenantAccessSettings_13_PatchContract(t *testing.T) {
	for _, code := range []int{200, 204} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			var patch map[string]any
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				require.Equal(t, "/networkAccess/settings/crossTenantAccess", q.URL.Path)
				if q.Method == http.MethodPatch {
					var reader io.Reader = q.Body
					if q.Header.Get("Content-Encoding") == "gzip" {
						gz, err := gzip.NewReader(q.Body)
						require.NoError(t, err)
						defer gz.Close()
						reader = gz
					}
					require.NoError(t, json.NewDecoder(reader).Decode(&patch))
					if code == 204 {
						reply(w, 204, "")
					} else {
						reply(w, 200, `{"networkPacketTaggingStatus":"enabled"}`)
					}
					return
				}
				require.Equal(t, http.MethodGet, q.Method)
				reply(w, 200, `{"networkPacketTaggingStatus":"enabled","dataPlaneTaggingOptions":"none"}`)
			})
			old := testState(t, r, testModel())
			m := testModel()
			m.NetworkPacketTaggingStatus = types.StringValue("enabled")
			plan := testState(t, r, m)
			resp := resource.UpdateResponse{State: old}
			r.Update(context.Background(), resource.UpdateRequest{State: old, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			delete(patch, "@odata.type")
			require.Equal(t, map[string]any{"networkPacketTaggingStatus": "enabled"}, patch)
		})
	}
}

func TestUnitResourceNetworkCrossTenantAccessSettings_14_DeleteMakesNoRequest(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) { t.Error("destroy must not send HTTP") })
	state := testState(t, r, testModel())
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	require.True(t, resp.State.Raw.IsNull())
}

func TestUnitResourceNetworkCrossTenantAccessSettings_15_ImportID(t *testing.T) {
	r := NewNetworkCrossTenantAccessSettingsResource().(*NetworkCrossTenantAccessSettingsResource)
	for _, id := range []string{singletonID, "", "another-tenant", "random-id"} {
		resp := resource.ImportStateResponse{State: testState(t, r, testModel())}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		require.Equal(t, id != singletonID, resp.Diagnostics.HasError())
	}
}

func TestUnitResourceNetworkCrossTenantAccessSettings_16_Consistency(t *testing.T) {
	r := NewNetworkCrossTenantAccessSettingsResource().(*NetworkCrossTenantAccessSettingsResource)
	expected := testModel()
	predicate := consistencyPredicate(&expected)
	require.True(t, predicate(context.Background(), testState(t, r, expected)))
	stale := testModel()
	stale.NetworkPacketTaggingStatus = types.StringValue("enabled")
	require.False(t, predicate(context.Background(), testState(t, r, stale)))
}

func TestUnitResourceNetworkCrossTenantAccessSettings_17_TimeoutOnlyUpdateMakesNoRequest(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) { t.Error("timeout-only update must not send HTTP") })
	old := testState(t, r, testModel())
	m := testModel()
	m.Timeouts.Object = types.ObjectValueMust(m.Timeouts.AttributeTypes(context.Background()), map[string]attr.Value{"create": types.StringNull(), "read": types.StringValue("4m"), "update": types.StringNull(), "delete": types.StringNull()})
	plan := testState(t, r, m)
	resp := resource.UpdateResponse{State: old}
	r.Update(context.Background(), resource.UpdateRequest{State: old, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	require.True(t, plan.Raw.Equal(resp.State.Raw))
}

func TestUnitResourceNetworkCrossTenantAccessSettings_18_ObservedReadResponse(t *testing.T) {
	fixture, err := os.ReadFile("tests/responses/validate_get/get_cross_tenant_access_settings_success.json")
	require.NoError(t, err)
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) { reply(w, 200, string(fixture)) })
	state := testState(t, r, testModel())
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	var model NetworkCrossTenantAccessSettingsResourceModel
	require.False(t, resp.State.Get(context.Background(), &model).HasError())
	require.Equal(t, "crossTenantAccess", model.ID.ValueString())
	require.Equal(t, "disabled", model.NetworkPacketTaggingStatus.ValueString())
}
