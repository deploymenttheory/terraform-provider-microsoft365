package graphBetaAuthorizationPolicy

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

func testModel() AuthorizationPolicyResourceModel {
	return AuthorizationPolicyResourceModel{
		ID:                                     types.StringValue(singletonID),
		AllowedToSignUpEmailBasedSubscriptions: types.BoolValue(false),
		AllowedToUseSSPR:                       types.BoolValue(true),
		AllowEmailVerifiedUsersToJoinOrganization:         types.BoolValue(false),
		AllowUserConsentForRiskyApps:                      types.BoolValue(false),
		BlockMsolPowerShell:                               types.BoolValue(true),
		AllowInvitesFrom:                                  types.StringValue("none"),
		EnabledPreviewFeatures:                            types.SetNull(types.StringType),
		PermissionGrantPolicyIdsAssignedToDefaultUserRole: types.SetNull(types.StringType),
		DefaultUserRolePermissions: types.ObjectValueMust(
			defaultUserRolePermissionsAttributeTypes(),
			map[string]attr.Value{
				"allowed_to_create_agent_identity_blueprints":     types.BoolValue(false),
				"allowed_to_create_apps":                          types.BoolValue(false),
				"allowed_to_create_security_groups":               types.BoolValue(false),
				"allowed_to_create_tenants":                       types.BoolValue(false),
				"allowed_to_read_bitlocker_keys_for_owned_device": types.BoolValue(false),
				"allowed_to_read_other_users":                     types.BoolValue(true),
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
	r *AuthorizationPolicyResource,
	m AuthorizationPolicyResourceModel,
) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), m).HasError())
	return state
}

func testClient(t *testing.T, h http.HandlerFunc) *AuthorizationPolicyResource {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(
		&authentication.AnonymousAuthenticationProvider{},
	)
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewAuthorizationPolicyResource().(*AuthorizationPolicyResource)
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

const (
	readPath  = "/policies/authorizationPolicy"
	patchPath = readPath + "/authorizationPolicy"
)

// TestUnitResourceAuthorizationPolicy_10_RequestContract verifies the documented paths and PATCH payload.
func TestUnitResourceAuthorizationPolicy_10_RequestContract(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			var patch map[string]any
			requests := []string{}
			r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
				requests = append(requests, q.Method+" "+q.URL.Path)
				if q.Method == http.MethodPatch {
					require.Equal(t, patchPath, q.URL.Path)
					var reader io.Reader = q.Body
					if q.Header.Get("Content-Encoding") == "gzip" {
						compressed, err := gzip.NewReader(q.Body)
						require.NoError(t, err)
						defer compressed.Close()
						reader = compressed
					}
					require.NoError(t, json.NewDecoder(reader).Decode(&patch))
					reply(w, 204, "")
					return
				}
				require.Equal(t, http.MethodGet, q.Method)
				require.Equal(t, readPath, q.URL.Path)
				remote := map[string]any{"id": singletonID}
				for key, value := range patch {
					remote[key] = value
				}
				body, err := json.Marshal(map[string]any{"value": []any{remote}})
				require.NoError(t, err)
				reply(w, 200, string(body))
			})
			model := testModel()
			model.GuestUserRoleId = types.StringValue("2af84b1e-32c8-42b7-82bc-daa82404023b")
			model.EnabledPreviewFeatures = types.SetValueMust(types.StringType, []attr.Value{})
			model.PermissionGrantPolicyIdsAssignedToDefaultUserRole = types.SetValueMust(
				types.StringType,
				[]attr.Value{},
			)
			state := testState(t, r, model)
			plan := tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}
			if operation == "create" {
				resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
				r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
				require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
				require.False(t, resp.State.Get(context.Background(), &model).HasError())
			} else {
				resp := resource.UpdateResponse{State: state}
				r.Update(
					context.Background(),
					resource.UpdateRequest{State: state, Plan: plan},
					&resp,
				)
				require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
				require.False(t, resp.State.Get(context.Background(), &model).HasError())
			}
			require.Equal(t, singletonID, model.ID.ValueString())
			require.Equal(
				t,
				[]string{
					"PATCH " + patchPath,
					"GET " + readPath,
					"GET " + readPath,
					"GET " + readPath,
				},
				requests,
			)
			require.Equal(t, false, patch["allowedToSignUpEmailBasedSubscriptions"])
			require.Equal(t, false, patch["allowUserConsentForRiskyApps"])
			require.Equal(t, true, patch["allowedToUseSSPR"])
			require.Equal(t, true, patch["blockMsolPowerShell"])
			require.Equal(t, false, patch["allowEmailVerifiedUsersToJoinOrganization"])
			require.Equal(t, "none", patch["allowInvitesFrom"])
			require.Equal(t, "2af84b1e-32c8-42b7-82bc-daa82404023b", patch["guestUserRoleId"])
			require.Equal(t, []any{}, patch["enabledPreviewFeatures"])
			require.Equal(t, []any{}, patch["permissionGrantPolicyIdsAssignedToDefaultUserRole"])
			require.Equal(
				t,
				false,
				patch["defaultUserRolePermissions"].(map[string]any)["allowedToCreateApps"],
			)
			require.NotContains(t, patch, "id")
			require.NotContains(t, patch, "guestUserRole")
			require.NotContains(t, patch, "displayName")
			require.NotContains(t, patch, "description")
		})
	}
}

