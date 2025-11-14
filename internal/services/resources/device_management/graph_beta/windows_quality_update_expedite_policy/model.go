package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsQualityUpdateExpeditePolicyResourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	DisplayName                  types.String   `tfsdk:"display_name"`
	Description                  types.String   `tfsdk:"description"`
	ExpeditedUpdateSettings      types.Object   `tfsdk:"expedited_update_settings"`
	CreatedDateTime              types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime         types.String   `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds              types.Set      `tfsdk:"role_scope_tag_ids"`
	ReleaseDateDisplayName       types.String   `tfsdk:"release_date_display_name"`
	DeployableContentDisplayName types.String   `tfsdk:"deployable_content_display_name"`
	Assignments                  types.Set      `tfsdk:"assignments"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}
