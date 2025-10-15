package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteResourceStateToTerraform maps the remote conditional access policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource map[string]any) {
	tflog.Debug(ctx, "Starting MapRemoteResourceStateToTerraform", map[string]any{
		"remoteResource": remoteResource,
	})

	// Basic properties
	if id, ok := remoteResource["id"].(string); ok {
		tflog.Debug(ctx, "Mapping ID", map[string]any{"id": id})
		data.ID = types.StringValue(id)
	} else {
		tflog.Debug(ctx, "ID not found or not a string")
		data.ID = types.StringNull()
	}

	if displayName, ok := remoteResource["displayName"].(string); ok {
		tflog.Debug(ctx, "Mapping displayName", map[string]any{"displayName": displayName})
		data.DisplayName = types.StringValue(displayName)
	} else {
		tflog.Debug(ctx, "displayName not found or not a string")
		data.DisplayName = types.StringNull()
	}

	if state, ok := remoteResource["state"].(string); ok {
		tflog.Debug(ctx, "Mapping state", map[string]any{"state": state})
		data.State = types.StringValue(state)
	} else {
		tflog.Debug(ctx, "state not found or not a string")
		data.State = types.StringNull()
	}

	if createdDateTime, ok := remoteResource["createdDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping createdDateTime", map[string]any{"createdDateTime": createdDateTime})
		data.CreatedDateTime = types.StringValue(createdDateTime)
	} else {
		tflog.Debug(ctx, "createdDateTime not found or not a string")
		data.CreatedDateTime = types.StringNull()
	}

	if modifiedDateTime, ok := remoteResource["modifiedDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping modifiedDateTime", map[string]any{"modifiedDateTime": modifiedDateTime})
		data.ModifiedDateTime = types.StringValue(modifiedDateTime)
	} else {
		tflog.Debug(ctx, "modifiedDateTime not found or not a string")
		data.ModifiedDateTime = types.StringNull()
	}

	if deletedDateTime, ok := remoteResource["deletedDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping deletedDateTime", map[string]any{"deletedDateTime": deletedDateTime})
		data.DeletedDateTime = types.StringValue(deletedDateTime)
	} else {
		tflog.Debug(ctx, "deletedDateTime not found or not a string")
		data.DeletedDateTime = types.StringNull()
	}

	if templateId, ok := remoteResource["templateId"].(string); ok {
		tflog.Debug(ctx, "Mapping templateId", map[string]any{"templateId": templateId})
		data.TemplateId = types.StringValue(templateId)
	} else {
		tflog.Debug(ctx, "templateId not found or not a string")
		data.TemplateId = types.StringNull()
	}

	if partialEnablementStrategy, ok := remoteResource["partialEnablementStrategy"].(string); ok {
		tflog.Debug(ctx, "Mapping partialEnablementStrategy", map[string]any{"partialEnablementStrategy": partialEnablementStrategy})
		data.PartialEnablementStrategy = types.StringValue(partialEnablementStrategy)
	} else {
		tflog.Debug(ctx, "partialEnablementStrategy not found or not a string")
		data.PartialEnablementStrategy = types.StringNull()
	}

	// Map conditions
	if conditionsRaw, ok := remoteResource["conditions"]; ok {
		tflog.Debug(ctx, "Mapping conditions", map[string]any{"conditions": conditionsRaw})
		data.Conditions = mapConditions(ctx, conditionsRaw)
	} else {
		tflog.Debug(ctx, "conditions not found")
		data.Conditions = nil
	}

	// Map grant controls
	if grantControlsRaw, ok := remoteResource["grantControls"]; ok {
		tflog.Debug(ctx, "Mapping grantControls", map[string]any{"grantControls": grantControlsRaw})
		data.GrantControls = mapGrantControls(ctx, grantControlsRaw)
	} else {
		tflog.Debug(ctx, "grantControls not found")
		data.GrantControls = nil
	}

	// Map session controls
	if sessionControlsRaw, ok := remoteResource["sessionControls"]; ok {
		tflog.Debug(ctx, "Mapping sessionControls", map[string]any{"sessionControls": sessionControlsRaw})
		data.SessionControls = mapSessionControls(ctx, sessionControlsRaw)
	} else {
		tflog.Debug(ctx, "sessionControls not found")
		data.SessionControls = nil
	}

	tflog.Debug(ctx, "Completed MapRemoteResourceStateToTerraform")
}

func mapConditions(ctx context.Context, conditionsRaw any) *ConditionalAccessConditions {
	conditions, ok := conditionsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "conditions is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessConditions{}

	// Map clientAppTypes (set)
	if clientAppTypesRaw, ok := conditions["clientAppTypes"]; ok {
		tflog.Debug(ctx, "Mapping clientAppTypes", map[string]any{"clientAppTypes": clientAppTypesRaw})
		result.ClientAppTypes = mapStringSliceToSet(ctx, clientAppTypesRaw, "clientAppTypes")
	} else {
		tflog.Debug(ctx, "clientAppTypes not found, setting to null")
		result.ClientAppTypes = types.SetNull(types.StringType)
	}

	// Map signInRiskLevels (set)
	if signInRiskLevelsRaw, ok := conditions["signInRiskLevels"]; ok {
		tflog.Debug(ctx, "Mapping signInRiskLevels", map[string]any{"signInRiskLevels": signInRiskLevelsRaw})
		result.SignInRiskLevels = mapStringSliceToSet(ctx, signInRiskLevelsRaw, "signInRiskLevels")
	} else {
		tflog.Debug(ctx, "signInRiskLevels not found, setting to null")
		result.SignInRiskLevels = types.SetNull(types.StringType)
	}

	// Map userRiskLevels (set)
	if userRiskLevelsRaw, ok := conditions["userRiskLevels"]; ok {
		tflog.Debug(ctx, "Mapping userRiskLevels", map[string]any{"userRiskLevels": userRiskLevelsRaw})
		result.UserRiskLevels = mapStringSliceToSet(ctx, userRiskLevelsRaw, "userRiskLevels")
	} else {
		tflog.Debug(ctx, "userRiskLevels not found, setting to null")
		result.UserRiskLevels = types.SetNull(types.StringType)
	}

	// Map servicePrincipalRiskLevels (set)
	if servicePrincipalRiskLevelsRaw, ok := conditions["servicePrincipalRiskLevels"]; ok {
		tflog.Debug(ctx, "Mapping servicePrincipalRiskLevels", map[string]any{"servicePrincipalRiskLevels": servicePrincipalRiskLevelsRaw})
		result.ServicePrincipalRiskLevels = mapStringSliceToSet(ctx, servicePrincipalRiskLevelsRaw, "servicePrincipalRiskLevels")
	} else {
		tflog.Debug(ctx, "servicePrincipalRiskLevels not found, setting to null")
		result.ServicePrincipalRiskLevels = types.SetNull(types.StringType)
	}

	// Map applications
	if applicationsRaw, ok := conditions["applications"]; ok {
		tflog.Debug(ctx, "Mapping applications", map[string]any{"applications": applicationsRaw})
		result.Applications = mapApplications(ctx, applicationsRaw)
	} else {
		tflog.Debug(ctx, "applications not found")
		result.Applications = nil
	}

	// Map users
	if usersRaw, ok := conditions["users"]; ok {
		tflog.Debug(ctx, "Mapping users", map[string]any{"users": usersRaw})
		result.Users = mapUsers(ctx, usersRaw)
	} else {
		tflog.Debug(ctx, "users not found")
		result.Users = nil
	}

	// Map locations
	if locationsRaw, ok := conditions["locations"]; ok {
		tflog.Debug(ctx, "Mapping locations", map[string]any{"locations": locationsRaw})
		result.Locations = mapLocations(ctx, locationsRaw)
	} else {
		tflog.Debug(ctx, "locations not found")
		result.Locations = nil
	}

	// Map platforms
	if platformsRaw, ok := conditions["platforms"]; ok {
		tflog.Debug(ctx, "Mapping platforms", map[string]any{"platforms": platformsRaw})
		result.Platforms = mapPlatforms(ctx, platformsRaw)
	} else {
		tflog.Debug(ctx, "platforms not found")
		result.Platforms = nil
	}

	// Map devices
	if devicesRaw, ok := conditions["devices"]; ok {
		tflog.Debug(ctx, "Mapping devices", map[string]any{"devices": devicesRaw})
		result.Devices = mapDevices(ctx, devicesRaw)
	} else {
		tflog.Debug(ctx, "devices not found")
		result.Devices = nil
	}

	// Map clientApplications
	if clientApplicationsRaw, ok := conditions["clientApplications"]; ok {
		tflog.Debug(ctx, "Mapping clientApplications", map[string]any{"clientApplications": clientApplicationsRaw})
		result.ClientApplications = mapClientApplications(ctx, clientApplicationsRaw)
	} else {
		tflog.Debug(ctx, "clientApplications not found")
		result.ClientApplications = nil
	}

	// Map times
	if timesRaw, ok := conditions["times"]; ok {
		tflog.Debug(ctx, "Mapping times", map[string]any{"times": timesRaw})
		result.Times = mapTimes(ctx, timesRaw)
	} else {
		tflog.Debug(ctx, "times not found")
		result.Times = nil
	}

	// Map deviceStates
	if deviceStatesRaw, ok := conditions["deviceStates"]; ok {
		tflog.Debug(ctx, "Mapping deviceStates", map[string]any{"deviceStates": deviceStatesRaw})
		result.DeviceStates = mapDeviceStates(ctx, deviceStatesRaw)
	} else {
		tflog.Debug(ctx, "deviceStates not found")
		result.DeviceStates = nil
	}

	return result
}

