package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns a ConditionalAccessPolicy that can be serialized by Kiota
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *ConditionalAccessPolicyResourceModel) (graphmodels.ConditionalAccessPolicyable, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewConditionalAccessPolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if err := convert.FrameworkToGraphEnum(data.State,
		graphmodels.ParseConditionalAccessPolicyState, requestBody.SetState); err != nil {
		return nil, fmt.Errorf("failed to set state: %w", err)
	}

	// Build conditions
	if data.Conditions != nil {
		conditions, err := constructConditions(ctx, data.Conditions)
		if err != nil {
			return nil, fmt.Errorf("failed to construct conditions: %w", err)
		}
		requestBody.SetConditions(conditions)
	}

	// Build grant controls
	if data.GrantControls != nil {
		grantControls, err := constructGrantControls(ctx, data.GrantControls)
		if err != nil {
			return nil, fmt.Errorf("failed to construct grant controls: %w", err)
		}
		requestBody.SetGrantControls(grantControls)
	}

	// Build session controls
	if data.SessionControls != nil {
		sessionControls, err := constructSessionControls(ctx, data.SessionControls)
		if err != nil {
			return nil, fmt.Errorf("failed to construct session controls: %w", err)
		}
		requestBody.SetSessionControls(sessionControls)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructConditions builds the conditions object using SDK models
func constructConditions(ctx context.Context, data *ConditionalAccessConditions) (graphmodels.ConditionalAccessConditionSetable, error) {
	conditions := graphmodels.NewConditionalAccessConditionSet()

	// Client App Types
	if err := convert.FrameworkToGraphEnumCollection(ctx, data.ClientAppTypes,
		graphmodels.ParseConditionalAccessClientApp, conditions.SetClientAppTypes); err != nil {
		return nil, fmt.Errorf("failed to set client app types: %w", err)
	}

	// Sign-in Risk Levels
	if err := convert.FrameworkToGraphEnumCollection(ctx, data.SignInRiskLevels,
		graphmodels.ParseRiskLevel, conditions.SetSignInRiskLevels); err != nil {
		return nil, fmt.Errorf("failed to set sign in risk levels: %w", err)
	}

	// User Risk Levels
	if err := convert.FrameworkToGraphEnumCollection(ctx, data.UserRiskLevels,
		graphmodels.ParseRiskLevel, conditions.SetUserRiskLevels); err != nil {
		return nil, fmt.Errorf("failed to set user risk levels: %w", err)
	}

	// Service Principal Risk Levels
	if err := convert.FrameworkToGraphEnumCollection(ctx, data.ServicePrincipalRiskLevels,
		graphmodels.ParseRiskLevel, conditions.SetServicePrincipalRiskLevels); err != nil {
		return nil, fmt.Errorf("failed to set service principal risk levels: %w", err)
	}

	// Agent ID Risk Levels - bitmask enum
	if err := convert.FrameworkToGraphBitmaskEnumFromSet(
		ctx,
		data.AgentIdRiskLevels,
		graphmodels.ParseConditionalAccessAgentIdRiskLevels,
		conditions.SetAgentIdRiskLevels,
	); err != nil {
		return nil, fmt.Errorf("failed to convert agent id risk levels: %w", err)
	}

	// Insider Risk Levels - bitmask enum
	if err := convert.FrameworkToGraphBitmaskEnumFromSet(
		ctx,
		data.InsiderRiskLevels,
		graphmodels.ParseConditionalAccessInsiderRiskLevels,
		conditions.SetInsiderRiskLevels,
	); err != nil {
		return nil, fmt.Errorf("failed to convert insider risk levels: %w", err)
	}

	// Applications
	if data.Applications != nil {
		applications, err := constructApplications(ctx, data.Applications)
		if err != nil {
			return nil, fmt.Errorf("failed to construct applications: %w", err)
		}
		conditions.SetApplications(applications)
	}

	// Users
	if data.Users != nil {
		users, err := constructUsers(ctx, data.Users)
		if err != nil {
			return nil, fmt.Errorf("failed to construct users: %w", err)
		}
		conditions.SetUsers(users)
	}

	// Platforms
	if data.Platforms != nil {
		platforms, err := constructPlatforms(ctx, data.Platforms)
		if err != nil {
			return nil, fmt.Errorf("failed to construct platforms: %w", err)
		}
		conditions.SetPlatforms(platforms)
	}

	// Locations
	if data.Locations != nil {
		locations, err := constructLocations(ctx, data.Locations)
		if err != nil {
			return nil, fmt.Errorf("failed to construct locations: %w", err)
		}
		conditions.SetLocations(locations)
	}

	// Devices
	if data.Devices != nil {
		devices, err := constructDevices(ctx, data.Devices)
		if err != nil {
			return nil, fmt.Errorf("failed to construct devices: %w", err)
		}
		conditions.SetDevices(devices)
	}

	// Client Applications
	if data.ClientApplications != nil {
		clientApplications, err := constructClientApplications(ctx, data.ClientApplications)
		if err != nil {
			return nil, fmt.Errorf("failed to construct client applications: %w", err)
		}
		conditions.SetClientApplications(clientApplications)
	}

	// Times - not available in SDK, would need to use AdditionalData if required
	// Skipping for now as it's not in the SDK model

	// Device States
	if data.DeviceStates != nil {
		deviceStates, err := constructDeviceStates(ctx, data.DeviceStates)
		if err != nil {
			return nil, fmt.Errorf("failed to construct device states: %w", err)
		}
		conditions.SetDeviceStates(deviceStates)
	}

	// Authentication Flows
	if data.AuthenticationFlows != nil {
		authFlows, err := constructAuthenticationFlows(ctx, data.AuthenticationFlows)
		if err != nil {
			return nil, fmt.Errorf("failed to construct authentication flows: %w", err)
		}
		conditions.SetAuthenticationFlows(authFlows)
	}

	return conditions, nil
}

// constructApplications builds the applications object using SDK models
func constructApplications(ctx context.Context, data *ConditionalAccessApplications) (graphmodels.ConditionalAccessApplicationsable, error) {
	applications := graphmodels.NewConditionalAccessApplications()

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeApplications, applications.SetIncludeApplications); err != nil {
		return nil, fmt.Errorf("failed to convert include applications: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeApplications, applications.SetExcludeApplications); err != nil {
		return nil, fmt.Errorf("failed to convert exclude applications: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeUserActions, applications.SetIncludeUserActions); err != nil {
		return nil, fmt.Errorf("failed to convert include user actions: %w", err)
	}

	// Handle authentication context class references with predefined value mapping
	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeAuthenticationContextClassReferences, func(values []string) {
		mappedValues := make([]string, len(values))
		for i, value := range values {
			switch value {
			case "require_trusted_device":
				mappedValues[i] = "c1"
			case "require_terms_of_use":
				mappedValues[i] = "c2"
			case "require_trusted_location":
				mappedValues[i] = "c3"
			case "require_strong_authentication":
				mappedValues[i] = "c4"
			case "required_trust_type:azure_ad_joined":
				mappedValues[i] = "c5"
			case "require_access_from_an_approved_app":
				mappedValues[i] = "c6"
			case "required_trust_type:hybrid_azure_ad_joined":
				mappedValues[i] = "c7"
			default:
				mappedValues[i] = value
			}
		}
		applications.SetIncludeAuthenticationContextClassReferences(mappedValues)
	}); err != nil {
		return nil, fmt.Errorf("failed to convert include auth context class refs: %w", err)
	}

	// Application Filter
	if data.ApplicationFilter != nil {
		filter := graphmodels.NewConditionalAccessFilter()
		if err := convert.FrameworkToGraphEnum(data.ApplicationFilter.Mode,
			graphmodels.ParseFilterMode, filter.SetMode); err != nil {
			return nil, fmt.Errorf("failed to set application filter mode: %w", err)
		}
		convert.FrameworkToGraphString(data.ApplicationFilter.Rule, filter.SetRule)
		applications.SetApplicationFilter(filter)
	}

	return applications, nil
}

// constructUsers builds the users object using SDK models
func constructUsers(ctx context.Context, data *ConditionalAccessUsers) (graphmodels.ConditionalAccessUsersable, error) {
	users := graphmodels.NewConditionalAccessUsers()

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeUsers, users.SetIncludeUsers); err != nil {
		return nil, fmt.Errorf("failed to convert include users: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeUsers, users.SetExcludeUsers); err != nil {
		return nil, fmt.Errorf("failed to convert exclude users: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeGroups, users.SetIncludeGroups); err != nil {
		return nil, fmt.Errorf("failed to convert include groups: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeGroups, users.SetExcludeGroups); err != nil {
		return nil, fmt.Errorf("failed to convert exclude groups: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeRoles, users.SetIncludeRoles); err != nil {
		return nil, fmt.Errorf("failed to convert include roles: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeRoles, users.SetExcludeRoles); err != nil {
		return nil, fmt.Errorf("failed to convert exclude roles: %w", err)
	}

	// Include Guests or External Users
	if !data.IncludeGuestsOrExternalUsers.IsNull() && !data.IncludeGuestsOrExternalUsers.IsUnknown() {
		includeGuests, err := constructGuestsOrExternalUsers(ctx, data.IncludeGuestsOrExternalUsers)
		if err != nil {
			return nil, fmt.Errorf("failed to construct include guests or external users: %w", err)
		}
		users.SetIncludeGuestsOrExternalUsers(includeGuests)
	}

	// Exclude Guests or External Users
	if !data.ExcludeGuestsOrExternalUsers.IsNull() && !data.ExcludeGuestsOrExternalUsers.IsUnknown() {
		excludeGuests, err := constructGuestsOrExternalUsers(ctx, data.ExcludeGuestsOrExternalUsers)
		if err != nil {
			return nil, fmt.Errorf("failed to construct exclude guests or external users: %w", err)
		}
		users.SetExcludeGuestsOrExternalUsers(excludeGuests)
	}

	return users, nil
}

// constructGuestsOrExternalUsers builds the guests or external users object
func constructGuestsOrExternalUsers(ctx context.Context, obj types.Object) (graphmodels.ConditionalAccessGuestsOrExternalUsersable, error) {
	guestsOrExternalUsers := graphmodels.NewConditionalAccessGuestsOrExternalUsers()

	attrs := obj.Attributes()

	// Handle guest_or_external_user_types - bitmask enum
	if guestTypesAttr, ok := attrs["guest_or_external_user_types"]; ok {
		if guestTypesSet, ok := guestTypesAttr.(types.Set); ok {
			if err := convert.FrameworkToGraphBitmaskEnumFromSet(
				ctx,
				guestTypesSet,
				graphmodels.ParseConditionalAccessGuestOrExternalUserTypes,
				guestsOrExternalUsers.SetGuestOrExternalUserTypes,
			); err != nil {
				return nil, fmt.Errorf("failed to convert guest or external user types: %w", err)
			}
		}
	}

	// Handle external_tenants
	if externalTenantsAttr, ok := attrs["external_tenants"]; ok {
		if externalTenantsObj, ok := externalTenantsAttr.(types.Object); ok && !externalTenantsObj.IsNull() {
			externalTenants, err := constructExternalTenants(ctx, externalTenantsObj)
			if err != nil {
				return nil, fmt.Errorf("failed to construct external tenants: %w", err)
			}
			guestsOrExternalUsers.SetExternalTenants(externalTenants)
		}
	}

	return guestsOrExternalUsers, nil
}

// constructExternalTenants builds the external tenants object
func constructExternalTenants(ctx context.Context, obj types.Object) (graphmodels.ConditionalAccessExternalTenantsable, error) {
	attrs := obj.Attributes()

	var membershipKind string
	if membershipKindAttr, ok := attrs["membership_kind"]; ok {
		if membershipKindStr, ok := membershipKindAttr.(types.String); ok {
			membershipKind = membershipKindStr.ValueString()
		}
	}

	// Create the appropriate type based on membership kind
	switch membershipKind {
	case "all":
		externalTenants := graphmodels.NewConditionalAccessAllExternalTenants()
		return externalTenants, nil

	case "enumerated":
		externalTenants := graphmodels.NewConditionalAccessEnumeratedExternalTenants()

		// Handle members
		if membersAttr, ok := attrs["members"]; ok {
			if membersSet, ok := membersAttr.(types.Set); ok {
				if err := convert.FrameworkToGraphStringSet(ctx, membersSet, externalTenants.SetMembers); err != nil {
					return nil, fmt.Errorf("failed to convert external tenants members: %w", err)
				}
			}
		}

		return externalTenants, nil

	default:
		return nil, fmt.Errorf("invalid membership_kind: %s", membershipKind)
	}
}

// constructPlatforms builds the platforms object using SDK models
func constructPlatforms(ctx context.Context, data *ConditionalAccessPlatforms) (graphmodels.ConditionalAccessPlatformsable, error) {
	platforms := graphmodels.NewConditionalAccessPlatforms()

	if err := convert.FrameworkToGraphEnumCollection(ctx, data.IncludePlatforms,
		graphmodels.ParseConditionalAccessDevicePlatform, platforms.SetIncludePlatforms); err != nil {
		return nil, fmt.Errorf("failed to set include platforms: %w", err)
	}

	if err := convert.FrameworkToGraphEnumCollection(ctx, data.ExcludePlatforms,
		graphmodels.ParseConditionalAccessDevicePlatform, platforms.SetExcludePlatforms); err != nil {
		return nil, fmt.Errorf("failed to set exclude platforms: %w", err)
	}

	return platforms, nil
}

// constructLocations builds the locations object using SDK models
func constructLocations(ctx context.Context, data *ConditionalAccessLocations) (graphmodels.ConditionalAccessLocationsable, error) {
	locations := graphmodels.NewConditionalAccessLocations()

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeLocations, locations.SetIncludeLocations); err != nil {
		return nil, fmt.Errorf("failed to convert include locations: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeLocations, locations.SetExcludeLocations); err != nil {
		return nil, fmt.Errorf("failed to convert exclude locations: %w", err)
	}

	return locations, nil
}

// constructDevices builds the devices object using SDK models
func constructDevices(ctx context.Context, data *ConditionalAccessDevices) (graphmodels.ConditionalAccessDevicesable, error) {
	devices := graphmodels.NewConditionalAccessDevices()

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeDevices, devices.SetIncludeDevices); err != nil {
		return nil, fmt.Errorf("failed to convert include devices: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeDevices, devices.SetExcludeDevices); err != nil {
		return nil, fmt.Errorf("failed to convert exclude devices: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeDeviceStates, devices.SetIncludeDeviceStates); err != nil {
		return nil, fmt.Errorf("failed to convert include device states: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeDeviceStates, devices.SetExcludeDeviceStates); err != nil {
		return nil, fmt.Errorf("failed to convert exclude device states: %w", err)
	}

	// Device Filter
	if data.DeviceFilter != nil {
		filter := graphmodels.NewConditionalAccessFilter()
		if err := convert.FrameworkToGraphEnum(data.DeviceFilter.Mode,
			graphmodels.ParseFilterMode, filter.SetMode); err != nil {
			return nil, fmt.Errorf("failed to set device filter mode: %w", err)
		}
		convert.FrameworkToGraphString(data.DeviceFilter.Rule, filter.SetRule)
		devices.SetDeviceFilter(filter)
	}

	return devices, nil
}

// constructClientApplications builds the client applications object using SDK models
func constructClientApplications(ctx context.Context, data *ConditionalAccessClientApplications) (graphmodels.ConditionalAccessClientApplicationsable, error) {
	clientApplications := graphmodels.NewConditionalAccessClientApplications()

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeServicePrincipals, clientApplications.SetIncludeServicePrincipals); err != nil {
		return nil, fmt.Errorf("failed to convert include service principals: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeServicePrincipals, clientApplications.SetExcludeServicePrincipals); err != nil {
		return nil, fmt.Errorf("failed to convert exclude service principals: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeAgentIdServicePrincipals, clientApplications.SetIncludeAgentIdServicePrincipals); err != nil {
		return nil, fmt.Errorf("failed to convert include agent id service principals: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeAgentIdServicePrincipals, clientApplications.SetExcludeAgentIdServicePrincipals); err != nil {
		return nil, fmt.Errorf("failed to convert exclude agent id service principals: %w", err)
	}

	if data.AgentIdServicePrincipalFilter != nil {
		filter := graphmodels.NewConditionalAccessFilter()
		if err := convert.FrameworkToGraphEnum(data.AgentIdServicePrincipalFilter.Mode,
			graphmodels.ParseFilterMode, filter.SetMode); err != nil {
			return nil, fmt.Errorf("failed to set agent id service principal filter mode: %w", err)
		}
		convert.FrameworkToGraphString(data.AgentIdServicePrincipalFilter.Rule, filter.SetRule)
		clientApplications.SetAgentIdServicePrincipalFilter(filter)
	}

	if data.ServicePrincipalFilter != nil {
		filter := graphmodels.NewConditionalAccessFilter()
		if err := convert.FrameworkToGraphEnum(data.ServicePrincipalFilter.Mode,
			graphmodels.ParseFilterMode, filter.SetMode); err != nil {
			return nil, fmt.Errorf("failed to set service principal filter mode: %w", err)
		}
		convert.FrameworkToGraphString(data.ServicePrincipalFilter.Rule, filter.SetRule)
		clientApplications.SetServicePrincipalFilter(filter)
	}

	return clientApplications, nil
}

// constructTimes - Times is not available in the SDK, skipping
// If needed, this could be implemented via AdditionalData

// constructDeviceStates builds the device states object using SDK models
func constructDeviceStates(ctx context.Context, data *ConditionalAccessDeviceStates) (graphmodels.ConditionalAccessDeviceStatesable, error) {
	deviceStates := graphmodels.NewConditionalAccessDeviceStates()

	if err := convert.FrameworkToGraphStringSet(ctx, data.IncludeStates, deviceStates.SetIncludeStates); err != nil {
		return nil, fmt.Errorf("failed to convert include states: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ExcludeStates, deviceStates.SetExcludeStates); err != nil {
		return nil, fmt.Errorf("failed to convert exclude states: %w", err)
	}

	return deviceStates, nil
}

// constructAuthenticationFlows builds the authentication flows object using SDK models
func constructAuthenticationFlows(ctx context.Context, data *ConditionalAccessAuthenticationFlows) (graphmodels.ConditionalAccessAuthenticationFlowsable, error) {
	authFlows := graphmodels.NewConditionalAccessAuthenticationFlows()

	if err := convert.FrameworkToGraphEnum(data.TransferMethods,
		graphmodels.ParseConditionalAccessTransferMethods, authFlows.SetTransferMethods); err != nil {
		return nil, fmt.Errorf("failed to set transfer methods: %w", err)
	}

	return authFlows, nil
}

// constructGrantControls builds the grant controls object using SDK models
func constructGrantControls(ctx context.Context, data *ConditionalAccessGrantControls) (graphmodels.ConditionalAccessGrantControlsable, error) {
	grantControls := graphmodels.NewConditionalAccessGrantControls()

	convert.FrameworkToGraphString(data.Operator, grantControls.SetOperator)

	if err := convert.FrameworkToGraphEnumCollection(ctx, data.BuiltInControls,
		graphmodels.ParseConditionalAccessGrantControl, grantControls.SetBuiltInControls); err != nil {
		return nil, fmt.Errorf("failed to set built in controls: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.CustomAuthenticationFactors, grantControls.SetCustomAuthenticationFactors); err != nil {
		return nil, fmt.Errorf("failed to convert custom authentication factors: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.TermsOfUse, grantControls.SetTermsOfUse); err != nil {
		return nil, fmt.Errorf("failed to convert terms of use: %w", err)
	}

	// Authentication Strength
	if data.AuthenticationStrength != nil {
		authStrength := graphmodels.NewAuthenticationStrengthPolicy()
		convert.FrameworkToGraphString(data.AuthenticationStrength.ID, authStrength.SetId)
		grantControls.SetAuthenticationStrength(authStrength)
	}

	return grantControls, nil
}

// constructSessionControls builds the session controls object using SDK models
func constructSessionControls(ctx context.Context, data *ConditionalAccessSessionControls) (graphmodels.ConditionalAccessSessionControlsable, error) {
	sessionControls := graphmodels.NewConditionalAccessSessionControls()

	convert.FrameworkToGraphBool(data.DisableResilienceDefaults, sessionControls.SetDisableResilienceDefaults)

	// Application Enforced Restrictions
	if data.ApplicationEnforcedRestrictions != nil {
		appEnforcedRestrictions := graphmodels.NewApplicationEnforcedRestrictionsSessionControl()
		convert.FrameworkToGraphBool(data.ApplicationEnforcedRestrictions.IsEnabled, appEnforcedRestrictions.SetIsEnabled)
		sessionControls.SetApplicationEnforcedRestrictions(appEnforcedRestrictions)
	}

	// Cloud App Security
	if data.CloudAppSecurity != nil {
		cloudAppSecurity := graphmodels.NewCloudAppSecuritySessionControl()
		convert.FrameworkToGraphBool(data.CloudAppSecurity.IsEnabled, cloudAppSecurity.SetIsEnabled)

		if err := convert.FrameworkToGraphEnum(data.CloudAppSecurity.CloudAppSecurityType,
			graphmodels.ParseCloudAppSecuritySessionControlType, cloudAppSecurity.SetCloudAppSecurityType); err != nil {
			return nil, fmt.Errorf("failed to set cloud app security type: %w", err)
		}

		sessionControls.SetCloudAppSecurity(cloudAppSecurity)
	}

	// Sign In Frequency
	if data.SignInFrequency != nil {
		signInFrequency := graphmodels.NewSignInFrequencySessionControl()
		convert.FrameworkToGraphBool(data.SignInFrequency.IsEnabled, signInFrequency.SetIsEnabled)

		if err := convert.FrameworkToGraphEnum(data.SignInFrequency.Type,
			graphmodels.ParseSigninFrequencyType, signInFrequency.SetTypeEscaped); err != nil {
			return nil, fmt.Errorf("failed to set sign in frequency type: %w", err)
		}

		convert.FrameworkToGraphInt32(data.SignInFrequency.Value, signInFrequency.SetValue)

		if err := convert.FrameworkToGraphEnum(data.SignInFrequency.AuthenticationType,
			graphmodels.ParseSignInFrequencyAuthenticationType, signInFrequency.SetAuthenticationType); err != nil {
			return nil, fmt.Errorf("failed to set authentication type: %w", err)
		}

		if err := convert.FrameworkToGraphEnum(data.SignInFrequency.FrequencyInterval,
			graphmodels.ParseSignInFrequencyInterval, signInFrequency.SetFrequencyInterval); err != nil {
			return nil, fmt.Errorf("failed to set frequency interval: %w", err)
		}

		sessionControls.SetSignInFrequency(signInFrequency)
	}

	// Persistent Browser
	if data.PersistentBrowser != nil {
		persistentBrowser := graphmodels.NewPersistentBrowserSessionControl()
		convert.FrameworkToGraphBool(data.PersistentBrowser.IsEnabled, persistentBrowser.SetIsEnabled)

		if err := convert.FrameworkToGraphEnum(data.PersistentBrowser.Mode,
			graphmodels.ParsePersistentBrowserSessionMode, persistentBrowser.SetMode); err != nil {
			return nil, fmt.Errorf("failed to set persistent browser mode: %w", err)
		}

		sessionControls.SetPersistentBrowser(persistentBrowser)
	}

	// Continuous Access Evaluation
	if data.ContinuousAccessEvaluation != nil {
		cae := graphmodels.NewContinuousAccessEvaluationSessionControl()

		if err := convert.FrameworkToGraphEnum(data.ContinuousAccessEvaluation.Mode,
			graphmodels.ParseContinuousAccessEvaluationMode, cae.SetMode); err != nil {
			return nil, fmt.Errorf("failed to set continuous access evaluation mode: %w", err)
		}

		sessionControls.SetContinuousAccessEvaluation(cae)
	}

	// Secure Sign In Session
	if data.SecureSignInSession != nil {
		secureSignInSession := graphmodels.NewSecureSignInSessionControl()
		convert.FrameworkToGraphBool(data.SecureSignInSession.IsEnabled, secureSignInSession.SetIsEnabled)
		sessionControls.SetSecureSignInSession(secureSignInSession)
	}

	// Global Secure Access Filtering Profile
	if data.GlobalSecureAccessFilteringProfile != nil {
		gsaProfile := graphmodels.NewGlobalSecureAccessFilteringProfileSessionControl()
		convert.FrameworkToGraphBool(data.GlobalSecureAccessFilteringProfile.IsEnabled, gsaProfile.SetIsEnabled)
		convert.FrameworkToGraphString(data.GlobalSecureAccessFilteringProfile.ProfileId, gsaProfile.SetProfileId)
		sessionControls.SetGlobalSecureAccessFilteringProfile(gsaProfile)
	}

	return sessionControls, nil
}
