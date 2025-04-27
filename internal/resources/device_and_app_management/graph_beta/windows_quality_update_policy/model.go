package graphBetaWindowsQualityUpdatePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsQualityUpdatePolicyResourceModel struct {
	ID                           types.String                           `tfsdk:"id"`
	DisplayName                  types.String                           `tfsdk:"display_name"`
	Description                  types.String                           `tfsdk:"description"`
	QualityUpdateType            types.String                           `tfsdk:"quality_update_type"`
	ExpeditedUpdateSettings      *ExpeditedWindowsQualityUpdateSettings `tfsdk:"expedited_update_settings"`
	CreatedDateTime              types.String                           `tfsdk:"created_date_time"`
	LastModifiedDateTime         types.String                           `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds              types.Set                              `tfsdk:"role_scope_tag_ids"`
	ReleaseDateDisplayName       types.String                           `tfsdk:"release_date_display_name"`
	DeployableContentDisplayName types.String                           `tfsdk:"deployable_content_display_name"`
	HotpatchEnabled              types.Bool                             `tfsdk:"hotpatch_enabled"`
	Timeouts                     timeouts.Value                         `tfsdk:"timeouts"`
}

type ExpeditedWindowsQualityUpdateSettings struct {
	QualityUpdateRelease  types.String `tfsdk:"quality_update_release"`
	DaysUntilForcedReboot types.Int32  `tfsdk:"days_until_forced_reboot"`
}
