package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteResourceStateToTerraform maps the remote conditional access policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource map[string]interface{}) {
	// Basic properties using helpers
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	data.DisplayName = convert.GraphToFrameworkString(getStringPtr(remoteResource, "displayName"))
	data.State = convert.GraphToFrameworkString(getStringPtr(remoteResource, "state"))
	data.CreatedDateTime = convert.GraphToFrameworkString(getStringPtr(remoteResource, "createdDateTime"))
	data.ModifiedDateTime = convert.GraphToFrameworkString(getStringPtr(remoteResource, "modifiedDateTime"))
	data.DeletedDateTime = convert.GraphToFrameworkString(getStringPtr(remoteResource, "deletedDateTime"))
	data.TemplateId = convert.GraphToFrameworkString(getStringPtr(remoteResource, "templateId"))
	data.PartialEnablementStrategy = convert.GraphToFrameworkString(getStringPtr(remoteResource, "partialEnablementStrategy"))

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
	if sessionControlsData, ok := remoteResource["sessionControls"].(map[string]interface{}); ok {
		data.SessionControls = mapSessionControlsFromRemote(ctx, sessionControlsData)
	} else {
		data.SessionControls = nil
	}
}

// Helper function to get string pointer from map
func getStringPtr(data map[string]interface{}, key string) *string {
	if value, ok := data[key].(string); ok {
		return &value
	}
	return nil
}

// mapConditionsFromRemote maps conditions from the API response
func mapConditionsFromRemote(ctx context.Context, conditionsData map[string]interface{}) *ConditionalAccessConditions {
	conditions := &ConditionalAccessConditions{}

	// Client app types
	if clientAppTypes, ok := conditionsData["clientAppTypes"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(clientAppTypes))
		for _, item := range clientAppTypes {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.ClientAppTypes = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		conditions.ClientAppTypes = types.SetNull(types.StringType)
	}

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
			IncludeGuestsOrExternalUsers: nil,
			ExcludeGuestsOrExternalUsers: nil,
		}
	}

	// Platforms - only set if present in API response
	if platformsData, ok := conditionsData["platforms"].(map[string]interface{}); ok {
		conditions.Platforms = mapPlatformsFromRemote(ctx, platformsData)
	} else {
		// Set to nil if not present in API response
		conditions.Platforms = nil
	}

	// Locations - only set if present in API response
	if locationsData, ok := conditionsData["locations"].(map[string]interface{}); ok {
		conditions.Locations = mapLocationsFromRemote(ctx, locationsData)
	} else {
		conditions.Locations = nil
	}

	// Devices - only set if present in API response
	if devicesData, ok := conditionsData["devices"].(map[string]interface{}); ok {
		conditions.Devices = mapDevicesFromRemote(ctx, devicesData)
	} else {
		// Set to nil if not present in API response
		conditions.Devices = nil
	}

	// Risk levels - use PreserveEmpty for fields configured as empty arrays in Terraform
	if signInRiskLevels, ok := conditionsData["signInRiskLevels"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(signInRiskLevels))
		for _, item := range signInRiskLevels {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.SignInRiskLevels = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		conditions.SignInRiskLevels = types.SetNull(types.StringType)
	}

	if userRiskLevels, ok := conditionsData["userRiskLevels"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(userRiskLevels))
		for _, item := range userRiskLevels {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.UserRiskLevels = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		conditions.UserRiskLevels = types.SetNull(types.StringType)
	}

	if servicePrincipalRiskLevels, ok := conditionsData["servicePrincipalRiskLevels"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(servicePrincipalRiskLevels))
		for _, item := range servicePrincipalRiskLevels {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.ServicePrincipalRiskLevels = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		conditions.ServicePrincipalRiskLevels = types.SetNull(types.StringType)
	}

	return conditions
}

