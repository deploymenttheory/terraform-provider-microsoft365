package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"
	"math"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the graph beta SDK model
func constructResource(ctx context.Context, typeName string, data *ConditionalAccessPolicyResourceModel) (*models.ConditionalAccessPolicy, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", typeName))

	requestBody := models.NewConditionalAccessPolicy()

	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.State.IsNull() {
		stateStr := data.State.ValueString()
		stateAny, err := models.ParseConditionalAccessPolicyState(stateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid state: %s", err)
		}
		if stateAny != nil {
			state, ok := stateAny.(*models.ConditionalAccessPolicyState)
			if !ok {
				return nil, fmt.Errorf("unexpected type for state: %T", stateAny)
			}
			requestBody.SetState(state)
		}
	}

	if data.Conditions != nil {
		conditions, err := constructConditions(data.Conditions)
		if err != nil {
			return nil, fmt.Errorf("error constructing conditions: %s", err)
		}
		requestBody.SetConditions(conditions)
	}

	if data.GrantControls != nil {
		grantControls, err := constructGrantControls(data.GrantControls)
		if err != nil {
			return nil, fmt.Errorf("error constructing grant controls: %s", err)
		}
		requestBody.SetGrantControls(grantControls)
	}

	if data.SessionControls != nil {
		sessionControls, err := constructSessionControls(data.SessionControls)
		if err != nil {
			return nil, fmt.Errorf("error constructing session controls: %s", err)
		}
		requestBody.SetSessionControls(sessionControls)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", typeName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", typeName))

	return requestBody, nil
}

