package graphBetaAgentCardManifest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func testManifestModel() AgentCardManifestResourceModel {
	return AgentCardManifestResourceModel{
		AgentCardManifestID: types.StringValue("employee-assistant"),
		DisplayName:         types.StringValue("managed"),
		OwnerIDs:            types.SetNull(types.StringType),
		DefaultInputModes:   types.SetNull(types.StringType),
		DefaultOutputModes:  types.SetNull(types.StringType),
		Provider:            types.ObjectNull(providerAttrTypes),
		Capabilities:        types.ObjectNull(capabilitiesAttrTypes),
		Skills:              types.ListNull(types.ObjectType{AttrTypes: skillAttrTypes}),
		SecuritySchemesJSON: jsontypes.NewNormalizedNull(),
		SecurityJSON:        jsontypes.NewNormalizedNull(),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}),
		},
	}
}

func testManifestState(t *testing.T, r *AgentCardManifestResource, m AgentCardManifestResourceModel) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), &m).HasError())
	return state
}

func testManifestClient(t *testing.T, h http.HandlerFunc) *AgentCardManifestResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	adapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		&authentication.AnonymousAuthenticationProvider{},
		nil,
		nil,
		server.Client(),
	)
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewAgentCardManifestResource().(*AgentCardManifestResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}

func TestUnitAgentCardManifest_AdoptUpdateImportAndForget(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	current := map[string]any{"id": "employee-assistant", "displayName": "existing", "description": "external"}
	var previous map[string]any
	reads, writes, forbidden := 0, 0, 0
	r := testManifestClient(t, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		require.Equal(t, "/agentRegistry/agentCardManifests/employee-assistant", req.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case http.MethodGet:
			response := current
			if writes > 0 {
				reads++
				if reads == 2 {
					response = previous
				}
			}
			require.NoError(t, json.NewEncoder(w).Encode(response))
		case http.MethodPatch:
			var body map[string]any
			require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
			require.NotContains(t, body, "id")
			require.NotContains(t, body, "createdDateTime")
			require.NotContains(t, body, "description", "unconfigured computed state must not be written back")
			previous = current
			current = make(map[string]any)
			for k, v := range previous {
				current[k] = v
			}
			for k, v := range body {
				current[k] = v
			}
			writes++
			reads = 0
			require.NoError(t, json.NewEncoder(w).Encode(current))
		default:
			forbidden++
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	configured := testManifestModel()
	planModel := configured
	planModel.ID = types.StringUnknown()
	planModel.Description = types.StringUnknown()
	planned := testManifestState(t, r, planModel)
	config := testManifestState(t, r, configured)
	response := resource.CreateResponse{State: testManifestState(t, r, testManifestModel())}
	r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}, Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw}}, &response)
	require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
	require.GreaterOrEqual(t, reads, 5, "a fresh, stale, then three fresh reads must complete")
	var created AgentCardManifestResourceModel
	require.False(t, response.State.Get(ctx, &created).HasError())
	require.Equal(t, "external", created.Description.ValueString())
	require.Equal(t, "employee-assistant", created.ID.ValueString())

	mu.Lock()
	current["description"] = "external change"
	mu.Unlock()
	configured.DisplayName = types.StringValue("updated")
	configured.OwnerIDs = types.SetValueMust(types.StringType, []attr.Value{})
	configured.SecuritySchemesJSON = jsontypes.NewNormalizedValue("{ \"oauth\" : { \"type\" : \"oauth2\" } }")
	configured.SecurityJSON = jsontypes.NewNormalizedValue("[ { \"oauth\" : [\"read\"] } ]")
	planModel = created
	planModel.DisplayName = configured.DisplayName
	planModel.OwnerIDs = configured.OwnerIDs
	planModel.SecuritySchemesJSON = configured.SecuritySchemesJSON
	planModel.SecurityJSON = configured.SecurityJSON
	planned = testManifestState(t, r, planModel)
	config = testManifestState(t, r, configured)
	update := resource.UpdateResponse{State: response.State}
	r.Update(
		ctx,
		resource.UpdateRequest{Plan: tfsdk.Plan{Schema: planned.Schema, Raw: planned.Raw}, Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw}, State: response.State},
		&update,
	)
	require.False(t, update.Diagnostics.HasError(), update.Diagnostics)
	var updated AgentCardManifestResourceModel
	require.False(t, update.State.Get(ctx, &updated).HasError())
	require.Equal(t, "external change", updated.Description.ValueString())
	require.Equal(t, configured.SecurityJSON.ValueString(), updated.SecurityJSON.ValueString())
	require.Empty(t, updated.OwnerIDs.Elements())

	imported := resource.ImportStateResponse{State: testManifestState(t, r, testManifestModel())}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "employee-assistant"}, &imported)
	require.False(t, imported.Diagnostics.HasError())
	var target types.String
	require.False(t, imported.State.GetAttribute(ctx, path.Root("agent_card_manifest_id"), &target).HasError())
	require.Equal(t, "employee-assistant", target.ValueString())
	countBefore := reads
	deletion := resource.DeleteResponse{State: update.State}
	r.Delete(ctx, resource.DeleteRequest{State: update.State}, &deletion)
	require.False(t, deletion.Diagnostics.HasError())
	require.True(t, deletion.State.Raw.IsNull())
	require.Equal(t, countBefore, reads)
	require.Zero(t, forbidden, "standalone manifest POST/DELETE are not documented")
	require.Equal(t, "updated", current["displayName"])
}

