package graphBetaNetworkThreatIntelligencePolicyRule

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
	"time"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"

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

func testModel() NetworkThreatIntelligencePolicyRuleResourceModel {
	return NetworkThreatIntelligencePolicyRuleResourceModel{
		ID: types.StringValue(
			"11111111-1111-1111-1111-111111111111",
		),
		Name:                       types.StringValue("test"),
		Description:                types.StringNull(),
		ThreatIntelligencePolicyID: types.StringValue("22222222-2222-2222-2222-222222222222"),
		Action:                     types.StringValue("allow"),
		Priority:                   types.Int32Value(1000),
		Enabled:                    types.BoolValue(true),
		Status:                     types.StringValue("enabled"),
		Severity:                   types.StringValue("high"),
		Destinations: types.ListValueMust(
			destinationObjectType,
			[]attr.Value{
				types.ObjectValueMust(
					destinationObjectType.AttrTypes,
					map[string]attr.Value{
						"type": types.StringValue("fqdn"),
						"values": types.ListValueMust(
							types.StringType,
							[]attr.Value{types.StringValue("example.com")},
						),
					},
				),
			},
		),
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
	r *NetworkThreatIntelligencePolicyRuleResource,
	model NetworkThreatIntelligencePolicyRuleResourceModel,
) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), model).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *NetworkThreatIntelligencePolicyRuleResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	// Exercise the normal provider HTTP pipeline, including configured retry and compression.
	t.Setenv("TF_ACC", "")
	require.NoError(t, os.Unsetenv("TF_ACC"))
	httpClient, err := providerclient.ConfigureGraphClientOptions(context.Background(), &providerclient.ProviderData{
		ClientOptions: &providerclient.ClientOptions{
			EnableRetry: true, MaxRetries: 3, RetryDelaySeconds: 1, EnableCompression: true,
		},
	})
	require.NoError(t, err)
	adapter, err := kiotahttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		&authentication.AnonymousAuthenticationProvider{}, nil, nil, httpClient,
	)
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewNetworkThreatIntelligencePolicyRuleResource().(*NetworkThreatIntelligencePolicyRuleResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}

