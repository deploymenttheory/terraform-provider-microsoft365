package graphBetaNetworkForwardingOptions

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func testModel() NetworkForwardingOptionsResourceModel {
	return NetworkForwardingOptionsResourceModel{
		ID:                 types.StringValue(singletonID),
		SkipDNSLookupState: types.StringValue("disabled"),
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
	r *NetworkForwardingOptionsResource,
	m NetworkForwardingOptionsResourceModel,
) tfsdk.State {
	t.Helper()
	var s resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &s)
	state := tfsdk.State{Schema: s.Schema}
	require.False(t, state.Set(context.Background(), m).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *NetworkForwardingOptionsResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	a, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	a.SetBaseUrl(server.URL)
	r := NewNetworkForwardingOptionsResource().(*NetworkForwardingOptionsResource)
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

func decodeRequest(t *testing.T, q *http.Request) map[string]any {
	t.Helper()
	var reader io.Reader = q.Body
	if q.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(q.Body)
		require.NoError(t, err)
		defer gz.Close()
		reader = gz
	}
	var body map[string]any
	require.NoError(t, json.NewDecoder(reader).Decode(&body))
	return body
}

func TestUnitResourceNetworkForwardingOptions_10_AdoptUpdateDestroy(t *testing.T) {
	value := "enabled"
	patches := 0
	gets := 0
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, "/networkAccess/settings/forwardingOptions", q.URL.Path)
		switch q.Method {
		case "PATCH":
			patches++
			body := decodeRequest(t, q)
			require.Equal(t, map[string]any{"skipDnsLookupState": body["skipDnsLookupState"]}, body)
			value = body["skipDnsLookupState"].(string)
			respond(t, w, 204, "")
		case "GET":
			gets++
			respond(
				t,
				w,
				200,
				fmt.Sprintf(
					`{"skipDnsLookupState":%q,"unmanagedFutureSetting":"preserve-me"}`,
					value,
				),
			)
		default:
			t.Errorf("unexpected HTTP method %s", q.Method)
			respond(t, w, 405, "")
		}
	})
	ctx := context.Background()
	s := testState(t, r, testModel())
	create := resource.CreateResponse{State: tfsdk.State{Schema: s.Schema}}
	r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: s.Schema, Raw: s.Raw}}, &create)
	require.False(t, create.Diagnostics.HasError(), "%v", create.Diagnostics)
	require.Equal(t, 1, patches)
	require.Equal(t, 1, gets)
	var id string
	require.False(t, create.State.GetAttribute(ctx, path.Root("id"), &id).HasError())
	require.Equal(t, singletonID, id)
	model := testModel()
	model.SkipDNSLookupState = types.StringValue("enabled")
	planned := testState(t, r, model)
	update := resource.UpdateResponse{State: create.State}
	r.Update(
		ctx,
		resource.UpdateRequest{
			State: create.State,
			Plan:  tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw},
		},
		&update,
	)
	require.False(t, update.Diagnostics.HasError(), "%v", update.Diagnostics)
	require.Equal(t, 2, patches)
	require.Equal(t, "enabled", value)
	noChange := resource.UpdateResponse{State: update.State}
	r.Update(
		ctx,
		resource.UpdateRequest{
			State: update.State,
			Plan:  tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw},
		},
		&noChange,
	)
	require.False(t, noChange.Diagnostics.HasError())
	require.Equal(t, 2, patches)
	beforeGets := gets
	destroy := resource.DeleteResponse{State: noChange.State}
	r.Delete(ctx, resource.DeleteRequest{State: noChange.State}, &destroy)
	require.False(t, destroy.Diagnostics.HasError())
	require.True(t, destroy.State.Raw.IsNull())
	require.Equal(t, 2, patches)
	require.Equal(t, beforeGets, gets)
	require.Equal(t, "enabled", value)
}

func TestUnitResourceNetworkForwardingOptions_11_ReadErrorsPreserveState(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				respond(t, w, code, `{"error":{"code":"TestError","message":"Synthetic failure"}}`)
			})
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkForwardingOptions_12_InvalidResponse(t *testing.T) {
	for _, body := range []string{`{}`, `{"skipDnsLookupState":null}`, `{"skipDnsLookupState":"futureValue"}`} {
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

func TestUnitResourceNetworkForwardingOptions_13_FailedReadbackRetainsID(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			respond(t, w, 204, "")
		} else {
			respond(t, w, 403, `{"error":{"code":"Forbidden","message":"Synthetic failure"}}`)
		}
	})
	s := testState(t, r, testModel())
	resp := resource.CreateResponse{State: tfsdk.State{Schema: s.Schema}}
	r.Create(
		context.Background(),
		resource.CreateRequest{Plan: tfsdk.Plan{Schema: s.Schema, Raw: s.Raw}},
		&resp,
	)
	require.True(t, resp.Diagnostics.HasError())
	var id string
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, singletonID, id)
}

func TestUnitResourceNetworkForwardingOptions_14_Import(t *testing.T) {
	r := NewNetworkForwardingOptionsResource().(*NetworkForwardingOptionsResource)
	for _, id := range []string{"", "wrong", singletonID} {
		s := testState(t, r, testModel())
		resp := resource.ImportStateResponse{State: s}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		require.Equal(t, id != singletonID, resp.Diagnostics.HasError())
	}
}

func TestUnitResourceNetworkForwardingOptions_15_ConsistencyPredicate(t *testing.T) {
	r := NewNetworkForwardingOptionsResource().(*NetworkForwardingOptionsResource)
	model := testModel()
	p := readbackOptions(model, "Create").ConsistencyPredicate
	require.True(t, p(context.Background(), testState(t, r, model)))
	model.SkipDNSLookupState = types.StringValue("enabled")
	require.False(t, p(context.Background(), testState(t, r, model)))
}

func TestUnitResourceNetworkForwardingOptions_16_InvalidWriteSendsNothing(t *testing.T) {
	r := testClient(
		t,
		func(w http.ResponseWriter, q *http.Request) { t.Error("invalid write reached Graph") },
	)
	for _, v := range []types.String{types.StringNull(), types.StringUnknown(), types.StringValue(""), types.StringValue("invalid")} {
		m := testModel()
		m.SkipDNSLookupState = v
		require.Error(t, r.patch(context.Background(), m))
	}
}

func TestUnitResourceNetworkForwardingOptions_17_FailedUpdateRetainsState(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method != http.MethodPatch {
			t.Errorf("unexpected HTTP method after rejected update: %s", q.Method)
		}
		respond(t, w, 403, `{"error":{"code":"Forbidden","message":"Synthetic failure"}}`)
	})
	before := testState(t, r, testModel())
	model := testModel()
	model.SkipDNSLookupState = types.StringValue("enabled")
	after := testState(t, r, model)
	resp := resource.UpdateResponse{State: before}
	r.Update(
		context.Background(),
		resource.UpdateRequest{
			State: before,
			Plan:  tfsdk.Plan{Schema: after.Schema, Raw: after.Raw},
		},
		&resp,
	)
	require.True(t, resp.Diagnostics.HasError())
	require.True(t, before.Raw.Equal(resp.State.Raw))
}
