package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteResourceStateToTerraform maps the remote conditional access policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource map[string]interface{}) {
	// Basic properties using helpers
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	data.DisplayName = convert.MapToFrameworkString(remoteResource, "displayName")
	data.State = convert.MapToFrameworkString(remoteResource, "state")
	data.CreatedDateTime = convert.MapToFrameworkString(remoteResource, "createdDateTime")
	data.ModifiedDateTime = convert.MapToFrameworkString(remoteResource, "modifiedDateTime")
	data.DeletedDateTime = convert.MapToFrameworkString(remoteResource, "deletedDateTime")
	data.TemplateId = convert.MapToFrameworkString(remoteResource, "templateId")
	data.PartialEnablementStrategy = convert.MapToFrameworkString(remoteResource, "partialEnablementStrategy")

	// Map conditions
	if conditionsData, ok := remoteResource["conditions"].(map[string]interface{}); ok {
		data.Conditions = mapConditionsFromRemote(ctx, conditionsData)
	}

	// Map grant controls - only set if present in API response
	if grantControlsData, ok := remoteResource["grantControls"].(map[string]interface{}); ok {
		data.GrantControls = mapGrantControlsFromRemote(ctx, grantControlsData)
	} else {
		data.GrantControls = nil
	}

	// Map session controls - only set if present in API response
	if sessionControlsData, ok := remoteResource["sessionControls"].(map[string]interface{}); ok && len(sessionControlsData) > 0 {
		data.SessionControls = mapSessionControlsFromRemote(ctx, sessionControlsData)
	} else {
		data.SessionControls = nil
	}
}

// mapConditionsFromRemote maps conditions from the API response
func mapConditionsFromRemote(ctx context.Context, conditionsData map[string]interface{}) *ConditionalAccessConditions {
	conditions := &ConditionalAccessConditions{}

	// Client app types
	conditions.ClientAppTypes = convert.MapToFrameworkStringSet(ctx, conditionsData, "clientAppTypes")

	// Applications - always create the object since it's required in schema
	if applicationsData, ok := conditionsData["applications"].(map[string]interface{}); ok {
		conditions.Applications = mapApplicationsFromRemote(ctx, applicationsData)
	} else {
		// Create empty applications object to maintain consistency
		conditions.Applications = &ConditionalAccessApplications{
			IncludeApplications:                         types.SetNull(types.StringType),
			ExcludeApplications:                         types.SetNull(types.StringType),
			IncludeUserActions:                          types.SetNull(types.StringType),
			IncludeAuthenticationContextClassReferences: types.SetNull(types.StringType),
			ApplicationFilter:                           nil,
		}
	}

	// Users - always create the object since it's required in schema
	if usersData, ok := conditionsData["users"].(map[string]interface{}); ok {
		conditions.Users = mapUsersFromRemote(ctx, usersData)
	} else {
		// Create empty users object to maintain consistency
		conditions.Users = &ConditionalAccessUsers{
			IncludeUsers:                 types.SetNull(types.StringType),
			ExcludeUsers:                 types.SetNull(types.StringType),
			IncludeGroups:                types.SetNull(types.StringType),
			ExcludeGroups:                types.SetNull(types.StringType),
			IncludeRoles:                 types.SetNull(types.StringType),
			ExcludeRoles:                 types.SetNull(types.StringType),
			IncludeGuestsOrExternalUsers: types.ObjectNull(getGuestsOrExternalUsersAttributeTypes()),
			ExcludeGuestsOrExternalUsers: types.ObjectNull(getGuestsOrExternalUsersAttributeTypes()),
		}
	}

	// Platforms - preserve the structure if any platform data exists
	if platformsData, ok := conditionsData["platforms"].(map[string]interface{}); ok {
		conditions.Platforms = mapPlatformsFromRemote(ctx, platformsData)
	} else {
		// Set to nil if not present in API response
		conditions.Platforms = nil
	}

	// Locations - preserve the structure if any location data exists
	if locationsData, ok := conditionsData["locations"].(map[string]interface{}); ok {
		conditions.Locations = mapLocationsFromRemote(ctx, locationsData)
	} else {
		conditions.Locations = nil
	}

	// Devices - preserve the structure if any device data exists
	if devicesData, ok := conditionsData["devices"].(map[string]interface{}); ok {
		conditions.Devices = mapDevicesFromRemote(ctx, devicesData)
	} else {
		// Set to nil if not present in API response
		conditions.Devices = nil
	}

	// Risk levels - preserve empty sets for fields that can be configured as empty arrays
	conditions.SignInRiskLevels = convert.MapToFrameworkStringSet(ctx, conditionsData, "signInRiskLevels")
	conditions.UserRiskLevels = convert.MapToFrameworkStringSet(ctx, conditionsData, "userRiskLevels")

	if _, ok := conditionsData["servicePrincipalRiskLevels"]; ok {
		conditions.ServicePrincipalRiskLevels = convert.MapToFrameworkStringSet(ctx, conditionsData, "servicePrincipalRiskLevels")
	} else {
		// Always return a value (empty set or populated set) instead of null
		conditions.ServicePrincipalRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}

	return conditions
}