// Helper functions to construct nested objects
func constructConditions(data *ConditionalAccessConditionsModel) (*models.ConditionalAccessConditionSet, error) {
	if data == nil {
		return nil, nil
	}

	conditions := models.NewConditionalAccessConditionSet()

	// Applications
	if data.Applications != nil {
		applications, err := constructApplications(data.Applications)
		if err != nil {
			return nil, fmt.Errorf("error constructing applications: %v", err)
		}
		conditions.SetApplications(applications)
	}

	// Authentication Flows
	if data.AuthenticationFlows != nil {
		authFlows, err := constructAuthenticationFlows(data.AuthenticationFlows)
		if err != nil {
			return nil, fmt.Errorf("error constructing authentication flows: %v", err)
		}
		conditions.SetAuthenticationFlows(authFlows)
	}

	// Client Applications
	if data.ClientApplications != nil {
		clientApps, err := constructClientApplications(data.ClientApplications)
		if err != nil {
			return nil, fmt.Errorf("error constructing client applications: %v", err)
		}
		conditions.SetClientApplications(clientApps)
	}

	// Client App Types
	if len(data.ClientAppTypes) > 0 {
		clientAppTypes := make([]models.ConditionalAccessClientApp, 0, len(data.ClientAppTypes))
		for _, appType := range data.ClientAppTypes {
			clientAppTypeAny, err := models.ParseConditionalAccessClientApp(appType.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing client app type: %v", err)
			}
			if clientAppTypeAny != nil {
				clientAppType, ok := clientAppTypeAny.(*models.ConditionalAccessClientApp)
				if !ok {
					return nil, fmt.Errorf("unexpected type for client app type: %T", clientAppTypeAny)
				}
				clientAppTypes = append(clientAppTypes, *clientAppType)
			}
		}
		if len(clientAppTypes) > 0 {
			conditions.SetClientAppTypes(clientAppTypes)
		}
	}

	// Devices
	if data.Devices != nil {
		devices, err := constructDevices(data.Devices)
		if err != nil {
			return nil, fmt.Errorf("error constructing devices: %v", err)
		}
		conditions.SetDevices(devices)
	}

	// Device States (deprecated)
	if data.DeviceStates != nil {
		deviceStates, err := constructDeviceStates(data.DeviceStates)
		if err != nil {
			return nil, fmt.Errorf("error constructing device states: %v", err)
		}
		conditions.SetDeviceStates(deviceStates)
	}

	// Insider Risk Levels
	if !data.InsiderRiskLevels.IsNull() {
		insiderRiskLevelAny, err := models.ParseConditionalAccessInsiderRiskLevels(data.InsiderRiskLevels.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing insider risk level: %v", err)
		}
		if insiderRiskLevelAny != nil {
			insiderRiskLevel, ok := insiderRiskLevelAny.(*models.ConditionalAccessInsiderRiskLevels)
			if !ok {
				return nil, fmt.Errorf("unexpected type for insider risk level: %T", insiderRiskLevelAny)
			}
			conditions.SetInsiderRiskLevels(insiderRiskLevel)
		}
	}

	// Locations
	if data.Locations != nil {
		locations, err := constructLocations(data.Locations)
		if err != nil {
			return nil, fmt.Errorf("error constructing locations: %v", err)
		}
		conditions.SetLocations(locations)
	}

	// Platforms
	if data.Platforms != nil {
		platforms, err := constructPlatforms(data.Platforms)
		if err != nil {
			return nil, fmt.Errorf("error constructing platforms: %v", err)
		}
		conditions.SetPlatforms(platforms)
	}

	// Service Principal Risk Levels
	if len(data.ServicePrincipalRiskLevels) > 0 {
		riskLevels := make([]models.RiskLevel, 0, len(data.ServicePrincipalRiskLevels))
		for _, level := range data.ServicePrincipalRiskLevels {
			riskLevelAny, err := models.ParseRiskLevel(level.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing service principal risk level: %v", err)
			}
			if riskLevelAny != nil {
				riskLevel, ok := riskLevelAny.(*models.RiskLevel)
				if !ok {
					return nil, fmt.Errorf("unexpected type for risk level: %T", riskLevelAny)
				}
				riskLevels = append(riskLevels, *riskLevel)
			}
		}
		if len(riskLevels) > 0 {
			conditions.SetServicePrincipalRiskLevels(riskLevels)
		}
	}

	// Sign-in Risk Levels
	if len(data.SignInRiskLevels) > 0 {
		signInRiskLevels := make([]models.RiskLevel, 0, len(data.SignInRiskLevels))
		for _, level := range data.SignInRiskLevels {
			riskLevelAny, err := models.ParseRiskLevel(level.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing sign-in risk level: %v", err)
			}
			if riskLevelAny != nil {
				riskLevel, ok := riskLevelAny.(*models.RiskLevel)
				if !ok {
					return nil, fmt.Errorf("unexpected type for sign-in risk level: %T", riskLevelAny)
				}
				signInRiskLevels = append(signInRiskLevels, *riskLevel)
			}
		}
		if len(signInRiskLevels) > 0 {
			conditions.SetSignInRiskLevels(signInRiskLevels)
		}
	}

	// User Risk Levels
	if len(data.UserRiskLevels) > 0 {
		userRiskLevels := make([]models.RiskLevel, 0, len(data.UserRiskLevels))
		for _, level := range data.UserRiskLevels {
			riskLevelAny, err := models.ParseRiskLevel(level.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing user risk level: %v", err)
			}
			if riskLevelAny != nil {
				riskLevel, ok := riskLevelAny.(*models.RiskLevel)
				if !ok {
					return nil, fmt.Errorf("unexpected type for user risk level: %T", riskLevelAny)
				}
				userRiskLevels = append(userRiskLevels, *riskLevel)
			}
		}
		if len(userRiskLevels) > 0 {
			conditions.SetUserRiskLevels(userRiskLevels)
		}
	}

	// Users
	if data.Users != nil {
		users, err := constructUsers(data.Users)
		if err != nil {
			return nil, fmt.Errorf("error constructing users: %v", err)
		}
		conditions.SetUsers(users)
	}

	return conditions, nil
}

func constructApplications(data *ConditionalAccessApplicationsModel) (models.ConditionalAccessApplicationsable, error) {
	if data == nil {
		return nil, nil
	}

	applications := models.NewConditionalAccessApplications()

	if len(data.IncludeApplications) > 0 {
		includeApps := make([]string, len(data.IncludeApplications))
		for i, app := range data.IncludeApplications {
			includeApps[i] = app.ValueString()
		}
		applications.SetIncludeApplications(includeApps)
	}

	if len(data.ExcludeApplications) > 0 {
		excludeApps := make([]string, len(data.ExcludeApplications))
		for i, app := range data.ExcludeApplications {
			excludeApps[i] = app.ValueString()
		}
		applications.SetExcludeApplications(excludeApps)
	}

	if len(data.IncludeUserActions) > 0 {
		userActions := make([]string, len(data.IncludeUserActions))
		for i, action := range data.IncludeUserActions {
			userActions[i] = action.ValueString()
		}
		applications.SetIncludeUserActions(userActions)
	}

	return applications, nil
}

