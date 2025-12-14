package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote conditional access policy from Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data ConditionalAccessPolicyResourceModel, remoteResource graphmodels.ConditionalAccessPolicyable) ConditionalAccessPolicyResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.State = convert.GraphToFrameworkEnum(remoteResource.GetState())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.ModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetModifiedDateTime())
	data.DeletedDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeletedDateTime())

	// templateId is not directly available in the SDK, check AdditionalData
	if additionalData := remoteResource.GetAdditionalData(); additionalData != nil {
		if templateId, ok := additionalData["templateId"].(string); ok {
			data.TemplateId = types.StringValue(templateId)
		} else {
			data.TemplateId = types.StringNull()
		}
	} else {
		data.TemplateId = types.StringNull()
	}

	// partialEnablementStrategy is not in the SDK, check AdditionalData
	if additionalData := remoteResource.GetAdditionalData(); additionalData != nil {
		if partialEnablementStrategy, ok := additionalData["partialEnablementStrategy"].(string); ok {
			data.PartialEnablementStrategy = types.StringValue(partialEnablementStrategy)
		} else {
			data.PartialEnablementStrategy = types.StringNull()
		}
	} else {
		data.PartialEnablementStrategy = types.StringNull()
	}

	// Map conditions
	if conditions := remoteResource.GetConditions(); conditions != nil {
		tflog.Debug(ctx, "Mapping conditions")
		data.Conditions = stateConditions(ctx, conditions)
	} else {
		tflog.Debug(ctx, "conditions not found")
		data.Conditions = nil
	}

	// Map grant controls (Required field - must always be present)
	if grantControls := remoteResource.GetGrantControls(); grantControls != nil {
		tflog.Debug(ctx, "Mapping grantControls")
		data.GrantControls = mapGrantControls(ctx, grantControls)
	} else {
		tflog.Debug(ctx, "grantControls not found or null, creating empty grant controls for required field")
		// grant_controls is Required in schema, so we must return an empty object when API returns null
		data.GrantControls = &ConditionalAccessGrantControls{
			Operator:                    types.StringValue("OR"), // Default operator
			BuiltInControls:             types.SetValueMust(types.StringType, []attr.Value{}),
			CustomAuthenticationFactors: types.SetValueMust(types.StringType, []attr.Value{}),
			TermsOfUse:                  types.SetValueMust(types.StringType, []attr.Value{}),
			AuthenticationStrength:      nil,
		}
	}

	// Map session controls
	if sessionControls := remoteResource.GetSessionControls(); sessionControls != nil {
		tflog.Debug(ctx, "Mapping sessionControls")
		data.SessionControls = mapSessionControls(ctx, sessionControls)
	} else {
		tflog.Debug(ctx, "sessionControls not found")
		data.SessionControls = nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}

func stateConditions(ctx context.Context, conditions graphmodels.ConditionalAccessConditionSetable) *ConditionalAccessConditions {
	result := &ConditionalAccessConditions{}

	result.ClientAppTypes = mapEnumCollectionToSet(ctx, conditions.GetClientAppTypes(), "clientAppTypes")
	result.SignInRiskLevels = mapEnumCollectionToSet(ctx, conditions.GetSignInRiskLevels(), "signInRiskLevels")
	// UserRiskLevels: Preserve empty arrays to match Terraform config
	result.UserRiskLevels = mapEnumCollectionToSet(ctx, conditions.GetUserRiskLevels(), "userRiskLevels")
	// ServicePrincipalRiskLevels: Preserve empty arrays to match Terraform config
	result.ServicePrincipalRiskLevels = mapEnumCollectionToSet(ctx, conditions.GetServicePrincipalRiskLevels(), "servicePrincipalRiskLevels")
	// AgentIdRiskLevels: Bitmask enum, will be empty set if not present
	if agentRiskEnum := conditions.GetAgentIdRiskLevels(); agentRiskEnum != nil {
		result.AgentIdRiskLevels = convert.GraphToFrameworkBitmaskEnumAsSet(ctx, agentRiskEnum)
	} else {
		result.AgentIdRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}
	// InsiderRiskLevels: Bitmask enum, will be empty set if not present
	if insiderRiskEnum := conditions.GetInsiderRiskLevels(); insiderRiskEnum != nil {
		result.InsiderRiskLevels = convert.GraphToFrameworkBitmaskEnumAsSet(ctx, insiderRiskEnum)
	} else {
		result.InsiderRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}

	if applications := conditions.GetApplications(); applications != nil {
		result.Applications = stateApplications(ctx, applications)
	}

	if users := conditions.GetUsers(); users != nil {
		result.Users = stateUsers(ctx, users)
	}

	if locations := conditions.GetLocations(); locations != nil {
		result.Locations = stateLocations(ctx, locations)
	}

	if platforms := conditions.GetPlatforms(); platforms != nil {
		result.Platforms = statePlatforms(ctx, platforms)
	}

	if devices := conditions.GetDevices(); devices != nil {
		result.Devices = stateDevices(ctx, devices)
	}

	if clientApplications := conditions.GetClientApplications(); clientApplications != nil {
		result.ClientApplications = stateClientApplications(ctx, clientApplications)
	}

	if authenticationFlows := conditions.GetAuthenticationFlows(); authenticationFlows != nil {
		result.AuthenticationFlows = stateAuthenticationFlows(ctx, authenticationFlows)
	}

	// Times and DeviceStates may be in AdditionalData if not directly available
	// For now, leave as nil if not in SDK
	result.Times = nil
	result.DeviceStates = nil

	return result
}