// mapApplicationsFromRemote maps applications from the API response
func mapApplicationsFromRemote(ctx context.Context, applicationsData map[string]interface{}) *ConditionalAccessApplications {
	applications := &ConditionalAccessApplications{}

	applications.IncludeApplications = convert.MapToFrameworkStringSet(ctx, applicationsData, "includeApplications")

	// These fields are configured as empty arrays in Terraform, so preserve empty sets
	applications.ExcludeApplications = convert.MapToFrameworkStringSet(ctx, applicationsData, "excludeApplications")
	applications.IncludeUserActions = convert.MapToFrameworkStringSet(ctx, applicationsData, "includeUserActions")
	applications.IncludeAuthenticationContextClassReferences = convert.MapToFrameworkStringSet(ctx, applicationsData, "includeAuthenticationContextClassReferences")

	if applicationFilterData, ok := applicationsData["applicationFilter"].(map[string]interface{}); ok {
		filter := &ConditionalAccessFilter{
			Mode: convert.MapToFrameworkString(applicationFilterData, "mode"),
			Rule: convert.MapToFrameworkString(applicationFilterData, "rule"),
		}
		applications.ApplicationFilter = filter
	}

	return applications
}

// mapUsersFromRemote maps users from the API response
func mapUsersFromRemote(ctx context.Context, usersData map[string]interface{}) *ConditionalAccessUsers {
	users := &ConditionalAccessUsers{}

	users.IncludeUsers = convert.MapToFrameworkStringSet(ctx, usersData, "includeUsers")

	// These fields are configured as empty arrays in Terraform, so preserve empty sets
	users.ExcludeUsers = convert.MapToFrameworkStringSet(ctx, usersData, "excludeUsers")
	users.IncludeGroups = convert.MapToFrameworkStringSet(ctx, usersData, "includeGroups")
	users.ExcludeGroups = convert.MapToFrameworkStringSet(ctx, usersData, "excludeGroups")
	users.IncludeRoles = convert.MapToFrameworkStringSet(ctx, usersData, "includeRoles")

	users.ExcludeRoles = convert.MapToFrameworkStringSet(ctx, usersData, "excludeRoles")

	// Handle guests or external users - preserve the structure if configured in Terraform but API doesn't return it
	if includeGuestsData, ok := usersData["includeGuestsOrExternalUsers"].(map[string]interface{}); ok {
		users.IncludeGuestsOrExternalUsers = mapGuestsOrExternalUsersToObject(ctx, includeGuestsData)
	} else {
		// Set to null object when not present in API response
		users.IncludeGuestsOrExternalUsers = types.ObjectNull(getGuestsOrExternalUsersAttributeTypes())
	}

	if excludeGuestsData, ok := usersData["excludeGuestsOrExternalUsers"].(map[string]interface{}); ok {
		users.ExcludeGuestsOrExternalUsers = mapGuestsOrExternalUsersToObject(ctx, excludeGuestsData)
	} else {
		// Set to null object when not present in API response
		users.ExcludeGuestsOrExternalUsers = types.ObjectNull(getGuestsOrExternalUsersAttributeTypes())
	}

	return users
}

