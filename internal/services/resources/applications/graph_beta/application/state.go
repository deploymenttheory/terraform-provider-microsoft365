package graphBetaApplication

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the Microsoft Graph API application response to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *ApplicationResourceModel, remoteResource graphmodels.Applicationable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AppId = convert.GraphToFrameworkString(remoteResource.GetAppId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.SignInAudience = convert.GraphToFrameworkString(remoteResource.GetSignInAudience())
	
	// Map SignInAudienceRestrictions - complex polymorphic type, return null for now
	data.SignInAudienceRestrictions = types.ObjectNull(SignInAudienceRestrictionsAttrTypes)
	
	data.IdentifierUris = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetIdentifierUris())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.IsDeviceOnlyAuthSupported = convert.GraphToFrameworkBool(remoteResource.GetIsDeviceOnlyAuthSupported())
	data.IsFallbackPublicClient = convert.GraphToFrameworkBool(remoteResource.GetIsFallbackPublicClient())
	data.ServiceManagementReference = convert.GraphToFrameworkString(remoteResource.GetServiceManagementReference())
	data.Tags = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetTags())
	data.DisabledByMicrosoftStatus = convert.GraphToFrameworkString(remoteResource.GetDisabledByMicrosoftStatus())
	data.PublisherDomain = convert.GraphToFrameworkString(remoteResource.GetPublisherDomain())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.DeletedDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeletedDateTime())

	// Map GroupMembershipClaims
	if groupClaims := remoteResource.GetGroupMembershipClaims(); groupClaims != nil {
		data.GroupMembershipClaims = convert.GraphToFrameworkStringSet(ctx, []string{*groupClaims})
	} else {
		data.GroupMembershipClaims = types.SetNull(types.StringType)
	}

	// Map Api configuration
	if api := remoteResource.GetApi(); api != nil {
		data.Api = mapApiToTerraform(ctx, api)
	} else {
		data.Api = types.ObjectNull(ApplicationApiAttrTypes)
	}

	// Map AppRoles
	if appRoles := remoteResource.GetAppRoles(); appRoles != nil {
		data.AppRoles = mapAppRolesToTerraform(ctx, appRoles)
	} else {
		data.AppRoles = types.SetNull(types.ObjectType{AttrTypes: AppRoleAttrTypes})
	}

	// Map Info (InformationalUrl)
	if info := remoteResource.GetInfo(); info != nil {
		data.Info = mapInfoToTerraform(ctx, info)
	} else {
		data.Info = types.ObjectNull(ApplicationInformationalUrlAttrTypes)
	}

	// Map KeyCredentials
	if keyCredentials := remoteResource.GetKeyCredentials(); keyCredentials != nil {
		data.KeyCredentials = mapKeyCredentialsToTerraform(ctx, keyCredentials)
	} else {
		data.KeyCredentials = types.SetNull(types.ObjectType{AttrTypes: KeyCredentialAttrTypes})
	}

	// Map PasswordCredentials
	if passwordCredentials := remoteResource.GetPasswordCredentials(); passwordCredentials != nil {
		data.PasswordCredentials = mapPasswordCredentialsToTerraform(ctx, passwordCredentials)
	} else {
		data.PasswordCredentials = types.SetNull(types.ObjectType{AttrTypes: PasswordCredentialAttrTypes})
	}

	// Map OptionalClaims
	if optionalClaims := remoteResource.GetOptionalClaims(); optionalClaims != nil {
		data.OptionalClaims = mapOptionalClaimsToTerraform(ctx, optionalClaims)
	} else {
		data.OptionalClaims = types.ObjectNull(ApplicationOptionalClaimsAttrTypes)
	}

	// Map ParentalControlSettings
	if parentalControl := remoteResource.GetParentalControlSettings(); parentalControl != nil {
		data.ParentalControlSettings = mapParentalControlSettingsToTerraform(ctx, parentalControl)
	} else {
		data.ParentalControlSettings = types.ObjectNull(ApplicationParentalControlSettingsAttrTypes)
	}

	// Map PublicClient
	if publicClient := remoteResource.GetPublicClient(); publicClient != nil {
		data.PublicClient = mapPublicClientToTerraform(ctx, publicClient)
	} else {
		data.PublicClient = types.ObjectNull(ApplicationPublicClientAttrTypes)
	}

	// Map RequiredResourceAccess
	if requiredResourceAccess := remoteResource.GetRequiredResourceAccess(); requiredResourceAccess != nil {
		data.RequiredResourceAccess = mapRequiredResourceAccessToTerraform(ctx, requiredResourceAccess)
	} else {
		data.RequiredResourceAccess = types.SetNull(types.ObjectType{AttrTypes: RequiredResourceAccessAttrTypes})
	}

	// Map Spa
	if spa := remoteResource.GetSpa(); spa != nil {
		data.Spa = mapSpaToTerraform(ctx, spa)
	} else {
		data.Spa = types.ObjectNull(ApplicationSpaAttrTypes)
	}

	// Map Web
	if web := remoteResource.GetWeb(); web != nil {
		data.Web = mapWebToTerraform(ctx, web)
	} else {
		data.Web = types.ObjectNull(ApplicationWebAttrTypes)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// MapRemoteOwnersToTerraform maps the fetched owners to Terraform state.
// It filters the API response to only include owners that are explicitly configured
// in Terraform, ignoring auto-added owners like the app registration caller.
func MapRemoteOwnersToTerraform(ctx context.Context, data *ApplicationResourceModel, owners graphmodels.DirectoryObjectCollectionResponseable) {
	if owners == nil {
		tflog.Debug(ctx, "No owners response received, preserving existing state")
		return
	}

	ownerObjects := owners.GetValue()
	if len(ownerObjects) == 0 {
		tflog.Debug(ctx, "No owners found for application")
		data.OwnerUserIds = types.SetNull(types.StringType)
		return
	}

	var configuredOwnerIds []string
	if !data.OwnerUserIds.IsNull() && !data.OwnerUserIds.IsUnknown() {
		diags := data.OwnerUserIds.ElementsAs(ctx, &configuredOwnerIds, false)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to extract configured owner IDs, falling back to all owners")
			configuredOwnerIds = nil
		}
	}

	configuredSet := make(map[string]bool)
	for _, id := range configuredOwnerIds {
		configuredSet[id] = true
	}

	filteredOwnerIds := make([]string, 0)
	for _, owner := range ownerObjects {
		if owner != nil && owner.GetId() != nil {
			ownerId := *owner.GetId()
			if len(configuredSet) == 0 || configuredSet[ownerId] {
				filteredOwnerIds = append(filteredOwnerIds, ownerId)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d owners to Terraform state (filtered from %d total)", len(filteredOwnerIds), len(ownerObjects)))

	if len(filteredOwnerIds) > 0 {
		data.OwnerUserIds = convert.GraphToFrameworkStringSet(ctx, filteredOwnerIds)
	} else {
		data.OwnerUserIds = types.SetNull(types.StringType)
	}
}

// Helper functions for mapping nested structures
func mapApiToTerraform(ctx context.Context, api graphmodels.ApiApplicationable) types.Object {
	if api == nil {
		return types.ObjectNull(ApplicationApiAttrTypes)
	}

	result := ApplicationApi{
		AcceptMappedClaims:          convert.GraphToFrameworkBool(api.GetAcceptMappedClaims()),
		RequestedAccessTokenVersion: convert.GraphToFrameworkInt32(api.GetRequestedAccessTokenVersion()),
	}

	// Map KnownClientApplications
	if knownClients := api.GetKnownClientApplications(); knownClients != nil {
		clientIds := make([]string, 0, len(knownClients))
		for _, uuid := range knownClients {
			clientIds = append(clientIds, uuid.String())
		}
		result.KnownClientApplications = convert.GraphToFrameworkStringSet(ctx, clientIds)
	} else {
		result.KnownClientApplications = types.SetNull(types.StringType)
	}

	// Map OAuth2PermissionScopes
	if scopes := api.GetOauth2PermissionScopes(); scopes != nil {
		result.OAuth2PermissionScopes = mapOAuth2PermissionScopesToTerraform(ctx, scopes)
	} else {
		result.OAuth2PermissionScopes = types.SetNull(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes})
	}

	// Map PreAuthorizedApplications
	if preAuthApps := api.GetPreAuthorizedApplications(); preAuthApps != nil {
		result.PreAuthorizedApplications = mapPreAuthorizedApplicationsToTerraform(ctx, preAuthApps)
	} else {
		result.PreAuthorizedApplications = types.SetNull(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes})
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationApiAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert Api to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationApiAttrTypes)
	}

	return obj
}

