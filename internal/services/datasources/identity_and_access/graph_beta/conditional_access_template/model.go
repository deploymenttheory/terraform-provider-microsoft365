// REF: https://learn.microsoft.com/en-us/graph/api/conditionalaccessroot-list-templates?view=graph-rest-beta

package graphBetaConditionalAccessTemplate

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConditionalAccessTemplateDataSourceModel struct {
	ID          types.String                      `tfsdk:"id"`
	Name        types.String                      `tfsdk:"name"`
	TemplateID  types.String                      `tfsdk:"template_id"`
	Description types.String                      `tfsdk:"description"`
	Scenarios   types.Set                         `tfsdk:"scenarios"`
	Details     *ConditionalAccessTemplateDetails `tfsdk:"details"`
	Timeouts    timeouts.Value                    `tfsdk:"timeouts"`
}

type ConditionalAccessTemplateDetails struct {
	SessionControls *ConditionalAccessSessionControls `tfsdk:"session_controls"`
	Conditions      *ConditionalAccessConditions      `tfsdk:"conditions"`
	GrantControls   *ConditionalAccessGrantControls   `tfsdk:"grant_controls"`
}

type ConditionalAccessConditions struct {
	UserRiskLevels             types.List                           `tfsdk:"user_risk_levels"`
	SignInRiskLevels           types.List                           `tfsdk:"sign_in_risk_levels"`
	ClientAppTypes             types.List                           `tfsdk:"client_app_types"`
	ServicePrincipalRiskLevels types.List                           `tfsdk:"service_principal_risk_levels"`
	AgentIdRiskLevels          types.Set                            `tfsdk:"agent_id_risk_levels"`
	InsiderRiskLevels          types.Set                            `tfsdk:"insider_risk_levels"`
	Platforms                  *ConditionalAccessPlatforms          `tfsdk:"platforms"`
	Locations                  *ConditionalAccessLocations          `tfsdk:"locations"`
	Devices                    *ConditionalAccessDevices            `tfsdk:"devices"`
	ClientApplications         *ConditionalAccessClientApplications `tfsdk:"client_applications"`
	Applications               *ConditionalAccessApplications       `tfsdk:"applications"`
	Users                      *ConditionalAccessUsers              `tfsdk:"users"`
}

type ConditionalAccessApplications struct {
	IncludeApplications                         types.List `tfsdk:"include_applications"`
	ExcludeApplications                         types.List `tfsdk:"exclude_applications"`
	IncludeUserActions                          types.List `tfsdk:"include_user_actions"`
	IncludeAuthenticationContextClassReferences types.List `tfsdk:"include_authentication_context_class_references"`
}

type ConditionalAccessUsers struct {
	IncludeUsers                 types.List                              `tfsdk:"include_users"`
	ExcludeUsers                 types.List                              `tfsdk:"exclude_users"`
	IncludeGroups                types.List                              `tfsdk:"include_groups"`
	ExcludeGroups                types.List                              `tfsdk:"exclude_groups"`
	IncludeRoles                 types.List                              `tfsdk:"include_roles"`
	ExcludeRoles                 types.List                              `tfsdk:"exclude_roles"`
	IncludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsers `tfsdk:"include_guests_or_external_users"`
	ExcludeGuestsOrExternalUsers *ConditionalAccessGuestsOrExternalUsers `tfsdk:"exclude_guests_or_external_users"`
}

type ConditionalAccessPlatforms struct {
	IncludePlatforms types.List `tfsdk:"include_platforms"`
	ExcludePlatforms types.List `tfsdk:"exclude_platforms"`
}

type ConditionalAccessLocations struct {
	IncludeLocations types.List `tfsdk:"include_locations"`
	ExcludeLocations types.List `tfsdk:"exclude_locations"`
}

type ConditionalAccessDevices struct {
	IncludeDeviceStates types.List               `tfsdk:"include_device_states"`
	ExcludeDeviceStates types.List               `tfsdk:"exclude_device_states"`
	IncludeDevices      types.List               `tfsdk:"include_devices"`
	ExcludeDevices      types.List               `tfsdk:"exclude_devices"`
	DeviceFilter        *ConditionalAccessFilter `tfsdk:"device_filter"`
}

