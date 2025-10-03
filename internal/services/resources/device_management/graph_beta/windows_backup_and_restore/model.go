// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsrestoredeviceenrollmentconfiguration?view=graph-rest-beta
package graphBetaWindowsBackupAndRestore

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsBackupAndRestoreResourceModel struct {
	ID                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	Description                       types.String   `tfsdk:"description"`
	Priority                          types.Int32    `tfsdk:"priority"`
	CreatedDateTime                   types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime              types.String   `tfsdk:"last_modified_date_time"`
	Version                           types.Int32    `tfsdk:"version"`
	RoleScopeTagIds                   types.Set      `tfsdk:"role_scope_tag_ids"`
	DeviceEnrollmentConfigurationType types.String   `tfsdk:"device_enrollment_configuration_type"`
	State                             types.String   `tfsdk:"state"`
	Assignments                       types.Set      `tfsdk:"assignments"`
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
}
