package graphBetaWindowsQualityUpdatePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsQualityUpdatePolicyResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds      types.Set      `tfsdk:"role_scope_tag_ids"`
	HotpatchEnabled      types.Bool     `tfsdk:"hotpatch_enabled"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