func stateApplications(ctx context.Context, applications graphmodels.ConditionalAccessApplicationsable) *ConditionalAccessApplications {
	result := &ConditionalAccessApplications{}

	result.IncludeApplications = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, applications.GetIncludeApplications())
	result.ExcludeApplications = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, applications.GetExcludeApplications())
	result.IncludeUserActions = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, applications.GetIncludeUserActions())
	result.IncludeAuthenticationContextClassReferences = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, applications.GetIncludeAuthenticationContextClassReferences())

	if applicationFilter := applications.GetApplicationFilter(); applicationFilter != nil {
		result.ApplicationFilter = mapFilter(ctx, applicationFilter)
	}

	// TODO: GlobalSecureAccess field is not in the current SDK model
	// I think this is eitehr a future SDK addition or i don't have it enabled in my tenant
	// for testing yet some the api fields are not showing?

	return result
}

func stateUsers(ctx context.Context, users graphmodels.ConditionalAccessUsersable) *ConditionalAccessUsers {
	result := &ConditionalAccessUsers{}

	result.IncludeUsers = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, users.GetIncludeUsers())
	result.ExcludeUsers = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, users.GetExcludeUsers())
	result.IncludeGroups = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, users.GetIncludeGroups())
	result.ExcludeGroups = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, users.GetExcludeGroups())
	result.IncludeRoles = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, users.GetIncludeRoles())
	result.ExcludeRoles = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, users.GetExcludeRoles())

	if includeGuestsOrExternalUsers := users.GetIncludeGuestsOrExternalUsers(); includeGuestsOrExternalUsers != nil {
		result.IncludeGuestsOrExternalUsers = mapGuestsOrExternalUsersToObject(ctx, includeGuestsOrExternalUsers)
	} else {
		result.IncludeGuestsOrExternalUsers = types.ObjectNull(map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{AttrTypes: map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			}},
		})
	}

	if excludeGuestsOrExternalUsers := users.GetExcludeGuestsOrExternalUsers(); excludeGuestsOrExternalUsers != nil {
		result.ExcludeGuestsOrExternalUsers = mapGuestsOrExternalUsersToObject(ctx, excludeGuestsOrExternalUsers)
	} else {
		result.ExcludeGuestsOrExternalUsers = types.ObjectNull(map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{AttrTypes: map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			}},
		})
	}

	return result
}

func stateLocations(ctx context.Context, locations graphmodels.ConditionalAccessLocationsable) *ConditionalAccessLocations {
	result := &ConditionalAccessLocations{}

	result.IncludeLocations = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, locations.GetIncludeLocations())
	result.ExcludeLocations = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, locations.GetExcludeLocations())

	return result
}