func mapOAuth2PermissionScopesToTerraform(ctx context.Context, scopes []graphmodels.PermissionScopeable) types.Set {
	if len(scopes) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes})
	}

	elements := make([]attr.Value, 0, len(scopes))
	for _, scope := range scopes {
		if scope == nil {
			continue
		}

		scopeMap := map[string]attr.Value{
			"admin_consent_description":  convert.GraphToFrameworkString(scope.GetAdminConsentDescription()),
			"admin_consent_display_name": convert.GraphToFrameworkString(scope.GetAdminConsentDisplayName()),
			"id":                         convert.GraphToFrameworkUUID(scope.GetId()),
			"is_enabled":                 convert.GraphToFrameworkBool(scope.GetIsEnabled()),
			"type":                       convert.GraphToFrameworkString(scope.GetTypeEscaped()),
			"user_consent_description":   convert.GraphToFrameworkString(scope.GetUserConsentDescription()),
			"user_consent_display_name":  convert.GraphToFrameworkString(scope.GetUserConsentDisplayName()),
			"value":                      convert.GraphToFrameworkString(scope.GetValue()),
		}

		objVal, diags := types.ObjectValue(OAuth2PermissionScopeAttrTypes, scopeMap)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create OAuth2PermissionScope object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			continue
		}
		elements = append(elements, objVal)
	}

	if len(elements) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes})
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes}, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create OAuth2PermissionScopes set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: OAuth2PermissionScopeAttrTypes})
	}

	return setVal
}