type ConditionalAccessFilter struct {
	Mode types.String `tfsdk:"mode"`
	Rule types.String `tfsdk:"rule"`
}

type ConditionalAccessClientApplications struct {
	IncludeServicePrincipals        types.List `tfsdk:"include_service_principals"`
	IncludeAgentIdServicePrincipals types.List `tfsdk:"include_agent_id_service_principals"`
	ExcludeServicePrincipals        types.List `tfsdk:"exclude_service_principals"`
	ExcludeAgentIdServicePrincipals types.List `tfsdk:"exclude_agent_id_service_principals"`
}

type ConditionalAccessGrantControls struct {
	Operator                    types.String                             `tfsdk:"operator"`
	BuiltInControls             types.List                               `tfsdk:"built_in_controls"`
	CustomAuthenticationFactors types.List                               `tfsdk:"custom_authentication_factors"`
	TermsOfUse                  types.List                               `tfsdk:"terms_of_use"`
	AuthenticationStrength      *ConditionalAccessAuthenticationStrength `tfsdk:"authentication_strength"`
}

type ConditionalAccessAuthenticationStrength struct {
	ID                    types.String `tfsdk:"id"`
	CreatedDateTime       types.String `tfsdk:"created_date_time"`
	ModifiedDateTime      types.String `tfsdk:"modified_date_time"`
	DisplayName           types.String `tfsdk:"display_name"`
	Description           types.String `tfsdk:"description"`
	PolicyType            types.String `tfsdk:"policy_type"`
	RequirementsSatisfied types.String `tfsdk:"requirements_satisfied"`
	AllowedCombinations   types.List   `tfsdk:"allowed_combinations"`
}

type ConditionalAccessSessionControls struct {
	DisableResilienceDefaults          types.Bool                                                 `tfsdk:"disable_resilience_defaults"`
	ApplicationEnforcedRestrictions    *ConditionalAccessApplicationEnforcedRestrictions          `tfsdk:"application_enforced_restrictions"`
	CloudAppSecurity                   *ConditionalAccessCloudAppSecuritySessionControl           `tfsdk:"cloud_app_security"`
	PersistentBrowser                  *ConditionalAccessPersistentBrowser                        `tfsdk:"persistent_browser"`
	ContinuousAccessEvaluation         *ConditionalAccessContinuousAccessEvaluationSessionControl `tfsdk:"continuous_access_evaluation"`
	SecureSignInSession                *ConditionalAccessSecureSignInSessionControl               `tfsdk:"secure_sign_in_session"`
	GlobalSecureAccessFilteringProfile *ConditionalAccessGlobalSecureAccessFilteringProfile       `tfsdk:"global_secure_access_filtering_profile"`
	SignInFrequency                    *ConditionalAccessSignInFrequency                          `tfsdk:"sign_in_frequency"`
}

type ConditionalAccessCloudAppSecuritySessionControl struct {
	IsEnabled            types.Bool   `tfsdk:"is_enabled"`
	CloudAppSecurityType types.String `tfsdk:"cloud_app_security_type"`
}

type ConditionalAccessContinuousAccessEvaluationSessionControl struct {
	Mode types.String `tfsdk:"mode"`
}

type ConditionalAccessSecureSignInSessionControl struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

type ConditionalAccessGlobalSecureAccessFilteringProfile struct {
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
	ProfileId types.String `tfsdk:"profile_id"`
}

type ConditionalAccessSignInFrequency struct {
	Value              types.Int64  `tfsdk:"value"`
	Type               types.String `tfsdk:"type"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	FrequencyInterval  types.String `tfsdk:"frequency_interval"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
}

type ConditionalAccessPersistentBrowser struct {
	Mode      types.String `tfsdk:"mode"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
}

type ConditionalAccessApplicationEnforcedRestrictions struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

type ConditionalAccessGuestsOrExternalUsers struct {
	GuestOrExternalUserTypes types.Set                         `tfsdk:"guest_or_external_user_types"`
	ExternalTenants          *ConditionalAccessExternalTenants `tfsdk:"external_tenants"`
}

type ConditionalAccessExternalTenants struct {
	MembershipKind types.String `tfsdk:"membership_kind"`
}
