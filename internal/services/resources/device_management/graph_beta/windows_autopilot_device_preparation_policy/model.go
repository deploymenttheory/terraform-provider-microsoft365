// REF: https://cloudflow.be/windows-autopilot-device-perpetration-with-graph-api/blog-post/
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/autopilot/device-preparation/overview

package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsAutopilotDevicePreparationPolicyResourceModel holds the configuration for a Windows Autopilot Device Preparation policy.
// This aligns with DeviceManagementConfigurationPolicy in Graph API
type WindowsAutopilotDevicePreparationPolicyResourceModel struct {
	// Base policy fields from DeviceManagementConfigurationPolicy
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	IsAssigned           types.Bool   `tfsdk:"is_assigned"`
	SettingsCount        types.Int32  `tfsdk:"settings_count"`
	RoleScopeTagIds      types.Set    `tfsdk:"role_scope_tag_ids"`
	// Platforms and technologies are computed and not user-configurable for this policy
	Platforms    types.String `tfsdk:"platforms"`
	Technologies types.String `tfsdk:"technologies"`
	// Template reference is computed but included for completeness
	TemplateId     types.String `tfsdk:"template_id"`
	TemplateFamily types.String `tfsdk:"template_family"`
	// Windows Autopilot Device Preparation specific fields
	DeviceSecurityGroup types.String             `tfsdk:"device_security_group"`
	DeploymentSettings  *DeploymentSettingsModel `tfsdk:"deployment_settings"`
	OOBESettings        *OOBESettingsModel       `tfsdk:"oobe_settings"`
	AllowedApps         []AllowedAppModel        `tfsdk:"allowed_apps"`
	AllowedScripts      []types.String           `tfsdk:"allowed_scripts"`
	// Assignment and timeouts
	Assignments *WindowsAutopilotDevicePreparationAssignment `tfsdk:"assignments"`
	Timeouts    timeouts.Value                               `tfsdk:"timeouts"`
}

// DeploymentSettingsModel represents the deployment settings for a Windows Autopilot Device Preparation policy
type DeploymentSettingsModel struct {
	DeploymentMode types.String `tfsdk:"deployment_mode"`
	DeploymentType types.String `tfsdk:"deployment_type"`
	JoinType       types.String `tfsdk:"join_type"`
	AccountType    types.String `tfsdk:"account_type"`
}

// OOBESettingsModel represents the out-of-box experience settings for a Windows Autopilot Device Preparation policy
type OOBESettingsModel struct {
	TimeoutInMinutes   types.Int64  `tfsdk:"timeout_in_minutes"`
	CustomErrorMessage types.String `tfsdk:"custom_error_message"`
	AllowSkip          types.Bool   `tfsdk:"allow_skip"`
	AllowDiagnostics   types.Bool   `tfsdk:"allow_diagnostics"`
}

// WindowsAutopilotDevicePreparationAssignment represents the assignment settings for a Windows Autopilot Device Preparation policy
type WindowsAutopilotDevicePreparationAssignment struct {
	IncludeGroupIds []types.String `tfsdk:"include_group_ids"`
}

// AllowedAppModel represents an application that is allowed to be installed during Windows Autopilot Device Preparation
type AllowedAppModel struct {
	AppID   types.String `tfsdk:"app_id"`
	AppType types.String `tfsdk:"app_type"`
}
