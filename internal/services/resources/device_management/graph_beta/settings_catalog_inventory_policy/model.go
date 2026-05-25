package graphBetaSettingsCatalogInventoryPolicy

import (
	configPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InventoryPolicyResourceModel holds the configuration for a Settings Catalog Inventory Policy.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type InventoryPolicyResourceModel struct {
	ID                   types.String                                          `tfsdk:"id"`
	Name                 types.String                                          `tfsdk:"name"`
	Description          types.String                                          `tfsdk:"description"`
	Platforms            types.String                                          `tfsdk:"platforms"`
	Technologies         types.String                                          `tfsdk:"technologies"`
	RoleScopeTagIds      types.Set                                             `tfsdk:"role_scope_tag_ids"`
	SettingsCount        types.Int32                                           `tfsdk:"settings_count"`
	LastModifiedDateTime types.String                                          `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String                                          `tfsdk:"created_date_time"`
	ConfigurationPolicy  *configPolicy.DeviceConfigV2GraphServiceResourceModel `tfsdk:"configuration_policy"`
	Assignments          types.Set                                             `tfsdk:"assignments"`
	Timeouts             timeouts.Value                                        `tfsdk:"timeouts"`
}
