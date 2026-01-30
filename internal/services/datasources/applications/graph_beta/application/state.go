package graphBetaApplication

import (
	"context"
	"encoding/base64"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps an Application Graph API object to the Terraform data source model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.Applicationable, ownerUserIds []string, config ApplicationDataSourceModel) ApplicationDataSourceModel {
	result := ApplicationDataSourceModel{
		ID:                         convert.GraphToFrameworkString(data.GetId()),
		ObjectId:                   convert.GraphToFrameworkString(data.GetId()),
		AppId:                      convert.GraphToFrameworkString(data.GetAppId()),
		DisplayName:                convert.GraphToFrameworkString(data.GetDisplayName()),
		Description:                convert.GraphToFrameworkString(data.GetDescription()),
		SignInAudience:             convert.GraphToFrameworkString(data.GetSignInAudience()),
		IdentifierUris:             convert.GraphToFrameworkStringSet(ctx, data.GetIdentifierUris()),
		Notes:                      convert.GraphToFrameworkString(data.GetNotes()),
		IsDeviceOnlyAuthSupported:  convert.GraphToFrameworkBool(data.GetIsDeviceOnlyAuthSupported()),
		IsFallbackPublicClient:     convert.GraphToFrameworkBool(data.GetIsFallbackPublicClient()),
		ServiceManagementReference: convert.GraphToFrameworkString(data.GetServiceManagementReference()),
		Tags:                       convert.GraphToFrameworkStringSet(ctx, data.GetTags()),
		DisabledByMicrosoftStatus:  convert.GraphToFrameworkString(data.GetDisabledByMicrosoftStatus()),
		PublisherDomain:            convert.GraphToFrameworkString(data.GetPublisherDomain()),
		CreatedDateTime:            convert.GraphToFrameworkTime(data.GetCreatedDateTime()),
		DeletedDateTime:            convert.GraphToFrameworkTime(data.GetDeletedDateTime()),
		ODataQuery:                 config.ODataQuery,
		Timeouts:                   config.Timeouts,
	}

	// Map GroupMembershipClaims
	if groupClaims := data.GetGroupMembershipClaims(); groupClaims != nil {
		result.GroupMembershipClaims = convert.GraphToFrameworkStringSet(ctx, []string{*groupClaims})
	} else {
		result.GroupMembershipClaims = types.SetNull(types.StringType)
	}

	// Map SignInAudienceRestrictions
	if restrictions := data.GetSignInAudienceRestrictions(); restrictions != nil {
		result.SignInAudienceRestrictions = mapSignInAudienceRestrictions(ctx, restrictions)
	} else {
		result.SignInAudienceRestrictions = types.ObjectNull(signInAudienceRestrictionsAttrTypes)
	}

	// Map Api configuration
	if api := data.GetApi(); api != nil {
		result.Api = mapApi(ctx, api)
	} else {
		result.Api = types.ObjectNull(apiAttrTypes)
	}

	// Map AppRoles
	if appRoles := data.GetAppRoles(); appRoles != nil {
		result.AppRoles = mapAppRoles(ctx, appRoles)
	} else {
		result.AppRoles = types.SetNull(types.ObjectType{AttrTypes: appRoleAttrTypes})
	}

	// Map Info (InformationalUrl)
	if info := data.GetInfo(); info != nil {
		result.Info = mapInfo(ctx, info)
	} else {
		result.Info = types.ObjectNull(infoAttrTypes)
	}

	// Map KeyCredentials
	if keyCredentials := data.GetKeyCredentials(); keyCredentials != nil {
		result.KeyCredentials = mapKeyCredentials(ctx, keyCredentials)
	} else {
		result.KeyCredentials = types.SetNull(types.ObjectType{AttrTypes: keyCredentialAttrTypes})
	}

	// Map PasswordCredentials
	if passwordCredentials := data.GetPasswordCredentials(); passwordCredentials != nil {
		result.PasswordCredentials = mapPasswordCredentials(ctx, passwordCredentials)
	} else {
		result.PasswordCredentials = types.SetNull(types.ObjectType{AttrTypes: passwordCredentialAttrTypes})
	}

	// Map OptionalClaims
	if optionalClaims := data.GetOptionalClaims(); optionalClaims != nil {
		result.OptionalClaims = mapOptionalClaims(ctx, optionalClaims)
	} else {
		result.OptionalClaims = types.ObjectNull(optionalClaimsAttrTypes)
	}

	// Map ParentalControlSettings
	if parentalControl := data.GetParentalControlSettings(); parentalControl != nil {
		result.ParentalControlSettings = mapParentalControlSettings(ctx, parentalControl)
	} else {
		result.ParentalControlSettings = types.ObjectNull(parentalControlSettingsAttrTypes)
	}

	// Map PublicClient
	if publicClient := data.GetPublicClient(); publicClient != nil {
		result.PublicClient = mapPublicClient(ctx, publicClient)
	} else {
		result.PublicClient = types.ObjectNull(publicClientAttrTypes)
	}

	// Map RequiredResourceAccess
	if requiredResourceAccess := data.GetRequiredResourceAccess(); requiredResourceAccess != nil {
		result.RequiredResourceAccess = mapRequiredResourceAccess(ctx, requiredResourceAccess)
	} else {
		result.RequiredResourceAccess = types.SetNull(types.ObjectType{AttrTypes: requiredResourceAccessAttrTypes})
	}

	// Map Spa
	if spa := data.GetSpa(); spa != nil {
		result.Spa = mapSpa(ctx, spa)
	} else {
		result.Spa = types.ObjectNull(spaAttrTypes)
	}

	// Map Web
	if web := data.GetWeb(); web != nil {
		result.Web = mapWeb(ctx, web)
	} else {
		result.Web = types.ObjectNull(webAttrTypes)
	}

	// Map OwnerUserIds
	if len(ownerUserIds) > 0 {
		result.OwnerUserIds = convert.GraphToFrameworkStringSet(ctx, ownerUserIds)
	} else {
		result.OwnerUserIds = types.SetNull(types.StringType)
	}

	return result
}