func successFixture(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(
		"tests/responses/validate_get/get_threat_intelligence_policy_rule_success.json",
	)
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

func TestUnitResourceNetworkThreatIntelligencePolicyRule_10_ReadErrors(t *testing.T) {
	for _, code := range []int{400, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				writeResponse(
					t,
					w,
					code,
					[]byte(`{"error":{"code":"TestError","message":"Synthetic read failure"}}`),
				)
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

func TestUnitResourceNetworkThreatIntelligencePolicyRule_11_CreateReadbackFailureRetainsID(
	t *testing.T,
) {
	for _, code := range []int{400, 403, 404, 200} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			posts := 0
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				if q.Method == "POST" {
					posts++
					writeResponse(t, w, 201, successFixture(t))
					return
				}
				writeResponse(
					t,
					w,
					code,
					[]byte(`{"error":{"code":"TestError","message":"Readback failed"}}`),
				)
			})
			state := testState(t, r, testModel())
			var identitySchema resource.IdentitySchemaResponse
			r.IdentitySchema(
				context.Background(),
				resource.IdentitySchemaRequest{},
				&identitySchema,
			)
			identity := tfsdk.ResourceIdentity{Schema: identitySchema.IdentitySchema}
			resp := resource.CreateResponse{
				State:    tfsdk.State{Schema: state.Schema},
				Identity: &identity,
			}
			r.Create(
				context.Background(),
				resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}},
				&resp,
			)
			require.True(t, resp.Diagnostics.HasError())
			var id string
			require.False(
				t,
				resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError(),
			)
			require.Equal(t, testModel().ID.ValueString(), id)
			require.Equal(t, 1, posts)
			var identityID string
			require.False(
				t,
				resp.Identity.GetAttribute(context.Background(), path.Root("id"), &identityID).
					HasError(),
			)
			require.Equal(t, id, identityID)
		})
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_12_DeleteErrors(t *testing.T) {
	for _, code := range []int{204, 400, 403, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				if code == 204 {
					writeResponse(t, w, code, nil)
					return
				}
				writeResponse(
					t,
					w,
					code,
					[]byte(`{"error":{"code":"TestError","message":"Synthetic delete failure"}}`),
				)
			})
			state := testState(t, r, testModel())
			resp := resource.DeleteResponse{State: state}
			r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
			require.Equal(t, code == 400 || code == 403, resp.Diagnostics.HasError())
			require.Equal(t, code == 204 || code == 404, resp.State.Raw.IsNull())
		})
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_13_DiffPatchClearsDescription(
	t *testing.T,
) {
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
	r.Update(
		context.Background(),
		resource.UpdateRequest{
			State: state,
			Plan:  tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw},
		},
		&resp,
	)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Contains(t, patch, "description")
	require.Nil(t, patch["description"])
	require.NotContains(t, patch, "name")
	require.NotContains(t, patch, "settings")
	require.NotContains(t, patch, "policyRules")
	require.NotContains(t, patch, "version")
	require.Equal(t, "#microsoft.graph.networkaccess.threatIntelligenceRule", patch["@odata.type"])
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_14_UpdateFailureRetainsState(
	t *testing.T,
) {
	for _, method := range []string{"PATCH", "GET"} {
		for _, code := range []int{400, 403, 404} {
			t.Run(fmt.Sprintf("%s-%d", method, code), func(t *testing.T) {
				r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
					if q.Method == method {
						writeResponse(
							t,
							w,
							code,
							[]byte(
								`{"error":{"code":"TestError","message":"Synthetic update failure"}}`,
							),
						)
						return
					}
					writeResponse(t, w, 204, nil)
				})
				old := testModel()
				old.Description = types.StringValue("old")
				state := testState(t, r, old)
				resp := resource.UpdateResponse{State: state}
				r.Update(
					context.Background(),
					resource.UpdateRequest{
						State: state,
						Plan: tfsdk.Plan{
							Schema: state.Schema,
							Raw:    testState(t, r, testModel()).Raw,
						},
					},
					&resp,
				)
				require.True(t, resp.Diagnostics.HasError())
				require.True(t, state.Raw.Equal(resp.State.Raw))
			})
		}
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_15_UnknownRemoteValues(t *testing.T) {
	for _, scenario := range []string{"status", "destination", "system", "overflow", "missing"} {
		t.Run(scenario, func(t *testing.T) {
			var payload map[string]any
			require.NoError(t, json.Unmarshal(successFixture(t), &payload))
			switch scenario {
			case "status":
				payload["settings"] = map[string]any{"status": "unknownFutureValue"}
			case "destination":
				payload["matchingConditions"] = map[string]any{
					"severity": "high",
					"destinations": []any{
						map[string]any{
							"@odata.type": "#microsoft.graph.networkaccess.futureDestination",
							"values":      []string{"example.com"},
						},
					},
				}
			case "system":
				payload["priority"] = 65000
			case "overflow":
				payload["priority"] = int64(2147483648)
			case "missing":
				delete(payload, "id")
			}
			body, err := json.Marshal(payload)
			require.NoError(t, err)
			r := testClient(
				t,
				func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, 200, body) },
			)
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			if scenario == "status" {
				var actual NetworkThreatIntelligencePolicyRuleResourceModel
				require.False(t, resp.State.Get(context.Background(), &actual).HasError())
				require.Equal(t, "unknownFutureValue", actual.Status.ValueString())
				require.True(t, actual.Enabled.ValueBool())
			} else {
				require.True(t, state.Raw.Equal(resp.State.Raw))
			}
		})
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_16_EnabledOnlyPatch(t *testing.T) {
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
	r.Update(
		context.Background(),
		resource.UpdateRequest{
			State: state,
			Plan:  tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw},
		},
		&resp,
	)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Len(t, patch, 2)
	require.Equal(t, "#microsoft.graph.networkaccess.threatIntelligenceRule", patch["@odata.type"])
	require.Equal(t, map[string]any{"status": "disabled"}, patch["settings"])
	var result NetworkThreatIntelligencePolicyRuleResourceModel
	require.False(t, resp.State.Get(context.Background(), &result).HasError())
	require.False(t, result.Enabled.ValueBool())
	require.Equal(t, "disabled", result.Status.ValueString())
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_17_CreateErrors(t *testing.T) {
	for _, code := range []int{400, 403, 404, 201} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				writeResponse(
					t,
					w,
					code,
					[]byte(`{"error":{"code":"TestError","message":"Synthetic create failure"}}`),
				)
			})
			state := testState(t, r, testModel())
			resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
			r.Create(
				context.Background(),
				resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}},
				&resp,
			)
			require.True(t, resp.Diagnostics.HasError())
		})
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_18_InvalidImport(t *testing.T) {
	r := NewNetworkThreatIntelligencePolicyRuleResource().(*NetworkThreatIntelligencePolicyRuleResource)
	for _, id := range []string{"", "not-a-uuid", "/", "11111111-1111-1111-1111-111111111111/invalid"} {
		resp := resource.ImportStateResponse{State: testState(t, r, testModel())}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		require.True(t, resp.Diagnostics.HasError(), id)
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_19_NoAPIPatchForUnchangedFields(
	t *testing.T,
) {
	patches := 0
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			patches++
		}
		writeResponse(t, w, 200, successFixture(t))
	})
	state := testState(t, r, testModel())
	plan := testModel()
	plan.Timeouts = timeouts.Value{
		Object: types.ObjectValueMust(
			map[string]attr.Type{
				"create": types.StringType,
				"read":   types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			},
			map[string]attr.Value{
				"create": types.StringValue("5m"),
				"read":   types.StringNull(),
				"update": types.StringNull(),
				"delete": types.StringNull(),
			},
		),
	}
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
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Zero(t, patches)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_20_IdentityImport(t *testing.T) {
	r := NewNetworkThreatIntelligencePolicyRuleResource().(*NetworkThreatIntelligencePolicyRuleResource)
	var schema resource.IdentitySchemaResponse
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, &schema)
	identity := tfsdk.ResourceIdentity{Schema: schema.IdentitySchema}
	require.False(
		t,
		identity.Set(context.Background(), ThreatIntelligencePolicyRuleIdentity{ID: testModel().ID.ValueString(), ThreatIntelligencePolicyID: testModel().ThreatIntelligencePolicyID.ValueString()}).
			HasError(),
	)
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

