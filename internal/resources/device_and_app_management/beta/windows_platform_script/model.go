// https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementscript?view=graph-rest-beta

package graphBetaWindowsPlatformScript

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsPlatformScriptResourceModel struct {
	ID                    types.String                                                `tfsdk:"id"`
	DisplayName           types.String                                                `tfsdk:"display_name"`
	Description           types.String                                                `tfsdk:"description"`
	ScriptContent         types.String                                                `tfsdk:"script_content"`
	RunAsAccount          types.String                                                `tfsdk:"run_as_account"`
	EnforceSignatureCheck types.Bool                                                  `tfsdk:"enforce_signature_check"`
	FileName              types.String                                                `tfsdk:"file_name"`
	RoleScopeTagIds       []types.String                                              `tfsdk:"role_scope_tag_ids"`
	RunAs32Bit            types.Bool                                                  `tfsdk:"run_as_32_bit"`
	Assignments           *sharedmodels.DeviceManagementScriptAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts              timeouts.Value                                              `tfsdk:"timeouts"`
}
