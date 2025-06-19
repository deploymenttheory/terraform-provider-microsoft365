// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicy?view=graph-rest-1.0
package graphCloudPcProvisioningPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcProvisioningPolicyResourceModel struct {
	ID                       types.String                   `tfsdk:"id"`
	AlternateResourceUrl     types.String                   `tfsdk:"alternate_resource_url"`
	CloudPcGroupDisplayName  types.String                   `tfsdk:"cloud_pc_group_display_name"`
	CloudPcNamingTemplate    types.String                   `tfsdk:"cloud_pc_naming_template"`
	Description              types.String                   `tfsdk:"description"`
	DisplayName              types.String                   `tfsdk:"display_name"`
	DomainJoinConfigurations []DomainJoinConfigurationModel `tfsdk:"domain_join_configurations"`
	EnableSingleSignOn       types.Bool                     `tfsdk:"enable_single_sign_on"`
	GracePeriodInHours       types.Int32                    `tfsdk:"grace_period_in_hours"`
	ImageDisplayName         types.String                   `tfsdk:"image_display_name"`
	ImageId                  types.String                   `tfsdk:"image_id"`
	ImageType                types.String                   `tfsdk:"image_type"`
	LocalAdminEnabled        types.Bool                     `tfsdk:"local_admin_enabled"`
	MicrosoftManagedDesktop  *MicrosoftManagedDesktopModel  `tfsdk:"microsoft_managed_desktop"`
	ProvisioningType         types.String                   `tfsdk:"provisioning_type"`
	WindowsSetting           *WindowsSettingModel           `tfsdk:"windows_setting"`
	Timeouts                 timeouts.Value                 `tfsdk:"timeouts"`
}

type DomainJoinConfigurationModel struct {
	DomainJoinType         types.String `tfsdk:"domain_join_type"`
	OnPremisesConnectionId types.String `tfsdk:"on_premises_connection_id"`
	RegionName             types.String `tfsdk:"region_name"`
}

type MicrosoftManagedDesktopModel struct {
	ManagedType types.String `tfsdk:"managed_type"`
	Profile     types.String `tfsdk:"profile"`
}

type WindowsSettingModel struct {
	Locale types.String `tfsdk:"locale"`
}