func mapPreAuthorizedApplicationsToTerraform(ctx context.Context, preAuthApps []graphmodels.PreAuthorizedApplicationable) types.Set {
	if len(preAuthApps) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes})
	}

	elements := make([]attr.Value, 0, len(preAuthApps))
	for _, preAuthApp := range preAuthApps {
		if preAuthApp == nil {
			continue
		}

		delegatedPermissionIds := make([]string, 0)
		if permissions := preAuthApp.GetPermissionIds(); permissions != nil {
			delegatedPermissionIds = permissions
		}

		appMap := map[string]attr.Value{
			"app_id":                   convert.GraphToFrameworkString(preAuthApp.GetAppId()),
			"delegated_permission_ids": convert.GraphToFrameworkStringSet(ctx, delegatedPermissionIds),
		}

		objVal, diags := types.ObjectValue(PreAuthorizedApplicationAttrTypes, appMap)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create PreAuthorizedApplication object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			continue
		}
		elements = append(elements, objVal)
	}

	if len(elements) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes})
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes}, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create PreAuthorizedApplications set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: PreAuthorizedApplicationAttrTypes})
	}

	return setVal
}

func mapAppRolesToTerraform(ctx context.Context, appRoles []graphmodels.AppRoleable) types.Set {
	if len(appRoles) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: AppRoleAttrTypes})
	}

	elements := make([]attr.Value, 0, len(appRoles))
	for _, appRole := range appRoles {
		if appRole == nil {
			continue
		}

		allowedMemberTypes := make([]string, 0)
		if memberTypes := appRole.GetAllowedMemberTypes(); memberTypes != nil {
			allowedMemberTypes = memberTypes
		}

		roleMap := map[string]attr.Value{
			"allowed_member_types": convert.GraphToFrameworkStringSet(ctx, allowedMemberTypes),
			"description":          convert.GraphToFrameworkString(appRole.GetDescription()),
			"display_name":         convert.GraphToFrameworkString(appRole.GetDisplayName()),
			"id":                   convert.GraphToFrameworkUUID(appRole.GetId()),
			"is_enabled":           convert.GraphToFrameworkBool(appRole.GetIsEnabled()),
			"origin":               convert.GraphToFrameworkString(appRole.GetOrigin()),
			"value":                convert.GraphToFrameworkString(appRole.GetValue()),
		}

		objVal, diags := types.ObjectValue(AppRoleAttrTypes, roleMap)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create AppRole object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			continue
		}
		elements = append(elements, objVal)
	}

	if len(elements) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: AppRoleAttrTypes})
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: AppRoleAttrTypes}, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create AppRoles set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: AppRoleAttrTypes})
	}

	return setVal
}

