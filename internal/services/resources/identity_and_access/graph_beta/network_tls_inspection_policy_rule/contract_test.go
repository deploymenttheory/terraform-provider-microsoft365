package graphBetaNetworkTLSInspectionPolicyRule

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
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
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func testModel() NetworkTLSInspectionPolicyRuleResourceModel {
	return NetworkTLSInspectionPolicyRuleResourceModel{ID: types.StringValue("11111111-1111-1111-1111-111111111111"), Name: types.StringValue("test"), Description: types.StringNull(), TLSInspectionPolicyID: types.StringValue("22222222-2222-2222-2222-222222222222"), Action: types.StringValue("inspect"), Priority: types.Int32Value(1000), Enabled: types.BoolValue(true), Status: types.StringValue("enabled"), Destinations: types.ListValueMust(tlsInspectionPolicyRuleDestinationObjectType(), []attr.Value{types.ObjectValueMust(tlsInspectionPolicyRuleDestinationObjectType().AttrTypes, map[string]attr.Value{"type": types.StringValue("fqdn"), "values": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("example.com")})})}),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})}}
}
func testState(t *testing.T, r *NetworkTLSInspectionPolicyRuleResource, model NetworkTLSInspectionPolicyRuleResourceModel) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), model).HasError())
	return state
}
func testClient(t *testing.T, h http.HandlerFunc) *NetworkTLSInspectionPolicyRuleResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewNetworkTLSInspectionPolicyRuleResource().(*NetworkTLSInspectionPolicyRuleResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}
func successFixture(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("tests/responses/validate_get/get_tls_inspection_policy_rule_success.json")
	require.NoError(t, err)
	return b
}
func writeResponse(t *testing.T, w http.ResponseWriter, code int, body []byte) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if len(body) > 0 {
		_, err := w.Write(body)
		require.NoError(t, err)
	}
}
func TestUnitResourceNetworkTLSInspectionPolicyRule_10_ReadErrors(t *testing.T) {
	for _, code := range []int{400, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				writeResponse(t, w, code, []byte(`{"error":{"code":"TestError","message":"Synthetic read failure"}}`))
			})
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.Equal(t, code != 404, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			require.Equal(t, code == 404, resp.State.Raw.IsNull())
			if code != 404 {
				require.True(t, state.Raw.Equal(resp.State.Raw))
			}
		})
	}
}
func TestUnitResourceNetworkTLSInspectionPolicyRule_11_CreateReadbackFailureRetainsID(t *testing.T) {
	for _, code := range []int{400, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			posts := 0
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				if q.Method == "POST" {
					posts++
					writeResponse(t, w, 201, successFixture(t))
					return
				}
				writeResponse(t, w, code, []byte(`{"error":{"code":"TestError","message":"Readback failed"}}`))
			})
			state := testState(t, r, testModel())
			resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
			r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			var id string
			require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
			require.Equal(t, testModel().ID.ValueString(), id)
			require.Equal(t, 1, posts)
		})
	}
}
func TestUnitResourceNetworkTLSInspectionPolicyRule_12_DeleteErrors(t *testing.T) {
	for _, code := range []int{204, 400, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				if code == 204 {
					writeResponse(t, w, code, nil)
					return
				}
				writeResponse(t, w, code, []byte(`{"error":{"code":"TestError","message":"Synthetic delete failure"}}`))
			})
			state := testState(t, r, testModel())
			resp := resource.DeleteResponse{State: state}
			r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
			require.Equal(t, code == 400 || code == 403, resp.Diagnostics.HasError())
			require.Equal(t, code == 204 || code == 404, resp.State.Raw.IsNull())
		})
	}
}
func TestUnitResourceNetworkTLSInspectionPolicyRule_13_DiffPatchClearsDescription(t *testing.T) {
	var patch map[string]any
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			decodeRequest(t, q, &patch)
			writeResponse(t, w, 204, nil)
			return
		}
		writeResponse(t, w, 200, successFixture(t))
	})
	old := testModel()
	old.Description = types.StringValue("remove me")
	state := testState(t, r, old)
	planned := testState(t, r, testModel())
	resp := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{State: state, Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Contains(t, patch, "description")
	require.Nil(t, patch["description"])
	require.NotContains(t, patch, "name")
	require.NotContains(t, patch, "settings")
	require.NotContains(t, patch, "policyRules")
	require.NotContains(t, patch, "version")
	require.Equal(t, "#microsoft.graph.networkaccess.tlsInspectionRule", patch["@odata.type"])
}
func TestUnitResourceNetworkTLSInspectionPolicyRule_14_UpdateFailureRetainsState(t *testing.T) {
	for _, method := range []string{"PATCH", "GET"} {
		for _, code := range []int{400, 403, 404} {
			t.Run(fmt.Sprintf("%s-%d", method, code), func(t *testing.T) {
				r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
					if q.Method == method {
						writeResponse(t, w, code, []byte(`{"error":{"code":"TestError","message":"Synthetic update failure"}}`))
						return
					}
					writeResponse(t, w, 204, nil)
				})
				old := testModel()
				old.Description = types.StringValue("old")
				state := testState(t, r, old)
				resp := resource.UpdateResponse{State: state}
				r.Update(context.Background(), resource.UpdateRequest{State: state, Plan: tfsdk.Plan{Schema: state.Schema, Raw: testState(t, r, testModel()).Raw}}, &resp)
				require.True(t, resp.Diagnostics.HasError())
				require.True(t, state.Raw.Equal(resp.State.Raw))
			})
		}
	}
}