func constructAuthenticationFlows(data *ConditionalAccessAuthenticationFlowsModel) (models.ConditionalAccessAuthenticationFlowsable, error) {
	if data == nil {
		return nil, nil
	}

	authFlows := models.NewConditionalAccessAuthenticationFlows()

	if !data.TransferMethods.IsNull() {
		transferMethodsStr := data.TransferMethods.ValueString()
		transferMethods, err := models.ParseConditionalAccessTransferMethods(transferMethodsStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing transfer methods: %v", err)
		}
		if transferMethods != nil {
			authFlows.SetTransferMethods(transferMethods.(*models.ConditionalAccessTransferMethods))
		}
	}

	return authFlows, nil
}

func constructUsers(data *ConditionalAccessUsersModel) (models.ConditionalAccessUsersable, error) {
	if data == nil {
		return nil, nil
	}

	users := models.NewConditionalAccessUsers()

	if len(data.IncludeUsers) > 0 {
		includeUsers := make([]string, len(data.IncludeUsers))
		for i, user := range data.IncludeUsers {
			includeUsers[i] = user.ValueString()
		}
		users.SetIncludeUsers(includeUsers)
	}

	if len(data.ExcludeUsers) > 0 {
		excludeUsers := make([]string, len(data.ExcludeUsers))
		for i, user := range data.ExcludeUsers {
			excludeUsers[i] = user.ValueString()
		}
		users.SetExcludeUsers(excludeUsers)
	}

	if len(data.IncludeGroups) > 0 {
		includeGroups := make([]string, len(data.IncludeGroups))
		for i, group := range data.IncludeGroups {
			includeGroups[i] = group.ValueString()
		}
		users.SetIncludeGroups(includeGroups)
	}

	if len(data.ExcludeGroups) > 0 {
		excludeGroups := make([]string, len(data.ExcludeGroups))
		for i, group := range data.ExcludeGroups {
			excludeGroups[i] = group.ValueString()
		}
		users.SetExcludeGroups(excludeGroups)
	}

	if len(data.IncludeRoles) > 0 {
		includeRoles := make([]string, len(data.IncludeRoles))
		for i, role := range data.IncludeRoles {
			includeRoles[i] = role.ValueString()
		}
		users.SetIncludeRoles(includeRoles)
	}

	if len(data.ExcludeRoles) > 0 {
		excludeRoles := make([]string, len(data.ExcludeRoles))
		for i, role := range data.ExcludeRoles {
			excludeRoles[i] = role.ValueString()
		}
		users.SetExcludeRoles(excludeRoles)
	}

	if data.IncludeGuestsOrExternalUsers != nil {
		guestsOrExternalUsers, err := constructGuestsOrExternalUsers(data.IncludeGuestsOrExternalUsers)
		if err != nil {
			return nil, fmt.Errorf("error constructing include guests or external users: %v", err)
		}
		users.SetIncludeGuestsOrExternalUsers(guestsOrExternalUsers)
	}

	if data.ExcludeGuestsOrExternalUsers != nil {
		guestsOrExternalUsers, err := constructGuestsOrExternalUsers(data.ExcludeGuestsOrExternalUsers)
		if err != nil {
			return nil, fmt.Errorf("error constructing exclude guests or external users: %v", err)
		}
		users.SetExcludeGuestsOrExternalUsers(guestsOrExternalUsers)
	}

	return users, nil
}