func statePlatforms(ctx context.Context, platforms graphmodels.ConditionalAccessPlatformsable) *ConditionalAccessPlatforms {
	result := &ConditionalAccessPlatforms{}

	result.IncludePlatforms = mapEnumCollectionToSet(ctx, platforms.GetIncludePlatforms(), "includePlatforms")
	result.ExcludePlatforms = mapEnumCollectionToSet(ctx, platforms.GetExcludePlatforms(), "excludePlatforms")

	return result
}

func stateDevices(ctx context.Context, devices graphmodels.ConditionalAccessDevicesable) *ConditionalAccessDevices {
	result := &ConditionalAccessDevices{}

	if deviceFilter := devices.GetDeviceFilter(); deviceFilter != nil {
		result.DeviceFilter = mapFilter(ctx, deviceFilter)
	}

	// IncludeDevices/ExcludeDevices might be in AdditionalData
	if additionalData := devices.GetAdditionalData(); additionalData != nil {
		if includeDevices, ok := additionalData["includeDevices"].([]any); ok {
			strings := make([]string, len(includeDevices))
			for i, v := range includeDevices {
				if str, ok := v.(string); ok {
					strings[i] = str
				}
			}
			result.IncludeDevices = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, strings)
		} else {
			result.IncludeDevices = types.SetNull(types.StringType)
		}

		if excludeDevices, ok := additionalData["excludeDevices"].([]any); ok {
			strings := make([]string, len(excludeDevices))
			for i, v := range excludeDevices {
				if str, ok := v.(string); ok {
					strings[i] = str
				}
			}
			result.ExcludeDevices = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, strings)
		} else {
			result.ExcludeDevices = types.SetNull(types.StringType)
		}

		if includeDeviceStates, ok := additionalData["includeDeviceStates"].([]any); ok {
			strings := make([]string, len(includeDeviceStates))
			for i, v := range includeDeviceStates {
				if str, ok := v.(string); ok {
					strings[i] = str
				}
			}
			result.IncludeDeviceStates = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, strings)
		} else {
			result.IncludeDeviceStates = types.SetNull(types.StringType)
		}

		if excludeDeviceStates, ok := additionalData["excludeDeviceStates"].([]any); ok {
			strings := make([]string, len(excludeDeviceStates))
			for i, v := range excludeDeviceStates {
				if str, ok := v.(string); ok {
					strings[i] = str
				}
			}
			result.ExcludeDeviceStates = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, strings)
		} else {
			result.ExcludeDeviceStates = types.SetNull(types.StringType)
		}
	} else {
		result.IncludeDevices = types.SetNull(types.StringType)
		result.ExcludeDevices = types.SetNull(types.StringType)
		result.IncludeDeviceStates = types.SetNull(types.StringType)
		result.ExcludeDeviceStates = types.SetNull(types.StringType)
	}

	return result
}