func mapApplications(ctx context.Context, applicationsRaw any) *ConditionalAccessApplications {
	applications, ok := applicationsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "applications is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessApplications{}

	// Map includeApplications (set)
	if includeApplicationsRaw, ok := applications["includeApplications"]; ok {
		tflog.Debug(ctx, "Mapping includeApplications", map[string]any{"includeApplications": includeApplicationsRaw})
		result.IncludeApplications = mapStringSliceToSet(ctx, includeApplicationsRaw, "includeApplications")
	} else {
		tflog.Debug(ctx, "includeApplications not found, setting to null")
		result.IncludeApplications = types.SetNull(types.StringType)
	}

	// Map excludeApplications (set)
	if excludeApplicationsRaw, ok := applications["excludeApplications"]; ok {
		tflog.Debug(ctx, "Mapping excludeApplications", map[string]any{"excludeApplications": excludeApplicationsRaw})
		result.ExcludeApplications = mapStringSliceToSet(ctx, excludeApplicationsRaw, "excludeApplications")
	} else {
		tflog.Debug(ctx, "excludeApplications not found, setting to null")
		result.ExcludeApplications = types.SetNull(types.StringType)
	}

	// Map includeUserActions (set)
	if includeUserActionsRaw, ok := applications["includeUserActions"]; ok {
		tflog.Debug(ctx, "Mapping includeUserActions", map[string]any{"includeUserActions": includeUserActionsRaw})
		result.IncludeUserActions = mapStringSliceToSet(ctx, includeUserActionsRaw, "includeUserActions")
	} else {
		tflog.Debug(ctx, "includeUserActions not found, setting to null")
		result.IncludeUserActions = types.SetNull(types.StringType)
	}

	// Map includeAuthenticationContextClassReferences (set)
	if includeAuthContextRaw, ok := applications["includeAuthenticationContextClassReferences"]; ok {
		tflog.Debug(ctx, "Mapping includeAuthenticationContextClassReferences", map[string]any{"includeAuthenticationContextClassReferences": includeAuthContextRaw})
		result.IncludeAuthenticationContextClassReferences = mapAuthContextClassReferencesToSet(ctx, includeAuthContextRaw, "includeAuthenticationContextClassReferences")
	} else {
		tflog.Debug(ctx, "includeAuthenticationContextClassReferences not found, setting to null")
		result.IncludeAuthenticationContextClassReferences = types.SetNull(types.StringType)
	}

	// Map applicationFilter
	if applicationFilterRaw, ok := applications["applicationFilter"]; ok {
		tflog.Debug(ctx, "Mapping applicationFilter", map[string]any{"applicationFilter": applicationFilterRaw})
		result.ApplicationFilter = mapFilter(ctx, applicationFilterRaw)
	} else {
		tflog.Debug(ctx, "applicationFilter not found")
		result.ApplicationFilter = nil
	}

	// Map globalSecureAccess (typically null)
	if globalSecureAccessRaw, ok := applications["globalSecureAccess"]; ok && globalSecureAccessRaw != nil {
		tflog.Debug(ctx, "Mapping globalSecureAccess", map[string]any{"globalSecureAccess": globalSecureAccessRaw})
		// For now, since this field is typically null, we'll map it as null object
		result.GlobalSecureAccess = types.ObjectNull(map[string]attr.Type{})
	} else {
		tflog.Debug(ctx, "globalSecureAccess not found or is null")
		result.GlobalSecureAccess = types.ObjectNull(map[string]attr.Type{})
	}

	return result
}