// TestUnitResourceAuthorizationPolicy_11_ReadErrorsPreserveState checks failed and malformed singleton reads.
func TestUnitResourceAuthorizationPolicy_11_ReadErrorsPreserveState(t *testing.T) {
	for _, test := range []struct {
		code int
		body string
	}{
		{403, `{"error":{"code":"Forbidden","message":"Read failed"}}`},
		{404, `{"error":{"code":"NotFound","message":"Read failed"}}`},
		{200, `{}`},
		{200, `{"value":[]}`},
		{200, `{"value":[{}]}`},
		{200, `{"value":[{"id":"other"}]}`},
		{200, `{"value":[{"id":"authorizationPolicy"},{"id":"authorizationPolicy"}]}`},
	} {
		t.Run(fmt.Sprintf("%d_%s", test.code, test.body), func(t *testing.T) {
			r := testClient(
				t,
				func(w http.ResponseWriter, q *http.Request) { reply(w, test.code, test.body) },
			)
			state := testState(t, r, testModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, state.Raw.Equal(resp.State.Raw))
		})
	}
}

// TestUnitResourceAuthorizationPolicy_12_CreateReadFailure retains the singleton ID after a successful PATCH.
func TestUnitResourceAuthorizationPolicy_12_CreateReadFailure(t *testing.T) {
	patches := 0
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == http.MethodPatch {
			patches++
			reply(w, 204, "")
			return
		}
		reply(w, 403, `{"error":{"code":"Forbidden","message":"Read failed"}}`)
	})
	model := testModel()
	model.ID = types.StringUnknown()
	state := testState(t, r, model)
	resp := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
	r.Create(
		context.Background(),
		resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}},
		&resp,
	)
	require.True(t, resp.Diagnostics.HasError())
	var id string
	require.False(t, resp.State.GetAttribute(context.Background(), path.Root("id"), &id).HasError())
	require.Equal(t, singletonID, id)
	require.Equal(t, 1, patches)
}

// TestUnitResourceAuthorizationPolicy_13_DeleteMakesNoRequest verifies state-only deletion.
func TestUnitResourceAuthorizationPolicy_13_DeleteMakesNoRequest(t *testing.T) {
	r := testClient(
		t,
		func(w http.ResponseWriter, q *http.Request) { t.Error("destroy must not send HTTP") },
	)
	state := testState(t, r, testModel())
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	require.True(t, resp.State.Raw.IsNull())
}

// TestUnitResourceAuthorizationPolicy_14_Import validates the fixed import identifier.
func TestUnitResourceAuthorizationPolicy_14_Import(t *testing.T) {
	r := NewAuthorizationPolicyResource().(*AuthorizationPolicyResource)
	for _, id := range []string{singletonID, "", "random-id"} {
		resp := resource.ImportStateResponse{State: testState(t, r, testModel())}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		require.Equal(t, id != singletonID, resp.Diagnostics.HasError())
	}
}

