package graphBetaNetworkProxyAutoConfiguration

import (
	"context"
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

func testModel() NetworkProxyAutoConfigurationResourceModel {
	return NetworkProxyAutoConfigurationResourceModel{
		ID: types.StringValue("11111111-2222-3333-4444-555555555555"),
		Name: types.StringValue(
			"test",
		),
		Content: types.StringValue(
			"function FindProxyForURL(url, host) { return \"DIRECT\"; }",
		),
		IsEnabled:            types.BoolValue(false),
		CreatedDateTime:      types.StringNull(),
		LastModifiedDateTime: types.StringNull(),
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
	r *NetworkProxyAutoConfigurationResource,
	m NetworkProxyAutoConfigurationResourceModel,
) tfsdk.State {
	t.Helper()
	var s resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &s)
	state := tfsdk.State{Schema: s.Schema}
	require.False(t, state.Set(context.Background(), m).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *NetworkProxyAutoConfigurationResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	a, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	a.SetBaseUrl(server.URL)
	r := NewNetworkProxyAutoConfigurationResource().(*NetworkProxyAutoConfigurationResource)
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

func TestUnitResourceNetworkProxyAutoConfiguration_10_ReadErrors(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				respond(t, w, code, `{"error":{"code":"test","message":"synthetic error"}}`)
			})
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			if code == 404 {
				require.True(t, resp.State.Raw.IsNull())
				require.False(t, resp.Diagnostics.HasError())
			} else {
				require.True(t, resp.Diagnostics.HasError())
				require.True(t, state.Raw.Equal(resp.State.Raw))
			}
		})
	}
}

func TestUnitResourceNetworkProxyAutoConfiguration_11_InvalidResponsePreservesState(t *testing.T) {
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

func TestUnitResourceNetworkProxyAutoConfiguration_12_FailedUpdateRetainsState(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, "PATCH", q.Method)
		respond(t, w, 403, `{"error":{"code":"Forbidden","message":"synthetic"}}`)
	})
	state := testState(t, r, testModel())
	plan := testModel()
	plan.Content = types.StringValue("// update")
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

func TestUnitResourceNetworkProxyAutoConfiguration_13_ImportValidation(t *testing.T) {
	r := NewNetworkProxyAutoConfigurationResource().(*NetworkProxyAutoConfigurationResource)
	state := testState(t, r, testModel())
	resp := resource.ImportStateResponse{State: state}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "wrong"}, &resp)
	require.True(t, resp.Diagnostics.HasError())
}

func TestUnitResourceNetworkProxyAutoConfiguration_14_Payloads(t *testing.T) {
	state := testModel()
	plan := state
	plan.Content = types.StringValue(
		"// comment\nfunction FindProxyForURL(url, host) { return \"DIRECT\"; }",
	)
	body, err := pacBody(plan, &state)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"content": plan.Content.ValueString()}, body)
	same, err := pacBody(state, &state)
	require.NoError(t, err)
	require.Empty(t, same)
	plan.Content = types.StringNull()
	_, err = pacBody(plan, nil)
	require.Error(t, err)
}

func TestUnitResourceNetworkProxyAutoConfiguration_15_FailedReadbackRetainsCreatedID(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "POST" {
			respond(t, w, 201, `{"id":"11111111-2222-3333-4444-555555555555"}`)
			return
		}
		require.Equal(t, "GET", q.Method)
		respond(t, w, 403, `{"error":{"code":"Forbidden","message":"synthetic"}}`)
	})
	ctx := context.Background()
	planned := testState(t, r, testModel())
	resp := resource.CreateResponse{State: tfsdk.State{Schema: planned.Schema}}
	r.Create(
		ctx,
		resource.CreateRequest{Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}},
		&resp,
	)
	require.True(t, resp.Diagnostics.HasError())
	var actual NetworkProxyAutoConfigurationResourceModel
	require.False(t, resp.State.Get(ctx, &actual).HasError())
	require.Equal(t, "11111111-2222-3333-4444-555555555555", actual.ID.ValueString())
}

func TestUnitResourceNetworkProxyAutoConfiguration_16_FailedDeleteRetainsState(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, "DELETE", q.Method)
		respond(t, w, 403, `{"error":{"code":"Forbidden","message":"synthetic"}}`)
	})
	state := testState(t, r, testModel())
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	require.True(t, state.Raw.Equal(resp.State.Raw))
}
