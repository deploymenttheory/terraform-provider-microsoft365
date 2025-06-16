package graphBetaConditionalAccessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConditionalAccessPolicyResourceModel represents the main conditional access policy
// This is the primary resource model, not ConditionalAccessRoot
type ConditionalAccessPolicyResourceModel struct {
	Id                        types.String                                   `tfsdk:"id"`
	DisplayName               types.String                                   `tfsdk:"display_name"`
	State                     types.String                                   `tfsdk:"state"`                       // ConditionalAccessPolicyState: enabled, disabled, enabledForReportingButNotEnforced
	CreatedDateTime           types.String                                   `tfsdk:"created_date_time"`           // Read-only
	ModifiedDateTime          types.String                                   `tfsdk:"modified_date_time"`          // Read-only
	DeletedDateTime           types.String                                   `tfsdk:"deleted_date_time"`           // Read-only
	TemplateId                types.String                                   `tfsdk:"template_id"`                 // Read-only
	PartialEnablementStrategy types.String                                   `tfsdk:"partial_enablement_strategy"` // Read-only
	Conditions                *ConditionalAccessConditionSetResourceModel    `tfsdk:"conditions"`
	GrantControls             *ConditionalAccessGrantControlsResourceModel   `tfsdk:"grant_controls"`
	SessionControls           *ConditionalAccessSessionControlsResourceModel `tfsdk:"session_controls"`
}

// ConditionalAccessConditionSetResourceModel represents conditions for conditional access
// Updated with all properties from JSON data
type ConditionalAccessConditionSetResourceModel struct {
	Applications               *ConditionalAccessApplicationsResourceModel         `tfsdk:"applications"`
	AuthenticationFlows        *ConditionalAccessAuthenticationFlowsResourceModel  `tfsdk:"authentication_flows"`
	ClientApplications         *ConditionalAccessClientApplicationsResourceModel   `tfsdk:"client_applications"`
	Clients                    *ConditionalAccessClientsResourceModel              `tfsdk:"clients"`          // New: different from clientApplications
	ClientAppTypes             types.Set                                           `tfsdk:"client_app_types"` // Set of ConditionalAccessClientApp enums
	Devices                    *ConditionalAccessDevicesResourceModel              `tfsdk:"devices"`
	DeviceStates               *ConditionalAccessDeviceStatesResourceModel         `tfsdk:"device_states"`       // Deprecated
	InsiderRiskLevels          types.String                                        `tfsdk:"insider_risk_levels"` // ConditionalAccessInsiderRiskLevels enum
	Locations                  *ConditionalAccessLocationsResourceModel            `tfsdk:"locations"`
	Platforms                  *ConditionalAccessPlatformsResourceModel            `tfsdk:"platforms"`
	ServicePrincipalRiskLevels types.Set                                           `tfsdk:"service_principal_risk_levels"` // Set of RiskLevel enums
	SignInRiskLevels           types.Set                                           `tfsdk:"sign_in_risk_levels"`           // Set of RiskLevel enums
	SignInRiskDetections       *ConditionalAccessSignInRiskDetectionsResourceModel `tfsdk:"sign_in_risk_detections"`       // New in PATCH format
	Times                      *ConditionalAccessTimesResourceModel                `tfsdk:"times"`                         // New: time-based conditions
	UserRiskLevels             types.Set                                           `tfsdk:"user_risk_levels"`              // Set of RiskLevel enums
	Users                      *ConditionalAccessUsersResourceModel                `tfsdk:"users"`
	ODataType                  types.String                                        `tfsdk:"odata_type"`
}

// ConditionalAccessApplicationsResourceModel represents application conditions
// Updated with globalSecureAccess
type ConditionalAccessApplicationsResourceModel struct {
	ApplicationFilter                           *ConditionalAccessFilterResourceModel             `tfsdk:"application_filter"`
	ExcludeApplications                         types.Set                                         `tfsdk:"exclude_applications"`
	GlobalSecureAccess                          *ConditionalAccessGlobalSecureAccessResourceModel `tfsdk:"global_secure_access"` // New in PATCH format
	IncludeApplications                         types.Set                                         `tfsdk:"include_applications"`
	IncludeAuthenticationContextClassReferences types.Set                                         `tfsdk:"include_authentication_context_class_references"`
	IncludeUserActions                          types.Set                                         `tfsdk:"include_user_actions"`
	NetworkAccess                               *ConditionalAccessNetworkAccessResourceModel      `tfsdk:"network_access"` // Deprecated
	ODataType                                   types.String                                      `tfsdk:"odata_type"`
}

