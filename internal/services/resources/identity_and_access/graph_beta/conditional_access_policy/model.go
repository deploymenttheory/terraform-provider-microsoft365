// REF: https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta
package graphBetaConditionalAccessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConditionalAccessPolicyResourceModel represents the schema for the ConditionalAccessPolicy resource
type ConditionalAccessPolicyResourceModel struct {
	ID                        types.String                      `tfsdk:"id"`
	DisplayName               types.String                      `tfsdk:"display_name"`
	State                     types.String                      `tfsdk:"state"`
	CreatedDateTime           types.String                      `tfsdk:"created_date_time"`
	ModifiedDateTime          types.String                      `tfsdk:"modified_date_time"`
	DeletedDateTime           types.String                      `tfsdk:"deleted_date_time"`
	Conditions                *ConditionalAccessConditions      `tfsdk:"conditions"`
	GrantControls             *ConditionalAccessGrantControls   `tfsdk:"grant_controls"`
	SessionControls           *ConditionalAccessSessionControls `tfsdk:"session_controls"`
	TemplateId                types.String                      `tfsdk:"template_id"`
	PartialEnablementStrategy types.String                      `tfsdk:"partial_enablement_strategy"`
	Timeouts                  timeouts.Value                    `tfsdk:"timeouts"`
}

// ConditionalAccessConditions represents the conditions for the conditional access policy
type ConditionalAccessConditions struct {
	ClientAppTypes             types.Set                            `tfsdk:"client_app_types"`
	Applications               *ConditionalAccessApplications       `tfsdk:"applications"`
	Users                      *ConditionalAccessUsers              `tfsdk:"users"`
	Platforms                  *ConditionalAccessPlatforms          `tfsdk:"platforms"`
	Locations                  *ConditionalAccessLocations          `tfsdk:"locations"`
	Devices                    *ConditionalAccessDevices            `tfsdk:"devices"`
	SignInRiskLevels           types.Set                            `tfsdk:"sign_in_risk_levels"`
	UserRiskLevels             types.Set                            `tfsdk:"user_risk_levels"`
	ServicePrincipalRiskLevels types.Set                            `tfsdk:"service_principal_risk_levels"`
	ClientApplications         *ConditionalAccessClientApplications `tfsdk:"client_applications"`
	Times                      *ConditionalAccessTimes              `tfsdk:"times"`
	DeviceStates               *ConditionalAccessDeviceStates       `tfsdk:"device_states"`
}

// ConditionalAccessApplications represents the applications condition
type ConditionalAccessApplications struct {
	IncludeApplications                         types.Set                `tfsdk:"include_applications"`
	ExcludeApplications                         types.Set                `tfsdk:"exclude_applications"`
	IncludeUserActions                          types.Set                `tfsdk:"include_user_actions"`
	IncludeAuthenticationContextClassReferences types.Set                `tfsdk:"include_authentication_context_class_references"`
	ApplicationFilter                           *ConditionalAccessFilter `tfsdk:"application_filter"`
}

// ConditionalAccessUsers represents the users condition
type ConditionalAccessUsers struct {
	IncludeUsers                 types.Set                               `tfsdk:"include_users"`
	ExcludeUsers                 types.Set                               `tfsdk:"exclude_users"`
	IncludeGroups                types.Set                               `tfsdk:"include_groups"`
	ExcludeGroups                types.Set                               `tfsdk:"exclude_groups"`
	IncludeRoles                 types.Set                               `tfsdk:"include_roles"`
	ExcludeRoles                 types.Set                               `tfsdk:"exclude_roles"`
	IncludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsers `tfsdk:"include_guests_or_external_users"`
	ExcludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsers `tfsdk:"exclude_guests_or_external_users"`
}

// ConditionalAccessGuestsOrExternalUsers represents the guests or external users configuration
type ConditionalAccessGuestsOrExternalUsers struct {
	GuestOrExternalUserTypes types.String                      `tfsdk:"guest_or_external_user_types"`
	ExternalTenants          *ConditionalAccessExternalTenants `tfsdk:"external_tenants"`
}

// ConditionalAccessExternalTenants represents the external tenants configuration
type ConditionalAccessExternalTenants struct {
	MembershipKind types.String `tfsdk:"membership_kind"`
	Members        types.Set    `tfsdk:"members"`
}

// ConditionalAccessPlatforms represents the platforms condition
type ConditionalAccessPlatforms struct {
	IncludePlatforms types.Set `tfsdk:"include_platforms"`
	ExcludePlatforms types.Set `tfsdk:"exclude_platforms"`
}

// ConditionalAccessLocations represents the locations condition
type ConditionalAccessLocations struct {
	IncludeLocations types.Set `tfsdk:"include_locations"`
	ExcludeLocations types.Set `tfsdk:"exclude_locations"`
}

