// REF: https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta
package graphBetaConditionalAccessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConditionalAccessPolicyResourceModel struct {
	ID               types.String                           `tfsdk:"id"`
	Description      types.String                           `tfsdk:"description"`
	DisplayName      types.String                           `tfsdk:"display_name"`
	CreatedDateTime  types.String                           `tfsdk:"created_date_time"`
	ModifiedDateTime types.String                           `tfsdk:"modified_date_time"`
	State            types.String                           `tfsdk:"state"`
	Conditions       *ConditionalAccessConditionsModel      `tfsdk:"conditions"`
	GrantControls    *ConditionalAccessGrantControlsModel   `tfsdk:"grant_controls"`
	SessionControls  *ConditionalAccessSessionControlsModel `tfsdk:"session_controls"`
	Timeouts         timeouts.Value                         `tfsdk:"timeouts"`
}

type ConditionalAccessConditionsModel struct {
	Applications               *ConditionalAccessApplicationsModel        `tfsdk:"applications"`
	Users                      *ConditionalAccessUsersModel               `tfsdk:"users"`
	ClientApplications         *ConditionalAccessClientApplicationsModel  `tfsdk:"client_applications"`
	ClientAppTypes             []types.String                             `tfsdk:"client_app_types"`
	DeviceStates               *ConditionalAccessDeviceStatesModel        `tfsdk:"device_states"`
	Devices                    *ConditionalAccessDevicesModel             `tfsdk:"devices"`
	Locations                  *ConditionalAccessLocationsModel           `tfsdk:"locations"`
	Platforms                  *ConditionalAccessPlatformsModel           `tfsdk:"platforms"`
	ServicePrincipalRiskLevels []types.String                             `tfsdk:"service_principal_risk_levels"`
	SignInRiskLevels           []types.String                             `tfsdk:"sign_in_risk_levels"`
	UserRiskLevels             []types.String                             `tfsdk:"user_risk_levels"`
	AuthenticationFlows        *ConditionalAccessAuthenticationFlowsModel `tfsdk:"authentication_flows"`
	InsiderRiskLevels          types.String                               `tfsdk:"insider_risk_levels"`
}

type ConditionalAccessApplicationsModel struct {
	IncludeApplications []types.String                `tfsdk:"include_applications"`
	ExcludeApplications []types.String                `tfsdk:"exclude_applications"`
	ApplicationFilter   *ConditionalAccessFilterModel `tfsdk:"application_filter"`
	IncludeUserActions  []types.String                `tfsdk:"include_user_actions"`
}

type ConditionalAccessUsersModel struct {
	ExcludeGroups                []types.String                               `tfsdk:"exclude_groups"`
	ExcludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsersModel `tfsdk:"exclude_guests_or_external_users"`
	ExcludeRoles                 []types.String                               `tfsdk:"exclude_roles"`
	ExcludeUsers                 []types.String                               `tfsdk:"exclude_users"`
	IncludeGroups                []types.String                               `tfsdk:"include_groups"`
	IncludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsersModel `tfsdk:"include_guests_or_external_users"`
	IncludeRoles                 []types.String                               `tfsdk:"include_roles"`
	IncludeUsers                 []types.String                               `tfsdk:"include_users"`
}

type ConditionalAccessGuestsOrExternalUsersModel struct {
	ExternalTenants          *ConditionalAccessExternalTenantsModel `tfsdk:"external_tenants"`
	GuestOrExternalUserTypes types.String                           `tfsdk:"guest_or_external_user_types"`
}

type ConditionalAccessExternalTenantsModel struct {
	MembershipKind types.String `tfsdk:"membership_kind"`
}

type ConditionalAccessClientApplicationsModel struct {
	ExcludeServicePrincipals []types.String                `tfsdk:"exclude_service_principals"`
	IncludeServicePrincipals []types.String                `tfsdk:"include_service_principals"`
	ServicePrincipalFilter   *ConditionalAccessFilterModel `tfsdk:"service_principal_filter"`
}

type ConditionalAccessDeviceStatesModel struct {
	IncludeStates []types.String `tfsdk:"include_states"`
	ExcludeStates []types.String `tfsdk:"exclude_states"`
}

