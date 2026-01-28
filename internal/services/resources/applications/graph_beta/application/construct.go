package graphBetaApplication

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
// isCreate indicates whether this is a create operation (POST) or update operation (PATCH)
// currentID is the resource ID for update operations, empty string for create operations
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *ApplicationResourceModel, currentID string, isCreate bool) (graphmodels.Applicationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model (isCreate: %t)", ResourceName, isCreate))

	if err := validateRequest(ctx, client, data, currentID); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewApplication()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.SignInAudience, requestBody.SetSignInAudience)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.ServiceManagementReference, requestBody.SetServiceManagementReference)
	convert.FrameworkToGraphString(data.DisabledByMicrosoftStatus, requestBody.SetDisabledByMicrosoftStatus)
	convert.FrameworkToGraphBool(data.IsDeviceOnlyAuthSupported, requestBody.SetIsDeviceOnlyAuthSupported)
	convert.FrameworkToGraphBool(data.IsFallbackPublicClient, requestBody.SetIsFallbackPublicClient)

	if err := constructGroupMembershipClaims(ctx, data, requestBody); err != nil {
		return nil, err
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IdentifierUris, requestBody.SetIdentifierUris); err != nil {
		return nil, fmt.Errorf("failed to set identifier_uris: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.Tags, requestBody.SetTags); err != nil {
		return nil, fmt.Errorf("failed to set tags: %w", err)
	}

	if !data.SignInAudienceRestrictions.IsNull() && !data.SignInAudienceRestrictions.IsUnknown() {
		var restrictionsData SignInAudienceRestrictions
		diags := data.SignInAudienceRestrictions.As(ctx, &restrictionsData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract sign_in_audience_restrictions data: %s", diags.Errors()[0].Detail())
		}

		restrictions, err := constructSignInAudienceRestrictions(ctx, &restrictionsData)
		if err != nil {
			return nil, err
		}
		requestBody.SetSignInAudienceRestrictions(restrictions)
	}

	if !data.Api.IsNull() && !data.Api.IsUnknown() {
		var apiData ApplicationApi
		diags := data.Api.As(ctx, &apiData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract api data: %s", diags.Errors()[0].Detail())
		}

		api, err := constructApi(ctx, &apiData)
		if err != nil {
			return nil, err
		}
		requestBody.SetApi(api)
	}

	if !data.AppRoles.IsNull() && !data.AppRoles.IsUnknown() {
		appRoles, err := constructAppRoles(ctx, data.AppRoles)
		if err != nil {
			return nil, err
		}
		requestBody.SetAppRoles(appRoles)
	}

	if !data.Info.IsNull() && !data.Info.IsUnknown() {
		var infoData ApplicationInformationalUrl
		diags := data.Info.As(ctx, &infoData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract info data: %s", diags.Errors()[0].Detail())
		}

		info := constructInformationalUrl(&infoData)
		requestBody.SetInfo(info)
	}

	if !data.KeyCredentials.IsNull() && !data.KeyCredentials.IsUnknown() {
		keyCredentials, err := constructKeyCredentials(ctx, data.KeyCredentials)
		if err != nil {
			return nil, err
		}
		requestBody.SetKeyCredentials(keyCredentials)
	}

	if !data.PasswordCredentials.IsNull() && !data.PasswordCredentials.IsUnknown() {
		passwordCredentials, err := constructPasswordCredentials(ctx, data.PasswordCredentials)
		if err != nil {
			return nil, err
		}
		requestBody.SetPasswordCredentials(passwordCredentials)
	}

	if !data.OptionalClaims.IsNull() && !data.OptionalClaims.IsUnknown() {
		var optionalClaimsData ApplicationOptionalClaims
		diags := data.OptionalClaims.As(ctx, &optionalClaimsData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract optional_claims data: %s", diags.Errors()[0].Detail())
		}

		optionalClaims, err := constructOptionalClaims(ctx, &optionalClaimsData)
		if err != nil {
			return nil, err
		}
		requestBody.SetOptionalClaims(optionalClaims)
	}

	if !data.ParentalControlSettings.IsNull() && !data.ParentalControlSettings.IsUnknown() {
		var parentalControlData ApplicationParentalControlSettings
		diags := data.ParentalControlSettings.As(ctx, &parentalControlData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract parental_control_settings data: %s", diags.Errors()[0].Detail())
		}

		parentalControl, err := constructParentalControlSettings(ctx, &parentalControlData)
		if err != nil {
			return nil, err
		}
		requestBody.SetParentalControlSettings(parentalControl)
	}

	if !data.PublicClient.IsNull() && !data.PublicClient.IsUnknown() {
		var publicClientData ApplicationPublicClient
		diags := data.PublicClient.As(ctx, &publicClientData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract public_client data: %s", diags.Errors()[0].Detail())
		}

		publicClient, err := constructPublicClient(ctx, &publicClientData)
		if err != nil {
			return nil, err
		}
		requestBody.SetPublicClient(publicClient)
	}

	if !data.RequiredResourceAccess.IsNull() && !data.RequiredResourceAccess.IsUnknown() {
		requiredResourceAccess, err := constructRequiredResourceAccess(ctx, data.RequiredResourceAccess)
		if err != nil {
			return nil, err
		}
		requestBody.SetRequiredResourceAccess(requiredResourceAccess)
	}

	if !data.Spa.IsNull() && !data.Spa.IsUnknown() {
		var spaData ApplicationSpa
		diags := data.Spa.As(ctx, &spaData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract spa data: %s", diags.Errors()[0].Detail())
		}

		spa, err := constructSpa(ctx, &spaData)
		if err != nil {
			return nil, err
		}
		requestBody.SetSpa(spa)
	}

	if !data.Web.IsNull() && !data.Web.IsUnknown() {
		var webData ApplicationWeb
		diags := data.Web.As(ctx, &webData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract web data: %s", diags.Errors()[0].Detail())
		}

		web, err := constructWeb(ctx, &webData)
		if err != nil {
			return nil, err
		}
		requestBody.SetWeb(web)
	}

	// Owners - handled separately in create operation via OData bind
	if isCreate {
		if err := constructOwners(ctx, data, requestBody); err != nil {
			return nil, err
		}
	}

	return requestBody, nil
}

