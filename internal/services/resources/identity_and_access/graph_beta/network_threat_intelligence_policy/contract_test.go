package graphBetaNetworkThreatIntelligencePolicy

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

func testModel() NetworkThreatIntelligencePolicyResourceModel {
	return NetworkThreatIntelligencePolicyResourceModel{
		ID: types.StringValue(
			"11111111-1111-1111-1111-111111111111",
		),
		Name:                 types.StringValue("test"),
		Description:          types.StringNull(),
		DefaultAction:        types.StringValue("allow"),
		Version:              types.StringValue("1.0.0"),
		LastModifiedDateTime: types.StringValue("2026-09-05T01:00:00Z"),
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
	r *NetworkThreatIntelligencePolicyResource,
	model NetworkThreatIntelligencePolicyResourceModel,
) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), model).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *NetworkThreatIntelligencePolicyResource {
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
	r := NewNetworkThreatIntelligencePolicyResource().(*NetworkThreatIntelligencePolicyResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}

func successFixture(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(
		"tests/responses/validate_get/get_threat_intelligence_policy_success.json",
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

func TestUnitResourceNetworkThreatIntelligencePolicy_10_ReadErrors(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_11_CreateReadbackFailureRetainsID(
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

func TestUnitResourceNetworkThreatIntelligencePolicy_12_DeleteErrors(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_13_DiffPatchClearsDescription(t *testing.T) {
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
	require.Len(t, patch, 2)
}

func TestUnitResourceNetworkThreatIntelligencePolicy_14_UpdateFailureRetainsState(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_17_CreateErrors(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_18_InvalidImport(t *testing.T) {
	r := NewNetworkThreatIntelligencePolicyResource().(*NetworkThreatIntelligencePolicyResource)
	for _, id := range []string{"", "not-a-uuid", "/", "11111111-1111-1111-1111-111111111111/invalid"} {
		resp := resource.ImportStateResponse{State: testState(t, r, testModel())}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		require.True(t, resp.Diagnostics.HasError(), id)
	}
}

func TestUnitResourceNetworkThreatIntelligencePolicy_19_NoAPIPatchForUnchangedFields(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_20_IdentityImport(t *testing.T) {
	r := NewNetworkThreatIntelligencePolicyResource().(*NetworkThreatIntelligencePolicyResource)
	var schema resource.IdentitySchemaResponse
	r.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, &schema)
	identity := tfsdk.ResourceIdentity{Schema: schema.IdentitySchema}
	require.False(t, identity.Set(context.Background(), struct {
		ID string `tfsdk:"id"`
	}{ID: testModel().ID.ValueString()}).HasError())
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

func TestUnitResourceNetworkThreatIntelligencePolicy_21_InvalidReadResponsePreservesState(
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

func TestUnitResourceNetworkThreatIntelligencePolicy_27_CreateRetriesThrottling(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_28_CreateThrottlingRetryLimit(t *testing.T) {
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

func TestUnitResourceNetworkThreatIntelligencePolicy_29_UpdateRetriesThrottling(t *testing.T) {
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
	err = r.updateThreatIntelligencePolicy(context.Background(), model.ID.ValueString(), body)
	require.NoError(t, err)
	require.Equal(t, 2, patches)
}
