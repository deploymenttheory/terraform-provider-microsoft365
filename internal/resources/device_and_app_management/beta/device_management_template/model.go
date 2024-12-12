package graphBetaDeviceManagementTemplate

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceManagementTemplateProfileResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type DeviceManagementTemplateResourceModel struct {
	ID                   types.String                                                 `tfsdk:"id"`
	Name                 types.String                                                 `tfsdk:"name"`
	Description          types.String                                                 `tfsdk:"description"`
	Platforms            types.String                                                 `tfsdk:"platforms"`
	Technologies         []types.String                                               `tfsdk:"technologies"`
	RoleScopeTagIds      []types.String                                               `tfsdk:"role_scope_tag_ids"`
	SettingsCount        types.Int64                                                  `tfsdk:"settings_count"`
	IsAssigned           types.Bool                                                   `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String                                                 `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String                                                 `tfsdk:"created_date_time"`
	Settings             types.String                                                 `tfsdk:"settings"`
	Assignments          *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts             timeouts.Value                                               `tfsdk:"timeouts"`
}