func mapUsers(ctx context.Context, usersRaw any) *ConditionalAccessUsers {
	users, ok := usersRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "users is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessUsers{}

	// Map includeUsers (set)
	if includeUsersRaw, ok := users["includeUsers"]; ok {
		tflog.Debug(ctx, "Mapping includeUsers", map[string]any{"includeUsers": includeUsersRaw})
		result.IncludeUsers = mapStringSliceToSet(ctx, includeUsersRaw, "includeUsers")
	} else {
		tflog.Debug(ctx, "includeUsers not found, setting to null")
		result.IncludeUsers = types.SetNull(types.StringType)
	}

	// Map excludeUsers (set)
	if excludeUsersRaw, ok := users["excludeUsers"]; ok {
		tflog.Debug(ctx, "Mapping excludeUsers", map[string]any{"excludeUsers": excludeUsersRaw})
		result.ExcludeUsers = mapStringSliceToSet(ctx, excludeUsersRaw, "excludeUsers")
	} else {
		tflog.Debug(ctx, "excludeUsers not found, setting to null")
		result.ExcludeUsers = types.SetNull(types.StringType)
	}

	// Map includeGroups (set)
	if includeGroupsRaw, ok := users["includeGroups"]; ok {
		tflog.Debug(ctx, "Mapping includeGroups", map[string]any{"includeGroups": includeGroupsRaw})
		result.IncludeGroups = mapStringSliceToSet(ctx, includeGroupsRaw, "includeGroups")
	} else {
		tflog.Debug(ctx, "includeGroups not found, setting to null")
		result.IncludeGroups = types.SetNull(types.StringType)
	}

	// Map excludeGroups (set)
	if excludeGroupsRaw, ok := users["excludeGroups"]; ok {
		tflog.Debug(ctx, "Mapping excludeGroups", map[string]any{"excludeGroups": excludeGroupsRaw})
		result.ExcludeGroups = mapStringSliceToSet(ctx, excludeGroupsRaw, "excludeGroups")
	} else {
		tflog.Debug(ctx, "excludeGroups not found, setting to null")
		result.ExcludeGroups = types.SetNull(types.StringType)
	}

	// Map includeRoles (set)
	if includeRolesRaw, ok := users["includeRoles"]; ok {
		tflog.Debug(ctx, "Mapping includeRoles", map[string]any{"includeRoles": includeRolesRaw})
		result.IncludeRoles = mapStringSliceToSet(ctx, includeRolesRaw, "includeRoles")
	} else {
		tflog.Debug(ctx, "includeRoles not found, setting to null")
		result.IncludeRoles = types.SetNull(types.StringType)
	}

	// Map excludeRoles (set)
	if excludeRolesRaw, ok := users["excludeRoles"]; ok {
		tflog.Debug(ctx, "Mapping excludeRoles", map[string]any{"excludeRoles": excludeRolesRaw})
		result.ExcludeRoles = mapStringSliceToSet(ctx, excludeRolesRaw, "excludeRoles")
	} else {
		tflog.Debug(ctx, "excludeRoles not found, setting to null")
		result.ExcludeRoles = types.SetNull(types.StringType)
	}

	// Map includeGuestsOrExternalUsers (object)
	if includeGuestsOrExternalUsersRaw, ok := users["includeGuestsOrExternalUsers"]; ok {
		tflog.Debug(ctx, "Mapping includeGuestsOrExternalUsers", map[string]any{"includeGuestsOrExternalUsers": includeGuestsOrExternalUsersRaw})
		result.IncludeGuestsOrExternalUsers = mapGuestsOrExternalUsersToObject(ctx, includeGuestsOrExternalUsersRaw)
	} else {
		tflog.Debug(ctx, "includeGuestsOrExternalUsers not found, setting to null")
		result.IncludeGuestsOrExternalUsers = types.ObjectNull(map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"membership_kind": types.StringType,
					"members":         types.SetType{ElemType: types.StringType},
				},
			},
		})
	}

	// Map excludeGuestsOrExternalUsers (object)
	if excludeGuestsOrExternalUsersRaw, ok := users["excludeGuestsOrExternalUsers"]; ok {
		tflog.Debug(ctx, "Mapping excludeGuestsOrExternalUsers", map[string]any{"excludeGuestsOrExternalUsers": excludeGuestsOrExternalUsersRaw})
		result.ExcludeGuestsOrExternalUsers = mapGuestsOrExternalUsersToObject(ctx, excludeGuestsOrExternalUsersRaw)
	} else {
		tflog.Debug(ctx, "excludeGuestsOrExternalUsers not found, setting to null")
		result.ExcludeGuestsOrExternalUsers = types.ObjectNull(map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"membership_kind": types.StringType,
					"members":         types.SetType{ElemType: types.StringType},
				},
			},
		})
	}

	return result
}

func mapLocations(ctx context.Context, locationsRaw any) *ConditionalAccessLocations {
	locations, ok := locationsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "locations is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessLocations{}

	// Map includeLocations (set)
	if includeLocationsRaw, ok := locations["includeLocations"]; ok {
		tflog.Debug(ctx, "Mapping includeLocations", map[string]any{"includeLocations": includeLocationsRaw})
		result.IncludeLocations = mapStringSliceToSet(ctx, includeLocationsRaw, "includeLocations")
	} else {
		tflog.Debug(ctx, "includeLocations not found, setting to null")
		result.IncludeLocations = types.SetNull(types.StringType)
	}

	// Map excludeLocations (set)
	if excludeLocationsRaw, ok := locations["excludeLocations"]; ok {
		tflog.Debug(ctx, "Mapping excludeLocations", map[string]any{"excludeLocations": excludeLocationsRaw})
		result.ExcludeLocations = mapStringSliceToSet(ctx, excludeLocationsRaw, "excludeLocations")
	} else {
		tflog.Debug(ctx, "excludeLocations not found, setting to null")
		result.ExcludeLocations = types.SetNull(types.StringType)
	}

	return result
}

func mapPlatforms(ctx context.Context, platformsRaw any) *ConditionalAccessPlatforms {
	platforms, ok := platformsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "platforms is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessPlatforms{}

	// Map includePlatforms (set)
	if includePlatformsRaw, ok := platforms["includePlatforms"]; ok {
		tflog.Debug(ctx, "Mapping includePlatforms", map[string]any{"includePlatforms": includePlatformsRaw})
		result.IncludePlatforms = mapStringSliceToSet(ctx, includePlatformsRaw, "includePlatforms")
	} else {
		tflog.Debug(ctx, "includePlatforms not found, setting to null")
		result.IncludePlatforms = types.SetNull(types.StringType)
	}

	// Map excludePlatforms (set)
	if excludePlatformsRaw, ok := platforms["excludePlatforms"]; ok {
		tflog.Debug(ctx, "Mapping excludePlatforms", map[string]any{"excludePlatforms": excludePlatformsRaw})
		result.ExcludePlatforms = mapStringSliceToSet(ctx, excludePlatformsRaw, "excludePlatforms")
	} else {
		tflog.Debug(ctx, "excludePlatforms not found, setting to null")
		result.ExcludePlatforms = types.SetNull(types.StringType)
	}

	return result
}