// TestUnitResourceAuthorizationPolicy_15_ReadFixture checks collection decoding and service null values.
func TestUnitResourceAuthorizationPolicy_15_ReadFixture(t *testing.T) {
	fixture, err := os.ReadFile(
		"tests/responses/validate_get/get_authorization_policy_success.json",
	)
	require.NoError(t, err)
	r := testClient(
		t,
		func(w http.ResponseWriter, q *http.Request) { reply(w, 200, string(fixture)) },
	)
	state := testState(t, r, testModel())
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	var model AuthorizationPolicyResourceModel
	require.False(t, resp.State.Get(context.Background(), &model).HasError())
	require.Equal(t, singletonID, model.ID.ValueString())
	require.Equal(t, "adminsGuestInvitersAndAllMembers", model.AllowInvitesFrom.ValueString())
	require.True(t, model.AllowUserConsentForRiskyApps.IsNull())
	require.True(
		t,
		model.DefaultUserRolePermissions.Attributes()["allowed_to_create_agent_identity_blueprints"].Equal(
			types.BoolValue(true),
		),
	)
	require.False(t, model.EnabledPreviewFeatures.IsNull())
	require.Empty(t, model.EnabledPreviewFeatures.Elements())
	require.Equal(t, "10dae51f-b6af-4016-8d66-8c2a99b929b3", model.GuestUserRoleId.ValueString())
}

// TestUnitResourceAuthorizationPolicy_16_Consistency rejects stale Boolean, set, and nested permissions.
func TestUnitResourceAuthorizationPolicy_16_Consistency(t *testing.T) {
	r := NewAuthorizationPolicyResource().(*AuthorizationPolicyResource)
	expected := testModel()
	predicate := consistencyPredicate(&expected)
	require.True(t, predicate(context.Background(), testState(t, r, expected)))
	for _, mutate := range []func(*AuthorizationPolicyResourceModel){
		func(m *AuthorizationPolicyResourceModel) { m.AllowedToUseSSPR = types.BoolValue(false) },
		func(m *AuthorizationPolicyResourceModel) { m.AllowInvitesFrom = types.StringValue("everyone") },
		func(m *AuthorizationPolicyResourceModel) {
			attrs := m.DefaultUserRolePermissions.Attributes()
			attrs["allowed_to_create_apps"] = types.BoolValue(true)
			m.DefaultUserRolePermissions = types.ObjectValueMust(defaultUserRolePermissionsAttributeTypes(), attrs)
		},
	} {
		stale := testModel()
		mutate(&stale)
		require.False(t, predicate(context.Background(), testState(t, r, stale)))
	}
	expected.EnabledPreviewFeatures = types.SetValueMust(types.StringType, []attr.Value{})
	require.False(
		t,
		consistencyPredicate(&expected)(context.Background(), testState(t, r, testModel())),
	)
}

// TestUnitResourceAuthorizationPolicy_17_ConsentPolicyCasing handles observed prefix normalization without hiding drift.
func TestUnitResourceAuthorizationPolicy_17_ConsentPolicyCasing(t *testing.T) {
	configured := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("managePermissionGrantsForSelf.policy-id")},
	)
	require.True(
		t,
		configured.Equal(
			mapPermissionGrantPolicyIDs(
				context.Background(),
				[]string{"ManagePermissionGrantsForSelf.policy-id"},
				configured,
			),
		),
	)
	changed := mapPermissionGrantPolicyIDs(
		context.Background(),
		[]string{"ManagePermissionGrantsForSelf.different-id"},
		configured,
	)
	require.False(t, configured.Equal(changed))
	require.False(
		t,
		configured.Equal(
			mapPermissionGrantPolicyIDs(
				context.Background(),
				[]string{"ManagePermissionGrantsForSelf.POLICY-ID"},
				configured,
			),
		),
	)
	empty := mapPermissionGrantPolicyIDs(context.Background(), []string{}, configured)
	require.False(t, empty.IsNull())
	require.Empty(t, empty.Elements())
}

// TestUnitResourceAuthorizationPolicy_18_StableReads rejects an isolated matching replica.
func TestUnitResourceAuthorizationPolicy_18_StableReads(t *testing.T) {
	r := NewAuthorizationPolicyResource().(*AuthorizationPolicyResource)
	expected := testModel()
	current := testState(t, r, expected)
	staleModel := testModel()
	staleModel.BlockMsolPowerShell = types.BoolValue(false)
	stale := testState(t, r, staleModel)
	predicate := stableConsistencyPredicate(&expected)
	ctx := context.Background()
	require.False(t, predicate(ctx, current))
	require.False(t, predicate(ctx, current))
	require.False(t, predicate(ctx, stale))
	require.False(t, predicate(ctx, current))
	require.False(t, predicate(ctx, current))
	require.True(t, predicate(ctx, current))
}