// mapApplicationsFromRemote maps applications from the API response
func mapApplicationsFromRemote(ctx context.Context, applicationsData map[string]interface{}) *ConditionalAccessApplications {
	applications := &ConditionalAccessApplications{}

	if includeApplications, ok := applicationsData["includeApplications"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeApplications))
		for _, item := range includeApplications {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		applications.IncludeApplications = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		applications.IncludeApplications = types.SetNull(types.StringType)
	}

	// These fields are configured as empty arrays in Terraform, so preserve empty sets
	if excludeApplications, ok := applicationsData["excludeApplications"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeApplications))
		for _, item := range excludeApplications {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		applications.ExcludeApplications = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		applications.ExcludeApplications = types.SetNull(types.StringType)
	}

	if includeUserActions, ok := applicationsData["includeUserActions"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeUserActions))
		for _, item := range includeUserActions {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		applications.IncludeUserActions = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		applications.IncludeUserActions = types.SetNull(types.StringType)
	}

	// Only set this field if it's present in the API response
	if _, hasAuthContext := applicationsData["includeAuthenticationContextClassReferences"]; hasAuthContext {
		if authContextRefs, ok := applicationsData["includeAuthenticationContextClassReferences"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(authContextRefs))
			for _, item := range authContextRefs {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			applications.IncludeAuthenticationContextClassReferences = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		} else {
			applications.IncludeAuthenticationContextClassReferences = types.SetNull(types.StringType)
		}
	} else {
		applications.IncludeAuthenticationContextClassReferences = types.SetNull(types.StringType)
	}

	if applicationFilterData, ok := applicationsData["applicationFilter"].(map[string]interface{}); ok {
		filter := &ConditionalAccessFilter{
			Mode: convert.GraphToFrameworkString(getStringPtr(applicationFilterData, "mode")),
			Rule: convert.GraphToFrameworkString(getStringPtr(applicationFilterData, "rule")),
		}
		applications.ApplicationFilter = filter
	}

	return applications
}

// mapUsersFromRemote maps users from the API response
func mapUsersFromRemote(ctx context.Context, usersData map[string]interface{}) *ConditionalAccessUsers {
	users := &ConditionalAccessUsers{}

	if includeUsers, ok := usersData["includeUsers"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeUsers))
		for _, item := range includeUsers {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.IncludeUsers = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		users.IncludeUsers = types.SetNull(types.StringType)
	}

	// These fields are configured as empty arrays in Terraform, so preserve empty sets
	if excludeUsers, ok := usersData["excludeUsers"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeUsers))
		for _, item := range excludeUsers {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.ExcludeUsers = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		users.ExcludeUsers = types.SetNull(types.StringType)
	}

	if includeGroups, ok := usersData["includeGroups"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeGroups))
		for _, item := range includeGroups {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.IncludeGroups = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		users.IncludeGroups = types.SetNull(types.StringType)
	}

	if excludeGroups, ok := usersData["excludeGroups"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeGroups))
		for _, item := range excludeGroups {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.ExcludeGroups = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		users.ExcludeGroups = types.SetNull(types.StringType)
	}

	// Only set roles if they're present in the API response
	if _, hasIncludeRoles := usersData["includeRoles"]; hasIncludeRoles {
		if includeRoles, ok := usersData["includeRoles"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(includeRoles))
			for _, item := range includeRoles {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			users.IncludeRoles = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		} else {
			users.IncludeRoles = types.SetNull(types.StringType)
		}
	} else {
		users.IncludeRoles = types.SetNull(types.StringType)
	}

	if excludeRoles, ok := usersData["excludeRoles"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeRoles))
		for _, item := range excludeRoles {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.ExcludeRoles = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		users.ExcludeRoles = types.SetNull(types.StringType)
	}

	// Handle guests or external users
	if includeGuestsData, ok := usersData["includeGuestsOrExternalUsers"].(map[string]interface{}); ok {
		users.IncludeGuestsOrExternalUsers = mapGuestsOrExternalUsersFromRemote(ctx, includeGuestsData)
	}

	if excludeGuestsData, ok := usersData["excludeGuestsOrExternalUsers"].(map[string]interface{}); ok {
		users.ExcludeGuestsOrExternalUsers = mapGuestsOrExternalUsersFromRemote(ctx, excludeGuestsData)
	}

	return users
}

// mapGuestsOrExternalUsersFromRemote maps guests or external users configuration
func mapGuestsOrExternalUsersFromRemote(ctx context.Context, guestsData map[string]interface{}) *ConditionalAccessGuestsOrExternalUsers {
	guests := &ConditionalAccessGuestsOrExternalUsers{
		GuestOrExternalUserTypes: convert.GraphToFrameworkString(getStringPtr(guestsData, "guestOrExternalUserTypes")),
	}

	if externalTenantsData, ok := guestsData["externalTenants"].(map[string]interface{}); ok {
		externalTenants := &ConditionalAccessExternalTenants{
			MembershipKind: convert.GraphToFrameworkString(getStringPtr(externalTenantsData, "membershipKind")),
		}

		if members, ok := externalTenantsData["members"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(members))
			for _, item := range members {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			externalTenants.Members = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		} else {
			externalTenants.Members = types.SetNull(types.StringType)
		}

		guests.ExternalTenants = externalTenants
	}

	return guests
}

