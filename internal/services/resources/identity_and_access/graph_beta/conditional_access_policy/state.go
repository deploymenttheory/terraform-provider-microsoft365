package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteResourceStateToTerraform maps the remote conditional access policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource map[string]any) {
	tflog.Debug(ctx, "Starting MapRemoteResourceStateToTerraform", map[string]any{
		"remoteResource": remoteResource,
	})

	// Basic properties - using helper functions
	data.ID = convert.MapToFrameworkString(remoteResource, "id")
	data.DisplayName = convert.MapToFrameworkString(remoteResource, "displayName")
	data.State = convert.MapToFrameworkString(remoteResource, "state")
	data.CreatedDateTime = convert.MapToFrameworkString(remoteResource, "createdDateTime")
	data.ModifiedDateTime = convert.MapToFrameworkString(remoteResource, "modifiedDateTime")
	data.DeletedDateTime = convert.MapToFrameworkString(remoteResource, "deletedDateTime")
	data.TemplateId = convert.MapToFrameworkString(remoteResource, "templateId")
	data.PartialEnablementStrategy = convert.MapToFrameworkString(remoteResource, "partialEnablementStrategy")

	// Map conditions
	if conditionsRaw, ok := remoteResource["conditions"]; ok {
		tflog.Debug(ctx, "Mapping conditions", map[string]any{"conditions": conditionsRaw})
		data.Conditions = mapConditions(ctx, conditionsRaw)
	} else {
		tflog.Debug(ctx, "conditions not found")
		data.Conditions = nil
	}

	// Map grant controls (Required field - must always be present)
	if grantControlsRaw, ok := remoteResource["grantControls"]; ok && grantControlsRaw != nil {
		tflog.Debug(ctx, "Mapping grantControls", map[string]any{"grantControls": grantControlsRaw})
		data.GrantControls = mapGrantControls(ctx, grantControlsRaw)
	} else {
		tflog.Debug(ctx, "grantControls not found or null, creating empty grant controls for required field")
		// grant_controls is Required in schema, so we must return an empty object when API returns null
		data.GrantControls = &ConditionalAccessGrantControls{
			Operator:                    types.StringValue("OR"), // Default operator
			BuiltInControls:             types.SetValueMust(types.StringType, []attr.Value{}),
			CustomAuthenticationFactors: types.SetValueMust(types.StringType, []attr.Value{}),
			TermsOfUse:                  types.SetNull(types.StringType),
			AuthenticationStrength:      nil,
		}
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

	// Map sets using helper functions
	result.ClientAppTypes = mapStringSliceToSet(ctx, conditions["clientAppTypes"], "clientAppTypes")
	result.SignInRiskLevels = mapStringSliceToSet(ctx, conditions["signInRiskLevels"], "signInRiskLevels")
	result.UserRiskLevels = mapOptionalStringSliceToSet(ctx, conditions["userRiskLevels"], "userRiskLevels")
	result.ServicePrincipalRiskLevels = mapStringSliceToSet(ctx, conditions["servicePrincipalRiskLevels"], "servicePrincipalRiskLevels")

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

	// Map authenticationFlows
	if authenticationFlowsRaw, ok := conditions["authenticationFlows"]; ok {
		tflog.Debug(ctx, "Mapping authenticationFlows", map[string]any{"authenticationFlows": authenticationFlowsRaw})
		result.AuthenticationFlows = mapAuthenticationFlows(ctx, authenticationFlowsRaw)
	} else {
		tflog.Debug(ctx, "authenticationFlows not found")
		result.AuthenticationFlows = nil
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

	// Map sets
	result.IncludeApplications = mapStringSliceToSet(ctx, applications["includeApplications"], "includeApplications")
	result.ExcludeApplications = mapStringSliceToSet(ctx, applications["excludeApplications"], "excludeApplications")
	result.IncludeUserActions = mapStringSliceToSet(ctx, applications["includeUserActions"], "includeUserActions")
	result.IncludeAuthenticationContextClassReferences = mapAuthContextClassReferencesToSet(ctx, applications["includeAuthenticationContextClassReferences"], "includeAuthenticationContextClassReferences")

	// Map applicationFilter
	if applicationFilterRaw, ok := applications["applicationFilter"]; ok {
		tflog.Debug(ctx, "Mapping applicationFilter", map[string]any{"applicationFilter": applicationFilterRaw})
		result.ApplicationFilter = mapFilter(ctx, applicationFilterRaw)
	} else {
		tflog.Debug(ctx, "applicationFilter not found")
		result.ApplicationFilter = nil
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

	// Map sets
	result.IncludeUsers = mapStringSliceToSet(ctx, users["includeUsers"], "includeUsers")
	result.ExcludeUsers = mapStringSliceToSet(ctx, users["excludeUsers"], "excludeUsers")
	result.IncludeGroups = mapStringSliceToSet(ctx, users["includeGroups"], "includeGroups")
	result.ExcludeGroups = mapStringSliceToSet(ctx, users["excludeGroups"], "excludeGroups")
	result.IncludeRoles = mapStringSliceToSet(ctx, users["includeRoles"], "includeRoles")
	result.ExcludeRoles = mapStringSliceToSet(ctx, users["excludeRoles"], "excludeRoles")

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
	result.IncludeLocations = mapStringSliceToSet(ctx, locations["includeLocations"], "includeLocations")
	result.ExcludeLocations = mapStringSliceToSet(ctx, locations["excludeLocations"], "excludeLocations")
	return result
}

func mapPlatforms(ctx context.Context, platformsRaw any) *ConditionalAccessPlatforms {
	platforms, ok := platformsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "platforms is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessPlatforms{}
	result.IncludePlatforms = mapStringSliceToSet(ctx, platforms["includePlatforms"], "includePlatforms")
	result.ExcludePlatforms = mapStringSliceToSet(ctx, platforms["excludePlatforms"], "excludePlatforms")
	return result
}

func mapDevices(ctx context.Context, devicesRaw any) *ConditionalAccessDevices {
	devices, ok := devicesRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "devices is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessDevices{}
	result.IncludeDevices = mapOptionalStringSliceToSet(ctx, devices["includeDevices"], "includeDevices")
	result.ExcludeDevices = mapOptionalStringSliceToSet(ctx, devices["excludeDevices"], "excludeDevices")
	result.IncludeDeviceStates = mapOptionalStringSliceToSet(ctx, devices["includeDeviceStates"], "includeDeviceStates")
	result.ExcludeDeviceStates = mapOptionalStringSliceToSet(ctx, devices["excludeDeviceStates"], "excludeDeviceStates")

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
	result.IncludeServicePrincipals = mapStringSliceToSet(ctx, clientApplications["includeServicePrincipals"], "includeServicePrincipals")
	result.ExcludeServicePrincipals = mapStringSliceToSet(ctx, clientApplications["excludeServicePrincipals"], "excludeServicePrincipals")
	return result
}

func mapTimes(ctx context.Context, timesRaw any) *ConditionalAccessTimes {
	times, ok := timesRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "times is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessTimes{}
	result.IncludedRanges = mapStringSliceToSet(ctx, times["includedRanges"], "includedRanges")
	result.ExcludedRanges = mapStringSliceToSet(ctx, times["excludedRanges"], "excludedRanges")
	result.AllDay = convert.MapToFrameworkBool(times, "allDay")
	result.StartTime = convert.MapToFrameworkString(times, "startTime")
	result.EndTime = convert.MapToFrameworkString(times, "endTime")
	result.TimeZone = convert.MapToFrameworkString(times, "timeZone")
	return result
}

func mapDeviceStates(ctx context.Context, deviceStatesRaw any) *ConditionalAccessDeviceStates {
	deviceStates, ok := deviceStatesRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "deviceStates is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessDeviceStates{}
	result.IncludeStates = mapStringSliceToSet(ctx, deviceStates["includeStates"], "includeStates")
	result.ExcludeStates = mapStringSliceToSet(ctx, deviceStates["excludeStates"], "excludeStates")
	return result
}

func mapAuthenticationFlows(ctx context.Context, authenticationFlowsRaw any) *ConditionalAccessAuthenticationFlows {
	authenticationFlows, ok := authenticationFlowsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "authenticationFlows is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessAuthenticationFlows{}
	result.TransferMethods = convert.MapToFrameworkString(authenticationFlows, "transferMethods")
	return result
}

func mapGrantControls(ctx context.Context, grantControlsRaw any) *ConditionalAccessGrantControls {
	grantControls, ok := grantControlsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "grantControls is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessGrantControls{}
	result.Operator = convert.MapToFrameworkString(grantControls, "operator")
	result.BuiltInControls = mapStringSliceToSet(ctx, grantControls["builtInControls"], "builtInControls")
	result.CustomAuthenticationFactors = mapStringSliceToSet(ctx, grantControls["customAuthenticationFactors"], "customAuthenticationFactors")
	result.TermsOfUse = mapOptionalStringSliceToSet(ctx, grantControls["termsOfUse"], "termsOfUse")

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
	result.DisableResilienceDefaults = convert.MapToFrameworkBool(sessionControls, "disableResilienceDefaults")

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

	// Map globalSecureAccessFilteringProfile
	if globalSecureAccessFilteringProfileRaw, ok := sessionControls["globalSecureAccessFilteringProfile"]; ok {
		tflog.Debug(ctx, "Mapping globalSecureAccessFilteringProfile", map[string]any{"globalSecureAccessFilteringProfile": globalSecureAccessFilteringProfileRaw})
		result.GlobalSecureAccessFilteringProfile = mapGlobalSecureAccessFilteringProfile(ctx, globalSecureAccessFilteringProfileRaw)
	} else {
		tflog.Debug(ctx, "globalSecureAccessFilteringProfile not found")
		result.GlobalSecureAccessFilteringProfile = nil
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
	result.Mode = convert.MapToFrameworkString(filter, "mode")
	result.Rule = convert.MapToFrameworkString(filter, "rule")
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
	result.ID = convert.MapToFrameworkString(authenticationStrength, "id")
	result.DisplayName = convert.MapToFrameworkString(authenticationStrength, "displayName")
	result.Description = convert.MapToFrameworkString(authenticationStrength, "description")
	result.PolicyType = convert.MapToFrameworkString(authenticationStrength, "policyType")
	result.RequirementsSatisfied = convert.MapToFrameworkString(authenticationStrength, "requirementsSatisfied")
	result.AllowedCombinations = mapStringSliceToSet(ctx, authenticationStrength["allowedCombinations"], "allowedCombinations")
	result.CreatedDateTime = convert.MapToFrameworkString(authenticationStrength, "createdDateTime")
	result.ModifiedDateTime = convert.MapToFrameworkString(authenticationStrength, "modifiedDateTime")
	return result
}

func mapApplicationEnforcedRestrictions(ctx context.Context, applicationEnforcedRestrictionsRaw any) *ConditionalAccessApplicationEnforcedRestrictions {
	applicationEnforcedRestrictions, ok := applicationEnforcedRestrictionsRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "applicationEnforcedRestrictions is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessApplicationEnforcedRestrictions{}
	result.IsEnabled = convert.MapToFrameworkBool(applicationEnforcedRestrictions, "isEnabled")
	return result
}

func mapCloudAppSecurity(ctx context.Context, cloudAppSecurityRaw any) *ConditionalAccessCloudAppSecurity {
	cloudAppSecurity, ok := cloudAppSecurityRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "cloudAppSecurity is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessCloudAppSecurity{}
	result.IsEnabled = convert.MapToFrameworkBool(cloudAppSecurity, "isEnabled")
	result.CloudAppSecurityType = convert.MapToFrameworkString(cloudAppSecurity, "cloudAppSecurityType")
	return result
}

func mapSignInFrequency(ctx context.Context, signInFrequencyRaw any) *ConditionalAccessSignInFrequency {
	signInFrequency, ok := signInFrequencyRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "signInFrequency is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessSignInFrequency{}
	result.IsEnabled = convert.MapToFrameworkBool(signInFrequency, "isEnabled")
	result.Type = convert.MapToFrameworkString(signInFrequency, "type")
	result.Value = convert.MapToFrameworkInt64(signInFrequency, "value")
	result.AuthenticationType = convert.MapToFrameworkString(signInFrequency, "authenticationType")
	result.FrequencyInterval = convert.MapToFrameworkString(signInFrequency, "frequencyInterval")
	return result
}

func mapPersistentBrowser(ctx context.Context, persistentBrowserRaw any) *ConditionalAccessPersistentBrowser {
	persistentBrowser, ok := persistentBrowserRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "persistentBrowser is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessPersistentBrowser{}
	result.IsEnabled = convert.MapToFrameworkBool(persistentBrowser, "isEnabled")
	result.Mode = convert.MapToFrameworkString(persistentBrowser, "mode")
	return result
}

func mapContinuousAccessEvaluation(ctx context.Context, continuousAccessEvaluationRaw any) *ConditionalAccessContinuousAccessEvaluation {
	continuousAccessEvaluation, ok := continuousAccessEvaluationRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "continuousAccessEvaluation is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessContinuousAccessEvaluation{}
	result.Mode = convert.MapToFrameworkString(continuousAccessEvaluation, "mode")
	return result
}

func mapSecureSignInSession(ctx context.Context, secureSignInSessionRaw any) *ConditionalAccessSecureSignInSession {
	secureSignInSession, ok := secureSignInSessionRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "secureSignInSession is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessSecureSignInSession{}
	result.IsEnabled = convert.MapToFrameworkBool(secureSignInSession, "isEnabled")
	return result
}

func mapGlobalSecureAccessFilteringProfile(ctx context.Context, globalSecureAccessFilteringProfileRaw any) *ConditionalAccessGlobalSecureAccessFilteringProfile {
	globalSecureAccessFilteringProfile, ok := globalSecureAccessFilteringProfileRaw.(map[string]any)
	if !ok {
		tflog.Debug(ctx, "globalSecureAccessFilteringProfile is not a map[string]any")
		return nil
	}

	result := &ConditionalAccessGlobalSecureAccessFilteringProfile{}
	result.IsEnabled = convert.MapToFrameworkBool(globalSecureAccessFilteringProfile, "isEnabled")
	result.ProfileId = convert.MapToFrameworkString(globalSecureAccessFilteringProfile, "profileId")
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

// Helper function to map OPTIONAL string slices to Terraform sets
// For optional fields: null or empty arrays from API should remain null in Terraform state
func mapOptionalStringSliceToSet(ctx context.Context, raw any, fieldName string) types.Set {
	tflog.Debug(ctx, fmt.Sprintf("Processing optional field %s: %v (type: %T)", fieldName, raw, raw))

	// For optional fields, keep as null if empty
	if slice, ok := raw.([]any); ok && len(slice) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s is empty array, keeping as null for optional field", fieldName))
		return types.SetNull(types.StringType)
	} else if raw == nil {
		tflog.Debug(ctx, fmt.Sprintf("%s is null, keeping as null for optional field", fieldName))
		return types.SetNull(types.StringType)
	}

	// If it has values, process normally
	return mapStringSliceToSet(ctx, raw, fieldName)
}