func mapDevices(ctx context.Context, devicesRaw any) *ConditionalAccessDevices {
	devices, ok := devicesRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "devices is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessDevices{}

	// Map includeDevices (set)
	if includeDevicesRaw, ok := devices["includeDevices"]; ok {
		tflog.Debug(ctx, "Mapping includeDevices", map[string]any{"includeDevices": includeDevicesRaw})
		result.IncludeDevices = mapStringSliceToSet(ctx, includeDevicesRaw, "includeDevices")
	} else {
		tflog.Debug(ctx, "includeDevices not found, setting to null")
		result.IncludeDevices = types.SetNull(types.StringType)
	}

	// Map excludeDevices (set)
	if excludeDevicesRaw, ok := devices["excludeDevices"]; ok {
		tflog.Debug(ctx, "Mapping excludeDevices", map[string]any{"excludeDevices": excludeDevicesRaw})
		result.ExcludeDevices = mapStringSliceToSet(ctx, excludeDevicesRaw, "excludeDevices")
	} else {
		tflog.Debug(ctx, "excludeDevices not found, setting to null")
		result.ExcludeDevices = types.SetNull(types.StringType)
	}

	// Map includeDeviceStates (set)
	if includeDeviceStatesRaw, ok := devices["includeDeviceStates"]; ok {
		tflog.Debug(ctx, "Mapping includeDeviceStates", map[string]any{"includeDeviceStates": includeDeviceStatesRaw})
		result.IncludeDeviceStates = mapStringSliceToSet(ctx, includeDeviceStatesRaw, "includeDeviceStates")
	} else {
		tflog.Debug(ctx, "includeDeviceStates not found, setting to null")
		result.IncludeDeviceStates = types.SetNull(types.StringType)
	}

	// Map excludeDeviceStates (set)
	if excludeDeviceStatesRaw, ok := devices["excludeDeviceStates"]; ok {
		tflog.Debug(ctx, "Mapping excludeDeviceStates", map[string]any{"excludeDeviceStates": excludeDeviceStatesRaw})
		result.ExcludeDeviceStates = mapStringSliceToSet(ctx, excludeDeviceStatesRaw, "excludeDeviceStates")
	} else {
		tflog.Debug(ctx, "excludeDeviceStates not found, setting to null")
		result.ExcludeDeviceStates = types.SetNull(types.StringType)
	}

	// Map deviceFilter
	if deviceFilterRaw, ok := devices["deviceFilter"]; ok {
		tflog.Debug(ctx, "Mapping deviceFilter", map[string]any{"deviceFilter": deviceFilterRaw})
		result.DeviceFilter = mapFilter(ctx, deviceFilterRaw)
	} else {
		tflog.Debug(ctx, "deviceFilter not found")
		result.DeviceFilter = nil
	}

	return result
}

func mapClientApplications(ctx context.Context, clientApplicationsRaw any) *ConditionalAccessClientApplications {
	clientApplications, ok := clientApplicationsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "clientApplications is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessClientApplications{}

	// Map includeServicePrincipals (set)
	if includeServicePrincipalsRaw, ok := clientApplications["includeServicePrincipals"]; ok {
		tflog.Debug(ctx, "Mapping includeServicePrincipals", map[string]any{"includeServicePrincipals": includeServicePrincipalsRaw})
		result.IncludeServicePrincipals = mapStringSliceToSet(ctx, includeServicePrincipalsRaw, "includeServicePrincipals")
	} else {
		tflog.Debug(ctx, "includeServicePrincipals not found, setting to null")
		result.IncludeServicePrincipals = types.SetNull(types.StringType)
	}

	// Map excludeServicePrincipals (set)
	if excludeServicePrincipalsRaw, ok := clientApplications["excludeServicePrincipals"]; ok {
		tflog.Debug(ctx, "Mapping excludeServicePrincipals", map[string]any{"excludeServicePrincipals": excludeServicePrincipalsRaw})
		result.ExcludeServicePrincipals = mapStringSliceToSet(ctx, excludeServicePrincipalsRaw, "excludeServicePrincipals")
	} else {
		tflog.Debug(ctx, "excludeServicePrincipals not found, setting to null")
		result.ExcludeServicePrincipals = types.SetNull(types.StringType)
	}

	return result
}

func mapTimes(ctx context.Context, timesRaw any) *ConditionalAccessTimes {
	times, ok := timesRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "times is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessTimes{}

	// Map includedRanges (set)
	if includedRangesRaw, ok := times["includedRanges"]; ok {
		tflog.Debug(ctx, "Mapping includedRanges", map[string]any{"includedRanges": includedRangesRaw})
		result.IncludedRanges = mapStringSliceToSet(ctx, includedRangesRaw, "includedRanges")
	} else {
		tflog.Debug(ctx, "includedRanges not found, setting to null")
		result.IncludedRanges = types.SetNull(types.StringType)
	}

	// Map excludedRanges (set)
	if excludedRangesRaw, ok := times["excludedRanges"]; ok {
		tflog.Debug(ctx, "Mapping excludedRanges", map[string]any{"excludedRanges": excludedRangesRaw})
		result.ExcludedRanges = mapStringSliceToSet(ctx, excludedRangesRaw, "excludedRanges")
	} else {
		tflog.Debug(ctx, "excludedRanges not found, setting to null")
		result.ExcludedRanges = types.SetNull(types.StringType)
	}

	if allDay, ok := times["allDay"].(bool); ok {
		tflog.Debug(ctx, "Mapping allDay", map[string]any{"allDay": allDay})
		result.AllDay = types.BoolValue(allDay)
	} else {
		tflog.Debug(ctx, "allDay not found or not a bool")
		result.AllDay = types.BoolNull()
	}

	if startTime, ok := times["startTime"].(string); ok {
		tflog.Debug(ctx, "Mapping startTime", map[string]any{"startTime": startTime})
		result.StartTime = types.StringValue(startTime)
	} else {
		tflog.Debug(ctx, "startTime not found or not a string")
		result.StartTime = types.StringNull()
	}

	if endTime, ok := times["endTime"].(string); ok {
		tflog.Debug(ctx, "Mapping endTime", map[string]any{"endTime": endTime})
		result.EndTime = types.StringValue(endTime)
	} else {
		tflog.Debug(ctx, "endTime not found or not a string")
		result.EndTime = types.StringNull()
	}

	if timeZone, ok := times["timeZone"].(string); ok {
		tflog.Debug(ctx, "Mapping timeZone", map[string]any{"timeZone": timeZone})
		result.TimeZone = types.StringValue(timeZone)
	} else {
		tflog.Debug(ctx, "timeZone not found or not a string")
		result.TimeZone = types.StringNull()
	}

	return result
}

func mapDeviceStates(ctx context.Context, deviceStatesRaw any) *ConditionalAccessDeviceStates {
	deviceStates, ok := deviceStatesRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "deviceStates is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessDeviceStates{}

	// Map includeStates (set)
	if includeStatesRaw, ok := deviceStates["includeStates"]; ok {
		tflog.Debug(ctx, "Mapping includeStates", map[string]any{"includeStates": includeStatesRaw})
		result.IncludeStates = mapStringSliceToSet(ctx, includeStatesRaw, "includeStates")
	} else {
		tflog.Debug(ctx, "includeStates not found, setting to null")
		result.IncludeStates = types.SetNull(types.StringType)
	}

	// Map excludeStates (set)
	if excludeStatesRaw, ok := deviceStates["excludeStates"]; ok {
		tflog.Debug(ctx, "Mapping excludeStates", map[string]any{"excludeStates": excludeStatesRaw})
		result.ExcludeStates = mapStringSliceToSet(ctx, excludeStatesRaw, "excludeStates")
	} else {
		tflog.Debug(ctx, "excludeStates not found, setting to null")
		result.ExcludeStates = types.SetNull(types.StringType)
	}

	return result
}