// constructGroupMembershipClaims converts a set of group membership claims to a comma-separated string
func constructGroupMembershipClaims(ctx context.Context, data *ApplicationResourceModel, requestBody graphmodels.Applicationable) error {
	if data.GroupMembershipClaims.IsNull() || data.GroupMembershipClaims.IsUnknown() {
		return nil
	}

	var claims []string
	diags := data.GroupMembershipClaims.ElementsAs(ctx, &claims, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract group_membership_claims: %v", diags.Errors())
	}

	if len(claims) == 0 {
		return nil
	}

	// Join multiple claims with comma (Graph API expects single string with comma-separated values)
	claimsStr := claims[0]
	for i := 1; i < len(claims); i++ {
		claimsStr += "," + claims[i]
	}
	requestBody.SetGroupMembershipClaims(&claimsStr)

	return nil
}

// constructOwners sets owners using OData bind for create operations
func constructOwners(ctx context.Context, data *ApplicationResourceModel, requestBody graphmodels.Applicationable) error {
	if data.OwnerUserIds.IsNull() || data.OwnerUserIds.IsUnknown() {
		return nil
	}

	var owners []string
	diags := data.OwnerUserIds.ElementsAs(ctx, &owners, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract owner_user_ids: %v", diags.Errors())
	}

	if len(owners) == 0 {
		return nil
	}

	additionalData := requestBody.GetAdditionalData()
	if additionalData == nil {
		additionalData = make(map[string]any)
	}

	ownerUrls := make([]string, len(owners))
	for i, ownerId := range owners {
		ownerUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", ownerId)
	}

	additionalData["owners@odata.bind"] = ownerUrls
	requestBody.SetAdditionalData(additionalData)
	tflog.Debug(ctx, fmt.Sprintf("Adding %d owners to application", len(owners)))

	return nil
}

