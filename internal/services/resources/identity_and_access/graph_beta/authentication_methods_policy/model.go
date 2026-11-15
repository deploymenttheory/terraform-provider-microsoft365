// REF: https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodspolicy?view=graph-rest-beta
package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthenticationMethodsPolicyResourceModel represents the schema for the AuthenticationMethodsPolicy resource.
//
// Every nested block is a types.Object because each one is optional and computed: blocks that are not
// configured are read from the tenant and are never written back.
type AuthenticationMethodsPolicyResourceModel struct {
	ID                               types.String   `tfsdk:"id"`
	DisplayName                      types.String   `tfsdk:"display_name"`
	Description                      types.String   `tfsdk:"description"`
	PolicyMigrationState             types.String   `tfsdk:"policy_migration_state"`
	RegistrationEnforcement          types.Object   `tfsdk:"registration_enforcement"`
	ReportSuspiciousActivitySettings types.Object   `tfsdk:"report_suspicious_activity_settings"`
	SystemCredentialPreferences      types.Object   `tfsdk:"system_credential_preferences"`
	Email                            types.Object   `tfsdk:"email"`
	TemporaryAccessPass              types.Object   `tfsdk:"temporary_access_pass"`
	Sms                              types.Object   `tfsdk:"sms"`
	SoftwareOath                     types.Object   `tfsdk:"software_oath"`
	HardwareOath                     types.Object   `tfsdk:"hardware_oath"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}

// RegistrationEnforcementModel represents the registrationEnforcement complex type.
type RegistrationEnforcementModel struct {
	AuthenticationMethodsRegistrationCampaign types.Object `tfsdk:"authentication_methods_registration_campaign"`
}

// RegistrationCampaignModel represents the authenticationMethodsRegistrationCampaign complex type.
type RegistrationCampaignModel struct {
	SnoozeDurationInDays                   types.Int32  `tfsdk:"snooze_duration_in_days"`
	EnforceRegistrationAfterAllowedSnoozes types.Bool   `tfsdk:"enforce_registration_after_allowed_snoozes"`
	State                                  types.String `tfsdk:"state"`
	IncludeTargets                         types.Set    `tfsdk:"include_targets"`
	ExcludeTargets                         types.Set    `tfsdk:"exclude_targets"`
}

// RegistrationCampaignIncludeTargetModel represents the authenticationMethodsRegistrationCampaignIncludeTarget complex type.
type RegistrationCampaignIncludeTargetModel struct {
	ID                           types.String `tfsdk:"id"`
	TargetType                   types.String `tfsdk:"target_type"`
	TargetedAuthenticationMethod types.String `tfsdk:"targeted_authentication_method"`
}

// ReportSuspiciousActivitySettingsModel represents the reportSuspiciousActivitySettings complex type.
type ReportSuspiciousActivitySettingsModel struct {
	State              types.String `tfsdk:"state"`
	VoiceReportingCode types.Int32  `tfsdk:"voice_reporting_code"`
	IncludeTarget      types.Object `tfsdk:"include_target"`
}

// SystemCredentialPreferencesModel represents the systemCredentialPreferences complex type.
type SystemCredentialPreferencesModel struct {
	State          types.String `tfsdk:"state"`
	IncludeTargets types.Set    `tfsdk:"include_targets"`
	ExcludeTargets types.Set    `tfsdk:"exclude_targets"`
}

// TargetModel represents the includeTarget and excludeTarget complex types (a user or group reference).
type TargetModel struct {
	ID         types.String `tfsdk:"id"`
	TargetType types.String `tfsdk:"target_type"`
}

// MethodTargetModel represents the authenticationMethodTarget complex type.
type MethodTargetModel struct {
	ID                     types.String `tfsdk:"id"`
	TargetType             types.String `tfsdk:"target_type"`
	IsRegistrationRequired types.Bool   `tfsdk:"is_registration_required"`
}

// SmsMethodTargetModel represents the smsAuthenticationMethodTarget complex type.
type SmsMethodTargetModel struct {
	ID                     types.String `tfsdk:"id"`
	TargetType             types.String `tfsdk:"target_type"`
	IsRegistrationRequired types.Bool   `tfsdk:"is_registration_required"`
	IsUsableForSignIn      types.Bool   `tfsdk:"is_usable_for_sign_in"`
}

// EmailMethodModel represents the emailAuthenticationMethodConfiguration resource type.
type EmailMethodModel struct {
	State                        types.String `tfsdk:"state"`
	AllowExternalIdToUseEmailOtp types.String `tfsdk:"allow_external_id_to_use_email_otp"`
	IncludeTargets               types.Set    `tfsdk:"include_targets"`
	ExcludeTargets               types.Set    `tfsdk:"exclude_targets"`
}

// TemporaryAccessPassMethodModel represents the temporaryAccessPassAuthenticationMethodConfiguration resource type.
type TemporaryAccessPassMethodModel struct {
	State                    types.String `tfsdk:"state"`
	DefaultLifetimeInMinutes types.Int32  `tfsdk:"default_lifetime_in_minutes"`
	DefaultLength            types.Int32  `tfsdk:"default_length"`
	MinimumLifetimeInMinutes types.Int32  `tfsdk:"minimum_lifetime_in_minutes"`
	MaximumLifetimeInMinutes types.Int32  `tfsdk:"maximum_lifetime_in_minutes"`
	IsUsableOnce             types.Bool   `tfsdk:"is_usable_once"`
	IncludeTargets           types.Set    `tfsdk:"include_targets"`
	ExcludeTargets           types.Set    `tfsdk:"exclude_targets"`
}

// SmsMethodModel represents the smsAuthenticationMethodConfiguration resource type.
type SmsMethodModel struct {
	State          types.String `tfsdk:"state"`
	IncludeTargets types.Set    `tfsdk:"include_targets"`
	ExcludeTargets types.Set    `tfsdk:"exclude_targets"`
}

// OathMethodModel represents the softwareOathAuthenticationMethodConfiguration and
// hardwareOathAuthenticationMethodConfiguration resource types, which share the same shape.
type OathMethodModel struct {
	State          types.String `tfsdk:"state"`
	IncludeTargets types.Set    `tfsdk:"include_targets"`
	ExcludeTargets types.Set    `tfsdk:"exclude_targets"`
}

// targetAttrTypes are the attribute types of TargetModel.
var targetAttrTypes = map[string]attr.Type{
	"id":          types.StringType,
	"target_type": types.StringType,
}

// methodTargetAttrTypes are the attribute types of MethodTargetModel.
var methodTargetAttrTypes = map[string]attr.Type{
	"id":                       types.StringType,
	"target_type":              types.StringType,
	"is_registration_required": types.BoolType,
}

// smsMethodTargetAttrTypes are the attribute types of SmsMethodTargetModel.
var smsMethodTargetAttrTypes = map[string]attr.Type{
	"id":                       types.StringType,
	"target_type":              types.StringType,
	"is_registration_required": types.BoolType,
	"is_usable_for_sign_in":    types.BoolType,
}

// registrationCampaignIncludeTargetAttrTypes are the attribute types of RegistrationCampaignIncludeTargetModel.
var registrationCampaignIncludeTargetAttrTypes = map[string]attr.Type{
	"id":                             types.StringType,
	"target_type":                    types.StringType,
	"targeted_authentication_method": types.StringType,
}

// registrationCampaignAttrTypes are the attribute types of RegistrationCampaignModel.
var registrationCampaignAttrTypes = map[string]attr.Type{
	"snooze_duration_in_days":                    types.Int32Type,
	"enforce_registration_after_allowed_snoozes": types.BoolType,
	"state":           types.StringType,
	"include_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: registrationCampaignIncludeTargetAttrTypes}},
	"exclude_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
}

// registrationEnforcementAttrTypes are the attribute types of RegistrationEnforcementModel.
var registrationEnforcementAttrTypes = map[string]attr.Type{
	"authentication_methods_registration_campaign": types.ObjectType{AttrTypes: registrationCampaignAttrTypes},
}

// reportSuspiciousActivitySettingsAttrTypes are the attribute types of ReportSuspiciousActivitySettingsModel.
var reportSuspiciousActivitySettingsAttrTypes = map[string]attr.Type{
	"state":                types.StringType,
	"voice_reporting_code": types.Int32Type,
	"include_target":       types.ObjectType{AttrTypes: targetAttrTypes},
}

// systemCredentialPreferencesAttrTypes are the attribute types of SystemCredentialPreferencesModel.
var systemCredentialPreferencesAttrTypes = map[string]attr.Type{
	"state":           types.StringType,
	"include_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
	"exclude_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
}

// emailMethodAttrTypes are the attribute types of EmailMethodModel.
var emailMethodAttrTypes = map[string]attr.Type{
	"state":                              types.StringType,
	"allow_external_id_to_use_email_otp": types.StringType,
	"include_targets":                    types.SetType{ElemType: types.ObjectType{AttrTypes: methodTargetAttrTypes}},
	"exclude_targets":                    types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
}

// temporaryAccessPassMethodAttrTypes are the attribute types of TemporaryAccessPassMethodModel.
var temporaryAccessPassMethodAttrTypes = map[string]attr.Type{
	"state":                       types.StringType,
	"default_lifetime_in_minutes": types.Int32Type,
	"default_length":              types.Int32Type,
	"minimum_lifetime_in_minutes": types.Int32Type,
	"maximum_lifetime_in_minutes": types.Int32Type,
	"is_usable_once":              types.BoolType,
	"include_targets":             types.SetType{ElemType: types.ObjectType{AttrTypes: methodTargetAttrTypes}},
	"exclude_targets":             types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
}

// smsMethodAttrTypes are the attribute types of SmsMethodModel.
var smsMethodAttrTypes = map[string]attr.Type{
	"state":           types.StringType,
	"include_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: smsMethodTargetAttrTypes}},
	"exclude_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
}

// oathMethodAttrTypes are the attribute types of OathMethodModel.
var oathMethodAttrTypes = map[string]attr.Type{
	"state":           types.StringType,
	"include_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: methodTargetAttrTypes}},
	"exclude_targets": types.SetType{ElemType: types.ObjectType{AttrTypes: targetAttrTypes}},
}
