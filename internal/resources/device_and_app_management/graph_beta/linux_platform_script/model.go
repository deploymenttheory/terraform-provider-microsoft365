package graphBetaLinuxPlatformScript

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LinuxPlatformScriptResourceModel represents the Terraform resource schema for Linux script configurations
type LinuxPlatformScriptResourceModel struct {
	ID              types.String                                                 `tfsdk:"id"`
	Name            types.String                                                 `tfsdk:"name"`
	Description     types.String                                                 `tfsdk:"description"`
	Platforms       types.String                                                 `tfsdk:"platforms"`    // Always "linux"
	Technologies    types.List                                                   `tfsdk:"technologies"` // Always ["linuxMdm"]
	RoleScopeTagIds types.List                                                   `tfsdk:"role_scope_tag_ids"`
	ScriptContent   types.String                                                 `tfsdk:"script_content"`
	Settings        *sharedmodels.DeviceConfigV2GraphServiceResourceModel        `tfsdk:"settings"`
	Assignments     *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts        timeouts.Value                                               `tfsdk:"timeouts"`
}