// constructSignInAudienceRestrictions builds the sign-in audience restrictions (polymorphic type)
func constructSignInAudienceRestrictions(ctx context.Context, data *SignInAudienceRestrictions) (graphmodels.SignInAudienceRestrictionsBaseable, error) {
	if data == nil {
		return nil, nil
	}

	odataType := data.ODataType.ValueString()

	switch odataType {
	case "#microsoft.graph.allowedTenantsAudience":
		// Create AllowedTenantsAudience with specific fields
		restrictions := graphmodels.NewAllowedTenantsAudience()
		convert.FrameworkToGraphString(data.ODataType, restrictions.SetOdataType)
		convert.FrameworkToGraphBool(data.IsHomeTenantAllowed, restrictions.SetIsHomeTenantAllowed)

		if err := convert.FrameworkToGraphStringSet(ctx, data.AllowedTenantIds, restrictions.SetAllowedTenantIds); err != nil {
			return nil, fmt.Errorf("failed to set allowed_tenant_ids: %w", err)
		}

		return restrictions, nil

	case "#microsoft.graph.unrestrictedAudience":
		// Create UnrestrictedAudience (no additional fields beyond base)
		restrictions := graphmodels.NewUnrestrictedAudience()
		convert.FrameworkToGraphString(data.ODataType, restrictions.SetOdataType)
		return restrictions, nil

	default:
		return nil, fmt.Errorf("unsupported sign_in_audience_restrictions odata_type: %s (must be '#microsoft.graph.allowedTenantsAudience' or '#microsoft.graph.unrestrictedAudience')", odataType)
	}
}

// constructApi builds the API configuration
func constructApi(ctx context.Context, data *ApplicationApi) (graphmodels.ApiApplicationable, error) {
	api := graphmodels.NewApiApplication()

	convert.FrameworkToGraphBool(data.AcceptMappedClaims, api.SetAcceptMappedClaims)
	convert.FrameworkToGraphInt32(data.RequestedAccessTokenVersion, api.SetRequestedAccessTokenVersion)

	if !data.KnownClientApplications.IsNull() && !data.KnownClientApplications.IsUnknown() {
		var clientAppStrings []string
		diags := data.KnownClientApplications.ElementsAs(ctx, &clientAppStrings, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract known_client_applications: %v", diags.Errors())
		}
		clientAppUUIDs := make([]uuid.UUID, 0, len(clientAppStrings))
		for _, appIDStr := range clientAppStrings {
			appUUID, err := uuid.Parse(appIDStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse known client application ID '%s': %w", appIDStr, err)
			}
			clientAppUUIDs = append(clientAppUUIDs, appUUID)
		}
		api.SetKnownClientApplications(clientAppUUIDs)
	}

	if !data.OAuth2PermissionScopes.IsNull() && !data.OAuth2PermissionScopes.IsUnknown() {
		scopes, err := constructOAuth2PermissionScopes(ctx, data.OAuth2PermissionScopes)
		if err != nil {
			return nil, err
		}
		api.SetOauth2PermissionScopes(scopes)
	}

	if !data.PreAuthorizedApplications.IsNull() && !data.PreAuthorizedApplications.IsUnknown() {
		preAuthApps, err := constructPreAuthorizedApplications(ctx, data.PreAuthorizedApplications)
		if err != nil {
			return nil, err
		}
		api.SetPreAuthorizedApplications(preAuthApps)
	}

	return api, nil
}