func TestUnitAgentCardManifest_ReadMissingAndDenied(t *testing.T) {
	for _, status := range []int{404, 403, 502} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			r := testManifestClient(t, func(w http.ResponseWriter, req *http.Request) {
				require.Equal(t, http.MethodGet, req.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, err := io.WriteString(w, `{"error":{"code":"Request_ResourceNotFound","message":"test failure"}}`)
				require.NoError(t, err)
			})
			m := testManifestModel()
			m.ID = types.StringValue("employee-assistant")
			state := testManifestState(t, r, m)
			response := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &response)
			if status == 404 {
				require.True(t, response.State.Raw.IsNull())
				require.False(t, response.Diagnostics.HasError())
			} else {
				require.True(t, response.Diagnostics.HasError())
				require.False(t, response.State.Raw.IsNull())
			}
		})
	}
}

func TestUnitAgentCardManifest_ConstructPreservesTypesAndOmissions(t *testing.T) {
	ctx := context.Background()
	m := testManifestModel()
	m.OwnerIDs = types.SetValueMust(types.StringType, []attr.Value{})
	params := jsontypes.NewNormalizedValue(`{"enabled":true,"retries":3,"nested":{"names":["a","b"]}}`)
	extensions, diagnostics := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: extensionAttrTypes},
		[]ExtensionModel{{URI: types.StringValue("https://example.test/extension"), ParamsJSON: params}},
	)
	require.False(t, diagnostics.HasError(), diagnostics)
	m.Capabilities, diagnostics = types.ObjectValueFrom(ctx, capabilitiesAttrTypes, CapabilitiesModel{Streaming: types.BoolValue(false), Extensions: extensions})
	require.False(t, diagnostics.HasError(), diagnostics)
	m.SecuritySchemesJSON = jsontypes.NewNormalizedValue(`{"key":{"type":"apiKey","in":"header","name":"X-Key"}}`)
	m.SecurityJSON = jsontypes.NewNormalizedValue(`[{"key":[]}]`)
	body, err := constructResource(ctx, &m)
	require.NoError(t, err)
	writer := jsonserialization.NewJsonSerializationWriter()
	require.NoError(t, writer.WriteObjectValue("", body))
	raw, err := writer.GetSerializedContent()
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	require.Equal(t, []any{}, payload["ownerIds"])
	require.NotContains(t, payload, "defaultInputModes")
	require.NotContains(t, payload, "createdBy")
	require.Equal(t, false, payload["capabilities"].(map[string]any)["streaming"])
	extension := payload["capabilities"].(map[string]any)["extensions"].([]any)[0].(map[string]any)
	require.Equal(t, true, extension["params"].(map[string]any)["enabled"])
	require.Equal(t, float64(3), extension["params"].(map[string]any)["retries"])
	require.NoError(t, MapRemoteResourceStateToTerraform(ctx, &m, bodyWithID(body, "employee-assistant")))
	require.Equal(t, params.ValueString(), func() string {
		var caps CapabilitiesModel
		require.False(t, m.Capabilities.As(ctx, &caps, basetypes.ObjectAsOptions{}).HasError())
		var values []ExtensionModel
		require.False(t, caps.Extensions.ElementsAs(ctx, &values, false).HasError())
		return values[0].ParamsJSON.ValueString()
	}())
}

func bodyWithID(body graphmodels.AgentCardManifestable, id string) graphmodels.AgentCardManifestable {
	body.SetId(&id)
	return body
}
