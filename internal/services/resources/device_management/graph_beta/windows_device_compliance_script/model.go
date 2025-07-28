// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicecompliancescript?view=graph-rest-beta
package graphBetaWindowsDeviceComplianceScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceComplianceScriptResourceModel defines the schema for a Device Compliance Script.
type DeviceComplianceScriptResourceModel struct {
	ID                     types.String   `tfsdk:"id"`
	DisplayName            types.String   `tfsdk:"display_name"`
	Description            types.String   `tfsdk:"description"`
	Publisher              types.String   `tfsdk:"publisher"`
	Version                types.String   `tfsdk:"version"`
	RunAs32Bit             types.Bool     `tfsdk:"run_as_32_bit"`
	RunAsAccount           types.String   `tfsdk:"run_as_account"`
	EnforceSignatureCheck  types.Bool     `tfsdk:"enforce_signature_check"`
	DetectionScriptContent types.String   `tfsdk:"detection_script_content"`
	RoleScopeTagIds        types.Set      `tfsdk:"role_scope_tag_ids"`
	CreatedDateTime        types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime   types.String   `tfsdk:"last_modified_date_time"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}
