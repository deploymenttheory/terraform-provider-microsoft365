package graphBetaConditionalAccessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConditionalAccessRootResourceModel represents the main resource model for ConditionalAccessRoot
// This is a management container that sits above individual policies and templates
type ConditionalAccessRootResourceModel struct {
	Id                                   types.String                                        `tfsdk:"id"`
	AuthenticationContextClassReferences []*AuthenticationContextClassReferenceResourceModel `tfsdk:"authentication_context_class_references"`
	AuthenticationStrength               *AuthenticationStrengthRootResourceModel            `tfsdk:"authentication_strength"` // Current property

}

// AuthenticationContextClassReferenceResourceModel represents an authentication context class reference
type AuthenticationContextClassReferenceResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	IsAvailable types.Bool   `tfsdk:"is_available"`
}

// AuthenticationStrengthRootResourceModel represents the authentication strength configuration
type AuthenticationStrengthRootResourceModel struct {
	Id                         types.String                                   `tfsdk:"id"`
	AuthenticationCombinations types.Set                                      `tfsdk:"authentication_combinations"` // Set of strings representing AuthenticationMethodModes
	Combinations               types.Set                                      `tfsdk:"combinations"`                // Set of strings representing AuthenticationMethodModes
	AuthenticationMethodModes  []*AuthenticationMethodModeDetailResourceModel `tfsdk:"authentication_method_modes"`
	Policies                   []*AuthenticationStrengthPolicyResourceModel   `tfsdk:"policies"`
}

// AuthenticationMethodModeDetailResourceModel represents details of an authentication method mode
type AuthenticationMethodModeDetailResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	DisplayName          types.String `tfsdk:"display_name"`
	AuthenticationMethod types.String `tfsdk:"authentication_method"` // String representation of BaseAuthenticationMethod enum
}

// AuthenticationStrengthPolicyResourceModel represents an authentication strength policy
type AuthenticationStrengthPolicyResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	DisplayName               types.String `tfsdk:"display_name"`
	Description               types.String `tfsdk:"description"`
	PolicyType                types.String `tfsdk:"policy_type"`                // builtIn, custom, unknownFutureValue
	RequirementsSatisfied     types.String `tfsdk:"requirements_satisfied"`     // singleFactor, multiFactor, unknownFutureValue
	AllowedCombinations       types.Set    `tfsdk:"allowed_combinations"`       // Set of strings representing AuthenticationMethodModes
	CombinationConfigurations types.Set    `tfsdk:"combination_configurations"` // Set of combination configuration objects
	CreatedDateTime           types.String `tfsdk:"created_date_time"`
	ModifiedDateTime          types.String `tfsdk:"modified_date_time"`
}

// ConditionalAccessTemplateResourceModel represents a conditional access template
type ConditionalAccessTemplateResourceModel struct {
	Id          types.String                                `tfsdk:"id"`
	Name        types.String                                `tfsdk:"name"`
	Description types.String                                `tfsdk:"description"`
	Scenarios   types.String                                `tfsdk:"scenarios"` // String representation of TemplateScenarios enum
	Details     *ConditionalAccessPolicyDetailResourceModel `tfsdk:"details"`
}

// ConditionalAccessPolicyDetailResourceModel represents the details of a conditional access policy template
type ConditionalAccessPolicyDetailResourceModel struct {
	Id               types.String                                   `tfsdk:"id"`
	DisplayName      types.String                                   `tfsdk:"display_name"`
	Description      types.String                                   `tfsdk:"description"`
	State            types.String                                   `tfsdk:"state"` // enabled, disabled, enabledForReportingButNotEnforced
	Conditions       *ConditionalAccessConditionSetResourceModel    `tfsdk:"conditions"`
	GrantControls    *ConditionalAccessGrantControlsResourceModel   `tfsdk:"grant_controls"`
	SessionControls  *ConditionalAccessSessionControlsResourceModel `tfsdk:"session_controls"`
	CreatedDateTime  types.String                                   `tfsdk:"created_date_time"`
	ModifiedDateTime types.String                                   `tfsdk:"modified_date_time"`
}

// ConditionalAccessConditionSetResourceModel represents conditions for conditional access
type ConditionalAccessConditionSetResourceModel struct {
	ClientAppTypes             types.Set                                          `tfsdk:"client_app_types"`              // Set of ConditionalAccessClientApp enums
	SignInRiskLevels           types.Set                                          `tfsdk:"sign_in_risk_levels"`           // Set of RiskLevel enums
	UserRiskLevels             types.Set                                          `tfsdk:"user_risk_levels"`              // Set of RiskLevel enums
	ServicePrincipalRiskLevels types.Set                                          `tfsdk:"service_principal_risk_levels"` // Set of RiskLevel enums
	InsiderRiskLevels          types.Set                                          `tfsdk:"insider_risk_levels"`           // Set of ConditionalAccessInsiderRiskLevels enums
	Applications               *ConditionalAccessApplicationsResourceModel        `tfsdk:"applications"`
	AuthenticationFlows        *ConditionalAccessAuthenticationFlowsResourceModel `tfsdk:"authentication_flows"`
	Users                      *ConditionalAccessUsersResourceModel               `tfsdk:"users"`
	ClientApplications         *ConditionalAccessClientApplicationsResourceModel  `tfsdk:"client_applications"`
	DeviceStates               *ConditionalAccessDeviceStatesResourceModel        `tfsdk:"device_states"` // Deprecated
	Devices                    *ConditionalAccessDevicesResourceModel             `tfsdk:"devices"`
	Locations                  *ConditionalAccessLocationsResourceModel           `tfsdk:"locations"`
	Platforms                  *ConditionalAccessPlatformsResourceModel           `tfsdk:"platforms"`
}