// constructOAuth2PermissionScopes builds the OAuth2 permission scopes collection
func constructOAuth2PermissionScopes(ctx context.Context, data types.Set) ([]graphmodels.PermissionScopeable, error) {
	var scopes []ApplicationApiPermissionScope
	diags := data.ElementsAs(ctx, &scopes, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract oauth2_permission_scopes: %v", diags.Errors())
	}

	result := make([]graphmodels.PermissionScopeable, 0, len(scopes))
	for _, scope := range scopes {
		permScope := graphmodels.NewPermissionScope()

		if err := convert.FrameworkToGraphUUID(scope.ID, permScope.SetId); err != nil {
			return nil, fmt.Errorf("failed to parse oauth2 permission scope id: %w", err)
		}
		convert.FrameworkToGraphString(scope.AdminConsentDescription, permScope.SetAdminConsentDescription)
		convert.FrameworkToGraphString(scope.AdminConsentDisplayName, permScope.SetAdminConsentDisplayName)
		convert.FrameworkToGraphBool(scope.IsEnabled, permScope.SetIsEnabled)
		convert.FrameworkToGraphString(scope.Type, permScope.SetTypeEscaped)
		convert.FrameworkToGraphString(scope.UserConsentDescription, permScope.SetUserConsentDescription)
		convert.FrameworkToGraphString(scope.UserConsentDisplayName, permScope.SetUserConsentDisplayName)
		convert.FrameworkToGraphString(scope.Value, permScope.SetValue)

		result = append(result, permScope)
	}

	return result, nil
}

// constructPreAuthorizedApplications builds the pre-authorized applications collection
func constructPreAuthorizedApplications(ctx context.Context, data types.Set) ([]graphmodels.PreAuthorizedApplicationable, error) {
	var preAuthApps []ApplicationApiPreAuthorizedApplication
	diags := data.ElementsAs(ctx, &preAuthApps, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract pre_authorized_applications: %v", diags.Errors())
	}

	result := make([]graphmodels.PreAuthorizedApplicationable, 0, len(preAuthApps))
	for _, preAuthApp := range preAuthApps {
		preAuth := graphmodels.NewPreAuthorizedApplication()

		convert.FrameworkToGraphString(preAuthApp.AppId, preAuth.SetAppId)

		if !preAuthApp.DelegatedPermission.IsNull() && !preAuthApp.DelegatedPermission.IsUnknown() {
			var permissionIds []string
			diags := preAuthApp.DelegatedPermission.ElementsAs(ctx, &permissionIds, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract delegated_permission_ids: %v", diags.Errors())
			}
			preAuth.SetPermissionIds(permissionIds)
		}

		result = append(result, preAuth)
	}

	return result, nil
}

// constructAppRoles builds the app roles collection
func constructAppRoles(ctx context.Context, data types.Set) ([]graphmodels.AppRoleable, error) {
	var appRoles []ApplicationAppRole
	diags := data.ElementsAs(ctx, &appRoles, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract app_roles: %v", diags.Errors())
	}

	result := make([]graphmodels.AppRoleable, 0, len(appRoles))
	for _, role := range appRoles {
		appRole := graphmodels.NewAppRole()

		if err := convert.FrameworkToGraphUUID(role.ID, appRole.SetId); err != nil {
			return nil, fmt.Errorf("failed to parse app role id: %w", err)
		}
		convert.FrameworkToGraphString(role.Description, appRole.SetDescription)
		convert.FrameworkToGraphString(role.DisplayName, appRole.SetDisplayName)
		convert.FrameworkToGraphBool(role.IsEnabled, appRole.SetIsEnabled)
		convert.FrameworkToGraphString(role.Origin, appRole.SetOrigin)
		convert.FrameworkToGraphString(role.Value, appRole.SetValue)

		if err := convert.FrameworkToGraphStringSet(ctx, role.AllowedMemberTypes, appRole.SetAllowedMemberTypes); err != nil {
			return nil, fmt.Errorf("failed to set allowed_member_types: %w", err)
		}

		result = append(result, appRole)
	}

	return result, nil
}

// constructInformationalUrl builds the informational URL object
func constructInformationalUrl(data *ApplicationInformationalUrl) graphmodels.InformationalUrlable {
	info := graphmodels.NewInformationalUrl()

	convert.FrameworkToGraphString(data.LogoUrl, info.SetLogoUrl)
	convert.FrameworkToGraphString(data.MarketingUrl, info.SetMarketingUrl)
	convert.FrameworkToGraphString(data.PrivacyStatementUrl, info.SetPrivacyStatementUrl)
	convert.FrameworkToGraphString(data.SupportUrl, info.SetSupportUrl)
	convert.FrameworkToGraphString(data.TermsOfServiceUrl, info.SetTermsOfServiceUrl)

	return info
}