// mapPlatformsFromRemote maps platforms from the API response
func mapPlatformsFromRemote(ctx context.Context, platformsData map[string]interface{}) *ConditionalAccessPlatforms {
	platforms := &ConditionalAccessPlatforms{}

	if includePlatforms, ok := platformsData["includePlatforms"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includePlatforms))
		for _, item := range includePlatforms {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		platforms.IncludePlatforms = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		platforms.IncludePlatforms = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, []string{})
	}

	if excludePlatforms, ok := platformsData["excludePlatforms"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludePlatforms))
		for _, item := range excludePlatforms {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		platforms.ExcludePlatforms = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		platforms.ExcludePlatforms = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, []string{})
	}

	return platforms
}

// mapLocationsFromRemote maps locations from the API response
func mapLocationsFromRemote(ctx context.Context, locationsData map[string]interface{}) *ConditionalAccessLocations {
	locations := &ConditionalAccessLocations{}

	if includeLocations, ok := locationsData["includeLocations"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeLocations))
		for _, item := range includeLocations {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		locations.IncludeLocations = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		locations.IncludeLocations = types.SetNull(types.StringType)
	}

	if excludeLocations, ok := locationsData["excludeLocations"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeLocations))
		for _, item := range excludeLocations {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		locations.ExcludeLocations = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		locations.ExcludeLocations = types.SetNull(types.StringType)
	}

	return locations
}

// mapDevicesFromRemote maps devices from the API response
func mapDevicesFromRemote(ctx context.Context, devicesData map[string]interface{}) *ConditionalAccessDevices {
	devices := &ConditionalAccessDevices{}

	// These fields are configured as empty arrays in Terraform, so preserve empty sets
	if includeDevices, ok := devicesData["includeDevices"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeDevices))
		for _, item := range includeDevices {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.IncludeDevices = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		devices.IncludeDevices = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, []string{})
	}

	if excludeDevices, ok := devicesData["excludeDevices"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeDevices))
		for _, item := range excludeDevices {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.ExcludeDevices = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, stringSlice)
	} else {
		devices.ExcludeDevices = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, []string{})
	}

	if includeDeviceStates, ok := devicesData["includeDeviceStates"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeDeviceStates))
		for _, item := range includeDeviceStates {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.IncludeDeviceStates = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		devices.IncludeDeviceStates = types.SetNull(types.StringType)
	}

	if excludeDeviceStates, ok := devicesData["excludeDeviceStates"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeDeviceStates))
		for _, item := range excludeDeviceStates {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.ExcludeDeviceStates = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		devices.ExcludeDeviceStates = types.SetNull(types.StringType)
	}

	if deviceFilterData, ok := devicesData["deviceFilter"].(map[string]interface{}); ok {
		filter := &ConditionalAccessFilter{
			Mode: convert.GraphToFrameworkString(getStringPtr(deviceFilterData, "mode")),
			Rule: convert.GraphToFrameworkString(getStringPtr(deviceFilterData, "rule")),
		}
		devices.DeviceFilter = filter
	}

	return devices
}

// mapGrantControlsFromRemote maps grant controls from the API response
func mapGrantControlsFromRemote(ctx context.Context, grantControlsData map[string]interface{}) *ConditionalAccessGrantControls {
	grantControls := &ConditionalAccessGrantControls{
		Operator: convert.GraphToFrameworkString(getStringPtr(grantControlsData, "operator")),
	}

	if builtInControls, ok := grantControlsData["builtInControls"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(builtInControls))
		for _, item := range builtInControls {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		grantControls.BuiltInControls = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	} else {
		grantControls.BuiltInControls = types.SetNull(types.StringType)
	}

	// Only set these fields if they're present in the API response
	// This prevents null -> empty set inconsistencies
	if _, hasCustomAuth := grantControlsData["customAuthenticationFactors"]; hasCustomAuth {
		if customAuthFactors, ok := grantControlsData["customAuthenticationFactors"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(customAuthFactors))
			for _, item := range customAuthFactors {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			grantControls.CustomAuthenticationFactors = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		} else {
			grantControls.CustomAuthenticationFactors = types.SetNull(types.StringType)
		}
	} else {
		grantControls.CustomAuthenticationFactors = types.SetNull(types.StringType)
	}

	if _, hasTermsOfUse := grantControlsData["termsOfUse"]; hasTermsOfUse {
		if termsOfUse, ok := grantControlsData["termsOfUse"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(termsOfUse))
			for _, item := range termsOfUse {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			grantControls.TermsOfUse = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		} else {
			grantControls.TermsOfUse = types.SetNull(types.StringType)
		}
	} else {
		grantControls.TermsOfUse = types.SetNull(types.StringType)
	}

	if authStrengthData, ok := grantControlsData["authenticationStrength"].(map[string]interface{}); ok {
		authStrength := &ConditionalAccessAuthenticationStrength{
			ID:                    convert.GraphToFrameworkString(getStringPtr(authStrengthData, "id")),
			DisplayName:           convert.GraphToFrameworkString(getStringPtr(authStrengthData, "displayName")),
			Description:           convert.GraphToFrameworkString(getStringPtr(authStrengthData, "description")),
			PolicyType:            convert.GraphToFrameworkString(getStringPtr(authStrengthData, "policyType")),
			RequirementsSatisfied: convert.GraphToFrameworkString(getStringPtr(authStrengthData, "requirementsSatisfied")),
			CreatedDateTime:       convert.GraphToFrameworkString(getStringPtr(authStrengthData, "createdDateTime")),
			ModifiedDateTime:      convert.GraphToFrameworkString(getStringPtr(authStrengthData, "modifiedDateTime")),
		}

		if allowedCombinations, ok := authStrengthData["allowedCombinations"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(allowedCombinations))
			for _, item := range allowedCombinations {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			authStrength.AllowedCombinations = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		} else {
			authStrength.AllowedCombinations = types.SetNull(types.StringType)
		}

		grantControls.AuthenticationStrength = authStrength
	}

	return grantControls
}