func mapInfoToTerraform(ctx context.Context, info graphmodels.InformationalUrlable) types.Object {
	if info == nil {
		return types.ObjectNull(ApplicationInformationalUrlAttrTypes)
	}

	result := ApplicationInformationalUrl{
		LogoUrl:             convert.GraphToFrameworkString(info.GetLogoUrl()),
		MarketingUrl:        convert.GraphToFrameworkString(info.GetMarketingUrl()),
		PrivacyStatementUrl: convert.GraphToFrameworkString(info.GetPrivacyStatementUrl()),
		SupportUrl:          convert.GraphToFrameworkString(info.GetSupportUrl()),
		TermsOfServiceUrl:   convert.GraphToFrameworkString(info.GetTermsOfServiceUrl()),
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationInformationalUrlAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert InformationalUrl to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationInformationalUrlAttrTypes)
	}

	return obj
}

func mapKeyCredentialsToTerraform(ctx context.Context, keyCredentials []graphmodels.KeyCredentialable) types.Set {
	if len(keyCredentials) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: KeyCredentialAttrTypes})
	}

	elements := make([]attr.Value, 0, len(keyCredentials))
	for _, keyCred := range keyCredentials {
		if keyCred == nil {
			continue
		}

		keyValue := types.StringNull()
		if key := keyCred.GetKey(); key != nil {
			keyValue = types.StringValue(base64.StdEncoding.EncodeToString(key))
		}

		customKeyId := types.StringNull()
		if customKey := keyCred.GetCustomKeyIdentifier(); customKey != nil {
			customKeyId = types.StringValue(base64.StdEncoding.EncodeToString(customKey))
		}

		credMap := map[string]attr.Value{
			"custom_key_identifier": customKeyId,
			"display_name":          convert.GraphToFrameworkString(keyCred.GetDisplayName()),
			"end_date_time":         convert.GraphToFrameworkTime(keyCred.GetEndDateTime()),
			"key":                   keyValue,
			"key_id":                convert.GraphToFrameworkUUID(keyCred.GetKeyId()),
			"start_date_time":       convert.GraphToFrameworkTime(keyCred.GetStartDateTime()),
			"type":                  convert.GraphToFrameworkString(keyCred.GetTypeEscaped()),
			"usage":                 convert.GraphToFrameworkString(keyCred.GetUsage()),
		}

		objVal, diags := types.ObjectValue(KeyCredentialAttrTypes, credMap)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create KeyCredential object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			continue
		}
		elements = append(elements, objVal)
	}

	if len(elements) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: KeyCredentialAttrTypes})
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: KeyCredentialAttrTypes}, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create KeyCredentials set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: KeyCredentialAttrTypes})
	}

	return setVal
}

