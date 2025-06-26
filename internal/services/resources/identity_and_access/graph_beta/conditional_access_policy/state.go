package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteResourceStateToTerraform maps the remote conditional access policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource map[string]interface{}) {
	// Basic properties
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	if displayName, ok := remoteResource["displayName"].(string); ok {
		data.DisplayName = types.StringValue(displayName)
	}

	if state, ok := remoteResource["state"].(string); ok {
		data.State = types.StringValue(state)
	}

	if createdDateTime, ok := remoteResource["createdDateTime"].(string); ok {
		data.CreatedDateTime = types.StringValue(createdDateTime)
	}

	if modifiedDateTime, ok := remoteResource["modifiedDateTime"].(string); ok {
		data.ModifiedDateTime = types.StringValue(modifiedDateTime)
	}

	if deletedDateTime, ok := remoteResource["deletedDateTime"].(string); ok {
		data.DeletedDateTime = types.StringValue(deletedDateTime)
	}

	if templateId, ok := remoteResource["templateId"].(string); ok {
		data.TemplateId = types.StringValue(templateId)
	}

	// Map conditions
	if conditionsData, ok := remoteResource["conditions"].(map[string]interface{}); ok {
		data.Conditions = mapConditionsFromRemote(ctx, conditionsData)
	}

	// Map grant controls
	if grantControlsData, ok := remoteResource["grantControls"].(map[string]interface{}); ok {
		data.GrantControls = mapGrantControlsFromRemote(ctx, grantControlsData)
	}

	// Map session controls
	if sessionControlsData, ok := remoteResource["sessionControls"].(map[string]interface{}); ok {
		data.SessionControls = mapSessionControlsFromRemote(ctx, sessionControlsData)
	}
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
		conditions.ClientAppTypes = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	// Applications
	if applicationsData, ok := conditionsData["applications"].(map[string]interface{}); ok {
		conditions.Applications = mapApplicationsFromRemote(ctx, applicationsData)
	}

	// Users
	if usersData, ok := conditionsData["users"].(map[string]interface{}); ok {
		conditions.Users = mapUsersFromRemote(ctx, usersData)
	}

	// Platforms
	if platformsData, ok := conditionsData["platforms"].(map[string]interface{}); ok {
		conditions.Platforms = mapPlatformsFromRemote(ctx, platformsData)
	}

	// Locations
	if locationsData, ok := conditionsData["locations"].(map[string]interface{}); ok {
		conditions.Locations = mapLocationsFromRemote(ctx, locationsData)
	}

	// Devices
	if devicesData, ok := conditionsData["devices"].(map[string]interface{}); ok {
		conditions.Devices = mapDevicesFromRemote(ctx, devicesData)
	}

	// Sign-in risk levels
	if signInRiskLevels, ok := conditionsData["signInRiskLevels"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(signInRiskLevels))
		for _, item := range signInRiskLevels {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.SignInRiskLevels = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	// User risk levels
	if userRiskLevels, ok := conditionsData["userRiskLevels"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(userRiskLevels))
		for _, item := range userRiskLevels {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.UserRiskLevels = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	// Service principal risk levels
	if servicePrincipalRiskLevels, ok := conditionsData["servicePrincipalRiskLevels"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(servicePrincipalRiskLevels))
		for _, item := range servicePrincipalRiskLevels {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		conditions.ServicePrincipalRiskLevels = convert.GraphToFrameworkStringSet(ctx, stringSlice)
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
	}

	if excludeApplications, ok := applicationsData["excludeApplications"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeApplications))
		for _, item := range excludeApplications {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		applications.ExcludeApplications = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if includeUserActions, ok := applicationsData["includeUserActions"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeUserActions))
		for _, item := range includeUserActions {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		applications.IncludeUserActions = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if includeAuthContextClassRefs, ok := applicationsData["includeAuthenticationContextClassReferences"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeAuthContextClassRefs))
		for _, item := range includeAuthContextClassRefs {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		applications.IncludeAuthenticationContextClassReferences = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if applicationFilterData, ok := applicationsData["applicationFilter"].(map[string]interface{}); ok {
		filter := &ConditionalAccessFilter{}

		if mode, ok := applicationFilterData["mode"].(string); ok {
			filter.Mode = types.StringValue(mode)
		}

		if rule, ok := applicationFilterData["rule"].(string); ok {
			filter.Rule = types.StringValue(rule)
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
	}

	if excludeUsers, ok := usersData["excludeUsers"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeUsers))
		for _, item := range excludeUsers {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.ExcludeUsers = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if includeGroups, ok := usersData["includeGroups"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeGroups))
		for _, item := range includeGroups {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.IncludeGroups = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if excludeGroups, ok := usersData["excludeGroups"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeGroups))
		for _, item := range excludeGroups {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.ExcludeGroups = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if includeRoles, ok := usersData["includeRoles"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeRoles))
		for _, item := range includeRoles {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.IncludeRoles = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if excludeRoles, ok := usersData["excludeRoles"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeRoles))
		for _, item := range excludeRoles {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		users.ExcludeRoles = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	return users
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
		platforms.IncludePlatforms = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if excludePlatforms, ok := platformsData["excludePlatforms"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludePlatforms))
		for _, item := range excludePlatforms {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		platforms.ExcludePlatforms = convert.GraphToFrameworkStringSet(ctx, stringSlice)
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
	}

	if excludeLocations, ok := locationsData["excludeLocations"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeLocations))
		for _, item := range excludeLocations {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		locations.ExcludeLocations = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	return locations
}

// mapDevicesFromRemote maps devices from the API response
func mapDevicesFromRemote(ctx context.Context, devicesData map[string]interface{}) *ConditionalAccessDevices {
	devices := &ConditionalAccessDevices{}

	if includeDevices, ok := devicesData["includeDevices"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeDevices))
		for _, item := range includeDevices {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.IncludeDevices = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if excludeDevices, ok := devicesData["excludeDevices"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeDevices))
		for _, item := range excludeDevices {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.ExcludeDevices = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if includeDeviceStates, ok := devicesData["includeDeviceStates"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(includeDeviceStates))
		for _, item := range includeDeviceStates {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.IncludeDeviceStates = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if excludeDeviceStates, ok := devicesData["excludeDeviceStates"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(excludeDeviceStates))
		for _, item := range excludeDeviceStates {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		devices.ExcludeDeviceStates = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if deviceFilterData, ok := devicesData["deviceFilter"].(map[string]interface{}); ok {
		filter := &ConditionalAccessFilter{}

		if mode, ok := deviceFilterData["mode"].(string); ok {
			filter.Mode = types.StringValue(mode)
		}

		if rule, ok := deviceFilterData["rule"].(string); ok {
			filter.Rule = types.StringValue(rule)
		}

		devices.DeviceFilter = filter
	}

	return devices
}

// mapGrantControlsFromRemote maps grant controls from the API response
func mapGrantControlsFromRemote(ctx context.Context, grantControlsData map[string]interface{}) *ConditionalAccessGrantControls {
	grantControls := &ConditionalAccessGrantControls{}

	if operator, ok := grantControlsData["operator"].(string); ok {
		grantControls.Operator = types.StringValue(operator)
	}

	if builtInControls, ok := grantControlsData["builtInControls"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(builtInControls))
		for _, item := range builtInControls {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		grantControls.BuiltInControls = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if customAuthFactors, ok := grantControlsData["customAuthenticationFactors"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(customAuthFactors))
		for _, item := range customAuthFactors {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		grantControls.CustomAuthenticationFactors = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if termsOfUse, ok := grantControlsData["termsOfUse"].([]interface{}); ok {
		stringSlice := make([]string, 0, len(termsOfUse))
		for _, item := range termsOfUse {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			}
		}
		grantControls.TermsOfUse = convert.GraphToFrameworkStringSet(ctx, stringSlice)
	}

	if authStrengthData, ok := grantControlsData["authenticationStrength"].(map[string]interface{}); ok {
		authStrength := &ConditionalAccessAuthenticationStrength{}

		if id, ok := authStrengthData["id"].(string); ok {
			authStrength.ID = types.StringValue(id)
		}

		if displayName, ok := authStrengthData["displayName"].(string); ok {
			authStrength.DisplayName = types.StringValue(displayName)
		}

		if description, ok := authStrengthData["description"].(string); ok {
			authStrength.Description = types.StringValue(description)
		}

		if policyType, ok := authStrengthData["policyType"].(string); ok {
			authStrength.PolicyType = types.StringValue(policyType)
		}

		if requirementsSatisfied, ok := authStrengthData["requirementsSatisfied"].(string); ok {
			authStrength.RequirementsSatisfied = types.StringValue(requirementsSatisfied)
		}

		if allowedCombinations, ok := authStrengthData["allowedCombinations"].([]interface{}); ok {
			stringSlice := make([]string, 0, len(allowedCombinations))
			for _, item := range allowedCombinations {
				if str, ok := item.(string); ok {
					stringSlice = append(stringSlice, str)
				}
			}
			authStrength.AllowedCombinations = convert.GraphToFrameworkStringSet(ctx, stringSlice)
		}

		if createdDateTime, ok := authStrengthData["createdDateTime"].(string); ok {
			authStrength.CreatedDateTime = types.StringValue(createdDateTime)
		}

		if modifiedDateTime, ok := authStrengthData["modifiedDateTime"].(string); ok {
			authStrength.ModifiedDateTime = types.StringValue(modifiedDateTime)
		}

		grantControls.AuthenticationStrength = authStrength
	}

	return grantControls
}

// mapSessionControlsFromRemote maps session controls from the API response
func mapSessionControlsFromRemote(ctx context.Context, sessionControlsData map[string]interface{}) *ConditionalAccessSessionControls {
	sessionControls := &ConditionalAccessSessionControls{}

	if appEnforcedRestrictionsData, ok := sessionControlsData["applicationEnforcedRestrictions"].(map[string]interface{}); ok {
		appEnforcedRestrictions := &ConditionalAccessApplicationEnforcedRestrictions{}

		if isEnabled, ok := appEnforcedRestrictionsData["isEnabled"].(bool); ok {
			appEnforcedRestrictions.IsEnabled = types.BoolValue(isEnabled)
		}

		sessionControls.ApplicationEnforcedRestrictions = appEnforcedRestrictions
	}

	if cloudAppSecurityData, ok := sessionControlsData["cloudAppSecurity"].(map[string]interface{}); ok {
		cloudAppSecurity := &ConditionalAccessCloudAppSecurity{}

		if isEnabled, ok := cloudAppSecurityData["isEnabled"].(bool); ok {
			cloudAppSecurity.IsEnabled = types.BoolValue(isEnabled)
		}

		if cloudAppSecurityType, ok := cloudAppSecurityData["cloudAppSecurityType"].(string); ok {
			cloudAppSecurity.CloudAppSecurityType = types.StringValue(cloudAppSecurityType)
		}

		sessionControls.CloudAppSecurity = cloudAppSecurity
	}

	if signInFrequencyData, ok := sessionControlsData["signInFrequency"].(map[string]interface{}); ok {
		signInFrequency := &ConditionalAccessSignInFrequency{}

		if isEnabled, ok := signInFrequencyData["isEnabled"].(bool); ok {
			signInFrequency.IsEnabled = types.BoolValue(isEnabled)
		}

		if frequencyType, ok := signInFrequencyData["type"].(string); ok {
			signInFrequency.Type = types.StringValue(frequencyType)
		}

		if value, ok := signInFrequencyData["value"].(float64); ok {
			signInFrequency.Value = types.Int64Value(int64(value))
		}

		sessionControls.SignInFrequency = signInFrequency
	}

	if persistentBrowserData, ok := sessionControlsData["persistentBrowser"].(map[string]interface{}); ok {
		persistentBrowser := &ConditionalAccessPersistentBrowser{}

		if isEnabled, ok := persistentBrowserData["isEnabled"].(bool); ok {
			persistentBrowser.IsEnabled = types.BoolValue(isEnabled)
		}

		if mode, ok := persistentBrowserData["mode"].(string); ok {
			persistentBrowser.Mode = types.StringValue(mode)
		}

		sessionControls.PersistentBrowser = persistentBrowser
	}

	if disableResilienceDefaults, ok := sessionControlsData["disableResilienceDefaults"].(bool); ok {
		sessionControls.DisableResilienceDefaults = types.BoolValue(disableResilienceDefaults)
	}

	if continuousAccessEvalData, ok := sessionControlsData["continuousAccessEvaluation"].(map[string]interface{}); ok {
		continuousAccessEval := &ConditionalAccessContinuousAccessEvaluation{}

		if mode, ok := continuousAccessEvalData["mode"].(string); ok {
			continuousAccessEval.Mode = types.StringValue(mode)
		}

		sessionControls.ContinuousAccessEvaluation = continuousAccessEval
	}

	if secureSignInSessionData, ok := sessionControlsData["secureSignInSession"].(map[string]interface{}); ok {
		secureSignInSession := &ConditionalAccessSecureSignInSession{}

		if isEnabled, ok := secureSignInSessionData["isEnabled"].(bool); ok {
			secureSignInSession.IsEnabled = types.BoolValue(isEnabled)
		}

		sessionControls.SecureSignInSession = secureSignInSession
	}

	return sessionControls
}