// mapPlatformsFromRemote maps platforms from the API response
func mapPlatformsFromRemote(ctx context.Context, platformsData map[string]interface{}) *ConditionalAccessPlatforms {
	platforms := &ConditionalAccessPlatforms{}

	platforms.IncludePlatforms = convert.MapToFrameworkStringSet(ctx, platformsData, "includePlatforms")
	platforms.ExcludePlatforms = convert.MapToFrameworkStringSet(ctx, platformsData, "excludePlatforms")

	return platforms
}

// mapLocationsFromRemote maps locations from the API response
func mapLocationsFromRemote(ctx context.Context, locationsData map[string]interface{}) *ConditionalAccessLocations {
	locations := &ConditionalAccessLocations{}

	locations.IncludeLocations = convert.MapToFrameworkStringSet(ctx, locationsData, "includeLocations")
	locations.ExcludeLocations = convert.MapToFrameworkStringSet(ctx, locationsData, "excludeLocations")

	return locations
}

// mapDevicesFromRemote maps devices from the API response
func mapDevicesFromRemote(ctx context.Context, devicesData map[string]interface{}) *ConditionalAccessDevices {
	devices := &ConditionalAccessDevices{}

	// Fix: Use SetNull instead of empty sets for device fields when not present
	devices.IncludeDevices = convert.MapToFrameworkStringSet(ctx, devicesData, "includeDevices")

	devices.ExcludeDevices = convert.MapToFrameworkStringSet(ctx, devicesData, "excludeDevices")

	devices.IncludeDeviceStates = convert.MapToFrameworkStringSet(ctx, devicesData, "includeDeviceStates")

	devices.ExcludeDeviceStates = convert.MapToFrameworkStringSet(ctx, devicesData, "excludeDeviceStates")

	if deviceFilterData, ok := devicesData["deviceFilter"].(map[string]interface{}); ok {
		filter := &ConditionalAccessFilter{
			Mode: convert.MapToFrameworkString(deviceFilterData, "mode"),
			Rule: convert.MapToFrameworkString(deviceFilterData, "rule"),
		}
		devices.DeviceFilter = filter
	}

	return devices
}

// mapGrantControlsFromRemote maps grant controls from the API response
func mapGrantControlsFromRemote(ctx context.Context, grantControlsData map[string]interface{}) *ConditionalAccessGrantControls {
	grantControls := &ConditionalAccessGrantControls{
		Operator: convert.MapToFrameworkString(grantControlsData, "operator"),
	}

	grantControls.BuiltInControls = convert.MapToFrameworkStringSet(ctx, grantControlsData, "builtInControls")

	// These fields can be configured as empty arrays in Terraform, so preserve empty sets
	grantControls.CustomAuthenticationFactors = convert.MapToFrameworkStringSet(ctx, grantControlsData, "customAuthenticationFactors")
	grantControls.TermsOfUse = convert.MapToFrameworkStringSet(ctx, grantControlsData, "termsOfUse")

	if authStrengthData, ok := grantControlsData["authenticationStrength"].(map[string]interface{}); ok {
		authStrength := &ConditionalAccessAuthenticationStrength{
			ID:                    convert.MapToFrameworkString(authStrengthData, "id"),
			DisplayName:           convert.MapToFrameworkString(authStrengthData, "displayName"),
			Description:           convert.MapToFrameworkString(authStrengthData, "description"),
			PolicyType:            convert.MapToFrameworkString(authStrengthData, "policyType"),
			RequirementsSatisfied: convert.MapToFrameworkString(authStrengthData, "requirementsSatisfied"),
			CreatedDateTime:       convert.MapToFrameworkString(authStrengthData, "createdDateTime"),
			ModifiedDateTime:      convert.MapToFrameworkString(authStrengthData, "modifiedDateTime"),
		}

		authStrength.AllowedCombinations = convert.MapToFrameworkStringSet(ctx, authStrengthData, "allowedCombinations")

		grantControls.AuthenticationStrength = authStrength
	}

	return grantControls
}

