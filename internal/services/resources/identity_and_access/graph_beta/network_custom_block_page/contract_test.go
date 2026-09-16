package graphBetaNetworkCustomBlockPage

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

const markdownConfigurationType = "#microsoft.graph.networkaccess.markdownBlockMessageConfiguration"

func testModel() NetworkCustomBlockPageResourceModel {
	return NetworkCustomBlockPageResourceModel{ID: types.StringValue(singletonID), State: types.StringValue("enabled"), Configuration: testConfiguration("Access blocked. [Support](https://example.com)"), Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})}}
}
func testState(t *testing.T, r *NetworkCustomBlockPageResource, model NetworkCustomBlockPageResourceModel) tfsdk.State {
	t.Helper()
	var s resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &s)
	state := tfsdk.State{Schema: s.Schema}
	require.False(t, state.Set(context.Background(), model).HasError())
	return state
}
func testClient(t *testing.T, h http.HandlerFunc) *NetworkCustomBlockPageResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewNetworkCustomBlockPageResource().(*NetworkCustomBlockPageResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}

func TestUnitResourceNetworkCustomBlockPage_Lifecycle(t *testing.T) {
	ctx := context.Background()
	currentState, currentBody := "disabled", "Existing message"
	methods := []string{}
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		methods = append(methods, q.Method)
		require.Equal(t, "/networkAccess/settings/customBlockPage", q.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch q.Method {
		case "PATCH":
			var payload map[string]any
			var reader io.Reader = q.Body
			if q.Header.Get("Content-Encoding") == "gzip" {
				gz, err := gzip.NewReader(q.Body)
				require.NoError(t, err)
				defer gz.Close()
				reader = gz
			}
			require.NoError(t, json.NewDecoder(reader).Decode(&payload))
			require.NotContains(t, payload, "id")
			require.NotContains(t, payload, "@odata.type")
			currentState = payload["state"].(string)
			if config, ok := payload["configuration"]; ok {
				c := config.(map[string]any)
				require.Equal(t, markdownConfigurationType, c["@odata.type"])
				currentBody = c["body"].(string)
			}
			w.WriteHeader(http.StatusNoContent)
		case "GET":
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"state": currentState, "configuration": map[string]any{"@odata.type": "#microsoft.graph.networkaccess.markdownBlockMessageConfiguration", "body": currentBody}}))
		default:
			t.Errorf("unsupported operation issued: %s", q.Method)
			w.WriteHeader(405)
		}
	})
	model := testModel()
	model.ID = types.StringUnknown()
	model.Configuration = types.ObjectUnknown(configurationAttrTypes)
	planned := testState(t, r, model)
	create := resource.CreateResponse{State: tfsdk.State{Schema: planned.Schema}}
	r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}}, &create)
	require.False(t, create.Diagnostics.HasError(), "%v", create.Diagnostics)
	var actual NetworkCustomBlockPageResourceModel
	require.False(t, create.State.Get(ctx, &actual).HasError())
	require.Equal(t, singletonID, actual.ID.ValueString())
	require.Equal(t, "Existing message", actual.Configuration.Attributes()["body"].(types.String).ValueString())
	require.Equal(t, []string{"PATCH", "GET"}, methods)

	for _, body := range []string{"Updated [Support](https://example.com)", "Replacement message"} {
		actual.Configuration = testConfiguration(body)
		actual.State = types.StringValue("disabled")
		planned = testState(t, r, actual)
		update := resource.UpdateResponse{State: create.State}
		r.Update(ctx, resource.UpdateRequest{State: create.State, Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}}, &update)
		require.False(t, update.Diagnostics.HasError(), "%v", update.Diagnostics)
		require.Equal(t, body, currentBody)
		create.State = update.State
	}
	blank := testModel()
	blank.ID, blank.State = types.StringNull(), types.StringNull()
	blank.Configuration = types.ObjectNull(configurationAttrTypes)
	imported := resource.ImportStateResponse{State: testState(t, r, blank)}
	r.ImportState(ctx, resource.ImportStateRequest{ID: singletonID}, &imported)
	require.False(t, imported.Diagnostics.HasError(), "%v", imported.Diagnostics)
	read := resource.ReadResponse{State: imported.State}
	r.Read(ctx, resource.ReadRequest{State: imported.State}, &read)
	require.False(t, read.Diagnostics.HasError(), "%v", read.Diagnostics)
	require.False(t, read.State.Get(ctx, &actual).HasError())
	require.Equal(t, "Replacement message", actual.Configuration.Attributes()["body"].(types.String).ValueString())
	before := len(methods)
	deleted := resource.DeleteResponse{State: read.State}
	r.Delete(ctx, resource.DeleteRequest{State: read.State}, &deleted)
	require.False(t, deleted.Diagnostics.HasError())
	require.True(t, deleted.State.Raw.IsNull())
	require.Len(t, methods, before, "destroy must not send any HTTP request")
	require.Equal(t, "disabled", currentState)
}

