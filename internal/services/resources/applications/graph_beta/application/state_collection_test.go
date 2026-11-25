package graphBetaApplication

import (
	"context"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitApplicationImportedEmptyCollectionsValidate(t *testing.T) {
	ctx := context.Background()
	api := graphmodels.NewApiApplication()
	api.SetKnownClientApplications([]uuid.UUID{})
	api.SetOauth2PermissionScopes([]graphmodels.PermissionScopeable{})
	api.SetPreAuthorizedApplications([]graphmodels.PreAuthorizedApplicationable{})
	remote := graphmodels.NewApplication()
	remote.SetId(new("11111111-1111-4111-8111-111111111111"))
	remote.SetApi(api)
	remote.SetAppRoles([]graphmodels.AppRoleable{})
	remote.SetRequiredResourceAccess([]graphmodels.RequiredResourceAccessable{})
	var data ApplicationResourceModel
	MapRemoteStateToTerraform(ctx, &data, remote)
	apiValues := data.Api.Attributes()

	var response resource.SchemaResponse
	NewApplicationResource().Schema(ctx, resource.SchemaRequest{}, &response)
	apiSchema := response.Schema.Attributes["api"].(schema.SingleNestedAttribute)
	for _, tc := range []struct {
		name  string
		path  path.Path
		value types.Set
		rules []validator.Set
	}{
		{"api.known_client_applications", path.Root("api").AtName("known_client_applications"), apiValues["known_client_applications"].(types.Set), apiSchema.Attributes["known_client_applications"].(schema.SetAttribute).Validators},
		{"api.oauth2_permission_scopes", path.Root("api").AtName("oauth2_permission_scopes"), apiValues["oauth2_permission_scopes"].(types.Set), apiSchema.Attributes["oauth2_permission_scopes"].(schema.SetNestedAttribute).Validators},
		{"api.pre_authorized_applications", path.Root("api").AtName("pre_authorized_applications"), apiValues["pre_authorized_applications"].(types.Set), apiSchema.Attributes["pre_authorized_applications"].(schema.SetNestedAttribute).Validators},
		{"app_roles", path.Root("app_roles"), data.AppRoles, response.Schema.Attributes["app_roles"].(schema.SetNestedAttribute).Validators},
		{"required_resource_access", path.Root("required_resource_access"), data.RequiredResourceAccess, response.Schema.Attributes["required_resource_access"].(schema.SetNestedAttribute).Validators},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.False(t, tc.value.IsNull())
			require.Empty(t, tc.value.Elements())
			for _, rule := range tc.rules {
				var result validator.SetResponse
				rule.ValidateSet(ctx, validator.SetRequest{Path: tc.path, ConfigValue: tc.value}, &result)
				require.False(t, result.Diagnostics.HasError(), result.Diagnostics.Errors())
			}
		})
	}
}

