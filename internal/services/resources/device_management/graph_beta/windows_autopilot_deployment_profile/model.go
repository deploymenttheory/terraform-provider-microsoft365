// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-activedirectorywindowsautopilotdeploymentprofile?view=graph-rest-beta
package graphBetaWindowsAutopilotDeploymentProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsAutopilotDeploymentProfileResourceModel represents the values for Windows Autopilot deployment profiles
type WindowsAutopilotDeploymentProfileResourceModel struct {
	ID                                     types.String                                `tfsdk:"id"`
	DisplayName                            types.String                                `tfsdk:"display_name"`
	Description                            types.String                                `tfsdk:"description"`
	Language                               types.String                                `tfsdk:"language"`
	Locale                                 types.String                                `tfsdk:"locale"`
	CreatedDateTime                        types.String                                `tfsdk:"created_date_time"`
	LastModifiedDateTime                   types.String                                `tfsdk:"last_modified_date_time"`
	DeviceJoinType                         types.String                                `tfsdk:"device_join_type"`
	OutOfBoxExperienceSetting              *OutOfBoxExperienceSettingModel             `tfsdk:"out_of_box_experience_setting"`
	EnrollmentStatusScreenSettings         *WindowsEnrollmentStatusScreenSettingsModel `tfsdk:"enrollment_status_screen_settings"`
	HardwareHashExtractionEnabled          types.Bool                                  `tfsdk:"hardware_hash_extraction_enabled"`
	DeviceNameTemplate                     types.String                                `tfsdk:"device_name_template"`
	DeviceType                             types.String                                `tfsdk:"device_type"`
	PreprovisioningAllowed                 types.Bool                                  `tfsdk:"preprovisioning_allowed"`
	RoleScopeTagIds                        types.Set                                   `tfsdk:"role_scope_tag_ids"`
	ManagementServiceAppId                 types.String                                `tfsdk:"management_service_app_id"`
	HybridAzureADJoinSkipConnectivityCheck types.Bool                                  `tfsdk:"hybrid_azure_ad_join_skip_connectivity_check"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// OutOfBoxExperienceSettingModel represents the current out-of-box experience settings
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-outofboxexperiencesetting?view=graph-rest-beta
type OutOfBoxExperienceSettingModel struct {
	PrivacySettingsHidden        types.Bool   `tfsdk:"privacy_settings_hidden"`
	EulaHidden                   types.Bool   `tfsdk:"eula_hidden"`
	UserType                     types.String `tfsdk:"user_type"`
	DeviceUsageType              types.String `tfsdk:"device_usage_type"`
	KeyboardSelectionPageSkipped types.Bool   `tfsdk:"keyboard_selection_page_skipped"`
	EscapeLinkHidden             types.Bool   `tfsdk:"escape_link_hidden"`
}

// WindowsEnrollmentStatusScreenSettingsModel represents enrollment status screen settings
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-windowsenrollmentstatusscreensettings?view=graph-rest-beta
type WindowsEnrollmentStatusScreenSettingsModel struct {
	HideInstallationProgress                         types.Bool   `tfsdk:"hide_installation_progress"`
	AllowDeviceUseBeforeProfileAndAppInstallComplete types.Bool   `tfsdk:"allow_device_use_before_profile_and_app_install_complete"`
	BlockDeviceSetupRetryByUser                      types.Bool   `tfsdk:"block_device_setup_retry_by_user"`
	AllowLogCollectionOnInstallFailure               types.Bool   `tfsdk:"allow_log_collection_on_install_failure"`
	CustomErrorMessage                               types.String `tfsdk:"custom_error_message"`
	InstallProgressTimeoutInMinutes                  types.Int32  `tfsdk:"install_progress_timeout_in_minutes"`
	AllowDeviceUseOnInstallFailure                   types.Bool   `tfsdk:"allow_device_use_on_install_failure"`
}
