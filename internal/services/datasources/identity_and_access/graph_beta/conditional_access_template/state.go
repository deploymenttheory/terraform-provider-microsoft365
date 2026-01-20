package graphBetaConditionalAccessTemplate

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps the remote state to the data source model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.ConditionalAccessTemplateable) ConditionalAccessTemplateDataSourceModel {
	model := ConditionalAccessTemplateDataSourceModel{
		TemplateID:  convert.GraphToFrameworkString(data.GetId()),
		Name:        convert.GraphToFrameworkString(data.GetName()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
		Scenarios:   convert.GraphToFrameworkBitmaskEnumAsSet(ctx, data.GetScenarios()),
	}

	details := data.GetDetails()
	if details != nil {
		model.Details = stateDetails(ctx, details)
	}

	return model
}

// stateDetails maps the details to the data source model
func stateDetails(ctx context.Context, details graphmodels.ConditionalAccessPolicyDetailable) *ConditionalAccessTemplateDetails {
	result := &ConditionalAccessTemplateDetails{}

	if conditions := details.GetConditions(); conditions != nil {
		result.Conditions = stateConditions(ctx, conditions)
	}

	if grantControls := details.GetGrantControls(); grantControls != nil {
		result.GrantControls = stateGrantControls(ctx, grantControls)
	}

	if sessionControls := details.GetSessionControls(); sessionControls != nil {
		result.SessionControls = stateSessionControls(ctx, sessionControls)
	}

	return result
}

// stateConditions maps the conditions to the data source model
func stateConditions(ctx context.Context, conditions graphmodels.ConditionalAccessConditionSetable) *ConditionalAccessConditions {
	result := &ConditionalAccessConditions{}

	result.UserRiskLevels = mapEnumCollectionToList(ctx, conditions.GetUserRiskLevels())
	result.SignInRiskLevels = mapEnumCollectionToList(ctx, conditions.GetSignInRiskLevels())
	result.ClientAppTypes = mapEnumCollectionToList(ctx, conditions.GetClientAppTypes())
	result.ServicePrincipalRiskLevels = mapEnumCollectionToList(ctx, conditions.GetServicePrincipalRiskLevels())
	result.AgentIdRiskLevels = convert.GraphToFrameworkBitmaskEnumAsSet(ctx, conditions.GetAgentIdRiskLevels())
	result.InsiderRiskLevels = convert.GraphToFrameworkBitmaskEnumAsSet(ctx, conditions.GetInsiderRiskLevels())

	if platforms := conditions.GetPlatforms(); platforms != nil {
		result.Platforms = statePlatforms(ctx, platforms)
	}

	if locations := conditions.GetLocations(); locations != nil {
		result.Locations = stateLocations(ctx, locations)
	}

	if devices := conditions.GetDevices(); devices != nil {
		result.Devices = stateDevices(ctx, devices)
	}

	if clientApplications := conditions.GetClientApplications(); clientApplications != nil {
		result.ClientApplications = stateClientApplications(ctx, clientApplications)
	}

	if applications := conditions.GetApplications(); applications != nil {
		result.Applications = stateApplications(ctx, applications)
	}

	if users := conditions.GetUsers(); users != nil {
		result.Users = stateUsers(ctx, users)
	}

	return result
}

// stateApplications maps the applications to the data source model
func stateApplications(ctx context.Context, applications graphmodels.ConditionalAccessApplicationsable) *ConditionalAccessApplications {
	result := &ConditionalAccessApplications{}

	result.IncludeApplications = convert.GraphToFrameworkStringList(applications.GetIncludeApplications())
	result.ExcludeApplications = convert.GraphToFrameworkStringList(applications.GetExcludeApplications())
	result.IncludeUserActions = convert.GraphToFrameworkStringList(applications.GetIncludeUserActions())
	result.IncludeAuthenticationContextClassReferences = convert.GraphToFrameworkStringList(applications.GetIncludeAuthenticationContextClassReferences())

	return result
}

// stateUsers maps the users to the data source model
func stateUsers(ctx context.Context, users graphmodels.ConditionalAccessUsersable) *ConditionalAccessUsers {
	result := &ConditionalAccessUsers{}

	result.IncludeUsers = convert.GraphToFrameworkStringList(users.GetIncludeUsers())
	result.ExcludeUsers = convert.GraphToFrameworkStringList(users.GetExcludeUsers())
	result.IncludeGroups = convert.GraphToFrameworkStringList(users.GetIncludeGroups())
	result.ExcludeGroups = convert.GraphToFrameworkStringList(users.GetExcludeGroups())
	result.IncludeRoles = convert.GraphToFrameworkStringList(users.GetIncludeRoles())
	result.ExcludeRoles = convert.GraphToFrameworkStringList(users.GetExcludeRoles())

	if includeGuestsOrExternalUsers := users.GetIncludeGuestsOrExternalUsers(); includeGuestsOrExternalUsers != nil {
		result.IncludeGuestsOrExternalUsers = stateGuestsOrExternalUsers(ctx, includeGuestsOrExternalUsers)
	}

	if excludeGuestsOrExternalUsers := users.GetExcludeGuestsOrExternalUsers(); excludeGuestsOrExternalUsers != nil {
		result.ExcludeGuestsOrExternalUsers = stateGuestsOrExternalUsers(ctx, excludeGuestsOrExternalUsers)
	}

	return result
}

// stateGuestsOrExternalUsers maps the guests or external users to the data source model
func stateGuestsOrExternalUsers(ctx context.Context, guestsOrExternalUsers graphmodels.ConditionalAccessGuestsOrExternalUsersable) *ConditionalAccessGuestsOrExternalUsers {
	result := &ConditionalAccessGuestsOrExternalUsers{}

	result.GuestOrExternalUserTypes = convert.GraphToFrameworkBitmaskEnumAsSet(ctx, guestsOrExternalUsers.GetGuestOrExternalUserTypes())

	if externalTenants := guestsOrExternalUsers.GetExternalTenants(); externalTenants != nil {
		result.ExternalTenants = stateExternalTenants(ctx, externalTenants)
	}

	return result
}

// stateExternalTenants maps the external tenants to the data source model
func stateExternalTenants(ctx context.Context, externalTenants graphmodels.ConditionalAccessExternalTenantsable) *ConditionalAccessExternalTenants {
	result := &ConditionalAccessExternalTenants{}

	result.MembershipKind = convert.GraphToFrameworkEnum(externalTenants.GetMembershipKind())

	return result
}

// statePlatforms maps the platforms to the data source model
func statePlatforms(ctx context.Context, platforms graphmodels.ConditionalAccessPlatformsable) *ConditionalAccessPlatforms {
	result := &ConditionalAccessPlatforms{}

	result.IncludePlatforms = mapEnumCollectionToList(ctx, platforms.GetIncludePlatforms())
	result.ExcludePlatforms = mapEnumCollectionToList(ctx, platforms.GetExcludePlatforms())

	return result
}

// stateLocations maps the locations to the data source model
func stateLocations(ctx context.Context, locations graphmodels.ConditionalAccessLocationsable) *ConditionalAccessLocations {
	result := &ConditionalAccessLocations{}

	result.IncludeLocations = convert.GraphToFrameworkStringList(locations.GetIncludeLocations())
	result.ExcludeLocations = convert.GraphToFrameworkStringList(locations.GetExcludeLocations())

	return result
}

// stateDevices maps the devices to the data source model
func stateDevices(ctx context.Context, devices graphmodels.ConditionalAccessDevicesable) *ConditionalAccessDevices {
	result := &ConditionalAccessDevices{}

	result.IncludeDeviceStates = convert.GraphToFrameworkStringList(devices.GetIncludeDeviceStates())
	result.ExcludeDeviceStates = convert.GraphToFrameworkStringList(devices.GetExcludeDeviceStates())
	result.IncludeDevices = convert.GraphToFrameworkStringList(devices.GetIncludeDevices())
	result.ExcludeDevices = convert.GraphToFrameworkStringList(devices.GetExcludeDevices())

	if deviceFilter := devices.GetDeviceFilter(); deviceFilter != nil {
		result.DeviceFilter = stateDeviceFilter(ctx, deviceFilter)
	}

	return result
}

// stateDeviceFilter maps the device filter to the data source model
func stateDeviceFilter(ctx context.Context, deviceFilter graphmodels.ConditionalAccessFilterable) *ConditionalAccessFilter {
	result := &ConditionalAccessFilter{}

	result.Mode = convert.GraphToFrameworkEnum(deviceFilter.GetMode())
	result.Rule = convert.GraphToFrameworkString(deviceFilter.GetRule())

	return result
}

// stateClientApplications maps the client applications to the data source model
func stateClientApplications(ctx context.Context, clientApps graphmodels.ConditionalAccessClientApplicationsable) *ConditionalAccessClientApplications {
	result := &ConditionalAccessClientApplications{}

	result.IncludeServicePrincipals = convert.GraphToFrameworkStringList(clientApps.GetIncludeServicePrincipals())
	result.IncludeAgentIdServicePrincipals = convert.GraphToFrameworkStringList(clientApps.GetIncludeAgentIdServicePrincipals())
	result.ExcludeServicePrincipals = convert.GraphToFrameworkStringList(clientApps.GetExcludeServicePrincipals())
	result.ExcludeAgentIdServicePrincipals = convert.GraphToFrameworkStringList(clientApps.GetExcludeAgentIdServicePrincipals())

	return result
}

// stateGrantControls maps the grant controls to the data source model
func stateGrantControls(ctx context.Context, grantControls graphmodels.ConditionalAccessGrantControlsable) *ConditionalAccessGrantControls {
	result := &ConditionalAccessGrantControls{}

	result.Operator = convert.GraphToFrameworkString(grantControls.GetOperator())
	result.BuiltInControls = mapEnumCollectionToList(ctx, grantControls.GetBuiltInControls())
	result.CustomAuthenticationFactors = convert.GraphToFrameworkStringList(grantControls.GetCustomAuthenticationFactors())
	result.TermsOfUse = convert.GraphToFrameworkStringList(grantControls.GetTermsOfUse())

	if authStrength := grantControls.GetAuthenticationStrength(); authStrength != nil {
		result.AuthenticationStrength = stateAuthenticationStrength(ctx, authStrength)
	}

	return result
}

// stateAuthenticationStrength maps the authentication strength to the data source model
func stateAuthenticationStrength(ctx context.Context, authStrength graphmodels.AuthenticationStrengthPolicyable) *ConditionalAccessAuthenticationStrength {
	result := &ConditionalAccessAuthenticationStrength{}

	result.ID = convert.GraphToFrameworkString(authStrength.GetId())
	result.CreatedDateTime = convert.GraphToFrameworkTime(authStrength.GetCreatedDateTime())
	result.ModifiedDateTime = convert.GraphToFrameworkTime(authStrength.GetModifiedDateTime())
	result.DisplayName = convert.GraphToFrameworkString(authStrength.GetDisplayName())
	result.Description = convert.GraphToFrameworkString(authStrength.GetDescription())
	result.PolicyType = convert.GraphToFrameworkEnum(authStrength.GetPolicyType())
	result.RequirementsSatisfied = convert.GraphToFrameworkEnum(authStrength.GetRequirementsSatisfied())
	result.AllowedCombinations = mapEnumCollectionToList(ctx, authStrength.GetAllowedCombinations())

	return result
}

// stateSessionControls maps the session controls to the data source model
func stateSessionControls(ctx context.Context, sessionControls graphmodels.ConditionalAccessSessionControlsable) *ConditionalAccessSessionControls {
	result := &ConditionalAccessSessionControls{}

	result.DisableResilienceDefaults = convert.GraphToFrameworkBool(sessionControls.GetDisableResilienceDefaults())

	if cloudAppSecurity := sessionControls.GetCloudAppSecurity(); cloudAppSecurity != nil {
		result.CloudAppSecurity = stateCloudAppSecuritySessionControl(ctx, cloudAppSecurity)
	}

	if continuousAccessEvaluation := sessionControls.GetContinuousAccessEvaluation(); continuousAccessEvaluation != nil {
		result.ContinuousAccessEvaluation = stateContinuousAccessEvaluationSessionControl(ctx, continuousAccessEvaluation)
	}

	if secureSignInSession := sessionControls.GetSecureSignInSession(); secureSignInSession != nil {
		result.SecureSignInSession = stateSecureSignInSessionControl(ctx, secureSignInSession)
	}

	if globalSecureAccessFilteringProfile := sessionControls.GetGlobalSecureAccessFilteringProfile(); globalSecureAccessFilteringProfile != nil {
		result.GlobalSecureAccessFilteringProfile = stateGlobalSecureAccessFilteringProfile(ctx, globalSecureAccessFilteringProfile)
	}

	if appEnforcedRestrictions := sessionControls.GetApplicationEnforcedRestrictions(); appEnforcedRestrictions != nil {
		result.ApplicationEnforcedRestrictions = stateApplicationEnforcedRestrictions(ctx, appEnforcedRestrictions)
	}

	if persistentBrowser := sessionControls.GetPersistentBrowser(); persistentBrowser != nil {
		result.PersistentBrowser = statePersistentBrowser(ctx, persistentBrowser)
	}

	if signInFrequency := sessionControls.GetSignInFrequency(); signInFrequency != nil {
		result.SignInFrequency = stateSignInFrequency(ctx, signInFrequency)
	}

	return result
}

func stateCloudAppSecuritySessionControl(ctx context.Context, control graphmodels.CloudAppSecuritySessionControlable) *ConditionalAccessCloudAppSecuritySessionControl {
	result := &ConditionalAccessCloudAppSecuritySessionControl{}
	result.IsEnabled = convert.GraphToFrameworkBool(control.GetIsEnabled())
	result.CloudAppSecurityType = convert.GraphToFrameworkEnum(control.GetCloudAppSecurityType())
	return result
}

func stateContinuousAccessEvaluationSessionControl(ctx context.Context, control graphmodels.ContinuousAccessEvaluationSessionControlable) *ConditionalAccessContinuousAccessEvaluationSessionControl {
	result := &ConditionalAccessContinuousAccessEvaluationSessionControl{}
	result.Mode = convert.GraphToFrameworkEnum(control.GetMode())
	return result
}

func stateSecureSignInSessionControl(ctx context.Context, control graphmodels.SecureSignInSessionControlable) *ConditionalAccessSecureSignInSessionControl {
	result := &ConditionalAccessSecureSignInSessionControl{}
	result.IsEnabled = convert.GraphToFrameworkBool(control.GetIsEnabled())
	return result
}

func stateGlobalSecureAccessFilteringProfile(ctx context.Context, control graphmodels.GlobalSecureAccessFilteringProfileSessionControlable) *ConditionalAccessGlobalSecureAccessFilteringProfile {
	result := &ConditionalAccessGlobalSecureAccessFilteringProfile{}
	result.IsEnabled = convert.GraphToFrameworkBool(control.GetIsEnabled())
	result.ProfileId = convert.GraphToFrameworkString(control.GetProfileId())
	return result
}

// stateApplicationEnforcedRestrictions maps the application enforced restrictions to the data source model
func stateApplicationEnforcedRestrictions(ctx context.Context, appRestrictions graphmodels.ApplicationEnforcedRestrictionsSessionControlable) *ConditionalAccessApplicationEnforcedRestrictions {
	result := &ConditionalAccessApplicationEnforcedRestrictions{}

	result.IsEnabled = convert.GraphToFrameworkBool(appRestrictions.GetIsEnabled())

	return result
}

// statePersistentBrowser maps the persistent browser to the data source model
func statePersistentBrowser(ctx context.Context, persistentBrowser graphmodels.PersistentBrowserSessionControlable) *ConditionalAccessPersistentBrowser {
	result := &ConditionalAccessPersistentBrowser{}

	result.Mode = convert.GraphToFrameworkEnum(persistentBrowser.GetMode())
	result.IsEnabled = convert.GraphToFrameworkBool(persistentBrowser.GetIsEnabled())

	return result
}

func stateSignInFrequency(ctx context.Context, signInFreq graphmodels.SignInFrequencySessionControlable) *ConditionalAccessSignInFrequency {
	result := &ConditionalAccessSignInFrequency{}

	result.Value = convert.GraphToFrameworkInt32AsInt64(signInFreq.GetValue())
	result.Type = convert.GraphToFrameworkEnum(signInFreq.GetTypeEscaped())
	result.AuthenticationType = convert.GraphToFrameworkEnum(signInFreq.GetAuthenticationType())
	result.FrequencyInterval = convert.GraphToFrameworkEnum(signInFreq.GetFrequencyInterval())
	result.IsEnabled = convert.GraphToFrameworkBool(signInFreq.GetIsEnabled())

	return result
}

// mapEnumCollectionToList converts an enum slice to a Terraform Framework list of strings
// Preserves empty slices as empty lists to maintain Terraform state consistency
func mapEnumCollectionToList[T fmt.Stringer](ctx context.Context, enums []T) types.List {
	values := make([]string, len(enums))
	for i, enum := range enums {
		values[i] = enum.String()
	}

	elemType := types.StringType
	if len(values) == 0 {
		return types.ListValueMust(elemType, []attr.Value{})
	}

	elements := make([]attr.Value, len(values))
	for i, v := range values {
		elements[i] = types.StringValue(v)
	}

	return types.ListValueMust(elemType, elements)
}