func mapGrantControls(ctx context.Context, grantControlsRaw any) *ConditionalAccessGrantControls {
	grantControls, ok := grantControlsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "grantControls is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessGrantControls{}

	if operator, ok := grantControls["operator"].(string); ok {
		tflog.Debug(ctx, "Mapping grantControls operator", map[string]any{"operator": operator})
		result.Operator = types.StringValue(operator)
	} else {
		tflog.Debug(ctx, "grantControls operator not found or not a string")
		result.Operator = types.StringNull()
	}

	// Map builtInControls (set)
	if builtInControlsRaw, ok := grantControls["builtInControls"]; ok {
		tflog.Debug(ctx, "Mapping builtInControls", map[string]any{"builtInControls": builtInControlsRaw})
		result.BuiltInControls = mapStringSliceToSet(ctx, builtInControlsRaw, "builtInControls")
	} else {
		tflog.Debug(ctx, "builtInControls not found, setting to null")
		result.BuiltInControls = types.SetNull(types.StringType)
	}

	// Map customAuthenticationFactors (set)
	if customAuthenticationFactorsRaw, ok := grantControls["customAuthenticationFactors"]; ok {
		tflog.Debug(ctx, "Mapping customAuthenticationFactors", map[string]any{"customAuthenticationFactors": customAuthenticationFactorsRaw})
		result.CustomAuthenticationFactors = mapStringSliceToSet(ctx, customAuthenticationFactorsRaw, "customAuthenticationFactors")
	} else {
		tflog.Debug(ctx, "customAuthenticationFactors not found, setting to null")
		result.CustomAuthenticationFactors = types.SetNull(types.StringType)
	}

	// Map termsOfUse (set)
	if termsOfUseRaw, ok := grantControls["termsOfUse"]; ok {
		tflog.Debug(ctx, "Mapping termsOfUse", map[string]any{"termsOfUse": termsOfUseRaw})
		result.TermsOfUse = mapStringSliceToSet(ctx, termsOfUseRaw, "termsOfUse")
	} else {
		tflog.Debug(ctx, "termsOfUse not found, setting to null")
		result.TermsOfUse = types.SetNull(types.StringType)
	}

	// Map authenticationStrength
	if authenticationStrengthRaw, ok := grantControls["authenticationStrength"]; ok {
		tflog.Debug(ctx, "Mapping authenticationStrength", map[string]any{"authenticationStrength": authenticationStrengthRaw})
		result.AuthenticationStrength = mapAuthenticationStrength(ctx, authenticationStrengthRaw)
	} else {
		tflog.Debug(ctx, "authenticationStrength not found")
		result.AuthenticationStrength = nil
	}

	return result
}

func mapSessionControls(ctx context.Context, sessionControlsRaw any) *ConditionalAccessSessionControls {
	sessionControls, ok := sessionControlsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "sessionControls is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessSessionControls{}

	if disableResilienceDefaults, ok := sessionControls["disableResilienceDefaults"].(bool); ok {
		tflog.Debug(ctx, "Mapping disableResilienceDefaults", map[string]any{"disableResilienceDefaults": disableResilienceDefaults})
		result.DisableResilienceDefaults = types.BoolValue(disableResilienceDefaults)
	} else {
		tflog.Debug(ctx, "disableResilienceDefaults not found or not a bool")
		result.DisableResilienceDefaults = types.BoolNull()
	}

	// Map applicationEnforcedRestrictions
	if applicationEnforcedRestrictionsRaw, ok := sessionControls["applicationEnforcedRestrictions"]; ok {
		tflog.Debug(ctx, "Mapping applicationEnforcedRestrictions", map[string]any{"applicationEnforcedRestrictions": applicationEnforcedRestrictionsRaw})
		result.ApplicationEnforcedRestrictions = mapApplicationEnforcedRestrictions(ctx, applicationEnforcedRestrictionsRaw)
	} else {
		tflog.Debug(ctx, "applicationEnforcedRestrictions not found")
		result.ApplicationEnforcedRestrictions = nil
	}

	// Map cloudAppSecurity
	if cloudAppSecurityRaw, ok := sessionControls["cloudAppSecurity"]; ok {
		tflog.Debug(ctx, "Mapping cloudAppSecurity", map[string]any{"cloudAppSecurity": cloudAppSecurityRaw})
		result.CloudAppSecurity = mapCloudAppSecurity(ctx, cloudAppSecurityRaw)
	} else {
		tflog.Debug(ctx, "cloudAppSecurity not found")
		result.CloudAppSecurity = nil
	}

	// Map signInFrequency
	if signInFrequencyRaw, ok := sessionControls["signInFrequency"]; ok {
		tflog.Debug(ctx, "Mapping signInFrequency", map[string]any{"signInFrequency": signInFrequencyRaw})
		result.SignInFrequency = mapSignInFrequency(ctx, signInFrequencyRaw)
	} else {
		tflog.Debug(ctx, "signInFrequency not found")
		result.SignInFrequency = nil
	}

	// Map persistentBrowser
	if persistentBrowserRaw, ok := sessionControls["persistentBrowser"]; ok {
		tflog.Debug(ctx, "Mapping persistentBrowser", map[string]any{"persistentBrowser": persistentBrowserRaw})
		result.PersistentBrowser = mapPersistentBrowser(ctx, persistentBrowserRaw)
	} else {
		tflog.Debug(ctx, "persistentBrowser not found")
		result.PersistentBrowser = nil
	}

	// Map continuousAccessEvaluation
	if continuousAccessEvaluationRaw, ok := sessionControls["continuousAccessEvaluation"]; ok {
		tflog.Debug(ctx, "Mapping continuousAccessEvaluation", map[string]any{"continuousAccessEvaluation": continuousAccessEvaluationRaw})
		result.ContinuousAccessEvaluation = mapContinuousAccessEvaluation(ctx, continuousAccessEvaluationRaw)
	} else {
		tflog.Debug(ctx, "continuousAccessEvaluation not found")
		result.ContinuousAccessEvaluation = nil
	}

	// Map secureSignInSession
	if secureSignInSessionRaw, ok := sessionControls["secureSignInSession"]; ok {
		tflog.Debug(ctx, "Mapping secureSignInSession", map[string]any{"secureSignInSession": secureSignInSessionRaw})
		result.SecureSignInSession = mapSecureSignInSession(ctx, secureSignInSessionRaw)
	} else {
		tflog.Debug(ctx, "secureSignInSession not found")
		result.SecureSignInSession = nil
	}

	return result
}