func stateClientApplications(ctx context.Context, clientApplications graphmodels.ConditionalAccessClientApplicationsable) *ConditionalAccessClientApplications {
	result := &ConditionalAccessClientApplications{}

	result.IncludeServicePrincipals = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, clientApplications.GetIncludeServicePrincipals())
	result.ExcludeServicePrincipals = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, clientApplications.GetExcludeServicePrincipals())

	// Try to get from SDK methods first
	includeAgentId := clientApplications.GetIncludeAgentIdServicePrincipals()
	excludeAgentId := clientApplications.GetExcludeAgentIdServicePrincipals()

	// If SDK methods return data, use them
	if includeAgentId != nil {
		result.IncludeAgentIdServicePrincipals = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, includeAgentId)
	} else {
		// Otherwise check AdditionalData
		if additionalData := clientApplications.GetAdditionalData(); additionalData != nil {
			if includeAgentIdServicePrincipals, ok := additionalData["includeAgentIdServicePrincipals"].([]any); ok {
				strings := make([]string, len(includeAgentIdServicePrincipals))
				for i, v := range includeAgentIdServicePrincipals {
					if str, ok := v.(string); ok {
						strings[i] = str
					}
				}
				result.IncludeAgentIdServicePrincipals = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, strings)
			} else {
				result.IncludeAgentIdServicePrincipals = types.SetNull(types.StringType)
			}
		} else {
			result.IncludeAgentIdServicePrincipals = types.SetNull(types.StringType)
		}
	}

	if excludeAgentId != nil {
		result.ExcludeAgentIdServicePrincipals = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, excludeAgentId)
	} else {
		// Otherwise check AdditionalData
		if additionalData := clientApplications.GetAdditionalData(); additionalData != nil {
			if excludeAgentIdServicePrincipals, ok := additionalData["excludeAgentIdServicePrincipals"].([]any); ok {
				strings := make([]string, len(excludeAgentIdServicePrincipals))
				for i, v := range excludeAgentIdServicePrincipals {
					if str, ok := v.(string); ok {
						strings[i] = str
					}
				}
				result.ExcludeAgentIdServicePrincipals = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, strings)
			} else {
				// Preserve empty set to match Terraform config
				result.ExcludeAgentIdServicePrincipals = types.SetValueMust(types.StringType, []attr.Value{})
			}
		} else {
			// Preserve empty set to match Terraform config
			result.ExcludeAgentIdServicePrincipals = types.SetValueMust(types.StringType, []attr.Value{})
		}
	}

	if servicePrincipalFilter := clientApplications.GetServicePrincipalFilter(); servicePrincipalFilter != nil {
		result.ServicePrincipalFilter = mapFilter(ctx, servicePrincipalFilter)
	}

	// Try SDK method first for agent ID filter
	agentIdFilter := clientApplications.GetAgentIdServicePrincipalFilter()
	if agentIdFilter != nil {
		result.AgentIdServicePrincipalFilter = mapFilter(ctx, agentIdFilter)
	} else {
		// Otherwise check AdditionalData
		result.AgentIdServicePrincipalFilter = mapAgentIdServicePrincipalFilter(ctx, clientApplications)
	}

	return result
}

// mapAgentIdServicePrincipalFilter extracts agent ID service principal filter from AdditionalData
func mapAgentIdServicePrincipalFilter(ctx context.Context, clientApplications graphmodels.ConditionalAccessClientApplicationsable) *ConditionalAccessFilter {
	additionalData := clientApplications.GetAdditionalData()
	if additionalData == nil {
		return nil
	}

	agentIdFilterData, ok := additionalData["agentIdServicePrincipalFilter"].(map[string]any)
	if !ok {
		return nil
	}

	filter := &ConditionalAccessFilter{
		Mode: types.StringNull(),
		Rule: types.StringNull(),
	}

	if mode, ok := agentIdFilterData["mode"].(string); ok {
		filter.Mode = types.StringValue(mode)
	}

	if rule, ok := agentIdFilterData["rule"].(string); ok {
		filter.Rule = types.StringValue(rule)
	}

	return filter
}

func stateAuthenticationFlows(ctx context.Context, authenticationFlows graphmodels.ConditionalAccessAuthenticationFlowsable) *ConditionalAccessAuthenticationFlows {
	result := &ConditionalAccessAuthenticationFlows{}

	result.TransferMethods = convert.GraphToFrameworkEnum(authenticationFlows.GetTransferMethods())

	return result
}

func mapGrantControls(ctx context.Context, grantControls graphmodels.ConditionalAccessGrantControlsable) *ConditionalAccessGrantControls {
	result := &ConditionalAccessGrantControls{}

	result.Operator = convert.GraphToFrameworkString(grantControls.GetOperator())
	result.BuiltInControls = mapEnumCollectionToSet(ctx, grantControls.GetBuiltInControls(), "builtInControls")
	result.CustomAuthenticationFactors = mapStringSliceToSetPreserveEmpty(ctx, grantControls.GetCustomAuthenticationFactors())
	result.TermsOfUse = mapStringSliceToSetPreserveEmpty(ctx, grantControls.GetTermsOfUse())

	if authenticationStrength := grantControls.GetAuthenticationStrength(); authenticationStrength != nil {
		result.AuthenticationStrength = mapAuthenticationStrength(ctx, authenticationStrength)
	}

	return result
}

