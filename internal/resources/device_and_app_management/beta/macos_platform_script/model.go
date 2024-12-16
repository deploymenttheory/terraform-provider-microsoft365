// https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-deviceshellScript?view=graph-rest-beta

package graphBetaMacOSPlatformScript

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MacOSPlatformScriptResourceModel struct {
	ID                          types.String                                                `tfsdk:"id"`
	DisplayName                 types.String                                                `tfsdk:"display_name"`
	Description                 types.String                                                `tfsdk:"description"`
	ScriptContent               types.String                                                `tfsdk:"script_content"`
	CreatedDateTime             types.String                                                `tfsdk:"created_date_time"`
	LastModifiedDateTime        types.String                                                `tfsdk:"last_modified_date_time"`
	RunAsAccount                types.String                                                `tfsdk:"run_as_account"`
	FileName                    types.String                                                `tfsdk:"file_name"`
	RoleScopeTagIds             []types.String                                              `tfsdk:"role_scope_tag_ids"`
	BlockExecutionNotifications types.Bool                                                  `tfsdk:"block_execution_notifications"`
	ExecutionFrequency          types.String                                                `tfsdk:"execution_frequency"`
	RetryCount                  types.Int32                                                 `tfsdk:"retry_count"`
	Assignments                 *sharedmodels.DeviceManagementScriptAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts                    timeouts.Value                                              `tfsdk:"timeouts"`
}