// Helper function to map authentication context class references with predefined value mapping
func mapAuthContextClassReferencesToSet(ctx context.Context, raw any, fieldName string) types.Set {
	tflog.Debug(ctx, fmt.Sprintf("Processing %s: %v (type: %T)", fieldName, raw, raw))

	if raw == nil {
		tflog.Debug(ctx, fmt.Sprintf("%s is null, returning null set", fieldName))
		return types.SetNull(types.StringType)
	}

	// Handle []any from JSON unmarshaling
	if slice, ok := raw.([]any); ok {
		tflog.Debug(ctx, fmt.Sprintf("%s is []any with %d elements", fieldName, len(slice)))

		// Convert []any to []string with reverse mapping
		stringSlice := make([]string, len(slice))
		for i, v := range slice {
			if str, ok := v.(string); ok {
				// Map predefined IDs back to their user-friendly names
				switch str {
				case "c1":
					stringSlice[i] = "require_trusted_device"
				case "c2":
					stringSlice[i] = "require_terms_of_use"
				case "c3":
					stringSlice[i] = "require_trusted_location"
				case "c4":
					stringSlice[i] = "require_strong_authentication"
				case "c5":
					stringSlice[i] = "required_trust_type:azure_ad_joined"
				case "c6":
					stringSlice[i] = "require_access_from_an_approved_app"
				case "c7":
					stringSlice[i] = "required_trust_type:hybrid_azure_ad_joined"
				default:
					stringSlice[i] = str
				}
				tflog.Trace(ctx, fmt.Sprintf("Element %d in %s: %q", i, fieldName, stringSlice[i]))
			} else {
				// Convert non-string values to strings
				stringSlice[i] = fmt.Sprintf("%v", v)
				tflog.Debug(ctx, fmt.Sprintf("Converting element %d in %s from %T to string: %q", i, fieldName, v, stringSlice[i]))
			}
		}

		// Use types.SetValueFrom to convert []string to types.Set
		setValue, diags := types.SetValueFrom(ctx, types.StringType, stringSlice)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating set for %s", fieldName), map[string]any{"diags": diags})
			return types.SetNull(types.StringType)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully created set for %s with %d elements", fieldName, len(stringSlice)))
		return setValue
	}

	tflog.Error(ctx, fmt.Sprintf("Unexpected type for %s: %T", fieldName, raw))
	return types.SetNull(types.StringType)
}

// Helper function to map string slices to Terraform sets, handling null vs empty arrays
func mapStringSliceToSet(ctx context.Context, raw any, fieldName string) types.Set {
	tflog.Debug(ctx, fmt.Sprintf("Processing %s: %v (type: %T)", fieldName, raw, raw))

	if raw == nil {
		tflog.Debug(ctx, fmt.Sprintf("%s is null, returning null set", fieldName))
		return types.SetNull(types.StringType)
	}

	// Handle []any from JSON unmarshaling
	if slice, ok := raw.([]any); ok {
		tflog.Debug(ctx, fmt.Sprintf("%s is []any with %d elements", fieldName, len(slice)))

		// Convert []any to []string
		stringSlice := make([]string, len(slice))
		for i, v := range slice {
			if str, ok := v.(string); ok {
				stringSlice[i] = str
				tflog.Trace(ctx, fmt.Sprintf("Element %d in %s: %q", i, fieldName, str))
			} else {
				// Convert non-string values to strings
				stringSlice[i] = fmt.Sprintf("%v", v)
				tflog.Debug(ctx, fmt.Sprintf("Converting element %d in %s from %T to string: %q", i, fieldName, v, stringSlice[i]))
			}
		}

		// Use types.SetValueFrom to convert []string to types.Set
		setValue, diags := types.SetValueFrom(ctx, types.StringType, stringSlice)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating set for %s", fieldName), map[string]any{"diags": diags})
			return types.SetNull(types.StringType)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully created set for %s with %d elements", fieldName, len(stringSlice)))
		return setValue
	}

	// Handle []any (alternative any representation)
	if slice, ok := raw.([]any); ok {
		tflog.Debug(ctx, fmt.Sprintf("%s is []any with %d elements", fieldName, len(slice)))

		// Convert []any to []string
		stringSlice := make([]string, len(slice))
		for i, v := range slice {
			if str, ok := v.(string); ok {
				stringSlice[i] = str
				tflog.Trace(ctx, fmt.Sprintf("Element %d in %s: %q", i, fieldName, str))
			} else {
				// Convert non-string values to strings
				stringSlice[i] = fmt.Sprintf("%v", v)
				tflog.Debug(ctx, fmt.Sprintf("Converting element %d in %s from %T to string: %q", i, fieldName, v, stringSlice[i]))
			}
		}

		// Use types.SetValueFrom to convert []string to types.Set
		setValue, diags := types.SetValueFrom(ctx, types.StringType, stringSlice)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating set for %s", fieldName), map[string]any{"diags": diags})
			return types.SetNull(types.StringType)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully created set for %s with %d elements", fieldName, len(stringSlice)))
		return setValue
	}

	// Handle []string directly
	if strSlice, ok := raw.([]string); ok {
		tflog.Debug(ctx, fmt.Sprintf("%s is []string with %d elements", fieldName, len(strSlice)))

		// Use types.SetValueFrom to convert []string to types.Set
		setValue, diags := types.SetValueFrom(ctx, types.StringType, strSlice)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating set for %s", fieldName), map[string]any{"diags": diags})
			return types.SetNull(types.StringType)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully created set for %s with %d elements", fieldName, len(strSlice)))
		return setValue
	}

	tflog.Debug(ctx, fmt.Sprintf("%s is not a recognized slice type, returning null set", fieldName))
	return types.SetNull(types.StringType)
}

// Helper functions for complex objects
func mapFilter(ctx context.Context, filterRaw any) *ConditionalAccessFilter {
	filter, ok := filterRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "filter is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessFilter{}

	if mode, ok := filter["mode"].(string); ok {
		tflog.Debug(ctx, "Mapping filter mode", map[string]any{"mode": mode})
		result.Mode = types.StringValue(mode)
	} else {
		tflog.Debug(ctx, "filter mode not found or not a string")
		result.Mode = types.StringNull()
	}

	if rule, ok := filter["rule"].(string); ok {
		tflog.Debug(ctx, "Mapping filter rule", map[string]any{"rule": rule})
		result.Rule = types.StringValue(rule)
	} else {
		tflog.Debug(ctx, "filter rule not found or not a string")
		result.Rule = types.StringNull()
	}

	return result
}

