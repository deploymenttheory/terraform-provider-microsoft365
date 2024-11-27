package graphBetaLinuxPlatformScript

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinuxPlatformScriptResourceModel struct {
	ID                  types.String                                                 `tfsdk:"id"`
	DisplayName         types.String                                                 `tfsdk:"display_name"`
	Description         types.String                                                 `tfsdk:"description"`
	ScriptContent       types.String                                                 `tfsdk:"script_content"`
	RoleScopeTagIds     []types.String                                               `tfsdk:"role_scope_tag_ids"`
	Platforms           types.String                                                 `tfsdk:"platforms"`             // e.g., "LINUX"
	Technologies        types.String                                                 `tfsdk:"technologies"`          // e.g., "LINUXMDM"
	Settings            []LinuxPlatformScriptConfigurationSettingResourceModel       `tfsdk:"settings"`              // Nested settings list
	TemplateReferenceID types.String                                                 `tfsdk:"template_reference_id"` // e.g., "92439f26-2b30-4503-8429-6d40f7e172dd_1"
	Assignments         *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts            timeouts.Value                                               `tfsdk:"timeouts"`
}

type LinuxPlatformScriptConfigurationSettingResourceModel struct {
	SettingDefinitionID           types.String                                           `tfsdk:"setting_definition_id"`
	SettingValue                  types.String                                           `tfsdk:"setting_value"`                    // Value for simple/choice settings
	SettingValueTemplateReference types.String                                           `tfsdk:"setting_value_template_reference"` // Template reference ID
	SettingInstanceTemplateID     types.String                                           `tfsdk:"setting_instance_template_id"`     // Instance template ID
	Children                      []LinuxPlatformScriptConfigurationSettingResourceModel `tfsdk:"children"`                         // For nested settings
}