// constructKeyCredentials builds the key credentials collection
func constructKeyCredentials(ctx context.Context, data types.Set) ([]graphmodels.KeyCredentialable, error) {
	var keyCredentials []ApplicationKeyCredential
	diags := data.ElementsAs(ctx, &keyCredentials, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract key_credentials: %v", diags.Errors())
	}

	result := make([]graphmodels.KeyCredentialable, 0, len(keyCredentials))
	for _, keyCred := range keyCredentials {
		keyCredential := graphmodels.NewKeyCredential()

		convert.FrameworkToGraphString(keyCred.DisplayName, keyCredential.SetDisplayName)
		if err := convert.FrameworkToGraphUUID(keyCred.KeyId, keyCredential.SetKeyId); err != nil {
			return nil, fmt.Errorf("failed to parse key credential key_id: %w", err)
		}
		convert.FrameworkToGraphString(keyCred.Type, keyCredential.SetTypeEscaped)
		convert.FrameworkToGraphString(keyCred.Usage, keyCredential.SetUsage)

		if !keyCred.CustomKeyIdentifier.IsNull() && !keyCred.CustomKeyIdentifier.IsUnknown() {
			customKeyId := keyCred.CustomKeyIdentifier.ValueString()
			customKeyIdBytes := []byte(customKeyId)
			keyCredential.SetCustomKeyIdentifier(customKeyIdBytes)
		}

		if !keyCred.Key.IsNull() && !keyCred.Key.IsUnknown() {
			// Key is base64 encoded, decode to byte array
			keyValue := keyCred.Key.ValueString()
			keyBytes, err := base64.StdEncoding.DecodeString(keyValue)
			if err != nil {
				return nil, fmt.Errorf("failed to decode base64 key: %w", err)
			}
			keyCredential.SetKey(keyBytes)
		}

		if err := convert.FrameworkToGraphTime(keyCred.StartDateTime, keyCredential.SetStartDateTime); err != nil {
			return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
		}
		if err := convert.FrameworkToGraphTime(keyCred.EndDateTime, keyCredential.SetEndDateTime); err != nil {
			return nil, fmt.Errorf("failed to parse end_date_time: %w", err)
		}

		result = append(result, keyCredential)
	}

	return result, nil
}

// constructPasswordCredentials builds the password credentials collection
func constructPasswordCredentials(ctx context.Context, data types.Set) ([]graphmodels.PasswordCredentialable, error) {
	var passwordCredentials []ApplicationPasswordCredential
	diags := data.ElementsAs(ctx, &passwordCredentials, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract password_credentials: %v", diags.Errors())
	}

	result := make([]graphmodels.PasswordCredentialable, 0, len(passwordCredentials))
	for _, pwdCred := range passwordCredentials {
		passwordCredential := graphmodels.NewPasswordCredential()

		convert.FrameworkToGraphString(pwdCred.DisplayName, passwordCredential.SetDisplayName)
		if err := convert.FrameworkToGraphUUID(pwdCred.KeyId, passwordCredential.SetKeyId); err != nil {
			return nil, fmt.Errorf("failed to parse password credential key_id: %w", err)
		}

		if !pwdCred.CustomKeyIdentifier.IsNull() && !pwdCred.CustomKeyIdentifier.IsUnknown() {
			customKeyId := pwdCred.CustomKeyIdentifier.ValueString()
			customKeyIdBytes := []byte(customKeyId)
			passwordCredential.SetCustomKeyIdentifier(customKeyIdBytes)
		}

		if err := convert.FrameworkToGraphTime(pwdCred.StartDateTime, passwordCredential.SetStartDateTime); err != nil {
			return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
		}
		if err := convert.FrameworkToGraphTime(pwdCred.EndDateTime, passwordCredential.SetEndDateTime); err != nil {
			return nil, fmt.Errorf("failed to parse end_date_time: %w", err)
		}

		result = append(result, passwordCredential)
	}

	return result, nil
}