func constructGuestsOrExternalUsers(data *ConditionalAccessGuestsOrExternalUsersModel) (models.ConditionalAccessGuestsOrExternalUsersable, error) {
	if data == nil {
		return nil, nil
	}

	guestsOrExternalUsers := models.NewConditionalAccessGuestsOrExternalUsers()

	if !data.GuestOrExternalUserTypes.IsNull() {
		userTypesAny, err := models.ParseConditionalAccessGuestOrExternalUserTypes(data.GuestOrExternalUserTypes.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing guest or external user types: %v", err)
		}
		if userTypesAny != nil {
			userTypes, ok := userTypesAny.(*models.ConditionalAccessGuestOrExternalUserTypes)
			if !ok {
				return nil, fmt.Errorf("unexpected type for guest or external user types: %T", userTypesAny)
			}
			guestsOrExternalUsers.SetGuestOrExternalUserTypes(userTypes)
		}
	}

	if data.ExternalTenants != nil {
		externalTenants, err := constructConditionalAccessExternalTenants(data.ExternalTenants)
		if err != nil {
			return nil, fmt.Errorf("error constructing external tenants: %v", err)
		}
		guestsOrExternalUsers.SetExternalTenants(externalTenants)
	}

	return guestsOrExternalUsers, nil
}

func constructConditionalAccessExternalTenants(data *ConditionalAccessExternalTenantsModel) (models.ConditionalAccessExternalTenantsable, error) {
	if data == nil {
		return nil, nil
	}

	externalTenants := models.NewConditionalAccessExternalTenants()

	if !data.MembershipKind.IsNull() {
		membershipKindAny, err := models.ParseConditionalAccessExternalTenantsMembershipKind(data.MembershipKind.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing membership kind: %v", err)
		}
		if membershipKindAny != nil {
			membershipKind, ok := membershipKindAny.(*models.ConditionalAccessExternalTenantsMembershipKind)
			if !ok {
				return nil, fmt.Errorf("unexpected type for membership kind: %T", membershipKindAny)
			}
			externalTenants.SetMembershipKind(membershipKind)
		}
	}

	return externalTenants, nil
}

func constructClientApplications(data *ConditionalAccessClientApplicationsModel) (models.ConditionalAccessClientApplicationsable, error) {
	if data == nil {
		return nil, nil
	}

	clientApps := models.NewConditionalAccessClientApplications()

	if len(data.IncludeServicePrincipals) > 0 {
		includeServicePrincipals := make([]string, len(data.IncludeServicePrincipals))
		for i, sp := range data.IncludeServicePrincipals {
			includeServicePrincipals[i] = sp.ValueString()
		}
		clientApps.SetIncludeServicePrincipals(includeServicePrincipals)
	}

	if len(data.ExcludeServicePrincipals) > 0 {
		excludeServicePrincipals := make([]string, len(data.ExcludeServicePrincipals))
		for i, sp := range data.ExcludeServicePrincipals {
			excludeServicePrincipals[i] = sp.ValueString()
		}
		clientApps.SetExcludeServicePrincipals(excludeServicePrincipals)
	}

	return clientApps, nil
}

func constructDevices(data *ConditionalAccessDevicesModel) (models.ConditionalAccessDevicesable, error) {
	if data == nil {
		return nil, nil
	}

	devices := models.NewConditionalAccessDevices()

	if len(data.IncludeDevices) > 0 {
		includeDevices := make([]string, len(data.IncludeDevices))
		for i, device := range data.IncludeDevices {
			includeDevices[i] = device.ValueString()
		}
		devices.SetIncludeDevices(includeDevices)
	}

	if len(data.ExcludeDevices) > 0 {
		excludeDevices := make([]string, len(data.ExcludeDevices))
		for i, device := range data.ExcludeDevices {
			excludeDevices[i] = device.ValueString()
		}
		devices.SetExcludeDevices(excludeDevices)
	}

	if data.IncludeStates != nil {
		if len(data.IncludeStates) > 0 {
			includeStates := make([]string, len(data.IncludeStates))
			for i, state := range data.IncludeStates {
				includeStates[i] = state.ValueString()
			}
			devices.SetIncludeDeviceStates(includeStates)
		}

		if len(data.ExcludeStates) > 0 {
			excludeStates := make([]string, len(data.ExcludeStates))
			for i, state := range data.ExcludeStates {
				excludeStates[i] = state.ValueString()
			}
			devices.SetExcludeDeviceStates(excludeStates)
		}
	}

	if data.DeviceFilter != nil {
		filter := models.NewConditionalAccessFilter()

		if !data.DeviceFilter.Mode.IsNull() {
			modeStr := data.DeviceFilter.Mode.ValueString()
			modeAny, err := models.ParseFilterMode(modeStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing device filter mode: %v", err)
			}
			if modeAny != nil {
				mode := modeAny.(*models.FilterMode)
				filter.SetMode(mode)
			}
		}

		if !data.DeviceFilter.Rule.IsNull() {
			rule := data.DeviceFilter.Rule.ValueString()
			filter.SetRule(&rule)
		}

		devices.SetDeviceFilter(filter)
	}

	return devices, nil
}

