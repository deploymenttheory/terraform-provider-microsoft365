package graphBetaMacOSCustomAttributeScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCustomAttributeShellScriptResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	CustomAttributeType  types.String   `tfsdk:"custom_attribute_type"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	ScriptContent        types.String   `tfsdk:"script_content"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	RunAsAccount         types.String   `tfsdk:"run_as_account"`
	FileName             types.String   `tfsdk:"file_name"`
	RoleScopeTagIds      types.Set      `tfsdk:"role_scope_tag_ids"`
	Assignments          types.Set      `tfsdk:"assignments"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
