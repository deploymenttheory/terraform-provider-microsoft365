package graphBetaApplication

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitApplicationCollectionRequestPayloads(t *testing.T) {
	for _, form := range []string{"null", "unknown", "empty"} {
		for _, create := range []bool{true, false} {
			operation := "update"
			if create {
				operation = "create"
			}
			t.Run(operation+"/"+form, func(t *testing.T) {
				ctx := context.Background()
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				set := func(element attr.Type) types.Set {
					switch form {
					case "unknown":
						return types.SetUnknown(element)
					case "empty":
						return types.SetValueMust(element, []attr.Value{})
					default:
						return types.SetNull(element)
					}
				}
				api, diags := types.ObjectValueFrom(ctx, ApplicationApiAttrTypes, ApplicationApi{
					KnownClientApplications:   set(types.StringType),
					OAuth2PermissionScopes:    set(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes}),
					PreAuthorizedApplications: set(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes}),
				})
				require.False(t, diags.HasError(), diags)
				model := ApplicationResourceModel{
					DisplayName: types.StringValue("Collection payload regression"), Api: api,
					AppRoles:               set(types.ObjectType{AttrTypes: AppRoleAttrTypes}),
					RequiredResourceAccess: set(types.ObjectType{AttrTypes: RequiredResourceAccessAttrTypes}),
					ParentalControlSettings: types.ObjectValueMust(ApplicationParentalControlSettingsAttrTypes, map[string]attr.Value{
						"countries_blocked_for_minors": set(types.StringType), "legal_age_group_rule": types.StringValue("Allow"),
					}),
					PublicClient: types.ObjectValueMust(ApplicationPublicClientAttrTypes, map[string]attr.Value{"redirect_uris": set(types.StringType)}),
					Spa:          types.ObjectValueMust(ApplicationSpaAttrTypes, map[string]attr.Value{"redirect_uris": set(types.StringType)}),
					Web: types.ObjectValueMust(ApplicationWebAttrTypes, map[string]attr.Value{
						"redirect_uris": set(types.StringType), "home_page_url": types.StringNull(), "logout_url": types.StringNull(),
						"implicit_grant_settings": types.ObjectNull(ImplicitGrantSettingsAttrTypes),
					}),
				}
				body, err := constructResource(ctx, client, &model, "11111111-1111-4111-8111-111111111111", create)
				require.NoError(t, err)
				var request *abstractions.RequestInformation
				if create {
					request, err = client.Applications().ToPostRequestInformation(ctx, body, nil)
				} else {
					request, err = client.Applications().ByApplicationId("11111111-1111-4111-8111-111111111111").ToPatchRequestInformation(ctx, body, nil)
				}
				require.NoError(t, err)
				var payload map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(request.Content, &payload))
				var apiPayload map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(payload["api"], &apiPayload))
				assertCollection := func(payload map[string]json.RawMessage, field string) {
					t.Helper()
					if form == "empty" {
						require.JSONEq(t, "[]", string(payload[field]), field)
					} else {
						require.NotContains(t, payload, field)
					}
				}
				for _, field := range []string{"knownClientApplications", "oauth2PermissionScopes", "preAuthorizedApplications"} {
					assertCollection(apiPayload, field)
				}
				assertCollection(payload, "requiredResourceAccess")
				for _, nested := range []struct{ object, field string }{
					{"parentalControlSettings", "countriesBlockedForMinors"},
					{"publicClient", "redirectUris"}, {"spa", "redirectUris"}, {"web", "redirectUris"},
				} {
					var nestedPayload map[string]json.RawMessage
					require.NoError(t, json.Unmarshal(payload[nested.object], &nestedPayload))
					assertCollection(nestedPayload, nested.field)
				}
				if create {
					assertCollection(payload, "appRoles")
				} else {
					// Updates retain the existing separate app-role PATCH.
					require.NotContains(t, payload, "appRoles")
					roles, err := ConstructAppRolesForUpdate(ctx, model.AppRoles)
					require.NoError(t, err)
					if form != "empty" {
						require.Nil(t, roles)
						return
					}
					roleBody := graphmodels.NewApplication()
					roleBody.SetAppRoles(roles)
					roleRequest, err := client.Applications().ByApplicationId("11111111-1111-4111-8111-111111111111").ToPatchRequestInformation(ctx, roleBody, nil)
					require.NoError(t, err)
					var rolePayload map[string]json.RawMessage
					require.NoError(t, json.Unmarshal(roleRequest.Content, &rolePayload))
					assertCollection(rolePayload, "appRoles")
				}
			})
		}
	}
}

func TestUnitApplicationCreationDefaults(t *testing.T) {
	for _, tc := range []struct {
		name    string
		flag    types.Bool
		version types.Int32
	}{
		{"null", types.BoolNull(), types.Int32Null()},
		{"unknown", types.BoolUnknown(), types.Int32Unknown()},
		{"explicit", types.BoolValue(true), types.Int32Value(2)},
	} {
		for _, create := range []bool{true, false} {
			operation := "update"
			if create {
				operation = "create"
			}
			t.Run(operation+"/"+tc.name, func(t *testing.T) {
				ctx := context.Background()
				client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
				api, diags := types.ObjectValueFrom(ctx, ApplicationApiAttrTypes, ApplicationApi{
					RequestedAccessTokenVersion: tc.version,
					KnownClientApplications:     types.SetNull(types.StringType),
					OAuth2PermissionScopes:      types.SetNull(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes}),
					PreAuthorizedApplications:   types.SetNull(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes}),
				})
				require.False(t, diags.HasError(), diags)
				model := ApplicationResourceModel{DisplayName: types.StringValue("Creation defaults regression"), Api: api,
					IsDeviceOnlyAuthSupported: tc.flag, IsFallbackPublicClient: tc.flag}
				body, err := constructResource(ctx, client, &model, "11111111-1111-4111-8111-111111111111", create)
				require.NoError(t, err)
				if tc.name == "explicit" {
					require.Equal(t, new(true), body.GetIsDeviceOnlyAuthSupported())
					require.Equal(t, new(true), body.GetIsFallbackPublicClient())
					require.Equal(t, new(int32(2)), body.GetApi().GetRequestedAccessTokenVersion())
				} else if create {
					require.Equal(t, new(false), body.GetIsDeviceOnlyAuthSupported())
					require.Equal(t, new(false), body.GetIsFallbackPublicClient())
					require.Equal(t, new(int32(1)), body.GetApi().GetRequestedAccessTokenVersion())
				} else {
					require.Nil(t, body.GetIsDeviceOnlyAuthSupported())
					require.Nil(t, body.GetIsFallbackPublicClient())
					require.Nil(t, body.GetApi().GetRequestedAccessTokenVersion())
				}
			})
		}
	}
	t.Run("omitted API remains omitted on creation", func(t *testing.T) {
		ctx := context.Background()
		client := providerclient.NewMockGraphClients(&http.Client{}).MockBetaClient
		model := ApplicationResourceModel{DisplayName: types.StringValue("No API regression")}
		body, err := constructResource(ctx, client, &model, "", true)
		require.NoError(t, err)
		require.Nil(t, body.GetApi())
	})
}