func constructDeviceStates(data *ConditionalAccessDeviceStatesModel) (models.ConditionalAccessDeviceStatesable, error) {
	if data == nil {
		return nil, nil
	}

	deviceStates := models.NewConditionalAccessDeviceStates()

	if len(data.IncludeStates) > 0 {
		includeStates := make([]string, len(data.IncludeStates))
		for i, state := range data.IncludeStates {
			includeStates[i] = state.ValueString()
		}
		deviceStates.SetIncludeStates(includeStates)
	}

	if len(data.ExcludeStates) > 0 {
		excludeStates := make([]string, len(data.ExcludeStates))
		for i, state := range data.ExcludeStates {
			excludeStates[i] = state.ValueString()
		}
		deviceStates.SetExcludeStates(excludeStates)
	}

	return deviceStates, nil
}

func constructLocations(data *ConditionalAccessLocationsModel) (models.ConditionalAccessLocationsable, error) {
	if data == nil {
		return nil, nil
	}

	locations := models.NewConditionalAccessLocations()

	if len(data.IncludeLocations) > 0 {
		includeLocations := make([]string, len(data.IncludeLocations))
		for i, location := range data.IncludeLocations {
			includeLocations[i] = location.ValueString()
		}
		locations.SetIncludeLocations(includeLocations)
	}

	if len(data.ExcludeLocations) > 0 {
		excludeLocations := make([]string, len(data.ExcludeLocations))
		for i, location := range data.ExcludeLocations {
			excludeLocations[i] = location.ValueString()
		}
		locations.SetExcludeLocations(excludeLocations)
	}

	return locations, nil
}

func constructPlatforms(data *ConditionalAccessPlatformsModel) (models.ConditionalAccessPlatformsable, error) {
	if data == nil {
		return nil, nil
	}

	platforms := models.NewConditionalAccessPlatforms()

	if len(data.IncludePlatforms) > 0 {
		includePlatforms := make([]models.ConditionalAccessDevicePlatform, 0, len(data.IncludePlatforms))
		for _, platform := range data.IncludePlatforms {
			platformAny, err := models.ParseConditionalAccessDevicePlatform(platform.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing include platform: %v", err)
			}
			if platformAny != nil {
				includePlatforms = append(includePlatforms, *platformAny.(*models.ConditionalAccessDevicePlatform))
			}
		}
		platforms.SetIncludePlatforms(includePlatforms)
	}

	if len(data.ExcludePlatforms) > 0 {
		excludePlatforms := make([]models.ConditionalAccessDevicePlatform, 0, len(data.ExcludePlatforms))
		for _, platform := range data.ExcludePlatforms {
			platformAny, err := models.ParseConditionalAccessDevicePlatform(platform.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing exclude platform: %v", err)
			}
			if platformAny != nil {
				excludePlatforms = append(excludePlatforms, *platformAny.(*models.ConditionalAccessDevicePlatform))
			}
		}
		platforms.SetExcludePlatforms(excludePlatforms)
	}

	return platforms, nil
}