// mapSignInAudienceRestrictions converts Graph API polymorphic SignInAudienceRestrictions to Terraform types.Object
func mapSignInAudienceRestrictions(ctx context.Context, restrictions graphmodels.SignInAudienceRestrictionsBaseable) types.Object {
	if restrictions == nil {
		return types.ObjectNull(signInAudienceRestrictionsAttrTypes)
	}

	result := map[string]attr.Value{
		"odata_type":             types.StringNull(),
		"is_home_tenant_allowed": types.BoolNull(),
		"allowed_tenant_ids":     types.SetNull(types.StringType),
	}

	odataType := restrictions.GetOdataType()
	if odataType != nil {
		result["odata_type"] = types.StringValue(*odataType)

		switch *odataType {
		case "#microsoft.graph.allowedTenantsAudience":
			if allowedTenants, ok := restrictions.(graphmodels.AllowedTenantsAudienceable); ok {
				result["is_home_tenant_allowed"] = convert.GraphToFrameworkBool(allowedTenants.GetIsHomeTenantAllowed())
				result["allowed_tenant_ids"] = convert.GraphToFrameworkStringSet(ctx, allowedTenants.GetAllowedTenantIds())
			}
		case "#microsoft.graph.unrestrictedAudience":
			// No additional fields
		}
	}

	obj, diags := types.ObjectValue(signInAudienceRestrictionsAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(signInAudienceRestrictionsAttrTypes)
	}
	return obj
}

// mapApi maps the Api configuration to Terraform state
func mapApi(ctx context.Context, api graphmodels.ApiApplicationable) types.Object {
	if api == nil {
		return types.ObjectNull(apiAttrTypes)
	}

	knownClients := types.SetNull(types.StringType)
	if knownClientApps := api.GetKnownClientApplications(); len(knownClientApps) > 0 {
		clientIds := make([]string, 0, len(knownClientApps))
		for _, uuid := range knownClientApps {
			clientIds = append(clientIds, uuid.String())
		}
		knownClients = convert.GraphToFrameworkStringSet(ctx, clientIds)
	}

	result := map[string]attr.Value{
		"accept_mapped_claims":           convert.GraphToFrameworkBool(api.GetAcceptMappedClaims()),
		"known_client_applications":      knownClients,
		"oauth2_permission_scopes":       mapOAuth2PermissionScopes(ctx, api.GetOauth2PermissionScopes()),
		"pre_authorized_applications":    mapPreAuthorizedApplications(ctx, api.GetPreAuthorizedApplications()),
		"requested_access_token_version": convert.GraphToFrameworkInt32(api.GetRequestedAccessTokenVersion()),
	}

	obj, diags := types.ObjectValue(apiAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(apiAttrTypes)
	}
	return obj
}

