package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model for stating.
// It populates the ConditionalAccessPolicyResourceModel with data from the ConditionalAccessPolicy.
func mapRemoteStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource models.ConditionalAccessPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.ModifiedDateTime = state.TimeToString(remoteResource.GetModifiedDateTime())

	if state := remoteResource.GetState(); state != nil {
		data.State = types.StringValue(state.String())
	} else {
		data.State = types.StringNull()
	}

	// Map Conditions
	if conditions := remoteResource.GetConditions(); conditions != nil {
		tflog.Debug(ctx, "Mapping conditions", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Conditions = mapConditions(ctx, conditions)
	}

	// Map Grant Controls
	if grantControls := remoteResource.GetGrantControls(); grantControls != nil {
		tflog.Debug(ctx, "Mapping grant controls", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.GrantControls = mapGrantControls(ctx, grantControls)
	}

	// Map Session Controls
	if sessionControls := remoteResource.GetSessionControls(); sessionControls != nil {
		tflog.Debug(ctx, "Mapping session controls", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.SessionControls = mapSessionControls(ctx, sessionControls)
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

// Map Conditions

func mapConditions(ctx context.Context, conditions models.ConditionalAccessConditionSetable) *ConditionalAccessConditionsModel {
	if conditions == nil {
		tflog.Debug(ctx, "Conditions are nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map conditions")

	result := &ConditionalAccessConditionsModel{
		Applications:               mapApplications(ctx, conditions.GetApplications()),
		Users:                      mapUsers(ctx, conditions.GetUsers()),
		ClientApplications:         mapClientApplications(ctx, conditions.GetClientApplications()),
		ClientAppTypes:             state.EnumSliceToTypeStringSlice(conditions.GetClientAppTypes()),
		Devices:                    mapDevices(ctx, conditions.GetDevices()),
		DeviceStates:               mapDeviceStates(ctx, conditions.GetDeviceStates()),
		Locations:                  mapLocations(ctx, conditions.GetLocations()),
		Platforms:                  mapPlatforms(ctx, conditions.GetPlatforms()),
		ServicePrincipalRiskLevels: state.EnumSliceToTypeStringSlice(conditions.GetServicePrincipalRiskLevels()),
		SignInRiskLevels:           state.EnumSliceToTypeStringSlice(conditions.GetSignInRiskLevels()),
		UserRiskLevels:             state.EnumSliceToTypeStringSlice(conditions.GetUserRiskLevels()),
		AuthenticationFlows:        mapAuthenticationFlows(ctx, conditions.GetAuthenticationFlows()),
		InsiderRiskLevels:          types.StringValue(conditions.GetInsiderRiskLevels().String()),
	}

	tflog.Debug(ctx, "Finished mapping conditions", map[string]interface{}{
		"hasApplications":       result.Applications != nil,
		"hasUsers":              result.Users != nil,
		"hasClientApplications": result.ClientApplications != nil,
		"hasDevices":            result.Devices != nil,
		"hasLocations":          result.Locations != nil,
		"hasPlatforms":          result.Platforms != nil,
	})

	return result
}

func mapApplications(ctx context.Context, apps models.ConditionalAccessApplicationsable) *ConditionalAccessApplicationsModel {
	if apps == nil {
		return nil
	}

	tflog.Debug(ctx, "Mapping applications")

	result := &ConditionalAccessApplicationsModel{
		IncludeApplications: state.SliceToTypeStringSlice(apps.GetIncludeApplications()),
		ExcludeApplications: state.SliceToTypeStringSlice(apps.GetExcludeApplications()),
		IncludeUserActions:  state.SliceToTypeStringSlice(apps.GetIncludeUserActions()),
		ApplicationFilter:   mapFilter(ctx, apps.GetApplicationFilter()),
	}

	tflog.Debug(ctx, "Finished mapping applications", map[string]interface{}{
		"includeCount": len(result.IncludeApplications),
		"excludeCount": len(result.ExcludeApplications),
	})

	return result
}

func mapUsers(ctx context.Context, users models.ConditionalAccessUsersable) *ConditionalAccessUsersModel {
	if users == nil {
		tflog.Debug(ctx, "Users model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map users")
	result := &ConditionalAccessUsersModel{
		ExcludeGroups:                state.SliceToTypeStringSlice(users.GetExcludeGroups()),
		ExcludeRoles:                 state.SliceToTypeStringSlice(users.GetExcludeRoles()),
		ExcludeUsers:                 state.SliceToTypeStringSlice(users.GetExcludeUsers()),
		IncludeGroups:                state.SliceToTypeStringSlice(users.GetIncludeGroups()),
		IncludeRoles:                 state.SliceToTypeStringSlice(users.GetIncludeRoles()),
		IncludeUsers:                 state.SliceToTypeStringSlice(users.GetIncludeUsers()),
		ExcludeGuestsOrExternalUsers: mapGuestsOrExternalUsers(ctx, users.GetExcludeGuestsOrExternalUsers()),
		IncludeGuestsOrExternalUsers: mapGuestsOrExternalUsers(ctx, users.GetIncludeGuestsOrExternalUsers()),
	}
	tflog.Debug(ctx, "Finished mapping users", map[string]interface{}{
		"excludeGroupsCount": len(result.ExcludeGroups),
		"includeGroupsCount": len(result.IncludeGroups),
	})
	return result
}

func mapAuthenticationFlows(ctx context.Context, authFlows models.ConditionalAccessAuthenticationFlowsable) *ConditionalAccessAuthenticationFlowsModel {
	if authFlows == nil {
		tflog.Debug(ctx, "Authentication flows model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map authentication flows")

	var transferMethodsString string
	if authFlows.GetTransferMethods() != nil {
		transferMethodsString = authFlows.GetTransferMethods().String()
	} else {
		transferMethodsString = ""
	}

	result := &ConditionalAccessAuthenticationFlowsModel{
		TransferMethods: types.StringValue(transferMethodsString),
	}

	tflog.Debug(ctx, "Finished mapping authentication flows", map[string]interface{}{
		"transferMethods": result.TransferMethods.ValueString(),
	})

	return result
}

func mapGuestsOrExternalUsers(ctx context.Context, guestsOrExternalUsers models.ConditionalAccessGuestsOrExternalUsersable) *ConditionalAccessGuestsOrExternalUsersModel {
	if guestsOrExternalUsers == nil {
		tflog.Debug(ctx, "Guests or External Users model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map guests or external users")

	var guestOrExternalUserTypesString string
	if guestsOrExternalUsers.GetGuestOrExternalUserTypes() != nil {
		guestOrExternalUserTypesString = guestsOrExternalUsers.GetGuestOrExternalUserTypes().String()
	} else {
		guestOrExternalUserTypesString = ""
	}

	result := &ConditionalAccessGuestsOrExternalUsersModel{
		ExternalTenants:          mapExternalTenants(ctx, guestsOrExternalUsers.GetExternalTenants()),
		GuestOrExternalUserTypes: types.StringValue(guestOrExternalUserTypesString),
	}

	tflog.Debug(ctx, "Finished mapping guests or external users", map[string]interface{}{
		"guestOrExternalUserTypes": result.GuestOrExternalUserTypes.ValueString(),
		"hasExternalTenants":       result.ExternalTenants != nil,
	})

	return result
}

func mapExternalTenants(ctx context.Context, externalTenants models.ConditionalAccessExternalTenantsable) *ConditionalAccessExternalTenantsModel {
	if externalTenants == nil {
		tflog.Debug(ctx, "External Tenants model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map external tenants")

	var membershipKindString string
	if externalTenants.GetMembershipKind() != nil {
		membershipKindString = externalTenants.GetMembershipKind().String()
	} else {
		membershipKindString = ""
	}

	result := &ConditionalAccessExternalTenantsModel{
		MembershipKind: types.StringValue(membershipKindString),
	}

	tflog.Debug(ctx, "Finished mapping external tenants", map[string]interface{}{
		"membershipKind": result.MembershipKind.ValueString(),
	})

	return result
}

func mapClientApplications(ctx context.Context, clientApps models.ConditionalAccessClientApplicationsable) *ConditionalAccessClientApplicationsModel {
	if clientApps == nil {
		tflog.Debug(ctx, "Client applications are nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map client applications")

	result := &ConditionalAccessClientApplicationsModel{
		ExcludeServicePrincipals: state.SliceToTypeStringSlice(clientApps.GetExcludeServicePrincipals()),
		IncludeServicePrincipals: state.SliceToTypeStringSlice(clientApps.GetIncludeServicePrincipals()),
		ServicePrincipalFilter:   mapFilter(ctx, clientApps.GetServicePrincipalFilter()),
	}

	tflog.Debug(ctx, "Finished mapping client applications", map[string]interface{}{
		"excludeServicePrincipalsCount": len(result.ExcludeServicePrincipals),
		"includeServicePrincipalsCount": len(result.IncludeServicePrincipals),
		"hasServicePrincipalFilter":     result.ServicePrincipalFilter != nil,
	})

	return result
}

func mapDevices(ctx context.Context, devices models.ConditionalAccessDevicesable) *ConditionalAccessDevicesModel {
	if devices == nil {
		tflog.Debug(ctx, "Devices model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map devices")

	result := &ConditionalAccessDevicesModel{
		IncludeDevices: state.SliceToTypeStringSlice(devices.GetIncludeDevices()),
		ExcludeDevices: state.SliceToTypeStringSlice(devices.GetExcludeDevices()),
		IncludeStates:  state.SliceToTypeStringSlice(devices.GetIncludeDeviceStates()),
		ExcludeStates:  state.SliceToTypeStringSlice(devices.GetExcludeDeviceStates()),
		DeviceFilter:   mapFilter(ctx, devices.GetDeviceFilter()),
	}

	tflog.Debug(ctx, "Finished mapping devices", map[string]interface{}{
		"includeDevicesCount": len(result.IncludeDevices),
		"excludeDevicesCount": len(result.ExcludeDevices),
		"includeStatesCount":  len(result.IncludeStates),
		"excludeStatesCount":  len(result.ExcludeStates),
		"hasDeviceFilter":     result.DeviceFilter != nil,
	})

	return result
}

func mapDeviceStates(ctx context.Context, deviceStates models.ConditionalAccessDeviceStatesable) *ConditionalAccessDeviceStatesModel {
	if deviceStates == nil {
		tflog.Debug(ctx, "Device states model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map device states")

	result := &ConditionalAccessDeviceStatesModel{
		IncludeStates: state.SliceToTypeStringSlice(deviceStates.GetIncludeStates()),
		ExcludeStates: state.SliceToTypeStringSlice(deviceStates.GetExcludeStates()),
	}

	tflog.Debug(ctx, "Finished mapping device states", map[string]interface{}{
		"includeStatesCount": len(result.IncludeStates),
		"excludeStatesCount": len(result.ExcludeStates),
	})

	return result
}

func mapLocations(ctx context.Context, locations models.ConditionalAccessLocationsable) *ConditionalAccessLocationsModel {
	if locations == nil {
		tflog.Debug(ctx, "Locations model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map locations")

	result := &ConditionalAccessLocationsModel{
		ExcludeLocations: state.SliceToTypeStringSlice(locations.GetExcludeLocations()),
		IncludeLocations: state.SliceToTypeStringSlice(locations.GetIncludeLocations()),
	}

	tflog.Debug(ctx, "Finished mapping locations", map[string]interface{}{
		"excludeLocationsCount": len(result.ExcludeLocations),
		"includeLocationsCount": len(result.IncludeLocations),
	})

	return result
}

func mapPlatforms(ctx context.Context, platforms models.ConditionalAccessPlatformsable) *ConditionalAccessPlatformsModel {
	if platforms == nil {
		tflog.Debug(ctx, "Platforms model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map platforms")

	result := &ConditionalAccessPlatformsModel{
		ExcludePlatforms: state.EnumSliceToTypeStringSlice(platforms.GetExcludePlatforms()),
		IncludePlatforms: state.EnumSliceToTypeStringSlice(platforms.GetIncludePlatforms()),
	}

	tflog.Debug(ctx, "Finished mapping platforms", map[string]interface{}{
		"excludePlatformsCount": len(result.ExcludePlatforms),
		"includePlatformsCount": len(result.IncludePlatforms),
	})

	return result
}

func mapFilter(ctx context.Context, filter models.ConditionalAccessFilterable) *ConditionalAccessFilterModel {
	if filter == nil {
		tflog.Debug(ctx, "Filter model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map filter")

	result := &ConditionalAccessFilterModel{
		Mode: types.StringValue(filter.GetMode().String()),
		Rule: types.StringValue(*filter.GetRule()),
	}

	tflog.Debug(ctx, "Finished mapping filter", map[string]interface{}{
		"mode": result.Mode.ValueString(),
		"rule": result.Rule.ValueString(),
	})

	return result
}

// Map Grant Controls
func mapGrantControls(ctx context.Context, grantControls models.ConditionalAccessGrantControlsable) *ConditionalAccessGrantControlsModel {
	if grantControls == nil {
		tflog.Debug(ctx, "Grant controls model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map grant controls")

	result := &ConditionalAccessGrantControlsModel{}

	if operator := grantControls.GetOperator(); operator != nil {
		result.Operator = types.StringValue(*operator)
	}

	if builtInControls := grantControls.GetBuiltInControls(); builtInControls != nil {
		result.BuiltInControls = state.EnumSliceToTypeStringSlice(builtInControls)
	}

	if customAuthenticationFactors := grantControls.GetCustomAuthenticationFactors(); customAuthenticationFactors != nil {
		result.CustomAuthenticationFactors = state.SliceToTypeStringSlice(customAuthenticationFactors)
	}

	if termsOfUse := grantControls.GetTermsOfUse(); termsOfUse != nil {
		result.TermsOfUse = state.SliceToTypeStringSlice(termsOfUse)
	}

	if authenticationStrength := grantControls.GetAuthenticationStrength(); authenticationStrength != nil {
		result.AuthenticationStrength = &AuthenticationStrengthPolicyModel{
			DisplayName: types.StringValue(state.StringPtrToString(authenticationStrength.GetDisplayName())),
			Description: types.StringValue(state.StringPtrToString(authenticationStrength.GetDescription())),
			PolicyType:  types.StringValue(authenticationStrength.GetPolicyType().String()),
			RequirementsSatisfied: types.StringValue(
				authenticationStrength.GetRequirementsSatisfied().String(),
			),
			AllowedCombinations: state.EnumSliceToTypeStringSlice(authenticationStrength.GetAllowedCombinations()),
		}
	}

	tflog.Debug(ctx, "Finished mapping grant controls", map[string]interface{}{
		"operator":                    result.Operator.ValueString(),
		"builtInControlsCount":        len(result.BuiltInControls),
		"customAuthenticationFactors": len(result.CustomAuthenticationFactors),
		"termsOfUseCount":             len(result.TermsOfUse),
		"hasAuthenticationStrength":   result.AuthenticationStrength != nil,
	})

	return result
}

// mapSessionControls
func mapSessionControls(ctx context.Context, sessionControls models.ConditionalAccessSessionControlsable) *ConditionalAccessSessionControlsModel {
	if sessionControls == nil {
		tflog.Debug(ctx, "Session controls model is nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map session controls")

	result := &ConditionalAccessSessionControlsModel{}

	if appRestrictions := sessionControls.GetApplicationEnforcedRestrictions(); appRestrictions != nil {
		result.ApplicationEnforcedRestrictions = &ApplicationEnforcedRestrictionsSessionControlModel{
			IsEnabled: types.BoolValue(*appRestrictions.GetIsEnabled()),
		}
	}

	if cloudAppSecurity := sessionControls.GetCloudAppSecurity(); cloudAppSecurity != nil {
		result.CloudAppSecurity = &CloudAppSecuritySessionControlModel{
			IsEnabled:            types.BoolValue(*cloudAppSecurity.GetIsEnabled()),
			CloudAppSecurityType: types.StringValue(cloudAppSecurity.GetCloudAppSecurityType().String()),
		}
	}

	if continuousAccessEvaluation := sessionControls.GetContinuousAccessEvaluation(); continuousAccessEvaluation != nil {
		result.ContinuousAccessEvaluation = &ContinuousAccessEvaluationSessionControlModel{
			Mode: types.StringValue(continuousAccessEvaluation.GetMode().String()),
		}
	}

	if persistentBrowser := sessionControls.GetPersistentBrowser(); persistentBrowser != nil {
		result.PersistentBrowser = &PersistentBrowserSessionControlModel{
			IsEnabled: types.BoolValue(*persistentBrowser.GetIsEnabled()),
			Mode:      types.StringValue(persistentBrowser.GetMode().String()),
		}
	}

	if signInFrequency := sessionControls.GetSignInFrequency(); signInFrequency != nil {
		result.SignInFrequency = &SignInFrequencySessionControlModel{
			IsEnabled:          types.BoolValue(*signInFrequency.GetIsEnabled()),
			Type:               types.StringValue(signInFrequency.GetTypeEscaped().String()),
			Value:              types.Int64Value(int64(*signInFrequency.GetValue())),
			FrequencyInterval:  types.StringValue(signInFrequency.GetFrequencyInterval().String()),
			AuthenticationType: types.StringValue(signInFrequency.GetAuthenticationType().String()),
		}
	}

	if secureSignInSession := sessionControls.GetSecureSignInSession(); secureSignInSession != nil {
		result.SecureSignInSession = &SecureSignInSessionControlModel{
			IsEnabled: types.BoolValue(*secureSignInSession.GetIsEnabled()),
		}
	}

	if disableResilienceDefaults := sessionControls.GetDisableResilienceDefaults(); disableResilienceDefaults != nil {
		result.DisableResilienceDefaults = types.BoolValue(*disableResilienceDefaults)
	}

	tflog.Debug(ctx, "Finished mapping session controls", map[string]interface{}{
		"hasApplicationEnforcedRestrictions": result.ApplicationEnforcedRestrictions != nil,
		"hasCloudAppSecurity":                result.CloudAppSecurity != nil,
		"hasContinuousAccessEvaluation":      result.ContinuousAccessEvaluation != nil,
		"hasPersistentBrowser":               result.PersistentBrowser != nil,
		"hasSignInFrequency":                 result.SignInFrequency != nil,
		"hasSecureSignInSession":             result.SecureSignInSession != nil,
		"disableResilienceDefaults":          result.DisableResilienceDefaults.ValueBool(),
	})

	return result
}