type ConditionalAccessDevicesModel struct {
	IncludeDevices []types.String                `tfsdk:"include_devices"`
	ExcludeDevices []types.String                `tfsdk:"exclude_devices"`
	IncludeStates  []types.String                `tfsdk:"include_states"` // TODO - validate this. sdk different to msft docs
	ExcludeStates  []types.String                `tfsdk:"exclude_states"` // TODO - validate this. sdk different to msft docs
	DeviceFilter   *ConditionalAccessFilterModel `tfsdk:"device_filter"`
}

type ConditionalAccessLocationsModel struct {
	ExcludeLocations []types.String `tfsdk:"exclude_locations"`
	IncludeLocations []types.String `tfsdk:"include_locations"`
}

type ConditionalAccessPlatformsModel struct {
	ExcludePlatforms []types.String `tfsdk:"exclude_platforms"`
	IncludePlatforms []types.String `tfsdk:"include_platforms"`
}

type ConditionalAccessAuthenticationFlowsModel struct {
	TransferMethods types.String `tfsdk:"transfer_methods"`
}

type ConditionalAccessFilterModel struct {
	Mode types.String `tfsdk:"mode"`
	Rule types.String `tfsdk:"rule"`
}

type ConditionalAccessGrantControlsModel struct {
	BuiltInControls             []types.String                     `tfsdk:"built_in_controls"`
	CustomAuthenticationFactors []types.String                     `tfsdk:"custom_authentication_factors"`
	Operator                    types.String                       `tfsdk:"operator"`
	TermsOfUse                  []types.String                     `tfsdk:"terms_of_use"`
	AuthenticationStrength      *AuthenticationStrengthPolicyModel `tfsdk:"authentication_strength"`
}

type AuthenticationStrengthPolicyModel struct {
	ID                    types.String   `tfsdk:"id"`
	CreatedDateTime       types.String   `tfsdk:"created_date_time"`
	ModifiedDateTime      types.String   `tfsdk:"modified_date_time"`
	DisplayName           types.String   `tfsdk:"display_name"`
	Description           types.String   `tfsdk:"description"`
	PolicyType            types.String   `tfsdk:"policy_type"`
	RequirementsSatisfied types.String   `tfsdk:"requirements_satisfied"`
	AllowedCombinations   []types.String `tfsdk:"allowed_combinations"`
}

// Ref: https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesssessioncontrols?view=graph-rest-beta
type ConditionalAccessSessionControlsModel struct {
	ApplicationEnforcedRestrictions *ApplicationEnforcedRestrictionsSessionControlModel `tfsdk:"application_enforced_restrictions"`
	CloudAppSecurity                *CloudAppSecuritySessionControlModel                `tfsdk:"cloud_app_security"`
	ContinuousAccessEvaluation      *ContinuousAccessEvaluationSessionControlModel      `tfsdk:"continuous_access_evaluation"`
	PersistentBrowser               *PersistentBrowserSessionControlModel               `tfsdk:"persistent_browser"`
	SignInFrequency                 *SignInFrequencySessionControlModel                 `tfsdk:"sign_in_frequency"`
	DisableResilienceDefaults       types.Bool                                          `tfsdk:"disable_resilience_defaults"`
	SecureSignInSession             *SecureSignInSessionControlModel                    `tfsdk:"secure_sign_in_session"`
}

type ApplicationEnforcedRestrictionsSessionControlModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

type CloudAppSecuritySessionControlModel struct {
	IsEnabled            types.Bool   `tfsdk:"is_enabled"`
	CloudAppSecurityType types.String `tfsdk:"cloud_app_security_type"`
}

type ContinuousAccessEvaluationSessionControlModel struct {
	Mode types.String `tfsdk:"mode"`
}

type PersistentBrowserSessionControlModel struct {
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
	Mode      types.String `tfsdk:"mode"`
}

type SignInFrequencySessionControlModel struct {
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Type               types.String `tfsdk:"type"`
	Value              types.Int64  `tfsdk:"value"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	FrequencyInterval  types.String `tfsdk:"frequency_interval"`
}

type SecureSignInSessionControlModel struct {
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Type               types.String `tfsdk:"type"`
	Value              types.Int64  `tfsdk:"value"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	FrequencyInterval  types.String `tfsdk:"frequency_interval"`
}
