// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementreusablepolicysetting?view=graph-rest-beta
package graphBetaLinuxDeviceComplianceScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LinuxDeviceComplianceScriptResourceModel defines the schema for a Linux Device Compliance Script.
type LinuxDeviceComplianceScriptResourceModel struct {
	ID                     types.String   `tfsdk:"id"`
	DisplayName            types.String   `tfsdk:"display_name"`
	Description            types.String   `tfsdk:"description"`
	DetectionScriptContent types.String   `tfsdk:"detection_script_content"`
	SettingDefinitionId    types.String   `tfsdk:"setting_definition_id"`
	Version                types.Int32    `tfsdk:"version"`
	LastModifiedDateTime   types.String   `tfsdk:"last_modified_date_time"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}