func mapGuestsOrExternalUsersToObject(ctx context.Context, raw any) types.Object {
	data, ok := raw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "guestsOrExternalUsers is not a map[string]any")
		return types.ObjectNull(map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"membership_kind": types.StringType,
					"members":         types.SetType{ElemType: types.StringType},
				},
			},
		})
	}

	attrs := map[string]attr.Value{}

	// Map guest_or_external_user_types
	if guestOrExternalUserTypesRaw, ok := data["guestOrExternalUserTypes"]; ok {
		attrs["guest_or_external_user_types"] = mapCommaSeparatedStringToSet(ctx, guestOrExternalUserTypesRaw, "guestOrExternalUserTypes")
	} else {
		attrs["guest_or_external_user_types"] = types.SetNull(types.StringType)
	}

	// Map external_tenants
	if externalTenantsRaw, ok := data["externalTenants"]; ok {
		if externalTenants, ok := externalTenantsRaw.(map[string]any); ok {
			externalTenantsAttrs := map[string]attr.Value{}

			if membershipKind, ok := externalTenants["membershipKind"].(string); ok {
				externalTenantsAttrs["membership_kind"] = types.StringValue(membershipKind)
			} else {
				externalTenantsAttrs["membership_kind"] = types.StringNull()
			}

			if membersRaw, ok := externalTenants["members"]; ok {
				externalTenantsAttrs["members"] = mapStringSliceToSet(ctx, membersRaw, "externalTenants.members")
			} else {
				externalTenantsAttrs["members"] = types.SetNull(types.StringType)
			}

			obj, diags := types.ObjectValue(map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			}, externalTenantsAttrs)

			if diags.HasError() {
				tflog.Error(ctx, "Error creating external_tenants object", map[string]any{"diags": diags})
				attrs["external_tenants"] = types.ObjectNull(map[string]attr.Type{
					"membership_kind": types.StringType,
					"members":         types.SetType{ElemType: types.StringType},
				})
			} else {
				attrs["external_tenants"] = obj
			}
		} else {
			attrs["external_tenants"] = types.ObjectNull(map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			})
		}
	} else {
		attrs["external_tenants"] = types.ObjectNull(map[string]attr.Type{
			"membership_kind": types.StringType,
			"members":         types.SetType{ElemType: types.StringType},
		})
	}

	obj, diags := types.ObjectValue(map[string]attr.Type{
		"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
		"external_tenants": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			},
		},
	}, attrs)

	if diags.HasError() {
		tflog.Error(ctx, "Error creating guestsOrExternalUsers object", map[string]any{"diags": diags})
		return types.ObjectNull(map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"membership_kind": types.StringType,
					"members":         types.SetType{ElemType: types.StringType},
				},
			},
		})
	}

	return obj
}

func mapAuthenticationStrength(ctx context.Context, authenticationStrengthRaw any) *ConditionalAccessAuthenticationStrength {
	authenticationStrength, ok := authenticationStrengthRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "authenticationStrength is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessAuthenticationStrength{}

	if id, ok := authenticationStrength["id"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength id", map[string]any{"id": id})
		// Map known GUIDs back to predefined string values for consistency
		switch id {
		case "00000000-0000-0000-0000-000000000002":
			result.ID = types.StringValue("multifactor_authentication")
		case "00000000-0000-0000-0000-000000000003":
			result.ID = types.StringValue("passwordless_mfa")
		case "00000000-0000-0000-0000-000000000004":
			result.ID = types.StringValue("phishing_resistant_mfa")
		default:
			result.ID = types.StringValue(id)
		}
	} else {
		tflog.Debug(ctx, "authenticationStrength id not found or not a string")
		result.ID = types.StringNull()
	}

	if displayName, ok := authenticationStrength["displayName"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength displayName", map[string]any{"displayName": displayName})
		result.DisplayName = types.StringValue(displayName)
	} else {
		tflog.Debug(ctx, "authenticationStrength displayName not found or not a string")
		result.DisplayName = types.StringNull()
	}

	if description, ok := authenticationStrength["description"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength description", map[string]any{"description": description})
		result.Description = types.StringValue(description)
	} else {
		tflog.Debug(ctx, "authenticationStrength description not found or not a string")
		result.Description = types.StringNull()
	}

	if policyType, ok := authenticationStrength["policyType"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength policyType", map[string]any{"policyType": policyType})
		result.PolicyType = types.StringValue(policyType)
	} else {
		tflog.Debug(ctx, "authenticationStrength policyType not found or not a string")
		result.PolicyType = types.StringNull()
	}

	if requirementsSatisfied, ok := authenticationStrength["requirementsSatisfied"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength requirementsSatisfied", map[string]any{"requirementsSatisfied": requirementsSatisfied})
		result.RequirementsSatisfied = types.StringValue(requirementsSatisfied)
	} else {
		tflog.Debug(ctx, "authenticationStrength requirementsSatisfied not found or not a string")
		result.RequirementsSatisfied = types.StringNull()
	}

	if allowedCombinationsRaw, ok := authenticationStrength["allowedCombinations"]; ok {
		tflog.Debug(ctx, "Mapping authenticationStrength allowedCombinations", map[string]any{"allowedCombinations": allowedCombinationsRaw})
		result.AllowedCombinations = mapStringSliceToSet(ctx, allowedCombinationsRaw, "allowedCombinations")
	} else {
		tflog.Debug(ctx, "authenticationStrength allowedCombinations not found, setting to null")
		result.AllowedCombinations = types.SetNull(types.StringType)
	}

	if createdDateTime, ok := authenticationStrength["createdDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength createdDateTime", map[string]any{"createdDateTime": createdDateTime})
		result.CreatedDateTime = types.StringValue(createdDateTime)
	} else {
		tflog.Debug(ctx, "authenticationStrength createdDateTime not found or not a string")
		result.CreatedDateTime = types.StringNull()
	}

	if modifiedDateTime, ok := authenticationStrength["modifiedDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping authenticationStrength modifiedDateTime", map[string]any{"modifiedDateTime": modifiedDateTime})
		result.ModifiedDateTime = types.StringValue(modifiedDateTime)
	} else {
		tflog.Debug(ctx, "authenticationStrength modifiedDateTime not found or not a string")
		result.ModifiedDateTime = types.StringNull()
	}

	return result
}

func mapApplicationEnforcedRestrictions(ctx context.Context, applicationEnforcedRestrictionsRaw any) *ConditionalAccessApplicationEnforcedRestrictions {
	applicationEnforcedRestrictions, ok := applicationEnforcedRestrictionsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "applicationEnforcedRestrictions is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessApplicationEnforcedRestrictions{}

	if isEnabled, ok := applicationEnforcedRestrictions["isEnabled"].(bool); ok {
		tflog.Debug(ctx, "Mapping applicationEnforcedRestrictions isEnabled", map[string]any{"isEnabled": isEnabled})
		result.IsEnabled = types.BoolValue(isEnabled)
	} else {
		tflog.Debug(ctx, "applicationEnforcedRestrictions isEnabled not found or not a bool")
		result.IsEnabled = types.BoolNull()
	}

	return result
}

func mapCloudAppSecurity(ctx context.Context, cloudAppSecurityRaw any) *ConditionalAccessCloudAppSecurity {
	cloudAppSecurity, ok := cloudAppSecurityRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "cloudAppSecurity is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessCloudAppSecurity{}

	if isEnabled, ok := cloudAppSecurity["isEnabled"].(bool); ok {
		tflog.Debug(ctx, "Mapping cloudAppSecurity isEnabled", map[string]any{"isEnabled": isEnabled})
		result.IsEnabled = types.BoolValue(isEnabled)
	} else {
		tflog.Debug(ctx, "cloudAppSecurity isEnabled not found or not a bool")
		result.IsEnabled = types.BoolNull()
	}

	if cloudAppSecurityType, ok := cloudAppSecurity["cloudAppSecurityType"].(string); ok {
		tflog.Debug(ctx, "Mapping cloudAppSecurity cloudAppSecurityType", map[string]any{"cloudAppSecurityType": cloudAppSecurityType})
		result.CloudAppSecurityType = types.StringValue(cloudAppSecurityType)
	} else {
		tflog.Debug(ctx, "cloudAppSecurity cloudAppSecurityType not found or not a string")
		result.CloudAppSecurityType = types.StringNull()
	}

	return result
}

