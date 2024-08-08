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
		ExcludeGroups:                helpers.SliceToTypeStringSlice(users.GetExcludeGroups()),
		ExcludeRoles:                 helpers.SliceToTypeStringSlice(users.GetExcludeRoles()),
		ExcludeUsers:                 helpers.SliceToTypeStringSlice(users.GetExcludeUsers()),
		IncludeGroups:                helpers.SliceToTypeStringSlice(users.GetIncludeGroups()),
		IncludeRoles:                 helpers.SliceToTypeStringSlice(users.GetIncludeRoles()),
		IncludeUsers:                 helpers.SliceToTypeStringSlice(users.GetIncludeUsers()),
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
		ExcludeServicePrincipals: helpers.SliceToTypeStringSlice(clientApps.GetExcludeServicePrincipals()),
		IncludeServicePrincipals: helpers.SliceToTypeStringSlice(clientApps.GetIncludeServicePrincipals()),
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
		IncludeDevices: helpers.SliceToTypeStringSlice(devices.GetIncludeDevices()),
		ExcludeDevices: helpers.SliceToTypeStringSlice(devices.GetExcludeDevices()),
		IncludeStates:  helpers.SliceToTypeStringSlice(devices.GetIncludeDeviceStates()),
		ExcludeStates:  helpers.SliceToTypeStringSlice(devices.GetExcludeDeviceStates()),
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
		IncludeStates: helpers.SliceToTypeStringSlice(deviceStates.GetIncludeStates()),
		ExcludeStates: helpers.SliceToTypeStringSlice(deviceStates.GetExcludeStates()),
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
		ExcludeLocations: helpers.SliceToTypeStringSlice(locations.GetExcludeLocations()),
		IncludeLocations: helpers.SliceToTypeStringSlice(locations.GetIncludeLocations()),
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
		ExcludePlatforms: helpers.EnumSliceToTypeStringSlice(platforms.GetExcludePlatforms()),
		IncludePlatforms: helpers.EnumSliceToTypeStringSlice(platforms.GetIncludePlatforms()),
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
// TODO

func mapGrantControls(ctx context.Context, grantControls models.ConditionalAccessGrantControlsable) *ConditionalAccessGrantControlsModel {

	return nil
}

func mapSessionControls(ctx context.Context, sessionControls models.ConditionalAccessSessionControlsable) *ConditionalAccessSessionControlsModel {

	return nil
}
