package graphBetaApplication

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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
	}

	// Map AppRoles
	if appRoles := remoteResource.GetAppRoles(); appRoles != nil {
		data.AppRoles = mapAppRolesToTerraform(ctx, appRoles)
	}

	// Map Info (InformationalUrl)
	if info := remoteResource.GetInfo(); info != nil {
		data.Info = mapInfoToTerraform(info)
	}

	// Map KeyCredentials
	if keyCredentials := remoteResource.GetKeyCredentials(); keyCredentials != nil {
		data.KeyCredentials = mapKeyCredentialsToTerraform(ctx, keyCredentials)
	}

	// Map PasswordCredentials
	if passwordCredentials := remoteResource.GetPasswordCredentials(); passwordCredentials != nil {
		data.PasswordCredentials = mapPasswordCredentialsToTerraform(ctx, passwordCredentials)
	}

	// Map OptionalClaims
	if optionalClaims := remoteResource.GetOptionalClaims(); optionalClaims != nil {
		data.OptionalClaims = mapOptionalClaimsToTerraform(ctx, optionalClaims)
	}

	// Map ParentalControlSettings
	if parentalControl := remoteResource.GetParentalControlSettings(); parentalControl != nil {
		data.ParentalControlSettings = mapParentalControlSettingsToTerraform(ctx, parentalControl)
	}

	// Map PublicClient
	if publicClient := remoteResource.GetPublicClient(); publicClient != nil {
		data.PublicClient = mapPublicClientToTerraform(ctx, publicClient)
	}

	// Map RequiredResourceAccess
	if requiredResourceAccess := remoteResource.GetRequiredResourceAccess(); requiredResourceAccess != nil {
		data.RequiredResourceAccess = mapRequiredResourceAccessToTerraform(ctx, requiredResourceAccess)
	}

	// Map Spa
	if spa := remoteResource.GetSpa(); spa != nil {
		data.Spa = mapSpaToTerraform(ctx, spa)
	}

	// Map Web
	if web := remoteResource.GetWeb(); web != nil {
		data.Web = mapWebToTerraform(ctx, web)
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
func mapApiToTerraform(ctx context.Context, api graphmodels.ApiApplicationable) *ApplicationApi {
	if api == nil {
		return nil
	}

	result := &ApplicationApi{
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
	}

	// Map OAuth2PermissionScopes
	if scopes := api.GetOauth2PermissionScopes(); scopes != nil {
		result.OAuth2PermissionScopes = mapOAuth2PermissionScopesToTerraform(ctx, scopes)
	}

	// Map PreAuthorizedApplications
	if preAuthApps := api.GetPreAuthorizedApplications(); preAuthApps != nil {
		result.PreAuthorizedApplications = mapPreAuthorizedApplicationsToTerraform(ctx, preAuthApps)
	}

	return result
}

func mapOAuth2PermissionScopesToTerraform(ctx context.Context, scopes []graphmodels.PermissionScopeable) types.Set {
	// Implementation for mapping OAuth2 permission scopes
	// This is a placeholder - implement based on your needs
	return types.SetNull(types.ObjectType{})
}

func mapPreAuthorizedApplicationsToTerraform(ctx context.Context, preAuthApps []graphmodels.PreAuthorizedApplicationable) types.Set {
	// Implementation for mapping pre-authorized applications
	// This is a placeholder - implement based on your needs
	return types.SetNull(types.ObjectType{})
}

func mapAppRolesToTerraform(ctx context.Context, appRoles []graphmodels.AppRoleable) types.Set {
	// Implementation for mapping app roles
	// This is a placeholder - implement based on your needs
	return types.SetNull(types.ObjectType{})
}

func mapInfoToTerraform(info graphmodels.InformationalUrlable) *ApplicationInformationalUrl {
	if info == nil {
		return nil
	}

	return &ApplicationInformationalUrl{
		LogoUrl:             convert.GraphToFrameworkString(info.GetLogoUrl()),
		MarketingUrl:        convert.GraphToFrameworkString(info.GetMarketingUrl()),
		PrivacyStatementUrl: convert.GraphToFrameworkString(info.GetPrivacyStatementUrl()),
		SupportUrl:          convert.GraphToFrameworkString(info.GetSupportUrl()),
		TermsOfServiceUrl:   convert.GraphToFrameworkString(info.GetTermsOfServiceUrl()),
	}
}

func mapKeyCredentialsToTerraform(ctx context.Context, keyCredentials []graphmodels.KeyCredentialable) types.Set {
	// Implementation for mapping key credentials
	// This is a placeholder - implement based on your needs
	return types.SetNull(types.ObjectType{})
}

func mapPasswordCredentialsToTerraform(ctx context.Context, passwordCredentials []graphmodels.PasswordCredentialable) types.Set {
	// Implementation for mapping password credentials
	// This is a placeholder - implement based on your needs
	return types.SetNull(types.ObjectType{})
}

func mapOptionalClaimsToTerraform(ctx context.Context, optionalClaims graphmodels.OptionalClaimsable) *ApplicationOptionalClaims {
	// Implementation for mapping optional claims
	// This is a placeholder - implement based on your needs
	return nil
}

func mapParentalControlSettingsToTerraform(ctx context.Context, parentalControl graphmodels.ParentalControlSettingsable) *ApplicationParentalControlSettings {
	if parentalControl == nil {
		return nil
	}

	return &ApplicationParentalControlSettings{
		CountriesBlockedForMinors: convert.GraphToFrameworkStringSet(ctx, parentalControl.GetCountriesBlockedForMinors()),
		LegalAgeGroupRule:         convert.GraphToFrameworkString(parentalControl.GetLegalAgeGroupRule()),
	}
}

func mapPublicClientToTerraform(ctx context.Context, publicClient graphmodels.PublicClientApplicationable) *ApplicationPublicClient {
	if publicClient == nil {
		return nil
	}

	return &ApplicationPublicClient{
		RedirectUris: convert.GraphToFrameworkStringSet(ctx, publicClient.GetRedirectUris()),
	}
}

func mapRequiredResourceAccessToTerraform(ctx context.Context, requiredResourceAccess []graphmodels.RequiredResourceAccessable) types.Set {
	// Implementation for mapping required resource access
	// This is a placeholder - implement based on your needs
	return types.SetNull(types.ObjectType{})
}

func mapSpaToTerraform(ctx context.Context, spa graphmodels.SpaApplicationable) *ApplicationSpa {
	if spa == nil {
		return nil
	}

	return &ApplicationSpa{
		RedirectUris: convert.GraphToFrameworkStringSet(ctx, spa.GetRedirectUris()),
	}
}

func mapWebToTerraform(ctx context.Context, web graphmodels.WebApplicationable) *ApplicationWeb {
	if web == nil {
		return nil
	}

	result := &ApplicationWeb{
		HomePageUrl:  convert.GraphToFrameworkString(web.GetHomePageUrl()),
		LogoutUrl:    convert.GraphToFrameworkString(web.GetLogoutUrl()),
		RedirectUris: convert.GraphToFrameworkStringSet(ctx, web.GetRedirectUris()),
	}

	// Map ImplicitGrantSettings
	if implicitGrant := web.GetImplicitGrantSettings(); implicitGrant != nil {
		result.ImplicitGrantSettings = &ApplicationWebImplicitGrantSettings{
			EnableAccessTokenIssuance: convert.GraphToFrameworkBool(implicitGrant.GetEnableAccessTokenIssuance()),
			EnableIdTokenIssuance:     convert.GraphToFrameworkBool(implicitGrant.GetEnableIdTokenIssuance()),
		}
	}

	return result
}