// ConditionalAccessApplicationsResourceModel represents application conditions
type ConditionalAccessApplicationsResourceModel struct {
	IncludeApplications                         types.Set                             `tfsdk:"include_applications"`
	ExcludeApplications                         types.Set                             `tfsdk:"exclude_applications"`
	IncludeUserActions                          types.Set                             `tfsdk:"include_user_actions"`
	IncludeAuthenticationContextClassReferences types.Set                             `tfsdk:"include_authentication_context_class_references"`
	ApplicationFilter                           *ConditionalAccessFilterResourceModel `tfsdk:"application_filter"`
}

// ConditionalAccessAuthenticationFlowsResourceModel represents authentication flow conditions
type ConditionalAccessAuthenticationFlowsResourceModel struct {
	TransferMethods types.Set `tfsdk:"transfer_methods"` // Set of ConditionalAccessTransferMethods enums
}

// ConditionalAccessUsersResourceModel represents user conditions
type ConditionalAccessUsersResourceModel struct {
	IncludeUsers                 types.Set                                            `tfsdk:"include_users"`
	ExcludeUsers                 types.Set                                            `tfsdk:"exclude_users"`
	IncludeGroups                types.Set                                            `tfsdk:"include_groups"`
	ExcludeGroups                types.Set                                            `tfsdk:"exclude_groups"`
	IncludeRoles                 types.Set                                            `tfsdk:"include_roles"`
	ExcludeRoles                 types.Set                                            `tfsdk:"exclude_roles"`
	IncludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsersResourceModel `tfsdk:"include_guests_or_external_users"`
	ExcludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsersResourceModel `tfsdk:"exclude_guests_or_external_users"`
}

// ConditionalAccessGuestsOrExternalUsersResourceModel represents guest/external user conditions
type ConditionalAccessGuestsOrExternalUsersResourceModel struct {
	GuestOrExternalUserTypes types.Set                                      `tfsdk:"guest_or_external_user_types"` // Set of ConditionalAccessGuestOrExternalUserTypes enums
	ExternalTenants          *ConditionalAccessExternalTenantsResourceModel `tfsdk:"external_tenants"`
}

// ConditionalAccessExternalTenantsResourceModel represents external tenant conditions
type ConditionalAccessExternalTenantsResourceModel struct {
	MembershipKind types.String `tfsdk:"membership_kind"` // ConditionalAccessExternalTenantsMembershipKind enum
	Members        types.Set    `tfsdk:"members"`         // Set of tenant IDs
}

// ConditionalAccessClientApplicationsResourceModel represents client application conditions
type ConditionalAccessClientApplicationsResourceModel struct {
	IncludeServicePrincipals types.Set                             `tfsdk:"include_service_principals"`
	ExcludeServicePrincipals types.Set                             `tfsdk:"exclude_service_principals"`
	ServicePrincipalFilter   *ConditionalAccessFilterResourceModel `tfsdk:"service_principal_filter"`
}

// ConditionalAccessDeviceStatesResourceModel represents device state conditions (deprecated)
type ConditionalAccessDeviceStatesResourceModel struct {
	IncludeStates types.Set `tfsdk:"include_states"`
	ExcludeStates types.Set `tfsdk:"exclude_states"`
}

// ConditionalAccessDevicesResourceModel represents device conditions
type ConditionalAccessDevicesResourceModel struct {
	IncludeDevices      types.Set                             `tfsdk:"include_devices"`
	ExcludeDevices      types.Set                             `tfsdk:"exclude_devices"`
	IncludeDeviceStates types.Set                             `tfsdk:"include_device_states"`
	ExcludeDeviceStates types.Set                             `tfsdk:"exclude_device_states"`
	DeviceFilter        *ConditionalAccessFilterResourceModel `tfsdk:"device_filter"`
}

// ConditionalAccessLocationsResourceModel represents location conditions
type ConditionalAccessLocationsResourceModel struct {
	IncludeLocations types.Set `tfsdk:"include_locations"`
	ExcludeLocations types.Set `tfsdk:"exclude_locations"`
}

