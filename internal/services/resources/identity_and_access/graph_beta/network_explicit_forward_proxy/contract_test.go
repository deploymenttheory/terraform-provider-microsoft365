package graphBetaNetworkExplicitForwardProxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func testModel() NetworkExplicitForwardProxyResourceModel {
	return NetworkExplicitForwardProxyResourceModel{
		ID: types.StringValue("explicitForwardProxyConfig"),
		InternetAccess: &InternetAccessResourceModel{
			IsEnabled:                        types.BoolValue(false),
			IsSourceIPSessionAffinityEnabled: types.BoolValue(true),
			SourceIPSessionAffinityOptions:   types.StringValue("useSessionId"),
		},
		ProxyAutoConfigurationFileURL: types.StringNull(),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(
				map[string]attr.Type{
					"create": types.StringType,
					"read":   types.StringType,
					"update": types.StringType,
					"delete": types.StringType,
				},
			),
		},
	}
}

func testState(
	t *testing.T,
	r *NetworkExplicitForwardProxyResource,
	m NetworkExplicitForwardProxyResourceModel,
) tfsdk.State {
	t.Helper()
	var s resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &s)
	state := tfsdk.State{Schema: s.Schema}
	require.False(t, state.Set(context.Background(), m).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *NetworkExplicitForwardProxyResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	a, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	a.SetBaseUrl(server.URL)
	r := NewNetworkExplicitForwardProxyResource().(*NetworkExplicitForwardProxyResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(a)
	return r
}

func respond(t *testing.T, w http.ResponseWriter, code int, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if body != "" {
		_, err := io.WriteString(w, body)
		require.NoError(t, err)
	}
}

func TestUnitResourceNetworkExplicitForwardProxy_10_ReadErrors(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				respond(t, w, code, `{"error":{"code":"test","message":"synthetic error"}}`)
			})
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkExplicitForwardProxy_11_InvalidResponsePreservesState(t *testing.T) {
	for _, body := range []string{`{}`, `null`, ``} {
		t.Run(body, func(t *testing.T) {
			r := testClient(
				t,
				func(w http.ResponseWriter, q *http.Request) { respond(t, w, 200, body) },
			)
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkExplicitForwardProxy_12_FailedUpdateRetainsState(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, "PATCH", q.Method)
		respond(t, w, 403, `{"error":{"code":"Forbidden","message":"synthetic"}}`)
	})
	state := testState(t, r, testModel())
	plan := testModel()
	plan.InternetAccess.IsEnabled = types.BoolValue(true)
	planned := testState(t, r, plan)
	resp := resource.UpdateResponse{State: state}
	r.Update(
		context.Background(),
		resource.UpdateRequest{
			State: state,
			Plan:  tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw},
		},
		&resp,
	)
	require.True(t, resp.Diagnostics.HasError())
	require.True(t, state.Raw.Equal(resp.State.Raw))
}

func TestUnitResourceNetworkExplicitForwardProxy_13_ImportValidation(t *testing.T) {
	r := NewNetworkExplicitForwardProxyResource().(*NetworkExplicitForwardProxyResource)
	state := testState(t, r, testModel())
	resp := resource.ImportStateResponse{State: state}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "wrong"}, &resp)
	require.True(t, resp.Diagnostics.HasError())
}

func TestUnitResourceNetworkExplicitForwardProxy_14_Payloads(t *testing.T) {
	for _, smart := range []bool{false, true} {
		for _, header := range []bool{false, true} {
			plan := testModel()
			plan.InternetAccess.IsSourceIPSessionAffinityEnabled = types.BoolValue(smart || header)
			options := "none"
			if smart && header {
				options = "useSessionId,useHttpHeader"
			} else if smart {
				options = "useSessionId"
			} else if header {
				options = "useHttpHeader"
			}
			plan.InternetAccess.SourceIPSessionAffinityOptions = types.StringValue(options)
			body, err := proxyPatch(plan, nil)
			require.NoError(t, err)
			access := body["internetAccess"].(map[string]any)
			require.Equal(t, smart || header, access["isSourceIpSessionAffinityEnabled"])
			if !smart && !header {
				require.NotContains(t, access, "sourceIpSessionAffinityOptions")
			} else {
				require.Contains(t, access, "sourceIpSessionAffinityOptions")
			}
			raw, _ := json.Marshal(body)
			var remote proxyResponse
			require.NoError(t, json.Unmarshal(raw, &remote))
			if !smart && !header {
				none := "none"
				remote.InternetAccess.SourceIPSessionAffinityOptions = &none
			}
			var actual NetworkExplicitForwardProxyResourceModel
			require.NoError(t, mapResponse(&actual, &remote))
			require.Equal(
				t,
				smart || header,
				actual.InternetAccess.IsSourceIPSessionAffinityEnabled.ValueBool(),
			)
			require.Equal(
				t,
				options,
				actual.InternetAccess.SourceIPSessionAffinityOptions.ValueString(),
			)
		}
	}
	state := testModel()
	plan := testModel()
	plan.InternetAccess.IsEnabled = types.BoolValue(true)
	body, err := proxyPatch(plan, &state)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"internetAccess": map[string]any{"isEnabled": true}}, body)
	same, err := proxyPatch(state, &state)
	require.NoError(t, err)
	require.Empty(t, same)
	plan.InternetAccess.IsEnabled = types.BoolNull()
	_, err = proxyPatch(plan, nil)
	require.Error(t, err)
}

func TestUnitResourceNetworkExplicitForwardProxy_15_AffinityValidation(t *testing.T) {
	for _, tc := range []struct {
		enabled bool
		options string
		valid   bool
	}{
		{false, "none", true},
		{true, "useSessionId", true},
		{true, "useHttpHeader", true},
		{true, "useSessionId,useHttpHeader", true},
		{true, "none", false},
		{false, "useSessionId", false},
		{false, "useHttpHeader", false},
		{true, "unsupported", false},
	} {
		t.Run(fmt.Sprintf("%t_%s", tc.enabled, tc.options), func(t *testing.T) {
			r := NewNetworkExplicitForwardProxyResource().(*NetworkExplicitForwardProxyResource)
			m := testModel()
			m.InternetAccess.IsSourceIPSessionAffinityEnabled = types.BoolValue(tc.enabled)
			m.InternetAccess.SourceIPSessionAffinityOptions = types.StringValue(tc.options)
			state := testState(t, r, m)
			var response resource.ValidateConfigResponse
			r.ValidateConfig(
				context.Background(),
				resource.ValidateConfigRequest{
					Config: tfsdk.Config{Schema: state.Schema, Raw: state.Raw},
				},
				&response,
			)
			require.Equal(t, !tc.valid, response.Diagnostics.HasError())
			_, err := proxyPatch(m, nil)
			require.Equal(t, !tc.valid, err != nil)
		})
	}
}