// mapSessionControlsFromRemote maps session controls from the API response
func mapSessionControlsFromRemote(ctx context.Context, sessionControlsData map[string]interface{}) *ConditionalAccessSessionControls {
	sessionControls := &ConditionalAccessSessionControls{
		DisableResilienceDefaults: convert.GraphToFrameworkBool(getBoolPtr(sessionControlsData, "disableResilienceDefaults")),
	}

	if appEnforcedData, ok := sessionControlsData["applicationEnforcedRestrictions"].(map[string]interface{}); ok {
		sessionControls.ApplicationEnforcedRestrictions = &ConditionalAccessApplicationEnforcedRestrictions{
			IsEnabled: convert.GraphToFrameworkBool(getBoolPtr(appEnforcedData, "isEnabled")),
		}
	}

	if cloudAppSecData, ok := sessionControlsData["cloudAppSecurity"].(map[string]interface{}); ok {
		sessionControls.CloudAppSecurity = &ConditionalAccessCloudAppSecurity{
			IsEnabled:            convert.GraphToFrameworkBool(getBoolPtr(cloudAppSecData, "isEnabled")),
			CloudAppSecurityType: convert.GraphToFrameworkString(getStringPtr(cloudAppSecData, "cloudAppSecurityType")),
		}
	}

	if signInFreqData, ok := sessionControlsData["signInFrequency"].(map[string]interface{}); ok {
		signInFreq := &ConditionalAccessSignInFrequency{
			IsEnabled:          convert.GraphToFrameworkBool(getBoolPtr(signInFreqData, "isEnabled")),
			Type:               convert.GraphToFrameworkString(getStringPtr(signInFreqData, "type")),
			AuthenticationType: convert.GraphToFrameworkString(getStringPtr(signInFreqData, "authenticationType")),
			FrequencyInterval:  convert.GraphToFrameworkString(getStringPtr(signInFreqData, "frequencyInterval")),
		}

		if value, ok := signInFreqData["value"].(float64); ok {
			signInFreq.Value = types.Int64Value(int64(value))
		}

		sessionControls.SignInFrequency = signInFreq
	}

	if persistentBrowserData, ok := sessionControlsData["persistentBrowser"].(map[string]interface{}); ok {
		sessionControls.PersistentBrowser = &ConditionalAccessPersistentBrowser{
			IsEnabled: convert.GraphToFrameworkBool(getBoolPtr(persistentBrowserData, "isEnabled")),
			Mode:      convert.GraphToFrameworkString(getStringPtr(persistentBrowserData, "mode")),
		}
	}

	if caeData, ok := sessionControlsData["continuousAccessEvaluation"].(map[string]interface{}); ok {
		sessionControls.ContinuousAccessEvaluation = &ConditionalAccessContinuousAccessEvaluation{
			Mode: convert.GraphToFrameworkString(getStringPtr(caeData, "mode")),
		}
	}

	if secureSignInData, ok := sessionControlsData["secureSignInSession"].(map[string]interface{}); ok {
		sessionControls.SecureSignInSession = &ConditionalAccessSecureSignInSession{
			IsEnabled: convert.GraphToFrameworkBool(getBoolPtr(secureSignInData, "isEnabled")),
		}
	}

	return sessionControls
}

// Helper function to get bool pointer from map
func getBoolPtr(data map[string]interface{}, key string) *bool {
	if value, ok := data[key].(bool); ok {
		return &value
	}
	return nil
}