func mapPasswordCredentialsToTerraform(ctx context.Context, passwordCredentials []graphmodels.PasswordCredentialable) types.Set {
	if len(passwordCredentials) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: PasswordCredentialAttrTypes})
	}

	elements := make([]attr.Value, 0, len(passwordCredentials))
	for _, passCred := range passwordCredentials {
		if passCred == nil {
			continue
		}

		customKeyId := types.StringNull()
		if customKey := passCred.GetCustomKeyIdentifier(); customKey != nil {
			customKeyId = types.StringValue(base64.StdEncoding.EncodeToString(customKey))
		}

		credMap := map[string]attr.Value{
			"custom_key_identifier": customKeyId,
			"display_name":          convert.GraphToFrameworkString(passCred.GetDisplayName()),
			"end_date_time":         convert.GraphToFrameworkTime(passCred.GetEndDateTime()),
			"hint":                  convert.GraphToFrameworkString(passCred.GetHint()),
			"key_id":                convert.GraphToFrameworkUUID(passCred.GetKeyId()),
			"secret_text":           convert.GraphToFrameworkString(passCred.GetSecretText()),
			"start_date_time":       convert.GraphToFrameworkTime(passCred.GetStartDateTime()),
		}

		objVal, diags := types.ObjectValue(PasswordCredentialAttrTypes, credMap)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create PasswordCredential object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			continue
		}
		elements = append(elements, objVal)
	}

	if len(elements) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: PasswordCredentialAttrTypes})
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: PasswordCredentialAttrTypes}, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create PasswordCredentials set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: PasswordCredentialAttrTypes})
	}

	return setVal
}

func mapOptionalClaimsToTerraform(ctx context.Context, optionalClaims graphmodels.OptionalClaimsable) types.Object {
	if optionalClaims == nil {
		return types.ObjectNull(ApplicationOptionalClaimsAttrTypes)
	}

	result := ApplicationOptionalClaims{
		AccessToken: types.SetNull(types.ObjectType{AttrTypes: OptionalClaimAttrTypes}),
		IdToken:     types.SetNull(types.ObjectType{AttrTypes: OptionalClaimAttrTypes}),
		Saml2Token:  types.SetNull(types.ObjectType{AttrTypes: OptionalClaimAttrTypes}),
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationOptionalClaimsAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert OptionalClaims to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationOptionalClaimsAttrTypes)
	}

	return obj
}

func mapParentalControlSettingsToTerraform(ctx context.Context, parentalControl graphmodels.ParentalControlSettingsable) types.Object {
	if parentalControl == nil {
		return types.ObjectNull(ApplicationParentalControlSettingsAttrTypes)
	}

	result := ApplicationParentalControlSettings{
		CountriesBlockedForMinors: convert.GraphToFrameworkStringSet(ctx, parentalControl.GetCountriesBlockedForMinors()),
		LegalAgeGroupRule:         convert.GraphToFrameworkString(parentalControl.GetLegalAgeGroupRule()),
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationParentalControlSettingsAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert ParentalControlSettings to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationParentalControlSettingsAttrTypes)
	}

	return obj
}

func mapPublicClientToTerraform(ctx context.Context, publicClient graphmodels.PublicClientApplicationable) types.Object {
	if publicClient == nil {
		return types.ObjectNull(ApplicationPublicClientAttrTypes)
	}

	result := ApplicationPublicClient{
		RedirectUris: convert.GraphToFrameworkStringSet(ctx, publicClient.GetRedirectUris()),
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationPublicClientAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert PublicClient to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationPublicClientAttrTypes)
	}

	return obj
}