// constructOptionalClaims builds the optional claims configuration
func constructOptionalClaims(ctx context.Context, data *ApplicationOptionalClaims) (graphmodels.OptionalClaimsable, error) {
	optionalClaims := graphmodels.NewOptionalClaims()

	if !data.AccessToken.IsNull() && !data.AccessToken.IsUnknown() {
		accessTokenClaims, err := constructOptionalClaimList(ctx, data.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to construct access token claims: %w", err)
		}
		optionalClaims.SetAccessToken(accessTokenClaims)
	}

	if !data.IdToken.IsNull() && !data.IdToken.IsUnknown() {
		idTokenClaims, err := constructOptionalClaimList(ctx, data.IdToken)
		if err != nil {
			return nil, fmt.Errorf("failed to construct id token claims: %w", err)
		}
		optionalClaims.SetIdToken(idTokenClaims)
	}

	if !data.Saml2Token.IsNull() && !data.Saml2Token.IsUnknown() {
		saml2TokenClaims, err := constructOptionalClaimList(ctx, data.Saml2Token)
		if err != nil {
			return nil, fmt.Errorf("failed to construct saml2 token claims: %w", err)
		}
		optionalClaims.SetSaml2Token(saml2TokenClaims)
	}

	return optionalClaims, nil
}

// constructOptionalClaimList builds a list of optional claims
func constructOptionalClaimList(ctx context.Context, data types.Set) ([]graphmodels.OptionalClaimable, error) {
	var claims []ApplicationOptionalClaim
	diags := data.ElementsAs(ctx, &claims, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract optional claims: %v", diags.Errors())
	}

	result := make([]graphmodels.OptionalClaimable, 0, len(claims))
	for _, claim := range claims {
		optionalClaim := graphmodels.NewOptionalClaim()

		convert.FrameworkToGraphString(claim.Name, optionalClaim.SetName)
		convert.FrameworkToGraphString(claim.Source, optionalClaim.SetSource)
		convert.FrameworkToGraphBool(claim.Essential, optionalClaim.SetEssential)

		if err := convert.FrameworkToGraphStringSet(ctx, claim.AdditionalProperties, optionalClaim.SetAdditionalProperties); err != nil {
			return nil, fmt.Errorf("failed to set additional_properties: %w", err)
		}

		result = append(result, optionalClaim)
	}

	return result, nil
}

// constructParentalControlSettings builds the parental control settings
func constructParentalControlSettings(ctx context.Context, data *ApplicationParentalControlSettings) (graphmodels.ParentalControlSettingsable, error) {
	settings := graphmodels.NewParentalControlSettings()

	convert.FrameworkToGraphString(data.LegalAgeGroupRule, settings.SetLegalAgeGroupRule)

	if err := convert.FrameworkToGraphStringSet(ctx, data.CountriesBlockedForMinors, settings.SetCountriesBlockedForMinors); err != nil {
		return nil, fmt.Errorf("failed to set countries_blocked_for_minors: %w", err)
	}

	return settings, nil
}

// constructPublicClient builds the public client configuration
func constructPublicClient(ctx context.Context, data *ApplicationPublicClient) (graphmodels.PublicClientApplicationable, error) {
	publicClient := graphmodels.NewPublicClientApplication()

	if err := convert.FrameworkToGraphStringSet(ctx, data.RedirectUris, publicClient.SetRedirectUris); err != nil {
		return nil, fmt.Errorf("failed to set redirect_uris: %w", err)
	}

	return publicClient, nil
}

