package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model.
// It populates the ConditionalAccessPolicyResourceModel with data from the ConditionalAccessPolicy.
func mapRemoteStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource models.ConditionalAccessPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform", map[string]interface{}{
		"resourceId": helpers.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(helpers.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(helpers.StringPtrToString(remoteResource.GetDescription()))
	data.CreatedDateTime = helpers.TimeToString(remoteResource.GetCreatedDateTime())
	data.ModifiedDateTime = helpers.TimeToString(remoteResource.GetModifiedDateTime())

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
		ClientAppTypes:             helpers.EnumSliceToTypeStringSlice(conditions.GetClientAppTypes()),
		Devices:                    mapDevices(ctx, conditions.GetDevices()),
		DeviceStates:               mapDeviceStates(ctx, conditions.GetDeviceStates()),
		Locations:                  mapLocations(ctx, conditions.GetLocations()),
		Platforms:                  mapPlatforms(ctx, conditions.GetPlatforms()),
		ServicePrincipalRiskLevels: helpers.EnumSliceToTypeStringSlice(conditions.GetServicePrincipalRiskLevels()),
		SignInRiskLevels:           helpers.EnumSliceToTypeStringSlice(conditions.GetSignInRiskLevels()),
		UserRiskLevels:             helpers.EnumSliceToTypeStringSlice(conditions.GetUserRiskLevels()),
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
		IncludeApplications: helpers.SliceToTypeStringSlice(apps.GetIncludeApplications()),
		ExcludeApplications: helpers.SliceToTypeStringSlice(apps.GetExcludeApplications()),
		IncludeUserActions:  helpers.SliceToTypeStringSlice(apps.GetIncludeUserActions()),
		ApplicationFilter:   mapFilter(apps.GetApplicationFilter()),
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
		ExcludeGroups:                helpers.SliceToTypeStringSlice(users.GetExcludeGroups()),
		ExcludeRoles:                 helpers.SliceToTypeStringSlice(users.GetExcludeRoles()),
		ExcludeUsers:                 helpers.SliceToTypeStringSlice(users.GetExcludeUsers()),
		IncludeGroups:                helpers.SliceToTypeStringSlice(users.GetIncludeGroups()),
		IncludeRoles:                 helpers.SliceToTypeStringSlice(users.GetIncludeRoles()),
		IncludeUsers:                 helpers.SliceToTypeStringSlice(users.GetIncludeUsers()),
		ExcludeGuestsOrExternalUsers: mapGuestsOrExternalUsers(users.GetExcludeGuestsOrExternalUsers()),
		IncludeGuestsOrExternalUsers: mapGuestsOrExternalUsers(users.GetIncludeGuestsOrExternalUsers()),
	}
	tflog.Debug(ctx, "Finished mapping users", map[string]interface{}{
		"excludeGroupsCount": len(result.ExcludeGroups),
		"includeGroupsCount": len(result.IncludeGroups),
	})
	return result
}