func mapRequiredResourceAccessToTerraform(ctx context.Context, requiredResourceAccess []graphmodels.RequiredResourceAccessable) types.Set {
	if len(requiredResourceAccess) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: RequiredResourceAccessAttrTypes})
	}

	elements := make([]attr.Value, 0, len(requiredResourceAccess))
	for _, rra := range requiredResourceAccess {
		if rra == nil {
			continue
		}

		resourceAccessElements := make([]attr.Value, 0)
		if resourceAccess := rra.GetResourceAccess(); resourceAccess != nil {
			for _, ra := range resourceAccess {
				if ra == nil {
					continue
				}

				raMap := map[string]attr.Value{
					"id":   convert.GraphToFrameworkUUID(ra.GetId()),
					"type": convert.GraphToFrameworkString(ra.GetTypeEscaped()),
				}

				raObj, diags := types.ObjectValue(ResourceAccessAttrTypes, raMap)
				if diags.HasError() {
					tflog.Error(ctx, "Failed to create ResourceAccess object", map[string]any{
						"error": diags.Errors()[0].Detail(),
					})
					continue
				}
				resourceAccessElements = append(resourceAccessElements, raObj)
			}
		}

		resourceAccessSet := types.SetNull(types.ObjectType{AttrTypes: ResourceAccessAttrTypes})
		if len(resourceAccessElements) > 0 {
			var diags diag.Diagnostics
			resourceAccessSet, diags = types.SetValue(types.ObjectType{AttrTypes: ResourceAccessAttrTypes}, resourceAccessElements)
			if diags.HasError() {
				tflog.Error(ctx, "Failed to create ResourceAccess set", map[string]any{
					"error": diags.Errors()[0].Detail(),
				})
				resourceAccessSet = types.SetNull(types.ObjectType{AttrTypes: ResourceAccessAttrTypes})
			}
		}

		rraMap := map[string]attr.Value{
			"resource_app_id": convert.GraphToFrameworkString(rra.GetResourceAppId()),
			"resource_access": resourceAccessSet,
		}

		objVal, diags := types.ObjectValue(RequiredResourceAccessAttrTypes, rraMap)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create RequiredResourceAccess object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			continue
		}
		elements = append(elements, objVal)
	}

	if len(elements) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: RequiredResourceAccessAttrTypes})
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: RequiredResourceAccessAttrTypes}, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create RequiredResourceAccess set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: RequiredResourceAccessAttrTypes})
	}

	return setVal
}

func mapSpaToTerraform(ctx context.Context, spa graphmodels.SpaApplicationable) types.Object {
	if spa == nil {
		return types.ObjectNull(ApplicationSpaAttrTypes)
	}

	result := ApplicationSpa{
		RedirectUris: convert.GraphToFrameworkStringSet(ctx, spa.GetRedirectUris()),
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationSpaAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert Spa to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationSpaAttrTypes)
	}

	return obj
}

func mapWebToTerraform(ctx context.Context, web graphmodels.WebApplicationable) types.Object {
	if web == nil {
		return types.ObjectNull(ApplicationWebAttrTypes)
	}

	implicitGrantObj := types.ObjectNull(ImplicitGrantSettingsAttrTypes)
	if implicitGrant := web.GetImplicitGrantSettings(); implicitGrant != nil {
		implicitGrantData := ApplicationWebImplicitGrantSettings{
			EnableAccessTokenIssuance: convert.GraphToFrameworkBool(implicitGrant.GetEnableAccessTokenIssuance()),
			EnableIdTokenIssuance:     convert.GraphToFrameworkBool(implicitGrant.GetEnableIdTokenIssuance()),
		}
		var diags2 diag.Diagnostics
		implicitGrantObj, diags2 = types.ObjectValueFrom(ctx, ImplicitGrantSettingsAttrTypes, implicitGrantData)
		if diags2.HasError() {
			tflog.Error(ctx, "Failed to convert ImplicitGrantSettings to types.Object", map[string]any{
				"error": diags2.Errors()[0].Detail(),
			})
			implicitGrantObj = types.ObjectNull(ImplicitGrantSettingsAttrTypes)
		}
	}

	result := ApplicationWeb{
		HomePageUrl:           convert.GraphToFrameworkString(web.GetHomePageUrl()),
		LogoutUrl:             convert.GraphToFrameworkString(web.GetLogoutUrl()),
		RedirectUris:          convert.GraphToFrameworkStringSet(ctx, web.GetRedirectUris()),
		ImplicitGrantSettings: implicitGrantObj,
		RedirectUriSettings:   types.SetNull(types.ObjectType{AttrTypes: RedirectUriSettingsAttrTypes}),
	}

	obj, diags := types.ObjectValueFrom(ctx, ApplicationWebAttrTypes, result)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert Web to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ObjectNull(ApplicationWebAttrTypes)
	}

	return obj
}