// ConditionalAccessDevices represents the devices condition
type ConditionalAccessDevices struct {
	IncludeDevices      types.Set                `tfsdk:"include_devices"`
	ExcludeDevices      types.Set                `tfsdk:"exclude_devices"`
	IncludeDeviceStates types.Set                `tfsdk:"include_device_states"`
	ExcludeDeviceStates types.Set                `tfsdk:"exclude_device_states"`
	DeviceFilter        *ConditionalAccessFilter `tfsdk:"device_filter"`
}

// ConditionalAccessFilter represents a filter for applications or devices
type ConditionalAccessFilter struct {
	Mode types.String `tfsdk:"mode"`
	Rule types.String `tfsdk:"rule"`
}

// ConditionalAccessTimes represents the times condition
type ConditionalAccessTimes struct {
	IncludedRanges types.Set    `tfsdk:"included_ranges"`
	ExcludedRanges types.Set    `tfsdk:"excluded_ranges"`
	AllDay         types.Bool   `tfsdk:"all_day"`
	StartTime      types.String `tfsdk:"start_time"`
	EndTime        types.String `tfsdk:"end_time"`
	TimeZone       types.String `tfsdk:"time_zone"`
}

// ConditionalAccessDeviceStates represents the device states condition
type ConditionalAccessDeviceStates struct {
	IncludeStates types.Set `tfsdk:"include_states"`
	ExcludeStates types.Set `tfsdk:"exclude_states"`
}

// ConditionalAccessGrantControls represents the grant controls for the conditional access policy
type ConditionalAccessGrantControls struct {
	Operator                    types.String                             `tfsdk:"operator"`
	BuiltInControls             types.Set                                `tfsdk:"built_in_controls"`
	CustomAuthenticationFactors types.Set                                `tfsdk:"custom_authentication_factors"`
	TermsOfUse                  types.Set                                `tfsdk:"terms_of_use"`
	AuthenticationStrength      *ConditionalAccessAuthenticationStrength `tfsdk:"authentication_strength"`
}

// ConditionalAccessAuthenticationStrength represents the authentication strength configuration
type ConditionalAccessAuthenticationStrength struct {
	ID                    types.String `tfsdk:"id"`
	DisplayName           types.String `tfsdk:"display_name"`
	Description           types.String `tfsdk:"description"`
	PolicyType            types.String `tfsdk:"policy_type"`
	RequirementsSatisfied types.String `tfsdk:"requirements_satisfied"`
	AllowedCombinations   types.Set    `tfsdk:"allowed_combinations"`
	CreatedDateTime       types.String `tfsdk:"created_date_time"`
	ModifiedDateTime      types.String `tfsdk:"modified_date_time"`
}

// ConditionalAccessSessionControls represents the session controls for the conditional access policy
type ConditionalAccessSessionControls struct {
	ApplicationEnforcedRestrictions *ConditionalAccessApplicationEnforcedRestrictions `tfsdk:"application_enforced_restrictions"`
	CloudAppSecurity                *ConditionalAccessCloudAppSecurity                `tfsdk:"cloud_app_security"`
	SignInFrequency                 *ConditionalAccessSignInFrequency                 `tfsdk:"sign_in_frequency"`
	PersistentBrowser               *ConditionalAccessPersistentBrowser               `tfsdk:"persistent_browser"`
	DisableResilienceDefaults       types.Bool                                        `tfsdk:"disable_resilience_defaults"`
	ContinuousAccessEvaluation      *ConditionalAccessContinuousAccessEvaluation      `tfsdk:"continuous_access_evaluation"`
	SecureSignInSession             *ConditionalAccessSecureSignInSession             `tfsdk:"secure_sign_in_session"`
}

// ConditionalAccessApplicationEnforcedRestrictions represents the application enforced restrictions configuration
type ConditionalAccessApplicationEnforcedRestrictions struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// ConditionalAccessCloudAppSecurity represents the cloud app security configuration
type ConditionalAccessCloudAppSecurity struct {
	IsEnabled            types.Bool   `tfsdk:"is_enabled"`
	CloudAppSecurityType types.String `tfsdk:"cloud_app_security_type"`
}

// ConditionalAccessSignInFrequency represents the sign-in frequency configuration
type ConditionalAccessSignInFrequency struct {
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Type               types.String `tfsdk:"type"`
	Value              types.Int64  `tfsdk:"value"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	FrequencyInterval  types.String `tfsdk:"frequency_interval"`
}

// ConditionalAccessPersistentBrowser represents the persistent browser configuration
type ConditionalAccessPersistentBrowser struct {
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
	Mode      types.String `tfsdk:"mode"`
}

// ConditionalAccessContinuousAccessEvaluation represents the continuous access evaluation configuration
type ConditionalAccessContinuousAccessEvaluation struct {
	Mode types.String `tfsdk:"mode"`
}

// ConditionalAccessSecureSignInSession represents the secure sign-in session configuration
type ConditionalAccessSecureSignInSession struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

// ConditionalAccessClientApplications represents the client applications configuration
type ConditionalAccessClientApplications struct {
	IncludeServicePrincipals types.Set `tfsdk:"include_service_principals"`
	ExcludeServicePrincipals types.Set `tfsdk:"exclude_service_principals"`
}