// ConditionalAccessGlobalSecureAccessResourceModel represents global secure access
type ConditionalAccessGlobalSecureAccessResourceModel struct {
	ODataType types.String `tfsdk:"odata_type"`
}

// ConditionalAccessNetworkAccessResourceModel represents network access (deprecated)
type ConditionalAccessNetworkAccessResourceModel struct {
	ODataType types.String `tfsdk:"odata_type"`
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
// Updated: guestOrExternalUserTypes is a string, not a Set (comma-separated values)
type ConditionalAccessGuestsOrExternalUsersResourceModel struct {
	GuestOrExternalUserTypes types.String                                   `tfsdk:"guest_or_external_user_types"` // Comma-separated string like "InternalGuest,B2bCollaborationGuest"
	ExternalTenants          *ConditionalAccessExternalTenantsResourceModel `tfsdk:"external_tenants"`
}

// ConditionalAccessExternalTenantsResourceModel represents external tenant conditions
// Updated: Added ODataType for proper deserialization
type ConditionalAccessExternalTenantsResourceModel struct {
	MembershipKind types.String `tfsdk:"membership_kind"` // ConditionalAccessExternalTenantsMembershipKind enum
	Members        types.Set    `tfsdk:"members"`         // Set of tenant IDs (only for enumerated type)
	ODataType      types.String `tfsdk:"odata_type"`      // @odata.type for discriminator (e.g., "#microsoft.graph.conditionalAccessAllExternalTenants")
}

// ConditionalAccessClientApplicationsResourceModel represents client application conditions
type ConditionalAccessClientApplicationsResourceModel struct {
	IncludeServicePrincipals types.Set                             `tfsdk:"include_service_principals"`
	ExcludeServicePrincipals types.Set                             `tfsdk:"exclude_service_principals"`
	ServicePrincipalFilter   *ConditionalAccessFilterResourceModel `tfsdk:"service_principal_filter"`
}

// ConditionalAccessClientsResourceModel represents client conditions (new in PATCH format)
type ConditionalAccessClientsResourceModel struct {
	// Add properties as they become available in the SDK
	ODataType types.String `tfsdk:"odata_type"`
}

// ConditionalAccessSignInRiskDetectionsResourceModel represents sign-in risk detection conditions (new in PATCH format)
type ConditionalAccessSignInRiskDetectionsResourceModel struct {
	// Add properties as they become available in the SDK
	ODataType types.String `tfsdk:"odata_type"`
}

// ConditionalAccessTimesResourceModel represents time-based conditions
type ConditionalAccessTimesResourceModel struct {
	// Add properties as they become available in the SDK
	ODataType types.String `tfsdk:"odata_type"`
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
// Updated: confirmed includePlatforms and excludePlatforms from JSON
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
// Updated with all session control properties from JSON
type ConditionalAccessSessionControlsResourceModel struct {
	DisableResilienceDefaults          types.Bool                                                     `tfsdk:"disable_resilience_defaults"`
	ApplicationEnforcedRestrictions    *ApplicationEnforcedRestrictionsSessionControlResourceModel    `tfsdk:"application_enforced_restrictions"`
	CloudAppSecurity                   *CloudAppSecuritySessionControlResourceModel                   `tfsdk:"cloud_app_security"`
	ContinuousAccessEvaluation         *ContinuousAccessEvaluationSessionControlResourceModel         `tfsdk:"continuous_access_evaluation"` // New in JSON
	SignInFrequency                    *SignInFrequencySessionControlResourceModel                    `tfsdk:"sign_in_frequency"`
	PersistentBrowser                  *PersistentBrowserSessionControlResourceModel                  `tfsdk:"persistent_browser"`
	SecureSignInSession                *SecureSignInSessionControlResourceModel                       `tfsdk:"secure_sign_in_session"`                 // New in JSON
	NetworkAccessSecurity              *NetworkAccessSecuritySessionControlResourceModel              `tfsdk:"network_access_security"`                // New in PATCH format
	GlobalSecureAccessFilteringProfile *GlobalSecureAccessFilteringProfileSessionControlResourceModel `tfsdk:"global_secure_access_filtering_profile"` // New in PATCH format
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

// ContinuousAccessEvaluationSessionControlResourceModel represents continuous access evaluation controls
type ContinuousAccessEvaluationSessionControlResourceModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// SignInFrequencySessionControlResourceModel represents sign-in frequency controls
// Confirmed properties from JSON data
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

// SecureSignInSessionControlResourceModel represents secure sign-in session controls
type SecureSignInSessionControlResourceModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// NetworkAccessSecuritySessionControlResourceModel represents network access security controls (new in PATCH format)
type NetworkAccessSecuritySessionControlResourceModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// GlobalSecureAccessFilteringProfileSessionControlResourceModel represents global secure access filtering profile controls (new in PATCH format)
type GlobalSecureAccessFilteringProfileSessionControlResourceModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// Legacy models for ConditionalAccessRoot (if still needed for other resources)
type ConditionalAccessRootResourceModel struct {
	Id                                   types.String                                        `tfsdk:"id"`
	AuthenticationContextClassReferences []*AuthenticationContextClassReferenceResourceModel `tfsdk:"authentication_context_class_references"`
	AuthenticationStrength               *AuthenticationStrengthRootResourceModel            `tfsdk:"authentication_strength"`
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
	AuthenticationMethodModes  []*AuthenticationMethodModeDetailResourceModel `tfsdk:"authentication_method_modes"`
	Combinations               types.Set                                      `tfsdk:"combinations"` // Set of strings representing AuthenticationMethodModes
	Policies                   []*AuthenticationStrengthPolicyResourceModel   `tfsdk:"policies"`
}

// AuthenticationMethodModeDetailResourceModel represents details of an authentication method mode
type AuthenticationMethodModeDetailResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	AuthenticationMethod types.String `tfsdk:"authentication_method"` // String representation of BaseAuthenticationMethod enum
	DisplayName          types.String `tfsdk:"display_name"`
}

// AuthenticationStrengthPolicyResourceModel represents an authentication strength policy
type AuthenticationStrengthPolicyResourceModel struct {
	Id                        types.String                                           `tfsdk:"id"`
	AllowedCombinations       types.Set                                              `tfsdk:"allowed_combinations"`       // Set of strings representing AuthenticationMethodModes
	CombinationConfigurations []*AuthenticationCombinationConfigurationResourceModel `tfsdk:"combination_configurations"` // Collection of combination configuration objects
	CreatedDateTime           types.String                                           `tfsdk:"created_date_time"`
	Description               types.String                                           `tfsdk:"description"`
	DisplayName               types.String                                           `tfsdk:"display_name"`
	ModifiedDateTime          types.String                                           `tfsdk:"modified_date_time"`
	PolicyType                types.String                                           `tfsdk:"policy_type"`            // AuthenticationStrengthPolicyType enum: builtIn, custom, unknownFutureValue
	RequirementsSatisfied     types.String                                           `tfsdk:"requirements_satisfied"` // AuthenticationStrengthRequirements enum: none, mfa, unknownFutureValue
}

// AuthenticationCombinationConfigurationResourceModel represents a combination configuration
type AuthenticationCombinationConfigurationResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	AppliesToCombinations types.Set    `tfsdk:"applies_to_combinations"` // Set of strings representing AuthenticationMethodModes
	ODataType             types.String `tfsdk:"odata_type"`              // @odata.type for discriminator
}

// ConditionalAccessPolicyDetailResourceModel represents the details of a conditional access policy template
type ConditionalAccessPolicyDetailResourceModel struct {
	Conditions      *ConditionalAccessConditionSetResourceModel    `tfsdk:"conditions"`
	GrantControls   *ConditionalAccessGrantControlsResourceModel   `tfsdk:"grant_controls"`
	SessionControls *ConditionalAccessSessionControlsResourceModel `tfsdk:"session_controls"`
	ODataType       types.String                                   `tfsdk:"odata_type"`
}

// Applied policy models (if needed for other resources)
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