func mapOAuth2PermissionScopes(ctx context.Context, scopes []graphmodels.PermissionScopeable) types.Set {
	if len(scopes) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: oAuth2PermissionScopeAttrTypes}, []attr.Value{})
		return emptySet
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

		objVal, diags := types.ObjectValue(oAuth2PermissionScopeAttrTypes, scopeMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: oAuth2PermissionScopeAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: oAuth2PermissionScopeAttrTypes})
	}
	return setVal
}

func mapPreAuthorizedApplications(ctx context.Context, preAuthApps []graphmodels.PreAuthorizedApplicationable) types.Set {
	if len(preAuthApps) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: preAuthorizedApplicationAttrTypes}, []attr.Value{})
		return emptySet
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

		objVal, diags := types.ObjectValue(preAuthorizedApplicationAttrTypes, appMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: preAuthorizedApplicationAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: preAuthorizedApplicationAttrTypes})
	}
	return setVal
}

func mapAppRoles(ctx context.Context, appRoles []graphmodels.AppRoleable) types.Set {
	if len(appRoles) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: appRoleAttrTypes}, []attr.Value{})
		return emptySet
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

		objVal, diags := types.ObjectValue(appRoleAttrTypes, roleMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: appRoleAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: appRoleAttrTypes})
	}
	return setVal
}

func mapInfo(ctx context.Context, info graphmodels.InformationalUrlable) types.Object {
	if info == nil {
		return types.ObjectNull(infoAttrTypes)
	}

	result := map[string]attr.Value{
		"logo_url":              convert.GraphToFrameworkString(info.GetLogoUrl()),
		"marketing_url":         convert.GraphToFrameworkString(info.GetMarketingUrl()),
		"privacy_statement_url": convert.GraphToFrameworkString(info.GetPrivacyStatementUrl()),
		"support_url":           convert.GraphToFrameworkString(info.GetSupportUrl()),
		"terms_of_service_url":  convert.GraphToFrameworkString(info.GetTermsOfServiceUrl()),
	}

	obj, diags := types.ObjectValue(infoAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(infoAttrTypes)
	}
	return obj
}

func mapKeyCredentials(ctx context.Context, keyCredentials []graphmodels.KeyCredentialable) types.Set {
	if len(keyCredentials) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: keyCredentialAttrTypes}, []attr.Value{})
		return emptySet
	}

	elements := make([]attr.Value, 0, len(keyCredentials))
	for _, keyCred := range keyCredentials {
		if keyCred == nil {
			continue
		}

		customKeyId := types.StringNull()
		if customKey := keyCred.GetCustomKeyIdentifier(); customKey != nil {
			customKeyId = types.StringValue(base64.StdEncoding.EncodeToString(customKey))
		}

		credMap := map[string]attr.Value{
			"custom_key_identifier": customKeyId,
			"display_name":          convert.GraphToFrameworkString(keyCred.GetDisplayName()),
			"end_date_time":         convert.GraphToFrameworkTime(keyCred.GetEndDateTime()),
			"key":                   types.StringNull(), // Key is never returned by API
			"key_id":                convert.GraphToFrameworkUUID(keyCred.GetKeyId()),
			"start_date_time":       convert.GraphToFrameworkTime(keyCred.GetStartDateTime()),
			"type":                  convert.GraphToFrameworkString(keyCred.GetTypeEscaped()),
			"usage":                 convert.GraphToFrameworkString(keyCred.GetUsage()),
		}

		objVal, diags := types.ObjectValue(keyCredentialAttrTypes, credMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: keyCredentialAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: keyCredentialAttrTypes})
	}
	return setVal
}