// constructRequiredResourceAccess builds the required resource access collection
func constructRequiredResourceAccess(ctx context.Context, data types.Set) ([]graphmodels.RequiredResourceAccessable, error) {
	var requiredResources []ApplicationRequiredResourceAccess
	diags := data.ElementsAs(ctx, &requiredResources, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract required_resource_access: %v", diags.Errors())
	}

	result := make([]graphmodels.RequiredResourceAccessable, 0, len(requiredResources))
	for _, resource := range requiredResources {
		requiredResourceAccess := graphmodels.NewRequiredResourceAccess()

		convert.FrameworkToGraphString(resource.ResourceAppId, requiredResourceAccess.SetResourceAppId)

		if !resource.ResourceAccess.IsNull() && !resource.ResourceAccess.IsUnknown() {
			var resourceAccessList []ApplicationResourceAccess
			diags := resource.ResourceAccess.ElementsAs(ctx, &resourceAccessList, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract resource_access: %v", diags.Errors())
			}

			accessList := make([]graphmodels.ResourceAccessable, 0, len(resourceAccessList))
			for _, access := range resourceAccessList {
				resourceAccess := graphmodels.NewResourceAccess()

				if err := convert.FrameworkToGraphUUID(access.ID, resourceAccess.SetId); err != nil {
					return nil, fmt.Errorf("failed to parse resource access id: %w", err)
				}
				convert.FrameworkToGraphString(access.Type, resourceAccess.SetTypeEscaped)

				accessList = append(accessList, resourceAccess)
			}
			requiredResourceAccess.SetResourceAccess(accessList)
		}

		result = append(result, requiredResourceAccess)
	}

	return result, nil
}

// constructSpa builds the SPA configuration
func constructSpa(ctx context.Context, data *ApplicationSpa) (graphmodels.SpaApplicationable, error) {
	spa := graphmodels.NewSpaApplication()

	if err := convert.FrameworkToGraphStringSet(ctx, data.RedirectUris, spa.SetRedirectUris); err != nil {
		return nil, fmt.Errorf("failed to set redirect_uris: %w", err)
	}

	return spa, nil
}

// constructWeb builds the web application configuration
func constructWeb(ctx context.Context, data *ApplicationWeb) (graphmodels.WebApplicationable, error) {
	web := graphmodels.NewWebApplication()

	convert.FrameworkToGraphString(data.HomePageUrl, web.SetHomePageUrl)
	convert.FrameworkToGraphString(data.LogoutUrl, web.SetLogoutUrl)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RedirectUris, web.SetRedirectUris); err != nil {
		return nil, fmt.Errorf("failed to set redirect_uris: %w", err)
	}

	if !data.ImplicitGrantSettings.IsNull() && !data.ImplicitGrantSettings.IsUnknown() {
		var implicitGrantData ApplicationWebImplicitGrantSettings
		diags := data.ImplicitGrantSettings.As(ctx, &implicitGrantData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract implicit_grant_settings data: %s", diags.Errors()[0].Detail())
		}

		implicitGrant := graphmodels.NewImplicitGrantSettings()
		convert.FrameworkToGraphBool(implicitGrantData.EnableAccessTokenIssuance, implicitGrant.SetEnableAccessTokenIssuance)
		convert.FrameworkToGraphBool(implicitGrantData.EnableIdTokenIssuance, implicitGrant.SetEnableIdTokenIssuance)
		web.SetImplicitGrantSettings(implicitGrant)
	}

	if !data.RedirectUriSettings.IsNull() && !data.RedirectUriSettings.IsUnknown() {
		redirectUriSettings, err := constructRedirectUriSettings(ctx, data.RedirectUriSettings)
		if err != nil {
			return nil, err
		}
		web.SetRedirectUriSettings(redirectUriSettings)
	}

	return web, nil
}

// constructRedirectUriSettings builds the redirect URI settings collection
func constructRedirectUriSettings(ctx context.Context, data types.Set) ([]graphmodels.RedirectUriSettingsable, error) {
	var settings []ApplicationWebRedirectUriSettings
	diags := data.ElementsAs(ctx, &settings, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract redirect_uri_settings: %v", diags.Errors())
	}

	result := make([]graphmodels.RedirectUriSettingsable, 0, len(settings))
	for _, setting := range settings {
		redirectUriSetting := graphmodels.NewRedirectUriSettings()

		convert.FrameworkToGraphString(setting.Uri, redirectUriSetting.SetUri)
		convert.FrameworkToGraphInt32(setting.Index, redirectUriSetting.SetIndex)

		result = append(result, redirectUriSetting)
	}

	return result, nil
}
