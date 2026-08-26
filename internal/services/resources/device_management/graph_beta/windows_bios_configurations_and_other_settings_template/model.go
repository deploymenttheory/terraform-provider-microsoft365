// https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-hardwareconfiguration?view=graph-rest-beta

package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel struct {
	ID                          types.String   `tfsdk:"id"`
	DisplayName                 types.String   `tfsdk:"display_name"`
	Description                 types.String   `tfsdk:"description"`
	FileName                    types.String   `tfsdk:"file_name"`
	ConfigurationFileContent    types.String   `tfsdk:"configuration_file_content"`
	HardwareConfigurationFormat types.String   `tfsdk:"hardware_configuration_format"`
	PerDevicePasswordDisabled   types.Bool     `tfsdk:"per_device_password_disabled"`
	RoleScopeTagIds             types.Set      `tfsdk:"role_scope_tag_ids"`
	Version                     types.Int32    `tfsdk:"version"`
	CreatedDateTime             types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime        types.String   `tfsdk:"last_modified_date_time"`
	Assignments                 types.Set      `tfsdk:"assignments"`
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
}