func mapSessionControls(ctx context.Context, sessionControls graphmodels.ConditionalAccessSessionControlsable) *ConditionalAccessSessionControls {
	result := &ConditionalAccessSessionControls{}

	if applicationEnforcedRestrictions := sessionControls.GetApplicationEnforcedRestrictions(); applicationEnforcedRestrictions != nil {
		result.ApplicationEnforcedRestrictions = mapApplicationEnforcedRestrictions(ctx, applicationEnforcedRestrictions)
	}

	if cloudAppSecurity := sessionControls.GetCloudAppSecurity(); cloudAppSecurity != nil {
		result.CloudAppSecurity = mapCloudAppSecurity(ctx, cloudAppSecurity)
	}

	if signInFrequency := sessionControls.GetSignInFrequency(); signInFrequency != nil {
		result.SignInFrequency = mapSignInFrequency(ctx, signInFrequency)
	}

	if persistentBrowser := sessionControls.GetPersistentBrowser(); persistentBrowser != nil {
		result.PersistentBrowser = mapPersistentBrowser(ctx, persistentBrowser)
	}

	if continuousAccessEvaluation := sessionControls.GetContinuousAccessEvaluation(); continuousAccessEvaluation != nil {
		result.ContinuousAccessEvaluation = mapContinuousAccessEvaluation(ctx, continuousAccessEvaluation)
	}

	if secureSignInSession := sessionControls.GetSecureSignInSession(); secureSignInSession != nil {
		result.SecureSignInSession = mapSecureSignInSession(ctx, secureSignInSession)
	}

	// GlobalSecureAccessFilteringProfile might be in AdditionalData or future SDK
	result.GlobalSecureAccessFilteringProfile = nil

	if disableResilienceDefaults := sessionControls.GetDisableResilienceDefaults(); disableResilienceDefaults != nil {
		result.DisableResilienceDefaults = convert.GraphToFrameworkBool(disableResilienceDefaults)
	} else {
		result.DisableResilienceDefaults = types.BoolNull()
	}

	return result
}

func mapFilter(ctx context.Context, filter graphmodels.ConditionalAccessFilterable) *ConditionalAccessFilter {
	result := &ConditionalAccessFilter{}

	result.Mode = convert.GraphToFrameworkEnum(filter.GetMode())
	result.Rule = convert.GraphToFrameworkString(filter.GetRule())

	return result
}

func mapGuestsOrExternalUsersToObject(ctx context.Context, guests graphmodels.ConditionalAccessGuestsOrExternalUsersable) types.Object {
	// GuestOrExternalUserTypes is a bitmask enum
	// Note: SDK parser requires strict comma-separated format (no spaces after commas)
	// Our GraphToFrameworkBitmaskEnumAsSet helper trims whitespace as a safety net
	guestOrExternalUserTypes := convert.GraphToFrameworkBitmaskEnumAsSet(ctx, guests.GetGuestOrExternalUserTypes())

	var externalTenantsObj types.Object
	if externalTenants := guests.GetExternalTenants(); externalTenants != nil {
		membershipKind := convert.GraphToFrameworkEnum(externalTenants.GetMembershipKind())

		// Check the @odata.type to determine the actual type
		// ConditionalAccessEnumeratedExternalTenants has Members, ConditionalAccessAllExternalTenants does not
		var members types.Set
		if additionalData := externalTenants.GetAdditionalData(); additionalData != nil {
			if odataType, ok := additionalData["@odata.type"].(string); ok {
				if odataType == "#microsoft.graph.conditionalAccessEnumeratedExternalTenants" {
					// Try to cast to enumerated type which has GetMembers
					if membersSlice, ok := additionalData["members"].([]any); ok {
						stringMembers := make([]string, len(membersSlice))
						for i, m := range membersSlice {
							if str, ok := m.(string); ok {
								stringMembers[i] = str
							}
						}
						members = convert.GraphToFrameworkStringSet(ctx, stringMembers)
					} else {
						members = types.SetNull(types.StringType)
					}
				} else {
					// For "all" type, members is not applicable
					members = types.SetNull(types.StringType)
				}
			} else {
				members = types.SetNull(types.StringType)
			}
		} else {
			members = types.SetNull(types.StringType)
		}

		externalTenantsObj, _ = types.ObjectValue(
			map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"membership_kind": membershipKind,
				"members":         members,
			},
		)
	} else {
		externalTenantsObj = types.ObjectNull(map[string]attr.Type{
			"membership_kind": types.StringType,
			"members":         types.SetType{ElemType: types.StringType},
		})
	}

	obj, _ := types.ObjectValue(
		map[string]attr.Type{
			"guest_or_external_user_types": types.SetType{ElemType: types.StringType},
			"external_tenants": types.ObjectType{AttrTypes: map[string]attr.Type{
				"membership_kind": types.StringType,
				"members":         types.SetType{ElemType: types.StringType},
			}},
		},
		map[string]attr.Value{
			"guest_or_external_user_types": guestOrExternalUserTypes,
			"external_tenants":             externalTenantsObj,
		},
	)

	return obj
}

