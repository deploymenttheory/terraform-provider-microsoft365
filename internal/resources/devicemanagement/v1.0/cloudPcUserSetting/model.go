// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcrestorepointsetting?view=graph-rest-1.0
package graphCloudPcUserSetting

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcUserSettingResourceModel struct {
	ID                   types.String                     `tfsdk:"id"`
	CreatedDateTime      types.String                     `tfsdk:"created_date_time"`
	DisplayName          types.String                     `tfsdk:"display_name"`
	LastModifiedDateTime types.String                     `tfsdk:"last_modified_date_time"`
	LocalAdminEnabled    types.Bool                       `tfsdk:"local_admin_enabled"`
	ResetEnabled         types.Bool                       `tfsdk:"reset_enabled"`
	RestorePointSetting  *CloudPcRestorePointSettingModel `tfsdk:"restore_point_setting"`
	Timeouts             timeouts.Value                   `tfsdk:"timeouts"`
}

type CloudPcRestorePointSettingModel struct {
	FrequencyType      types.String `tfsdk:"frequency_type"`
	UserRestoreEnabled types.Bool   `tfsdk:"user_restore_enabled"`
}