func TestUnitResourceNetworkCustomBlockPage_ReadFailurePreservesState(t *testing.T) {
	for _, code := range []int{403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
				fmt.Fprint(w, `{"error":{"code":"Unavailable","message":"Cannot read settings"}}`)
			})
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkCustomBlockPage_CreateReadFailureRetainsID(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			w.WriteHeader(204)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		fmt.Fprint(w, `{"error":{"code":"Forbidden","message":"Read unavailable"}}`)
	})
	model := testModel()
	model.ID = types.StringUnknown()
	plan := testState(t, r, model)
	resp := resource.CreateResponse{State: tfsdk.State{Schema: plan.Schema}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	var id string
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, singletonID, id)
}

func TestUnitResourceNetworkCustomBlockPage_ConsistencyAndImport(t *testing.T) {
	r := NewNetworkCustomBlockPageResource().(*NetworkCustomBlockPageResource)
	expected := testModel()
	actual := expected
	actual.Configuration = testConfiguration("stale")
	opts := readOptions(&expected, "Update")
	require.False(t, opts.ConsistencyPredicate(context.Background(), testState(t, r, actual)))
	require.True(t, opts.ConsistencyPredicate(context.Background(), testState(t, r, expected)))
	actual.Configuration = types.ObjectNull(configurationAttrTypes)
	require.False(t, opts.ConsistencyPredicate(context.Background(), testState(t, r, actual)))
	state := testState(t, r, expected)
	resp := resource.ImportStateResponse{State: state}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "wrong-id"}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	require.True(t, state.Raw.Equal(resp.State.Raw))
}

func testConfiguration(body string) types.Object {
	return types.ObjectValueMust(configurationAttrTypes, map[string]attr.Value{
		"body": types.StringValue(body),
	})
}

func TestUnitResourceNetworkCustomBlockPage_OmittedConfigurationPlan(t *testing.T) {
	r := NewNetworkCustomBlockPageResource().(*NetworkCustomBlockPageResource)
	state := testState(t, r, testModel())
	var s resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &s)
	attribute := s.Schema.Attributes["configuration"].(schema.SingleNestedAttribute)
	for _, unknownConfig := range []bool{false, true} {
		config := types.ObjectNull(configurationAttrTypes)
		if unknownConfig {
			config = types.ObjectUnknown(configurationAttrTypes)
		}
		req := planmodifier.ObjectRequest{
			State: state, StateValue: testModel().Configuration,
			ConfigValue: config, PlanValue: types.ObjectUnknown(configurationAttrTypes),
		}
		resp := planmodifier.ObjectResponse{PlanValue: req.PlanValue}
		for _, modifier := range attribute.PlanModifiers {
			modifier.PlanModifyObject(context.Background(), req, &resp)
		}
		require.False(t, resp.Diagnostics.HasError())
		if unknownConfig {
			require.True(t, resp.PlanValue.IsUnknown(), "preserve unresolved configuration expressions")
		} else {
			require.True(t, resp.PlanValue.Equal(testModel().Configuration), "omission must preserve the known message without an update diff")
		}
	}
}
