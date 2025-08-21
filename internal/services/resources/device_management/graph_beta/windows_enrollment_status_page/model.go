// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-windows10enrollmentcompletionpageconfiguration?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/intune-onboarding-windows10enrollmentcompletionpageconfiguration-create?view=graph-rest-beta
package graphBetaWindowsEnrollmentStatusPage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsEnrollmentStatusPageResourceModel struct {
	ID                                      types.String   `tfsdk:"id"`
	DisplayName                             types.String   `tfsdk:"display_name"`
	Description                             types.String   `tfsdk:"description"`
	Priority                                types.Int32    `tfsdk:"priority"`
	ShowInstallationProgress                types.Bool     `tfsdk:"show_installation_progress"`
	BlockDeviceSetupRetryByUser             types.Bool     `tfsdk:"block_device_setup_retry_by_user"`
	AllowDeviceResetOnInstallFailure        types.Bool     `tfsdk:"allow_device_reset_on_install_failure"`
	AllowLogCollectionOnInstallFailure      types.Bool     `tfsdk:"allow_log_collection_on_install_failure"`
	CustomErrorMessage                      types.String   `tfsdk:"custom_error_message"`
	InstallProgressTimeoutInMinutes         types.Int32    `tfsdk:"install_progress_timeout_in_minutes"`
	AllowDeviceUseOnInstallFailure          types.Bool     `tfsdk:"allow_device_use_on_install_failure"`
	SelectedMobileAppIds                    types.Set      `tfsdk:"selected_mobile_app_ids"`
	TrackInstallProgressForAutopilotOnly    types.Bool     `tfsdk:"track_install_progress_for_autopilot_only"`
	DisableUserStatusTrackingAfterFirstUser types.Bool     `tfsdk:"disable_user_status_tracking_after_first_user"`
	RoleScopeTagIds                         types.Set      `tfsdk:"role_scope_tag_ids"`
	Assignments                             types.Set      `tfsdk:"assignments"`
	Timeouts                                timeouts.Value `tfsdk:"timeouts"`
}