func TestUnitResourceNetworkTLSInspectionPolicyRule_15_UnknownRemoteValues(t *testing.T) {
	for _, scenario := range []string{"status", "destination", "system", "overflow", "missing"} {
		t.Run(scenario, func(t *testing.T) {
			var payload map[string]any
			require.NoError(t, json.Unmarshal(successFixture(t), &payload))
			switch scenario {
			case "status":
				payload["settings"] = map[string]any{"status": "unknownFutureValue"}
			case "destination":
				payload["matchingConditions"] = map[string]any{"destinations": []any{map[string]any{"@odata.type": "#microsoft.graph.networkaccess.futureDestination", "values": []string{"example.com"}}}}
			case "system":
				payload["priority"] = 50
				payload["matchingConditions"] = nil
			case "overflow":
				payload["priority"] = int64(2147483648)
			case "missing":
				delete(payload, "id")
			}
			body, err := json.Marshal(payload)
			require.NoError(t, err)
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, 200, body) })
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkTLSInspectionPolicyRule_16_EnabledOnlyPatch(t *testing.T) {
	var patch map[string]any
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			decodeRequest(t, q, &patch)
			writeResponse(t, w, 204, nil)
			return
		}
		var body map[string]any
		require.NoError(t, json.Unmarshal(successFixture(t), &body))
		body["settings"] = map[string]any{"status": "disabled"}
		b, err := json.Marshal(body)
		require.NoError(t, err)
		writeResponse(t, w, 200, b)
	})
	state := testState(t, r, testModel())
	plan := testModel()
	plan.Enabled = types.BoolValue(false)
	planned := testState(t, r, plan)
	resp := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{State: state, Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Len(t, patch, 2)
	require.Equal(t, "#microsoft.graph.networkaccess.tlsInspectionRule", patch["@odata.type"])
	require.Equal(t, map[string]any{"status": "disabled"}, patch["settings"])
	var result NetworkTLSInspectionPolicyRuleResourceModel
	require.False(t, resp.State.Get(context.Background(), &result).HasError())
	require.False(t, result.Enabled.ValueBool())
	require.Equal(t, "disabled", result.Status.ValueString())
}

func TestUnitResourceNetworkTLSInspectionPolicyRule_17_CreateErrors(t *testing.T) {
	for _, code := range []int{400, 403, 404, 201} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				writeResponse(t, w, code, []byte(`{"error":{"code":"TestError","message":"Synthetic create failure"}}`))
			})
			state := testState(t, r, testModel())
			resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
			r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &resp)
			require.True(t, resp.Diagnostics.HasError())
		})
	}
}
func TestUnitResourceNetworkTLSInspectionPolicyRule_18_InvalidImport(t *testing.T) {
	r := NewNetworkTLSInspectionPolicyRuleResource().(*NetworkTLSInspectionPolicyRuleResource)
	for _, id := range []string{"", "not-a-uuid", "/", "11111111-1111-1111-1111-111111111111/invalid"} {
		resp := resource.ImportStateResponse{State: testState(t, r, testModel())}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		require.True(t, resp.Diagnostics.HasError(), id)
	}
}

func TestUnitResourceNetworkTLSInspectionPolicyRule_19_NoAPIPatchForUnchangedFields(t *testing.T) {
	patches := 0
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			patches++
		}
		writeResponse(t, w, 200, successFixture(t))
	})
	state := testState(t, r, testModel())
	resp := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{State: state, Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Zero(t, patches)
}

func TestUnitResourceNetworkTLSInspectionPolicyRule_20_IdentityImport(t *testing.T) {
	r := NewNetworkTLSInspectionPolicyRuleResource().(*NetworkTLSInspectionPolicyRuleResource)
	var schema resource.IdentitySchemaResponse
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, &schema)
	identity := tfsdk.ResourceIdentity{Schema: schema.IdentitySchema}
	require.False(t, identity.Set(context.Background(), TLSInspectionPolicyRuleIdentity{ID: testModel().ID.ValueString(), TLSInspectionPolicyID: testModel().TLSInspectionPolicyID.ValueString()}).HasError())
	resp := resource.ImportStateResponse{State: testState(t, r, testModel())}
	r.ImportState(context.Background(), resource.ImportStateRequest{Identity: &identity}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	var id string
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, testModel().ID.ValueString(), id)
}

func decodeRequest(t *testing.T, q *http.Request, target any) {
	t.Helper()
	var body io.Reader = q.Body
	if q.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(q.Body)
		require.NoError(t, err)
		defer reader.Close()
		body = reader
	}
	require.NoError(t, json.NewDecoder(body).Decode(target))
}

func TestUnitResourceNetworkTLSInspectionPolicyRule_21_InvalidReadResponsePreservesState(t *testing.T) {
	for _, tc := range []struct {
		code int
		body string
	}{{204, ""}, {200, "{}"}, {200, `{"id":123}`}} {
		t.Run(fmt.Sprintf("%d-%s", tc.code, tc.body), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, tc.code, []byte(tc.body)) })
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}