func mapSignInFrequency(ctx context.Context, signInFrequencyRaw any) *ConditionalAccessSignInFrequency {
	signInFrequency, ok := signInFrequencyRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "signInFrequency is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessSignInFrequency{}

	if isEnabled, ok := signInFrequency["isEnabled"].(bool); ok {
		tflog.Debug(ctx, "Mapping signInFrequency isEnabled", map[string]any{"isEnabled": isEnabled})
		result.IsEnabled = types.BoolValue(isEnabled)
	} else {
		tflog.Debug(ctx, "signInFrequency isEnabled not found or not a bool")
		result.IsEnabled = types.BoolNull()
	}

	if signInFrequencyType, ok := signInFrequency["type"].(string); ok {
		tflog.Debug(ctx, "Mapping signInFrequency type", map[string]any{"type": signInFrequencyType})
		result.Type = types.StringValue(signInFrequencyType)
	} else {
		tflog.Debug(ctx, "signInFrequency type not found or not a string")
		result.Type = types.StringNull()
	}

	if value, ok := signInFrequency["value"]; ok {
		if intValue, ok := value.(int); ok {
			tflog.Debug(ctx, "Mapping signInFrequency value (int)", map[string]any{"value": intValue})
			result.Value = types.Int64Value(int64(intValue))
		} else if floatValue, ok := value.(float64); ok {
			tflog.Debug(ctx, "Mapping signInFrequency value (float64)", map[string]any{"value": floatValue})
			result.Value = types.Int64Value(int64(floatValue))
		} else {
			tflog.Debug(ctx, "signInFrequency value not a valid number type")
			result.Value = types.Int64Null()
		}
	} else {
		tflog.Debug(ctx, "signInFrequency value not found")
		result.Value = types.Int64Null()
	}

	if authenticationType, ok := signInFrequency["authenticationType"].(string); ok {
		tflog.Debug(ctx, "Mapping signInFrequency authenticationType", map[string]any{"authenticationType": authenticationType})
		result.AuthenticationType = types.StringValue(authenticationType)
	} else {
		tflog.Debug(ctx, "signInFrequency authenticationType not found or not a string")
		result.AuthenticationType = types.StringNull()
	}

	if frequencyInterval, ok := signInFrequency["frequencyInterval"].(string); ok {
		tflog.Debug(ctx, "Mapping signInFrequency frequencyInterval", map[string]any{"frequencyInterval": frequencyInterval})
		result.FrequencyInterval = types.StringValue(frequencyInterval)
	} else {
		tflog.Debug(ctx, "signInFrequency frequencyInterval not found or not a string")
		result.FrequencyInterval = types.StringNull()
	}

	return result
}

func mapPersistentBrowser(ctx context.Context, persistentBrowserRaw any) *ConditionalAccessPersistentBrowser {
	persistentBrowser, ok := persistentBrowserRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "persistentBrowser is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessPersistentBrowser{}

	if isEnabled, ok := persistentBrowser["isEnabled"].(bool); ok {
		tflog.Debug(ctx, "Mapping persistentBrowser isEnabled", map[string]any{"isEnabled": isEnabled})
		result.IsEnabled = types.BoolValue(isEnabled)
	} else {
		tflog.Debug(ctx, "persistentBrowser isEnabled not found or not a bool")
		result.IsEnabled = types.BoolNull()
	}

	if mode, ok := persistentBrowser["mode"].(string); ok {
		tflog.Debug(ctx, "Mapping persistentBrowser mode", map[string]any{"mode": mode})
		result.Mode = types.StringValue(mode)
	} else {
		tflog.Debug(ctx, "persistentBrowser mode not found or not a string")
		result.Mode = types.StringNull()
	}

	return result
}

func mapContinuousAccessEvaluation(ctx context.Context, continuousAccessEvaluationRaw any) *ConditionalAccessContinuousAccessEvaluation {
	continuousAccessEvaluation, ok := continuousAccessEvaluationRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "continuousAccessEvaluation is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessContinuousAccessEvaluation{}

	if mode, ok := continuousAccessEvaluation["mode"].(string); ok {
		tflog.Debug(ctx, "Mapping continuousAccessEvaluation mode", map[string]any{"mode": mode})
		result.Mode = types.StringValue(mode)
	} else {
		tflog.Debug(ctx, "continuousAccessEvaluation mode not found or not a string")
		result.Mode = types.StringNull()
	}

	return result
}

func mapSecureSignInSession(ctx context.Context, secureSignInSessionRaw any) *ConditionalAccessSecureSignInSession {
	secureSignInSession, ok := secureSignInSessionRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "secureSignInSession is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessSecureSignInSession{}

	if isEnabled, ok := secureSignInSession["isEnabled"].(bool); ok {
		tflog.Debug(ctx, "Mapping secureSignInSession isEnabled", map[string]any{"isEnabled": isEnabled})
		result.IsEnabled = types.BoolValue(isEnabled)
	} else {
		tflog.Debug(ctx, "secureSignInSession isEnabled not found or not a bool")
		result.IsEnabled = types.BoolNull()
	}

	return result
}

// Helper function to map comma-separated string to Terraform set
func mapCommaSeparatedStringToSet(ctx context.Context, raw any, fieldName string) types.Set {
	tflog.Debug(ctx, fmt.Sprintf("Processing comma-separated string %s: %v (type: %T)", fieldName, raw, raw))

	if raw == nil {
		tflog.Debug(ctx, fmt.Sprintf("%s is null, returning null set", fieldName))
		return types.SetNull(types.StringType)
	}

	// Handle string from JSON unmarshaling
	if str, ok := raw.(string); ok {
		tflog.Debug(ctx, fmt.Sprintf("%s is string: %s", fieldName, str))

		if str == "" {
			// Empty string should return empty set, not null
			setValue, diags := types.SetValueFrom(ctx, types.StringType, []string{})
			if diags.HasError() {
				tflog.Error(ctx, fmt.Sprintf("Error creating empty set for %s", fieldName), map[string]any{"diags": diags})
				return types.SetNull(types.StringType)
			}
			return setValue
		}

		// Split comma-separated string
		stringSlice := strings.Split(str, ",")
		// Trim whitespace from each value
		for i, value := range stringSlice {
			stringSlice[i] = strings.TrimSpace(value)
		}

		// Use types.SetValueFrom to convert []string to types.Set
		setValue, diags := types.SetValueFrom(ctx, types.StringType, stringSlice)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating set for %s", fieldName), map[string]any{"diags": diags})
			return types.SetNull(types.StringType)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully created set for %s with %d elements", fieldName, len(stringSlice)))
		return setValue
	}

	tflog.Debug(ctx, fmt.Sprintf("%s is not a string, returning null set", fieldName))
	return types.SetNull(types.StringType)
}