func TestUnitApplicationAPIConsistencyIgnoresUnconfiguredScalars(t *testing.T) {
	ctx := context.Background()
	api := graphmodels.NewApiApplication()
	api.SetKnownClientApplications([]uuid.UUID{})
	api.SetOauth2PermissionScopes([]graphmodels.PermissionScopeable{})
	api.SetPreAuthorizedApplications([]graphmodels.PreAuthorizedApplicationable{})
	api.SetRequestedAccessTokenVersion(new(int32(1)))
	remote := graphmodels.NewApplication()
	remote.SetId(new("11111111-1111-4111-8111-111111111111"))
	remote.SetAppId(new("22222222-2222-4222-8222-222222222222"))
	remote.SetApi(api)
	var actual ApplicationResourceModel
	MapRemoteStateToTerraform(ctx, &actual, remote)
	actual.OwnerUserIds = types.SetNull(types.StringType)
	actual.Timeouts = timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
		"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
	})}
	var schemaResponse resource.SchemaResponse
	NewApplicationResource().Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	state := tfsdk.State{Schema: schemaResponse.Schema}
	require.False(t, state.Set(ctx, &actual).HasError())
	expected := actual
	attributes := actual.Api.Attributes()
	attributes["accept_mapped_claims"] = types.BoolUnknown()
	attributes["requested_access_token_version"] = types.Int32Unknown()
	expected.Api = types.ObjectValueMust(ApplicationApiAttrTypes, attributes)
	require.True(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "unconfigured computed scalar values must not prevent accepting an otherwise matching API")
	attributes["requested_access_token_version"] = types.Int32Value(2)
	expected.Api = types.ObjectValueMust(ApplicationApiAttrTypes, attributes)
	require.False(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "a stale configured scalar must still trigger a retry")
	attributes["requested_access_token_version"] = types.Int32Value(1)
	attributes["known_client_applications"] = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("33333333-3333-4333-8333-333333333333")})
	expected.Api = types.ObjectValueMust(ApplicationApiAttrTypes, attributes)
	require.False(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "a stale collection must still trigger a retry")

	actual.ParentalControlSettings = types.ObjectValueMust(ApplicationParentalControlSettingsAttrTypes, map[string]attr.Value{
		"countries_blocked_for_minors": types.SetNull(types.StringType), "legal_age_group_rule": types.StringValue("Allow"),
	})
	require.False(t, state.Set(ctx, &actual).HasError())
	expected = actual
	parentalAttributes := actual.ParentalControlSettings.Attributes()
	parentalAttributes["countries_blocked_for_minors"] = types.SetUnknown(types.StringType)
	expected.ParentalControlSettings = types.ObjectValueMust(ApplicationParentalControlSettingsAttrTypes, parentalAttributes)
	require.True(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "unconfigured computed collections must not block a metadata update")
	parentalAttributes["legal_age_group_rule"] = types.StringValue("BlockMinors")
	expected.ParentalControlSettings = types.ObjectValueMust(ApplicationParentalControlSettingsAttrTypes, parentalAttributes)
	require.False(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "known nested values must still trigger a retry when stale")

	implicitGrant := types.ObjectValueMust(ImplicitGrantSettingsAttrTypes, map[string]attr.Value{
		"enable_access_token_issuance": types.BoolValue(false), "enable_id_token_issuance": types.BoolValue(false),
	})
	actual.Web = types.ObjectValueMust(ApplicationWebAttrTypes, map[string]attr.Value{
		"home_page_url": types.StringValue("https://example.com"), "logout_url": types.StringNull(),
		"redirect_uris": types.SetNull(types.StringType), "implicit_grant_settings": implicitGrant,
	})
	require.False(t, state.Set(ctx, &actual).HasError())
	expected = actual
	webAttributes := actual.Web.Attributes()
	grantAttributes := implicitGrant.Attributes()
	grantAttributes["enable_id_token_issuance"] = types.BoolUnknown()
	webAttributes["implicit_grant_settings"] = types.ObjectValueMust(ImplicitGrantSettingsAttrTypes, grantAttributes)
	expected.Web = types.ObjectValueMust(ApplicationWebAttrTypes, webAttributes)
	require.True(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "a computed field inside a nested object must not block consistency")
	grantAttributes["enable_access_token_issuance"] = types.BoolValue(true)
	webAttributes["implicit_grant_settings"] = types.ObjectValueMust(ImplicitGrantSettingsAttrTypes, grantAttributes)
	expected.Web = types.ObjectValueMust(ApplicationWebAttrTypes, webAttributes)
	require.False(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "a stale known field inside a nested object must trigger a retry")
	missingWebAttributes := actual.Web.Attributes()
	missingWebAttributes["implicit_grant_settings"] = types.ObjectNull(ImplicitGrantSettingsAttrTypes)
	actual.Web = types.ObjectValueMust(ApplicationWebAttrTypes, missingWebAttributes)
	require.False(t, state.Set(ctx, &actual).HasError())
	require.False(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "a missing nested object must trigger a retry")
	actual.Web = types.ObjectNull(ApplicationWebAttrTypes)
	require.False(t, state.Set(ctx, &actual).HasError())
	require.False(t, readbackAccepts(t, schemaResponse.Schema, expected, state), "a missing configured object must trigger a retry")
}

// readbackAccepts reports whether the read-back predicate accepts actual for the expected model.
func readbackAccepts(t *testing.T, s schema.Schema, expected ApplicationResourceModel, actual tfsdk.State) bool {
	t.Helper()
	expectedState := tfsdk.State{Schema: s}
	require.False(t, expectedState.Set(context.Background(), &expected).HasError())
	return crud.AttributesMatchPredicate(expectedState, readbackMatch)(context.Background(), actual)
}