func mapAuthenticationFlows(ctx context.Context, authFlows models.ConditionalAccessAuthenticationFlowsable) *ConditionalAccessAuthenticationFlowsModel {
	if authFlows == nil {
		tflog.Debug(ctx, "Authentication flows are nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map authentication flows")

	result := &ConditionalAccessAuthenticationFlowsModel{
		TransferMethods: types.StringValue(string(*authFlows.GetTransferMethods())),
	}

	tflog.Debug(ctx, "Finished mapping authentication flows", map[string]interface{}{
		"transferMethods": result.TransferMethods.ValueString(),
	})

	return result
}

func mapGuestsOrExternalUsers(guestsOrExternalUsers models.ConditionalAccessGuestsOrExternalUsersable) *ConditionalAccessGuestsOrExternalUsersModel {
	if guestsOrExternalUsers == nil {
		return nil
	}

	return &ConditionalAccessGuestsOrExternalUsersModel{
		ExternalTenants:          mapExternalTenants(guestsOrExternalUsers.GetExternalTenants()),
		GuestOrExternalUserTypes: types.StringValue(string(*guestsOrExternalUsers.GetGuestOrExternalUserTypes())),
	}
}

func mapExternalTenants(externalTenants models.ConditionalAccessExternalTenantsable) *ConditionalAccessExternalTenantsModel {
	if externalTenants == nil {
		return nil
	}

	return &ConditionalAccessExternalTenantsModel{
		MembershipKind: types.StringValue(string(*externalTenants.GetMembershipKind())),
	}
}

func mapClientApplications(ctx context.Context, clientApps models.ConditionalAccessClientApplicationsable) *ConditionalAccessClientApplicationsModel {
	if clientApps == nil {
		tflog.Debug(ctx, "Client applications are nil")
		return nil
	}

	tflog.Debug(ctx, "Starting to map client applications")

	result := &ConditionalAccessClientApplicationsModel{
		ExcludeServicePrincipals: helpers.SliceToTypeStringSlice(clientApps.GetExcludeServicePrincipals()),
		IncludeServicePrincipals: helpers.SliceToTypeStringSlice(clientApps.GetIncludeServicePrincipals()),
		ServicePrincipalFilter:   mapFilter(clientApps.GetServicePrincipalFilter()),
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
		return nil
	}

	return &ConditionalAccessDevicesModel{
		IncludeDevices: helpers.SliceToTypeStringSlice(devices.GetIncludeDevices()),
		ExcludeDevices: helpers.SliceToTypeStringSlice(devices.GetExcludeDevices()),
		IncludeStates:  helpers.SliceToTypeStringSlice(devices.GetIncludeDeviceStates()),
		ExcludeStates:  helpers.SliceToTypeStringSlice(devices.GetExcludeDeviceStates()),
		DeviceFilter:   mapFilter(devices.GetDeviceFilter()),
	}
}

func mapDeviceStates(ctx context.Context, deviceStates models.ConditionalAccessDeviceStatesable) *ConditionalAccessDeviceStatesModel {
	if deviceStates == nil {
		return nil
	}

	return &ConditionalAccessDeviceStatesModel{
		IncludeStates: helpers.SliceToTypeStringSlice(deviceStates.GetIncludeStates()),
		ExcludeStates: helpers.SliceToTypeStringSlice(deviceStates.GetExcludeStates()),
	}
}

func mapLocations(ctx context.Context, locations models.ConditionalAccessLocationsable) *ConditionalAccessLocationsModel {
	if locations == nil {
		return nil
	}

	return &ConditionalAccessLocationsModel{
		ExcludeLocations: helpers.SliceToTypeStringSlice(locations.GetExcludeLocations()),
		IncludeLocations: helpers.SliceToTypeStringSlice(locations.GetIncludeLocations()),
	}
}

func mapPlatforms(ctx context.Context, platforms models.ConditionalAccessPlatformsable) *ConditionalAccessPlatformsModel {
	if platforms == nil {
		return nil
	}

	return &ConditionalAccessPlatformsModel{
		ExcludePlatforms: helpers.EnumSliceToTypeStringSlice(platforms.GetExcludePlatforms()),
		IncludePlatforms: helpers.EnumSliceToTypeStringSlice(platforms.GetIncludePlatforms()),
	}
}

func mapFilter(filter models.ConditionalAccessFilterable) *ConditionalAccessFilterModel {
	if filter == nil {
		return nil
	}

	return &ConditionalAccessFilterModel{
		Mode: types.StringValue(string(*filter.GetMode())),
		Rule: types.StringValue(*filter.GetRule()),
	}
}

// Helper functions for mapping sub-components (e.g., mapApplications, mapUsers, etc.)

func mapGrantControls(ctx context.Context, grantControls models.ConditionalAccessGrantControlsable) *ConditionalAccessGrantControlsModel {
	if grantControls == nil {
		return nil
	}

	result := &ConditionalAccessGrantControlsModel{}

	if operator := grantControls.GetOperator(); operator != nil {
		result.Operator = types.StringValue(*operator)
	}

	if builtInControls := grantControls.GetBuiltInControls(); builtInControls != nil {
		result.BuiltInControls = helpers.EnumSliceToTypeStringSlice(builtInControls)
	}

	if customAuthenticationFactors := grantControls.GetCustomAuthenticationFactors(); customAuthenticationFactors != nil {
		result.CustomAuthenticationFactors = helpers.SliceToTypeStringSlice(customAuthenticationFactors)
	}

	if termsOfUse := grantControls.GetTermsOfUse(); termsOfUse != nil {
		result.TermsOfUse = helpers.SliceToTypeStringSlice(termsOfUse)
	}

	// Map AuthenticationStrength if needed

	return result
}

func mapSessionControls(ctx context.Context, sessionControls models.ConditionalAccessSessionControlsable) *ConditionalAccessSessionControlsModel {
	if sessionControls == nil {
		return nil
	}

	result := &ConditionalAccessSessionControlsModel{}

	// Map ApplicationEnforcedRestrictions
	if appRestrictions := sessionControls.GetApplicationEnforcedRestrictions(); appRestrictions != nil {
		result.ApplicationEnforcedRestrictions = &ApplicationEnforcedRestrictionsSessionControlModel{
			IsEnabled: types.BoolValue(*appRestrictions.GetIsEnabled()),
		}
	}

	// Map CloudAppSecurity
	if cloudAppSecurity := sessionControls.GetCloudAppSecurity(); cloudAppSecurity != nil {
		result.CloudAppSecurity = &CloudAppSecuritySessionControlModel{
			IsEnabled:            types.BoolValue(*cloudAppSecurity.GetIsEnabled()),
			CloudAppSecurityType: types.StringValue(string(*cloudAppSecurity.GetCloudAppSecurityType())),
		}
	}

	// Map other session controls (PersistentBrowser, SignInFrequency, etc.) similarly

	return result
}

// Additional helper functions for mapping other components would be defined here.