func mapPasswordCredentials(ctx context.Context, passwordCredentials []graphmodels.PasswordCredentialable) types.Set {
	if len(passwordCredentials) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: passwordCredentialAttrTypes}, []attr.Value{})
		return emptySet
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

		objVal, diags := types.ObjectValue(passwordCredentialAttrTypes, credMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: passwordCredentialAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: passwordCredentialAttrTypes})
	}
	return setVal
}

func mapOptionalClaims(ctx context.Context, optionalClaims graphmodels.OptionalClaimsable) types.Object {
	if optionalClaims == nil {
		return types.ObjectNull(optionalClaimsAttrTypes)
	}

	result := map[string]attr.Value{
		"access_token": types.SetNull(types.ObjectType{AttrTypes: optionalClaimAttrTypes}),
		"id_token":     types.SetNull(types.ObjectType{AttrTypes: optionalClaimAttrTypes}),
		"saml2_token":  types.SetNull(types.ObjectType{AttrTypes: optionalClaimAttrTypes}),
	}

	if accessToken := optionalClaims.GetAccessToken(); accessToken != nil {
		result["access_token"] = mapOptionalClaimArray(ctx, accessToken)
	}
	if idToken := optionalClaims.GetIdToken(); idToken != nil {
		result["id_token"] = mapOptionalClaimArray(ctx, idToken)
	}
	if saml2Token := optionalClaims.GetSaml2Token(); saml2Token != nil {
		result["saml2_token"] = mapOptionalClaimArray(ctx, saml2Token)
	}

	obj, diags := types.ObjectValue(optionalClaimsAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(optionalClaimsAttrTypes)
	}
	return obj
}

func mapOptionalClaimArray(ctx context.Context, claims []graphmodels.OptionalClaimable) types.Set {
	if len(claims) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: optionalClaimAttrTypes}, []attr.Value{})
		return emptySet
	}

	elements := make([]attr.Value, 0, len(claims))
	for _, claim := range claims {
		if claim == nil {
			continue
		}

		claimMap := map[string]attr.Value{
			"additional_properties": convert.GraphToFrameworkStringSet(ctx, claim.GetAdditionalProperties()),
			"essential":             convert.GraphToFrameworkBool(claim.GetEssential()),
			"name":                  convert.GraphToFrameworkString(claim.GetName()),
			"source":                convert.GraphToFrameworkString(claim.GetSource()),
		}

		objVal, diags := types.ObjectValue(optionalClaimAttrTypes, claimMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: optionalClaimAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: optionalClaimAttrTypes})
	}
	return setVal
}

func mapParentalControlSettings(ctx context.Context, parentalControl graphmodels.ParentalControlSettingsable) types.Object {
	if parentalControl == nil {
		return types.ObjectNull(parentalControlSettingsAttrTypes)
	}

	result := map[string]attr.Value{
		"countries_blocked_for_minors": convert.GraphToFrameworkStringSet(ctx, parentalControl.GetCountriesBlockedForMinors()),
		"legal_age_group_rule":         convert.GraphToFrameworkString(parentalControl.GetLegalAgeGroupRule()),
	}

	obj, diags := types.ObjectValue(parentalControlSettingsAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(parentalControlSettingsAttrTypes)
	}
	return obj
}

func mapPublicClient(ctx context.Context, publicClient graphmodels.PublicClientApplicationable) types.Object {
	if publicClient == nil {
		return types.ObjectNull(publicClientAttrTypes)
	}

	result := map[string]attr.Value{
		"redirect_uris": convert.GraphToFrameworkStringSet(ctx, publicClient.GetRedirectUris()),
	}

	obj, diags := types.ObjectValue(publicClientAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(publicClientAttrTypes)
	}
	return obj
}