func TestUnitResourceNetworkThreatIntelligencePolicyRule_21_InvalidReadResponsePreservesState(
	t *testing.T,
) {
	for _, tc := range []struct {
		code int
		body string
	}{{204, ""}, {200, "{}"}, {200, `{"id":123}`}} {
		t.Run(fmt.Sprintf("%d-%s", tc.code, tc.body), func(t *testing.T) {
			r := testClient(
				t,
				func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, tc.code, []byte(tc.body)) },
			)
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_22_ListRoundTrip(t *testing.T) {
	for _, values := range [][]string{{"Example.COM", "example.org", "Example.COM"}, {"example.org", "Example.COM", "example.org"}, {"example.org", "Example.COM"}} {
		t.Run(fmt.Sprint(values), func(t *testing.T) {
			var patch map[string]any
			var response map[string]any
			require.NoError(t, json.Unmarshal(successFixture(t), &response))
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				if q.Method == "PATCH" {
					decodeRequest(t, q, &patch)
					response["matchingConditions"] = patch["matchingConditions"]
					writeResponse(t, w, 204, nil)
					return
				}
				body, err := json.Marshal(response)
				require.NoError(t, err)
				writeResponse(t, w, 200, body)
			})
			plan := testModel()
			list, diags := types.ListValueFrom(context.Background(), types.StringType, values)
			require.False(t, diags.HasError())
			plan.Destinations = types.ListValueMust(
				destinationObjectType,
				[]attr.Value{
					types.ObjectValueMust(
						destinationObjectType.AttrTypes,
						map[string]attr.Value{"type": types.StringValue("fqdn"), "values": list},
					),
				},
			)
			state := testState(t, r, testModel())
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
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			var actual NetworkThreatIntelligencePolicyRuleResourceModel
			require.False(t, resp.State.Get(context.Background(), &actual).HasError())
			require.True(t, plan.Destinations.Equal(actual.Destinations))
		})
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_23_EmptyDestinations(t *testing.T) {
	plan := testModel()
	plan.Destinations = types.ListValueMust(destinationObjectType, []attr.Value{})
	b, err := constructResource(context.Background(), &plan)
	require.NoError(t, err)
	require.NotNil(t, b.GetMatchingConditions().GetDestinations())
	require.Empty(t, b.GetMatchingConditions().GetDestinations())
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_25_ReadImplicitLowSeverity(t *testing.T) {
	var payload map[string]any
	require.NoError(t, json.Unmarshal(successFixture(t), &payload))
	payload["matchingConditions"].(map[string]any)["severity"] = "low"
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	r := testClient(
		t,
		func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, 200, body) },
	)
	state := testState(t, r, testModel())
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	var actual NetworkThreatIntelligencePolicyRuleResourceModel
	require.False(t, resp.State.Get(context.Background(), &actual).HasError())
	require.Equal(t, "low", actual.Severity.ValueString())
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_26_UnknownStatusCannotInferImportedEnabled(
	t *testing.T,
) {
	var payload map[string]any
	require.NoError(t, json.Unmarshal(successFixture(t), &payload))
	payload["settings"] = map[string]any{"status": "unknownFutureValue"}
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	r := testClient(
		t,
		func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, 200, body) },
	)
	model := testModel()
	model.Enabled = types.BoolNull()
	state := testState(t, r, model)
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	var actual NetworkThreatIntelligencePolicyRuleResourceModel
	require.False(t, resp.State.Get(context.Background(), &actual).HasError())
	require.True(t, actual.Enabled.IsNull())
	require.Equal(t, "unknownFutureValue", actual.Status.ValueString())
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_27_CreateRetriesThrottling(t *testing.T) {
	posts, gets := 0, 0
	var firstBody map[string]any
	var throttledAt time.Time
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == http.MethodGet {
			gets++
			writeResponse(t, w, http.StatusOK, successFixture(t))
			return
		}
		require.Equal(t, http.MethodPost, q.Method)
		posts++
		var body map[string]any
		decodeRequest(t, q, &body)
		if posts == 1 {
			firstBody = body
			throttledAt = time.Now()
			w.Header().Set("Retry-After", "1")
			writeResponse(
				t,
				w,
				http.StatusTooManyRequests,
				[]byte(`{"error":{"code":"TooManyRequests","message":"Throttled"}}`),
			)
			return
		}
		require.GreaterOrEqual(t, time.Since(throttledAt), time.Second)
		require.Equal(t, firstBody, body)
		writeResponse(t, w, http.StatusCreated, successFixture(t))
	})
	state := testState(t, r, testModel())
	resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
	r.Create(
		context.Background(),
		resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}},
		&resp,
	)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Equal(t, 2, posts)
	require.Equal(t, 1, gets)
	var id string
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.NotEmpty(t, id)
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_28_CreateThrottlingRetryLimit(
	t *testing.T,
) {
	posts := 0
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, http.MethodPost, q.Method)
		posts++
		w.Header().Set("Retry-After", "0")
		writeResponse(
			t,
			w,
			http.StatusTooManyRequests,
			[]byte(`{"error":{"code":"TooManyRequests","message":"Throttled"}}`),
		)
	})
	state := testState(t, r, testModel())
	resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
	r.Create(
		context.Background(),
		resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}},
		&resp,
	)
	require.True(t, resp.Diagnostics.HasError())
	require.Equal(t, 4, posts) // Initial request plus the SDK's default three retries.
}

func TestUnitResourceNetworkThreatIntelligencePolicyRule_29_UpdateRetriesThrottling(t *testing.T) {
	patches := 0
	var original map[string]any
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		require.Equal(t, http.MethodPatch, q.Method)
		patches++
		var body map[string]any
		decodeRequest(t, q, &body)
		if patches == 1 {
			original = body
			w.Header().Set("Retry-After", "0")
			writeResponse(
				t,
				w,
				http.StatusTooManyRequests,
				[]byte(`{"error":{"code":"TooManyRequests","message":"Throttled"}}`),
			)
			return
		}
		require.Equal(t, original, body)
		writeResponse(t, w, http.StatusNoContent, nil)
	})
	model := testModel()
	body, err := constructResource(context.Background(), &model)
	require.NoError(t, err)
	err = r.updateThreatIntelligencePolicyRule(
		context.Background(),
		model.ThreatIntelligencePolicyID.ValueString(),
		model.ID.ValueString(),
		body,
	)
	require.NoError(t, err)
	require.Equal(t, 2, patches)
}