func constructGrantControls(data *ConditionalAccessGrantControlsModel) (*models.ConditionalAccessGrantControls, error) {
	if data == nil {
		return nil, nil
	}

	grantControls := models.NewConditionalAccessGrantControls()

	if !data.Operator.IsNull() {
		operator := data.Operator.ValueString()
		grantControls.SetOperator(&operator)
	}

	if len(data.BuiltInControls) > 0 {
		builtInControls := make([]models.ConditionalAccessGrantControl, 0, len(data.BuiltInControls))
		for _, control := range data.BuiltInControls {
			builtInControlAny, err := models.ParseConditionalAccessGrantControl(control.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing built-in control: %v", err)
			}
			if builtInControlAny != nil {
				builtInControl, ok := builtInControlAny.(*models.ConditionalAccessGrantControl)
				if !ok {
					return nil, fmt.Errorf("unexpected type for built-in control: %T", builtInControlAny)
				}
				builtInControls = append(builtInControls, *builtInControl)
			}
		}
		grantControls.SetBuiltInControls(builtInControls)
	}

	if len(data.CustomAuthenticationFactors) > 0 {
		customFactors := make([]string, len(data.CustomAuthenticationFactors))
		for i, factor := range data.CustomAuthenticationFactors {
			customFactors[i] = factor.ValueString()
		}
		grantControls.SetCustomAuthenticationFactors(customFactors)
	}

	if len(data.TermsOfUse) > 0 {
		termsOfUse := make([]string, len(data.TermsOfUse))
		for i, term := range data.TermsOfUse {
			termsOfUse[i] = term.ValueString()
		}
		grantControls.SetTermsOfUse(termsOfUse)
	}

	if data.AuthenticationStrength != nil {
		authStrength, err := constructAuthenticationStrength(data.AuthenticationStrength)
		if err != nil {
			return nil, fmt.Errorf("error constructing authentication strength: %v", err)
		}
		grantControls.SetAuthenticationStrength(authStrength)
	}

	return grantControls, nil
}

func constructAuthenticationStrength(data *AuthenticationStrengthPolicyModel) (*models.AuthenticationStrengthPolicy, error) {
	if data == nil {
		return nil, nil
	}

	authStrength := models.NewAuthenticationStrengthPolicy()

	if !data.DisplayName.IsNull() {
		displayName := data.DisplayName.ValueString()
		authStrength.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		authStrength.SetDescription(&description)
	}

	if !data.PolicyType.IsNull() {
		policyTypeStr := data.PolicyType.ValueString()
		policyTypeAny, err := models.ParseAuthenticationStrengthPolicyType(policyTypeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing policy type: %v", err)
		}
		if policyTypeAny != nil {
			policyType := policyTypeAny.(*models.AuthenticationStrengthPolicyType)
			authStrength.SetPolicyType(policyType)
		}
	}

	if !data.RequirementsSatisfied.IsNull() {
		requirementsSatisfiedStr := data.RequirementsSatisfied.ValueString()
		requirementsSatisfiedAny, err := models.ParseAuthenticationStrengthRequirements(requirementsSatisfiedStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing requirements satisfied: %v", err)
		}
		if requirementsSatisfiedAny != nil {
			requirementsSatisfied := requirementsSatisfiedAny.(*models.AuthenticationStrengthRequirements)
			authStrength.SetRequirementsSatisfied(requirementsSatisfied)
		}
	}
	if len(data.AllowedCombinations) > 0 {
		allowedCombinations := make([]models.AuthenticationMethodModes, 0, len(data.AllowedCombinations))
		for _, combination := range data.AllowedCombinations {
			combinationAny, err := models.ParseAuthenticationMethodModes(combination.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing allowed combination: %v", err)
			}
			if combinationAny != nil {
				authMethodMode := combinationAny.(*models.AuthenticationMethodModes)
				allowedCombinations = append(allowedCombinations, *authMethodMode)
			}
		}
		authStrength.SetAllowedCombinations(allowedCombinations)
	}

	return authStrength, nil
}

