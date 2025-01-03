package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema values to the graph beta SDK model and sets the values for  the request to the Graph API
func constructResource(ctx context.Context, data *ConditionalAccessPolicyResourceModel) (*models.ConditionalAccessPolicy, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewConditionalAccessPolicy()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	if err := constructors.SetEnumProperty(data.State, models.ParseConditionalAccessPolicyState, requestBody.SetState); err != nil {
		return nil, fmt.Errorf("failed to set state: %v", err)
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

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructConditions constructs the ConditionalAccessConditionSet object
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

// constructApplications constructs the ConditionalAccessApplicationsable object
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

	if len(data.IncludeAuthenticationContextClassReferences) > 0 {
		authRefs := make([]string, len(data.IncludeAuthenticationContextClassReferences))
		for i, ref := range data.IncludeAuthenticationContextClassReferences {
			authRefs[i] = ref.ValueString()
		}
		applications.SetIncludeAuthenticationContextClassReferences(authRefs)
	}

	if data.ApplicationFilter != nil {
		filter := models.NewConditionalAccessFilter()
		if !data.ApplicationFilter.Mode.IsNull() {
			modeStr := data.ApplicationFilter.Mode.ValueString()
			modeAny, err := models.ParseFilterMode(modeStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing filter mode: %v", err)
			}
			if modeAny != nil {
				mode := modeAny.(*models.FilterMode)
				filter.SetMode(mode)
			}
		}
		if !data.ApplicationFilter.Rule.IsNull() {
			rule := data.ApplicationFilter.Rule.ValueString()
			filter.SetRule(&rule)
		}
		applications.SetApplicationFilter(filter)
	}

	additionalData := make(map[string]interface{})
	applications.SetAdditionalData(additionalData)

	return applications, nil
}

// constructAuthenticationFlows constructs the ConditionalAccessAuthenticationFlowsable object
func constructAuthenticationFlows(data *ConditionalAccessAuthenticationFlowsModel) (models.ConditionalAccessAuthenticationFlowsable, error) {
	if data == nil {
		return nil, nil
	}

	authFlows := models.NewConditionalAccessAuthenticationFlows()

	if len(data.TransferMethods) > 0 {
		// Create a combined transfer methods value
		var combinedMethods models.ConditionalAccessTransferMethods
		for i, method := range data.TransferMethods {
			methodAny, err := models.ParseConditionalAccessTransferMethods(method.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing transfer methods: %v", err)
			}
			if methodAny != nil {
				if i == 0 {
					combinedMethods = *methodAny.(*models.ConditionalAccessTransferMethods)
				} else {
					combinedMethods |= *methodAny.(*models.ConditionalAccessTransferMethods)
				}
			}
		}
		authFlows.SetTransferMethods(&combinedMethods)
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

// constructGuestsOrExternalUsers constructs a ConditionalAccessGuestsOrExternalUsers object
// for the Microsoft Graph SDK based on the data provided from the Terraform model.
//
// This function processes the HCL input for `exclude_guests_or_external_users` or
// `include_guests_or_external_users` and maps the list of `guest_or_external_user_types`
// into a comma-separated string. The string is parsed using the Microsoft Graph SDK's
// `ParseConditionalAccessGuestOrExternalUserTypes` function, which combines the individual
// flags into a single bitmask value.
//
// Parameters:
//   - data: *ConditionalAccessGuestsOrExternalUsersModel
//     The Terraform model containing the input data for GuestOrExternalUserTypes
//     and ExternalTenants.
//
// Returns:
//   - models.ConditionalAccessGuestsOrExternalUsersable: The constructed object
//     that can be sent to the Microsoft Graph API.
//   - error: If an error occurs while parsing the guest or external user types.
//
// Behavior:
//   - If `guest_or_external_user_types` contains a list of strings (e.g., ["b2bCollaborationGuest", "b2bCollaborationMember"]),
//     they are combined into a single comma-separated string and parsed into the corresponding SDK enum bitmask.
//   - If `external_tenants` is provided, it is mapped using the helper function `constructConditionalAccessExternalTenants`.
//   - Handles cases where no data is provided for `guest_or_external_user_types` or `external_tenants` by safely returning nil.
//
// Example Input (HCL):
//
//	exclude_guests_or_external_users = {
//	  guest_or_external_user_types = ["b2bCollaborationGuest", "b2bCollaborationMember"]
//	  external_tenants = {
//	    membership_kind = "all"
//	  }
//	}
//
// Errors:
// - Returns an error if a value in `guest_or_external_user_types` cannot be parsed into a valid enum.
//
// Notes:
// - The Microsoft Graph SDK uses a bitmask-based enum for GuestOrExternalUserTypes.
// - This implementation ensures that Terraform's list input aligns correctly with the SDK's expected format.
func constructGuestsOrExternalUsers(data *ConditionalAccessGuestsOrExternalUsersModel) (models.ConditionalAccessGuestsOrExternalUsersable, error) {
	if data == nil {
		return nil, nil
	}

	guestsOrExternalUsers := models.NewConditionalAccessGuestsOrExternalUsers()

	if len(data.GuestOrExternalUserTypes) > 0 {
		var userTypesStrings []string
		for _, userType := range data.GuestOrExternalUserTypes {
			userTypesStrings = append(userTypesStrings, userType.ValueString())
		}

		combinedTypesString := strings.Join(userTypesStrings, ",")

		parsedAny, err := models.ParseConditionalAccessGuestOrExternalUserTypes(combinedTypesString)
		if err != nil {
			return nil, fmt.Errorf("error parsing guest or external user types: %v", err)
		}

		if parsedAny != nil {
			parsedType, ok := parsedAny.(*models.ConditionalAccessGuestOrExternalUserTypes)
			if !ok {
				return nil, fmt.Errorf("unexpected type for guest or external user types: %T", parsedAny)
			}
			guestsOrExternalUsers.SetGuestOrExternalUserTypes(parsedType)
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

	// if data.IncludeStates != nil {
	// 	if len(data.IncludeStates) > 0 {
	// 		includeStates := make([]string, len(data.IncludeStates))
	// 		for i, state := range data.IncludeStates {
	// 			includeStates[i] = state.ValueString()
	// 		}
	// 		devices.SetIncludeDeviceStates(includeStates)
	// 	}

	// 	if len(data.ExcludeStates) > 0 {
	// 		excludeStates := make([]string, len(data.ExcludeStates))
	// 		for i, state := range data.ExcludeStates {
	// 			excludeStates[i] = state.ValueString()
	// 		}
	// 		devices.SetExcludeDeviceStates(excludeStates)
	// 	}
	// }

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

// constructAuthenticationStrength constructs the AuthenticationStrengthPolicy object
func constructAuthenticationStrength(data *AuthenticationStrengthPolicyModel) (*models.AuthenticationStrengthPolicy, error) {
	if data == nil {
		return nil, nil
	}

	authStrength := models.NewAuthenticationStrengthPolicy()

	constructors.SetStringProperty(data.DisplayName, authStrength.SetDisplayName)
	constructors.SetStringProperty(data.Description, authStrength.SetDescription)

	if err := constructors.SetEnumProperty(data.PolicyType, models.ParseAuthenticationStrengthPolicyType, authStrength.SetPolicyType); err != nil {
		return nil, fmt.Errorf("error setting policy type: %v", err)
	}

	if err := constructors.SetEnumProperty(data.RequirementsSatisfied, models.ParseAuthenticationStrengthRequirements, authStrength.SetRequirementsSatisfied); err != nil {
		return nil, fmt.Errorf("error setting requirements satisfied: %v", err)
	}

	// Handle allowed combinations list
	if len(data.AllowedCombinations) > 0 {
		allowedCombinations := make([]models.AuthenticationMethodModes, 0, len(data.AllowedCombinations))
		for _, combination := range data.AllowedCombinations {
			var authMethodMode models.AuthenticationMethodModes
			if err := constructors.SetEnumProperty(
				combination,
				models.ParseAuthenticationMethodModes,
				func(mode models.AuthenticationMethodModes) {
					authMethodMode = mode
				},
			); err != nil {
				return nil, fmt.Errorf("error parsing allowed combination: %v", err)
			}
			allowedCombinations = append(allowedCombinations, authMethodMode)
		}
		authStrength.SetAllowedCombinations(allowedCombinations)
	}

	return authStrength, nil
}

// constructSessionControls constructs the ConditionalAccessSessionControlsable object
func constructSessionControls(data *ConditionalAccessSessionControlsModel) (models.ConditionalAccessSessionControlsable, error) {
	if data == nil {
		return nil, nil
	}

	sessionControls := models.NewConditionalAccessSessionControls()

	// Handle ApplicationEnforcedRestrictions
	if data.ApplicationEnforcedRestrictions != nil {
		appRestrictions := models.NewApplicationEnforcedRestrictionsSessionControl()
		constructors.SetBoolProperty(data.ApplicationEnforcedRestrictions.IsEnabled, appRestrictions.SetIsEnabled)
		sessionControls.SetApplicationEnforcedRestrictions(appRestrictions)
	}

	// Handle CloudAppSecurity
	if data.CloudAppSecurity != nil {
		cloudAppSecurity := models.NewCloudAppSecuritySessionControl()
		constructors.SetBoolProperty(data.CloudAppSecurity.IsEnabled, cloudAppSecurity.SetIsEnabled)

		if err := constructors.SetEnumProperty(
			data.CloudAppSecurity.CloudAppSecurityType,
			models.ParseCloudAppSecuritySessionControlType,
			cloudAppSecurity.SetCloudAppSecurityType,
		); err != nil {
			return nil, fmt.Errorf("error setting cloud app security type: %v", err)
		}

		sessionControls.SetCloudAppSecurity(cloudAppSecurity)
	}

	// Handle ContinuousAccessEvaluation
	if data.ContinuousAccessEvaluation != nil {
		continuousAccessEvaluation := models.NewContinuousAccessEvaluationSessionControl()

		if err := constructors.SetEnumProperty(
			data.ContinuousAccessEvaluation.Mode,
			models.ParseContinuousAccessEvaluationMode,
			continuousAccessEvaluation.SetMode,
		); err != nil {
			return nil, fmt.Errorf("error setting continuous access evaluation mode: %v", err)
		}

		sessionControls.SetContinuousAccessEvaluation(continuousAccessEvaluation)
	}

	// Handle PersistentBrowser
	if data.PersistentBrowser != nil {
		persistentBrowser := models.NewPersistentBrowserSessionControl()
		constructors.SetBoolProperty(data.PersistentBrowser.IsEnabled, persistentBrowser.SetIsEnabled)

		if err := constructors.SetEnumProperty(
			data.PersistentBrowser.Mode,
			models.ParsePersistentBrowserSessionMode,
			persistentBrowser.SetMode,
		); err != nil {
			return nil, fmt.Errorf("error setting persistent browser mode: %v", err)
		}

		sessionControls.SetPersistentBrowser(persistentBrowser)
	}

	// Handle SignInFrequency
	if data.SignInFrequency != nil {
		signInFrequency := models.NewSignInFrequencySessionControl()
		constructors.SetBoolProperty(data.SignInFrequency.IsEnabled, signInFrequency.SetIsEnabled)

		if err := constructors.SetEnumProperty(
			data.SignInFrequency.Type,
			models.ParseSigninFrequencyType,
			signInFrequency.SetTypeEscaped,
		); err != nil {
			return nil, fmt.Errorf("error setting sign-in frequency type: %v", err)
		}

		if !data.SignInFrequency.Value.IsNull() {
			value := data.SignInFrequency.Value.ValueInt64()
			if value > math.MaxInt32 || value < math.MinInt32 {
				return nil, fmt.Errorf("sign-in frequency value %d is out of range for int32", value)
			}
			int32Value := int32(value)
			signInFrequency.SetValue(&int32Value)
		}

		if err := constructors.SetEnumProperty(
			data.SignInFrequency.FrequencyInterval,
			models.ParseSignInFrequencyInterval,
			signInFrequency.SetFrequencyInterval,
		); err != nil {
			return nil, fmt.Errorf("error setting sign-in frequency interval: %v", err)
		}

		if err := constructors.SetEnumProperty(
			data.SignInFrequency.AuthenticationType,
			models.ParseSignInFrequencyAuthenticationType,
			signInFrequency.SetAuthenticationType,
		); err != nil {
			return nil, fmt.Errorf("error setting sign-in frequency authentication type: %v", err)
		}

		sessionControls.SetSignInFrequency(signInFrequency)
	}

	// Handle SecureSignInSession
	if data.SecureSignInSession != nil {
		secureSignInSession := models.NewSecureSignInSessionControl()
		constructors.SetBoolProperty(data.SecureSignInSession.IsEnabled, secureSignInSession.SetIsEnabled)
		sessionControls.SetSecureSignInSession(secureSignInSession)
	}

	// Handle DisableResilienceDefaults
	constructors.SetBoolProperty(data.DisableResilienceDefaults, sessionControls.SetDisableResilienceDefaults)

	return sessionControls, nil
}