func mapAuthenticationStrength(ctx context.Context, authenticationStrength graphmodels.AuthenticationStrengthPolicyable) *ConditionalAccessAuthenticationStrength {
	result := &ConditionalAccessAuthenticationStrength{}

	result.ID = convert.GraphToFrameworkString(authenticationStrength.GetId())
	result.DisplayName = convert.GraphToFrameworkString(authenticationStrength.GetDisplayName())
	result.Description = convert.GraphToFrameworkString(authenticationStrength.GetDescription())
	result.PolicyType = convert.GraphToFrameworkEnum(authenticationStrength.GetPolicyType())
	result.RequirementsSatisfied = convert.GraphToFrameworkEnum(authenticationStrength.GetRequirementsSatisfied())
	result.AllowedCombinations = mapEnumCollectionToSet(ctx, authenticationStrength.GetAllowedCombinations(), "allowedCombinations")

	return result
}

func mapApplicationEnforcedRestrictions(ctx context.Context, applicationEnforcedRestrictions graphmodels.ApplicationEnforcedRestrictionsSessionControlable) *ConditionalAccessApplicationEnforcedRestrictions {
	result := &ConditionalAccessApplicationEnforcedRestrictions{}

	result.IsEnabled = convert.GraphToFrameworkBool(applicationEnforcedRestrictions.GetIsEnabled())

	return result
}

func mapCloudAppSecurity(ctx context.Context, cloudAppSecurity graphmodels.CloudAppSecuritySessionControlable) *ConditionalAccessCloudAppSecurity {
	result := &ConditionalAccessCloudAppSecurity{}

	result.IsEnabled = convert.GraphToFrameworkBool(cloudAppSecurity.GetIsEnabled())
	result.CloudAppSecurityType = convert.GraphToFrameworkEnum(cloudAppSecurity.GetCloudAppSecurityType())

	return result
}

func mapSignInFrequency(ctx context.Context, signInFrequency graphmodels.SignInFrequencySessionControlable) *ConditionalAccessSignInFrequency {
	result := &ConditionalAccessSignInFrequency{}

	result.IsEnabled = convert.GraphToFrameworkBool(signInFrequency.GetIsEnabled())
	result.Type = convert.GraphToFrameworkEnum(signInFrequency.GetTypeEscaped())
	result.Value = convert.GraphToFrameworkInt32(signInFrequency.GetValue())
	result.AuthenticationType = convert.GraphToFrameworkEnum(signInFrequency.GetAuthenticationType())
	result.FrequencyInterval = convert.GraphToFrameworkEnum(signInFrequency.GetFrequencyInterval())

	return result
}

func mapPersistentBrowser(ctx context.Context, persistentBrowser graphmodels.PersistentBrowserSessionControlable) *ConditionalAccessPersistentBrowser {
	result := &ConditionalAccessPersistentBrowser{}

	result.IsEnabled = convert.GraphToFrameworkBool(persistentBrowser.GetIsEnabled())
	result.Mode = convert.GraphToFrameworkEnum(persistentBrowser.GetMode())

	return result
}

func mapContinuousAccessEvaluation(ctx context.Context, continuousAccessEvaluation graphmodels.ContinuousAccessEvaluationSessionControlable) *ConditionalAccessContinuousAccessEvaluation {
	result := &ConditionalAccessContinuousAccessEvaluation{}

	result.Mode = convert.GraphToFrameworkEnum(continuousAccessEvaluation.GetMode())

	return result
}

func mapSecureSignInSession(ctx context.Context, secureSignInSession graphmodels.SecureSignInSessionControlable) *ConditionalAccessSecureSignInSession {
	result := &ConditionalAccessSecureSignInSession{}

	result.IsEnabled = convert.GraphToFrameworkBool(secureSignInSession.GetIsEnabled())

	return result
}