// mapSessionControlsFromRemote maps session controls from the API response
func mapSessionControlsFromRemote(ctx context.Context, sessionControlsData map[string]interface{}) *ConditionalAccessSessionControls {
	sessionControls := &ConditionalAccessSessionControls{
		DisableResilienceDefaults: convert.MapToFrameworkBool(sessionControlsData, "disableResilienceDefaults"),
	}

	if appEnforcedData, ok := sessionControlsData["applicationEnforcedRestrictions"].(map[string]interface{}); ok {
		sessionControls.ApplicationEnforcedRestrictions = &ConditionalAccessApplicationEnforcedRestrictions{
			IsEnabled: convert.MapToFrameworkBool(appEnforcedData, "isEnabled"),
		}
	}

	if cloudAppSecData, ok := sessionControlsData["cloudAppSecurity"].(map[string]interface{}); ok {
		sessionControls.CloudAppSecurity = &ConditionalAccessCloudAppSecurity{
			IsEnabled:            convert.MapToFrameworkBool(cloudAppSecData, "isEnabled"),
			CloudAppSecurityType: convert.MapToFrameworkString(cloudAppSecData, "cloudAppSecurityType"),
		}
	}

	if signInFreqData, ok := sessionControlsData["signInFrequency"].(map[string]interface{}); ok {
		signInFreq := &ConditionalAccessSignInFrequency{
			IsEnabled: convert.MapToFrameworkBool(signInFreqData, "isEnabled"),
			Type:      convert.MapToFrameworkString(signInFreqData, "type"),
		}

		// Handle value conversion using helper
		signInFreq.Value = convert.MapToFrameworkInt64(signInFreqData, "value")

		// Always preserve authentication_type and frequency_interval from API response
		signInFreq.AuthenticationType = convert.MapToFrameworkString(signInFreqData, "authenticationType")
		signInFreq.FrequencyInterval = convert.MapToFrameworkString(signInFreqData, "frequencyInterval")

		sessionControls.SignInFrequency = signInFreq
	}

	if persistentBrowserData, ok := sessionControlsData["persistentBrowser"].(map[string]interface{}); ok {
		sessionControls.PersistentBrowser = &ConditionalAccessPersistentBrowser{
			IsEnabled: convert.MapToFrameworkBool(persistentBrowserData, "isEnabled"),
			Mode:      convert.MapToFrameworkString(persistentBrowserData, "mode"),
		}
	}

	if caeData, ok := sessionControlsData["continuousAccessEvaluation"].(map[string]interface{}); ok {
		sessionControls.ContinuousAccessEvaluation = &ConditionalAccessContinuousAccessEvaluation{
			Mode: convert.MapToFrameworkString(caeData, "mode"),
		}
	}

	if secureSignInData, ok := sessionControlsData["secureSignInSession"].(map[string]interface{}); ok {
		sessionControls.SecureSignInSession = &ConditionalAccessSecureSignInSession{
			IsEnabled: convert.MapToFrameworkBool(secureSignInData, "isEnabled"),
		}
	}

	return sessionControls
}

// getGuestsOrExternalUsersAttributeTypes returns the attribute types for the guests or external users object
func getGuestsOrExternalUsersAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
		"external_tenants": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			},
		},
	}
}

// mapGuestsOrExternalUsersToObject converts API data to types.Object
func mapGuestsOrExternalUsersToObject(ctx context.Context, guestsData map[string]interface{}) types.Object {
	attributes := make(map[string]attr.Value)

	// Map guest_or_external_user_types
	attributes["guest_or_external_user_types"] = convert.MapToFrameworkStringSet(ctx, guestsData, "guestOrExternalUserTypes")

	// Map external_tenants
	if externalTenantsData, ok := guestsData["externalTenants"].(map[string]interface{}); ok {
		externalTenantsAttrs := map[string]attr.Value{
			"membership_kind": convert.MapToFrameworkString(externalTenantsData, "membershipKind"),
			"members":         convert.MapToFrameworkStringSet(ctx, externalTenantsData, "members"),
		}
		attributes["external_tenants"] = types.ObjectValueMust(
			map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			},
			externalTenantsAttrs,
		)
	} else {
		attributes["external_tenants"] = types.ObjectNull(map[string]attr.Type{
			"membership_kind": types.StringType,
			"members":         types.SetType{ElemType: types.StringType},
		})
	}

	return types.ObjectValueMust(getGuestsOrExternalUsersAttributeTypes(), attributes)
}
