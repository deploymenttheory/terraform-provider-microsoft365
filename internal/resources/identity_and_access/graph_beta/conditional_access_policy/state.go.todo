// MapRemoteResourceStateToTerraform states the base properties of a ConditionalAccessPolicyResourceModel to a Terraform state
package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the base properties of a ConditionalAccessPolicyResourceModel to a Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource graphmodels.ConditionalAccessPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.State = state.EnumPtrToTypeString(remoteResource.GetState())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.ModifiedDateTime = state.TimeToString(remoteResource.GetModifiedDateTime())

	// Map Conditions
	if conditions := remoteResource.GetConditions(); conditions != nil {
		data.Conditions = mapConditionalAccessConditionSet(ctx, conditions)
	}

	// Map Grant Controls
	if grantControls := remoteResource.GetGrantControls(); grantControls != nil {
		data.GrantControls = mapConditionalAccessGrantControls(ctx, grantControls)
	}

	// Map Session Controls
	if sessionControls := remoteResource.GetSessionControls(); sessionControls != nil {
		data.SessionControls = mapConditionalAccessSessionControls(ctx, sessionControls)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapConditionalAccessConditionSet maps the condition set from API response to Terraform state
func mapConditionalAccessConditionSet(ctx context.Context, conditions graphmodels.ConditionalAccessConditionSetable) *ConditionalAccessConditionSetResourceModel {
	conditionSet := &ConditionalAccessConditionSetResourceModel{}

	// Map Client App Types
	if clientAppTypes := conditions.GetClientAppTypes(); clientAppTypes != nil {
		clientAppTypeStrings := make([]attr.Value, 0)
		for _, appType := range clientAppTypes {
			clientAppTypeStrings = append(clientAppTypeStrings, types.StringValue(string(appType)))
		}
		conditionSet.ClientAppTypes = types.SetValueMust(types.StringType, clientAppTypeStrings)
	} else {
		conditionSet.ClientAppTypes = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Map Sign-in Risk Levels
	if signInRiskLevels := conditions.GetSignInRiskLevels(); signInRiskLevels != nil {
		riskLevelStrings := make([]attr.Value, 0)
		for _, riskLevel := range signInRiskLevels {
			riskLevelStrings = append(riskLevelStrings, types.StringValue(string(riskLevel)))
		}
		conditionSet.SignInRiskLevels = types.SetValueMust(types.StringType, riskLevelStrings)
	} else {
		conditionSet.SignInRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Map User Risk Levels
	if userRiskLevels := conditions.GetUserRiskLevels(); userRiskLevels != nil {
		riskLevelStrings := make([]attr.Value, 0)
		for _, riskLevel := range userRiskLevels {
			riskLevelStrings = append(riskLevelStrings, types.StringValue(string(riskLevel)))
		}
		conditionSet.UserRiskLevels = types.SetValueMust(types.StringType, riskLevelStrings)
	} else {
		conditionSet.UserRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Map Service Principal Risk Levels
	if servicePrincipalRiskLevels := conditions.GetServicePrincipalRiskLevels(); servicePrincipalRiskLevels != nil {
		riskLevelStrings := make([]attr.Value, 0)
		for _, riskLevel := range servicePrincipalRiskLevels {
			riskLevelStrings = append(riskLevelStrings, types.StringValue(string(riskLevel)))
		}
		conditionSet.ServicePrincipalRiskLevels = types.SetValueMust(types.StringType, riskLevelStrings)
	} else {
		conditionSet.ServicePrincipalRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Map Insider Risk Levels
	if insiderRiskLevels := conditions.GetInsiderRiskLevels(); insiderRiskLevels != nil {
		riskLevelStrings := make([]attr.Value, 0)
		for _, riskLevel := range insiderRiskLevels {
			riskLevelStrings = append(riskLevelStrings, types.StringValue(string(riskLevel)))
		}
		conditionSet.InsiderRiskLevels = types.SetValueMust(types.StringType, riskLevelStrings)
	} else {
		conditionSet.InsiderRiskLevels = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Map Applications
	if applications := conditions.GetApplications(); applications != nil {
		conditionSet.Applications = &ConditionalAccessApplicationsResourceModel{
			IncludeApplications:                         state.StringSliceToSet(ctx, applications.GetIncludeApplications()),
			ExcludeApplications:                         state.StringSliceToSet(ctx, applications.GetExcludeApplications()),
			IncludeUserActions:                          state.StringSliceToSet(ctx, applications.GetIncludeUserActions()),
			IncludeAuthenticationContextClassReferences: state.StringSliceToSet(ctx, applications.GetIncludeAuthenticationContextClassReferences()),
		}

		// Map Application Filter
		if appFilter := applications.GetApplicationFilter(); appFilter != nil {
			conditionSet.Applications.ApplicationFilter = &ConditionalAccessFilterResourceModel{
				Mode: state.EnumPtrToTypeString(appFilter.GetMode()),
				Rule: types.StringPointerValue(appFilter.GetRule()),
			}
		}
	}

	// Map Authentication Flows
	if authFlows := conditions.GetAuthenticationFlows(); authFlows != nil {
		conditionSet.AuthenticationFlows = &ConditionalAccessAuthenticationFlowsResourceModel{}

		if transferMethods := authFlows.GetTransferMethods(); transferMethods != nil {
			methodStrings := make([]attr.Value, 0)
			for _, method := range transferMethods {
				methodStrings = append(methodStrings, types.StringValue(string(method)))
			}
			conditionSet.AuthenticationFlows.TransferMethods = types.SetValueMust(types.StringType, methodStrings)
		} else {
			conditionSet.AuthenticationFlows.TransferMethods = types.SetValueMust(types.StringType, []attr.Value{})
		}
	}

	// Map Users
	if users := conditions.GetUsers(); users != nil {
		conditionSet.Users = &ConditionalAccessUsersResourceModel{
			IncludeUsers:  state.StringSliceToSet(ctx, users.GetIncludeUsers()),
			ExcludeUsers:  state.StringSliceToSet(ctx, users.GetExcludeUsers()),
			IncludeGroups: state.StringSliceToSet(ctx, users.GetIncludeGroups()),
			ExcludeGroups: state.StringSliceToSet(ctx, users.GetExcludeGroups()),
			IncludeRoles:  state.StringSliceToSet(ctx, users.GetIncludeRoles()),
			ExcludeRoles:  state.StringSliceToSet(ctx, users.GetExcludeRoles()),
		}

		// Map Include Guests or External Users
		if includeGuests := users.GetIncludeGuestsOrExternalUsers(); includeGuests != nil {
			conditionSet.Users.IncludeGuestsOrExternalUsers = mapGuestsOrExternalUsers(ctx, includeGuests)
		}

		// Map Exclude Guests or External Users
		if excludeGuests := users.GetExcludeGuestsOrExternalUsers(); excludeGuests != nil {
			conditionSet.Users.ExcludeGuestsOrExternalUsers = mapGuestsOrExternalUsers(ctx, excludeGuests)
		}
	}

	// Map Client Applications
	if clientApps := conditions.GetClientApplications(); clientApps != nil {
		conditionSet.ClientApplications = &ConditionalAccessClientApplicationsResourceModel{
			IncludeServicePrincipals: state.StringSliceToSet(ctx, clientApps.GetIncludeServicePrincipals()),
			ExcludeServicePrincipals: state.StringSliceToSet(ctx, clientApps.GetExcludeServicePrincipals()),
		}

		// Map Service Principal Filter
		if spFilter := clientApps.GetServicePrincipalFilter(); spFilter != nil {
			conditionSet.ClientApplications.ServicePrincipalFilter = &ConditionalAccessFilterResourceModel{
				Mode: state.EnumPtrToTypeString(spFilter.GetMode()),
				Rule: types.StringPointerValue(spFilter.GetRule()),
			}
		}
	}

	// Map Device States (deprecated)
	if deviceStates := conditions.GetDeviceStates(); deviceStates != nil {
		conditionSet.DeviceStates = &ConditionalAccessDeviceStatesResourceModel{
			IncludeStates: state.StringSliceToSet(ctx, deviceStates.GetIncludeStates()),
			ExcludeStates: state.StringSliceToSet(ctx, deviceStates.GetExcludeStates()),
		}
	}

	// Map Devices
	if devices := conditions.GetDevices(); devices != nil {
		conditionSet.Devices = &ConditionalAccessDevicesResourceModel{
			IncludeDevices:      state.StringSliceToSet(ctx, devices.GetIncludeDevices()),
			ExcludeDevices:      state.StringSliceToSet(ctx, devices.GetExcludeDevices()),
			IncludeDeviceStates: state.StringSliceToSet(ctx, devices.GetIncludeDeviceStates()),
			ExcludeDeviceStates: state.StringSliceToSet(ctx, devices.GetExcludeDeviceStates()),
		}

		// Map Device Filter
		if deviceFilter := devices.GetDeviceFilter(); deviceFilter != nil {
			conditionSet.Devices.DeviceFilter = &ConditionalAccessFilterResourceModel{
				Mode: state.EnumPtrToTypeString(deviceFilter.GetMode()),
				Rule: types.StringPointerValue(deviceFilter.GetRule()),
			}

		}
	}

	// Map Locations
	if locations := conditions.GetLocations(); locations != nil {
		conditionSet.Locations = &ConditionalAccessLocationsResourceModel{
			IncludeLocations: state.StringSliceToSet(ctx, locations.GetIncludeLocations()),
			ExcludeLocations: state.StringSliceToSet(ctx, locations.GetExcludeLocations()),
		}
	}

	// Map Platforms
	if platforms := conditions.GetPlatforms(); platforms != nil {
		conditionSet.Platforms = &ConditionalAccessPlatformsResourceModel{}

		// Map Include Platforms
		if includePlatforms := platforms.GetIncludePlatforms(); includePlatforms != nil {
			platformStrings := make([]attr.Value, 0)
			for _, platform := range includePlatforms {
				platformStrings = append(platformStrings, types.StringValue(string(platform)))
			}
			conditionSet.Platforms.IncludePlatforms = types.SetValueMust(types.StringType, platformStrings)
		} else {
			conditionSet.Platforms.IncludePlatforms = types.SetValueMust(types.StringType, []attr.Value{})
		}

		// Map Exclude Platforms
		if excludePlatforms := platforms.GetExcludePlatforms(); excludePlatforms != nil {
			platformStrings := make([]attr.Value, 0)
			for _, platform := range excludePlatforms {
				platformStrings = append(platformStrings, types.StringValue(string(platform)))
			}
			conditionSet.Platforms.ExcludePlatforms = types.SetValueMust(types.StringType, platformStrings)
		} else {
			conditionSet.Platforms.ExcludePlatforms = types.SetValueMust(types.StringType, []attr.Value{})
		}
	}

	return conditionSet
}

// mapGuestsOrExternalUsers maps guests or external users from API response to Terraform state
func mapGuestsOrExternalUsers(ctx context.Context, guestsOrExternal graphmodels.ConditionalAccessGuestsOrExternalUsersable) *ConditionalAccessGuestsOrExternalUsersResourceModel {
	result := &ConditionalAccessGuestsOrExternalUsersResourceModel{}

	// Map Guest or External User Types
	if userTypes := guestsOrExternal.GetGuestOrExternalUserTypes(); userTypes != nil {
		userTypeStrings := make([]attr.Value, 0)
		for _, userType := range userTypes {
			userTypeStrings = append(userTypeStrings, types.StringValue(string(userType)))
		}
		result.GuestOrExternalUserTypes = types.SetValueMust(types.StringType, userTypeStrings)
	} else {
		result.GuestOrExternalUserTypes = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Map External Tenants
	if externalTenants := guestsOrExternal.GetExternalTenants(); externalTenants != nil {
		result.ExternalTenants = &ConditionalAccessExternalTenantsResourceModel{
			MembershipKind: state.EnumPtrToTypeString(externalTenants.GetMembershipKind()),
			Members:        state.StringSliceToSet(ctx, externalTenants.GetMembers()),
		}
	}

	return result
}

// mapConditionalAccessGrantControls maps the grant controls from API response to Terraform state
func mapConditionalAccessGrantControls(ctx context.Context, grantControls graphmodels.ConditionalAccessGrantControlsable) *ConditionalAccessGrantControlsResourceModel {
	grantControlsModel := &ConditionalAccessGrantControlsResourceModel{
		Operator:                    types.StringPointerValue(grantControls.GetOperator()),
		CustomAuthenticationFactors: state.StringSliceToSet(ctx, grantControls.GetCustomAuthenticationFactors()),
		TermsOfUse:                  state.StringSliceToSet(ctx, grantControls.GetTermsOfUse()),
	}

	// Map Built-in Controls
	if builtInControls := grantControls.GetBuiltInControls(); builtInControls != nil {
		controlStrings := make([]attr.Value, 0)
		for _, control := range builtInControls {
			controlStrings = append(controlStrings, types.StringValue(string(control)))
		}
		grantControlsModel.BuiltInControls = types.SetValueMust(types.StringType, controlStrings)
	} else {
		grantControlsModel.BuiltInControls = types.SetValueMust(types.StringType, []attr.Value{})
	}

	return grantControlsModel
}

// mapConditionalAccessSessionControls maps the session controls from API response to Terraform state
func mapConditionalAccessSessionControls(ctx context.Context, sessionControls graphmodels.ConditionalAccessSessionControlsable) *ConditionalAccessSessionControlsResourceModel {
	sessionControlsModel := &ConditionalAccessSessionControlsResourceModel{
		DisableResilienceDefaults: types.BoolPointerValue(sessionControls.GetDisableResilienceDefaults()),
	}

	// Map Application Enforced Restrictions
	if appRestrictions := sessionControls.GetApplicationEnforcedRestrictions(); appRestrictions != nil {
		sessionControlsModel.ApplicationEnforcedRestrictions = &ApplicationEnforcedRestrictionsSessionControlResourceModel{
			IsEnabled: types.BoolPointerValue(appRestrictions.GetIsEnabled()),
		}
	}

	// Map Cloud App Security
	if cloudAppSecurity := sessionControls.GetCloudAppSecurity(); cloudAppSecurity != nil {
		sessionControlsModel.CloudAppSecurity = &CloudAppSecuritySessionControlResourceModel{
			IsEnabled:            types.BoolPointerValue(cloudAppSecurity.GetIsEnabled()),
			CloudAppSecurityType: state.EnumPtrToTypeString(cloudAppSecurity.GetCloudAppSecurityType()),
		}
	}

	// Map Sign-in Frequency
	if signInFreq := sessionControls.GetSignInFrequency(); signInFreq != nil {
		sessionControlsModel.SignInFrequency = &SignInFrequencySessionControlResourceModel{
			IsEnabled:          types.BoolPointerValue(signInFreq.GetIsEnabled()),
			Type:               state.EnumPtrToTypeString(signInFreq.GetTypeEscaped()),
			Value:              types.Int32PointerValue(signInFreq.GetValue()),
			AuthenticationType: state.EnumPtrToTypeString(signInFreq.GetAuthenticationType()),
			FrequencyInterval:  state.EnumPtrToTypeString(signInFreq.GetFrequencyInterval()),
		}
	}

	// Map Persistent Browser
	if persistentBrowser := sessionControls.GetPersistentBrowser(); persistentBrowser != nil {
		sessionControlsModel.PersistentBrowser = &PersistentBrowserSessionControlResourceModel{
			IsEnabled: types.BoolPointerValue(persistentBrowser.GetIsEnabled()),
			Mode:      state.EnumPtrToTypeString(persistentBrowser.GetMode()),
		}
	}

	return sessionControlsModel
}