func constructSessionControls(data *ConditionalAccessSessionControlsModel) (models.ConditionalAccessSessionControlsable, error) {
	if data == nil {
		return nil, nil
	}

	sessionControls := models.NewConditionalAccessSessionControls()

	if data.ApplicationEnforcedRestrictions != nil {
		appRestrictions := models.NewApplicationEnforcedRestrictionsSessionControl()
		isEnabled := data.ApplicationEnforcedRestrictions.IsEnabled.ValueBool()
		appRestrictions.SetIsEnabled(&isEnabled)
		sessionControls.SetApplicationEnforcedRestrictions(appRestrictions)
	}

	if data.CloudAppSecurity != nil {
		cloudAppSecurity := models.NewCloudAppSecuritySessionControl()
		isEnabled := data.CloudAppSecurity.IsEnabled.ValueBool()
		cloudAppSecurity.SetIsEnabled(&isEnabled)

		if !data.CloudAppSecurity.CloudAppSecurityType.IsNull() {
			typeAny, err := models.ParseCloudAppSecuritySessionControlType(data.CloudAppSecurity.CloudAppSecurityType.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing cloud app security type: %v", err)
			}
			if typeAny != nil {
				cloudAppSecurityType := typeAny.(*models.CloudAppSecuritySessionControlType)
				cloudAppSecurity.SetCloudAppSecurityType(cloudAppSecurityType)
			}
		}

		sessionControls.SetCloudAppSecurity(cloudAppSecurity)
	}

	if data.ContinuousAccessEvaluation != nil {
		continuousAccessEvaluation := models.NewContinuousAccessEvaluationSessionControl()

		if !data.ContinuousAccessEvaluation.Mode.IsNull() {
			modeAny, err := models.ParseContinuousAccessEvaluationMode(data.ContinuousAccessEvaluation.Mode.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing continuous access evaluation mode: %v", err)
			}
			if modeAny != nil {
				mode := modeAny.(*models.ContinuousAccessEvaluationMode)
				continuousAccessEvaluation.SetMode(mode)
			}
		}

		sessionControls.SetContinuousAccessEvaluation(continuousAccessEvaluation)
	}

	if data.PersistentBrowser != nil {
		persistentBrowser := models.NewPersistentBrowserSessionControl()

		isEnabled := data.PersistentBrowser.IsEnabled.ValueBool()
		persistentBrowser.SetIsEnabled(&isEnabled)

		if !data.PersistentBrowser.Mode.IsNull() {
			modeAny, err := models.ParsePersistentBrowserSessionMode(data.PersistentBrowser.Mode.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing persistent browser session mode: %v", err)
			}
			if modeAny != nil {
				mode := modeAny.(*models.PersistentBrowserSessionMode)
				persistentBrowser.SetMode(mode)
			}
		}

		sessionControls.SetPersistentBrowser(persistentBrowser)
	}

	if data.SignInFrequency != nil {
		signInFrequency := models.NewSignInFrequencySessionControl()

		isEnabled := data.SignInFrequency.IsEnabled.ValueBool()
		signInFrequency.SetIsEnabled(&isEnabled)

		if !data.SignInFrequency.Type.IsNull() {
			typeAny, err := models.ParseSigninFrequencyType(data.SignInFrequency.Type.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing sign-in frequency type: %v", err)
			}
			if typeAny != nil {
				freqType := typeAny.(*models.SigninFrequencyType)
				signInFrequency.SetTypeEscaped(freqType)
			}
		}

		if !data.SignInFrequency.Value.IsNull() {
			value := data.SignInFrequency.Value.ValueInt64()
			if value > math.MaxInt32 || value < math.MinInt32 {
				return nil, fmt.Errorf("sign-in frequency value %d is out of range for int32", value)
			}
			int32Value := int32(value)
			signInFrequency.SetValue(&int32Value)
		}

		if !data.SignInFrequency.FrequencyInterval.IsNull() {
			intervalAny, err := models.ParseSignInFrequencyInterval(data.SignInFrequency.FrequencyInterval.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing sign-in frequency interval: %v", err)
			}
			if intervalAny != nil {
				interval := intervalAny.(*models.SignInFrequencyInterval)
				signInFrequency.SetFrequencyInterval(interval)
			}
		}

		if !data.SignInFrequency.AuthenticationType.IsNull() {
			authTypeAny, err := models.ParseSignInFrequencyAuthenticationType(data.SignInFrequency.AuthenticationType.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing sign-in frequency authentication type: %v", err)
			}
			if authTypeAny != nil {
				authType := authTypeAny.(*models.SignInFrequencyAuthenticationType)
				signInFrequency.SetAuthenticationType(authType)
			}
		}

		sessionControls.SetSignInFrequency(signInFrequency)
	}

	if data.SecureSignInSession != nil {
		secureSignInSession := models.NewSecureSignInSessionControl()
		isEnabled := data.SecureSignInSession.IsEnabled.ValueBool()
		secureSignInSession.SetIsEnabled(&isEnabled)
		sessionControls.SetSecureSignInSession(secureSignInSession)
	}

	if !data.DisableResilienceDefaults.IsNull() {
		disableResilienceDefaults := data.DisableResilienceDefaults.ValueBool()
		sessionControls.SetDisableResilienceDefaults(&disableResilienceDefaults)
	}

	return sessionControls, nil
}
