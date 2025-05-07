package graphBetaLinuxPlatformScript

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LinuxPlatformScriptResourceModel represents the Terraform resource schema for Linux script configurations
type LinuxPlatformScriptResourceModel struct {
	ID                 types.String                                                 `tfsdk:"id"`
	Name               types.String                                                 `tfsdk:"name"`
	Description        types.String                                                 `tfsdk:"description"`
	Platforms          types.String                                                 `tfsdk:"platforms"`
	Technologies       []types.String                                               `tfsdk:"technologies"`
	RoleScopeTagIds    types.Set                                                    `tfsdk:"role_scope_tag_ids"`
	ScriptContent      types.String                                                 `tfsdk:"script_content"`
	ExecutionContext   types.String                                                 `tfsdk:"execution_context"`
	ExecutionFrequency types.String                                                 `tfsdk:"execution_frequency"`
	ExecutionRetries   types.String                                                 `tfsdk:"execution_retries"`
	Assignments        *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts           timeouts.Value                                               `tfsdk:"timeouts"`
}