func mapRequiredResourceAccess(ctx context.Context, requiredResourceAccess []graphmodels.RequiredResourceAccessable) types.Set {
	if len(requiredResourceAccess) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: requiredResourceAccessAttrTypes}, []attr.Value{})
		return emptySet
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

				raObj, diags := types.ObjectValue(resourceAccessAttrTypes, raMap)
				if diags.HasError() {
					continue
				}
				resourceAccessElements = append(resourceAccessElements, raObj)
			}
		}

		resourceAccessSet := types.SetNull(types.ObjectType{AttrTypes: resourceAccessAttrTypes})
		if len(resourceAccessElements) > 0 {
			var diags diag.Diagnostics
			resourceAccessSet, diags = types.SetValue(types.ObjectType{AttrTypes: resourceAccessAttrTypes}, resourceAccessElements)
			if diags.HasError() {
				resourceAccessSet = types.SetNull(types.ObjectType{AttrTypes: resourceAccessAttrTypes})
			}
		}

		rraMap := map[string]attr.Value{
			"resource_app_id": convert.GraphToFrameworkString(rra.GetResourceAppId()),
			"resource_access": resourceAccessSet,
		}

		objVal, diags := types.ObjectValue(requiredResourceAccessAttrTypes, rraMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: requiredResourceAccessAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: requiredResourceAccessAttrTypes})
	}
	return setVal
}

func mapSpa(ctx context.Context, spa graphmodels.SpaApplicationable) types.Object {
	if spa == nil {
		return types.ObjectNull(spaAttrTypes)
	}

	result := map[string]attr.Value{
		"redirect_uris": convert.GraphToFrameworkStringSet(ctx, spa.GetRedirectUris()),
	}

	obj, diags := types.ObjectValue(spaAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(spaAttrTypes)
	}
	return obj
}

func mapWeb(ctx context.Context, web graphmodels.WebApplicationable) types.Object {
	if web == nil {
		return types.ObjectNull(webAttrTypes)
	}

	implicitGrantObj := types.ObjectNull(implicitGrantSettingsAttrTypes)
	if implicitGrant := web.GetImplicitGrantSettings(); implicitGrant != nil {
		implicitGrantData := map[string]attr.Value{
			"enable_access_token_issuance": convert.GraphToFrameworkBool(implicitGrant.GetEnableAccessTokenIssuance()),
			"enable_id_token_issuance":     convert.GraphToFrameworkBool(implicitGrant.GetEnableIdTokenIssuance()),
		}
		var diags2 diag.Diagnostics
		implicitGrantObj, diags2 = types.ObjectValue(implicitGrantSettingsAttrTypes, implicitGrantData)
		if diags2.HasError() {
			implicitGrantObj = types.ObjectNull(implicitGrantSettingsAttrTypes)
		}
	}

	redirectUriSettingsSet := types.SetNull(types.ObjectType{AttrTypes: redirectUriSettingsAttrTypes})
	if redirectUriSettings := web.GetRedirectUriSettings(); redirectUriSettings != nil {
		redirectUriSettingsSet = mapRedirectUriSettings(ctx, redirectUriSettings)
	}

	result := map[string]attr.Value{
		"home_page_url":           convert.GraphToFrameworkString(web.GetHomePageUrl()),
		"logout_url":              convert.GraphToFrameworkString(web.GetLogoutUrl()),
		"redirect_uris":           convert.GraphToFrameworkStringSet(ctx, web.GetRedirectUris()),
		"implicit_grant_settings": implicitGrantObj,
		"redirect_uri_settings":   redirectUriSettingsSet,
	}

	obj, diags := types.ObjectValue(webAttrTypes, result)
	if diags.HasError() {
		return types.ObjectNull(webAttrTypes)
	}
	return obj
}

func mapRedirectUriSettings(ctx context.Context, settings []graphmodels.RedirectUriSettingsable) types.Set {
	if len(settings) == 0 {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: redirectUriSettingsAttrTypes}, []attr.Value{})
		return emptySet
	}

	elements := make([]attr.Value, 0, len(settings))
	for _, setting := range settings {
		if setting == nil {
			continue
		}

		settingMap := map[string]attr.Value{
			"uri":   convert.GraphToFrameworkString(setting.GetUri()),
			"index": convert.GraphToFrameworkInt32(setting.GetIndex()),
		}

		objVal, diags := types.ObjectValue(redirectUriSettingsAttrTypes, settingMap)
		if diags.HasError() {
			continue
		}
		elements = append(elements, objVal)
	}

	setVal, diags := types.SetValue(types.ObjectType{AttrTypes: redirectUriSettingsAttrTypes}, elements)
	if diags.HasError() {
		return types.SetNull(types.ObjectType{AttrTypes: redirectUriSettingsAttrTypes})
	}
	return setVal
}