// ConditionalAccessPlatformsResourceModel represents platform conditions
type ConditionalAccessPlatformsResourceModel struct {
	IncludePlatforms types.Set `tfsdk:"include_platforms"` // Set of ConditionalAccessDevicePlatform enums
	ExcludePlatforms types.Set `tfsdk:"exclude_platforms"` // Set of ConditionalAccessDevicePlatform enums
}

// ConditionalAccessFilterResourceModel represents a filter condition
type ConditionalAccessFilterResourceModel struct {
	Mode types.String `tfsdk:"mode"` // FilterMode enum: include, exclude
	Rule types.String `tfsdk:"rule"` // Filter rule expression
}

// ConditionalAccessGrantControlsResourceModel represents grant controls
type ConditionalAccessGrantControlsResourceModel struct {
	Operator                    types.String `tfsdk:"operator"`                      // AND, OR
	BuiltInControls             types.Set    `tfsdk:"built_in_controls"`             // Set of ConditionalAccessGrantControl enums
	CustomAuthenticationFactors types.Set    `tfsdk:"custom_authentication_factors"` // Set of custom auth factor IDs
	TermsOfUse                  types.Set    `tfsdk:"terms_of_use"`                  // Set of terms of use IDs
}

// ConditionalAccessSessionControlsResourceModel represents session controls
type ConditionalAccessSessionControlsResourceModel struct {
	DisableResilienceDefaults       types.Bool                                                  `tfsdk:"disable_resilience_defaults"`
	ApplicationEnforcedRestrictions *ApplicationEnforcedRestrictionsSessionControlResourceModel `tfsdk:"application_enforced_restrictions"`
	CloudAppSecurity                *CloudAppSecuritySessionControlResourceModel                `tfsdk:"cloud_app_security"`
	SignInFrequency                 *SignInFrequencySessionControlResourceModel                 `tfsdk:"sign_in_frequency"`
	PersistentBrowser               *PersistentBrowserSessionControlResourceModel               `tfsdk:"persistent_browser"`
}

// ApplicationEnforcedRestrictionsSessionControlResourceModel represents app enforced restrictions
type ApplicationEnforcedRestrictionsSessionControlResourceModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// CloudAppSecuritySessionControlResourceModel represents cloud app security controls
type CloudAppSecuritySessionControlResourceModel struct {
	IsEnabled            types.Bool   `tfsdk:"is_enabled"`
	CloudAppSecurityType types.String `tfsdk:"cloud_app_security_type"` // CloudAppSecuritySessionControlType enum
}

// SignInFrequencySessionControlResourceModel represents sign-in frequency controls
type SignInFrequencySessionControlResourceModel struct {
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Type               types.String `tfsdk:"type"`                // SigninFrequencyType enum: days, hours
	Value              types.Int32  `tfsdk:"value"`               // Frequency value
	AuthenticationType types.String `tfsdk:"authentication_type"` // SigninFrequencyAuthenticationType enum
	FrequencyInterval  types.String `tfsdk:"frequency_interval"`  // SigninFrequencyInterval enum
}

// PersistentBrowserSessionControlResourceModel represents persistent browser controls
type PersistentBrowserSessionControlResourceModel struct {
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
	Mode      types.String `tfsdk:"mode"` // PersistentBrowserSessionControlMode enum: always, never
}

// Additional models for the applied policy results (if needed)
type AppliedConditionalAccessPolicyResourceModel struct {
	Id                          types.String                                   `tfsdk:"id"`
	DisplayName                 types.String                                   `tfsdk:"display_name"`
	Result                      types.String                                   `tfsdk:"result"`                         // AppliedConditionalAccessPolicyResult enum
	ConditionsSatisfied         types.String                                   `tfsdk:"conditions_satisfied"`           // ConditionalAccessConditions enum
	ConditionsNotSatisfied      types.String                                   `tfsdk:"conditions_not_satisfied"`       // ConditionalAccessConditions enum
	EnforcedGrantControls       types.Set                                      `tfsdk:"enforced_grant_controls"`        // Set of strings
	EnforcedSessionControls     types.Set                                      `tfsdk:"enforced_session_controls"`      // Set of strings
	SessionControlsNotSatisfied types.Set                                      `tfsdk:"session_controls_not_satisfied"` // Set of strings
	AuthenticationStrength      *AuthenticationStrengthResourceModel           `tfsdk:"authentication_strength"`
	IncludeRulesSatisfied       []*ConditionalAccessRuleSatisfiedResourceModel `tfsdk:"include_rules_satisfied"`
	ExcludeRulesSatisfied       []*ConditionalAccessRuleSatisfiedResourceModel `tfsdk:"exclude_rules_satisfied"`
}

// AuthenticationStrengthResourceModel represents authentication strength in applied policies
type AuthenticationStrengthResourceModel struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	PolicyType  types.String `tfsdk:"policy_type"`
}

// ConditionalAccessRuleSatisfiedResourceModel represents satisfied rules in applied policies
type ConditionalAccessRuleSatisfiedResourceModel struct {
	ConditionName types.String `tfsdk:"condition_name"`
	RuleSatisfied types.String `tfsdk:"rule_satisfied"`
}
