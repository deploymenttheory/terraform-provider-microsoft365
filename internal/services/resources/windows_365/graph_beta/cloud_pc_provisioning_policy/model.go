// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicy?view=graph-rest-beta
package graphBetaCloudPcProvisioningPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcProvisioningPolicyResourceModel struct {
	ID                       types.String                               `tfsdk:"id"`
	AlternateResourceUrl     types.String                               `tfsdk:"alternate_resource_url"`
	CloudPcGroupDisplayName  types.String                               `tfsdk:"cloud_pc_group_display_name"`
	CloudPcNamingTemplate    types.String                               `tfsdk:"cloud_pc_naming_template"`
	Description              types.String                               `tfsdk:"description"`
	DisplayName              types.String                               `tfsdk:"display_name"`
	DomainJoinConfigurations []DomainJoinConfigurationModel             `tfsdk:"domain_join_configurations"`
	EnableSingleSignOn       types.Bool                                 `tfsdk:"enable_single_sign_on"`
	GracePeriodInHours       types.Int32                                `tfsdk:"grace_period_in_hours"`
	ImageDisplayName         types.String                               `tfsdk:"image_display_name"`
	ImageId                  types.String                               `tfsdk:"image_id"`
	ImageType                types.String                               `tfsdk:"image_type"`
	LocalAdminEnabled        types.Bool                                 `tfsdk:"local_admin_enabled"`
	MicrosoftManagedDesktop  *MicrosoftManagedDesktopModel              `tfsdk:"microsoft_managed_desktop"`
	ProvisioningType         types.String                               `tfsdk:"provisioning_type"`
	WindowsSetting           *WindowsSettingModel                       `tfsdk:"windows_setting"`
	ApplyToExistingCloudPcs  *ApplyToExistingCloudPcsModel              `tfsdk:"apply_to_existing_cloud_pcs"`
	ManagedBy                types.String                               `tfsdk:"managed_by"`
	ScopeIds                 types.Set                                  `tfsdk:"scope_ids"`
	Autopatch                *AutopatchModel                            `tfsdk:"autopatch"`
	AutopilotConfiguration   *AutopilotConfigurationModel               `tfsdk:"autopilot_configuration"`
	Assignments              []CloudPcProvisioningPolicyAssignmentModel `tfsdk:"assignments"`
	Timeouts                 timeouts.Value                             `tfsdk:"timeouts"`
}

type DomainJoinConfigurationModel struct {
	DomainJoinType         types.String `tfsdk:"domain_join_type"`
	OnPremisesConnectionId types.String `tfsdk:"on_premises_connection_id"`
	RegionName             types.String `tfsdk:"region_name"`
	RegionGroup            types.String `tfsdk:"region_group"`
}

type MicrosoftManagedDesktopModel struct {
	ManagedType types.String `tfsdk:"managed_type"`
	Profile     types.String `tfsdk:"profile"`
}

type WindowsSettingModel struct {
	Locale types.String `tfsdk:"locale"`
}

type AutopatchModel struct {
	AutopatchGroupId types.String `tfsdk:"autopatch_group_id"`
}

type AutopilotConfigurationModel struct {
	DevicePreparationProfileId  types.String `tfsdk:"device_preparation_profile_id"`
	ApplicationTimeoutInMinutes types.Int32  `tfsdk:"application_timeout_in_minutes"`
	OnFailureDeviceAccessDenied types.Bool   `tfsdk:"on_failure_device_access_denied"`
}

type ApplyToExistingCloudPcsModel struct {
	MicrosoftEntraSingleSignOnForAllDevices        types.Bool `tfsdk:"microsoft_entra_single_sign_on_for_all_devices"`
	RegionOrAzureNetworkConnectionForAllDevices    types.Bool `tfsdk:"region_or_azure_network_connection_for_all_devices"`
	RegionOrAzureNetworkConnectionForSelectDevices types.Bool `tfsdk:"region_or_azure_network_connection_for_select_devices"`
}

// CloudPcProvisioningPolicyAssignmentModel represents an assignment of a Cloud PC provisioning policy to a group
type CloudPcProvisioningPolicyAssignmentModel struct {
	ID                    types.String `tfsdk:"id"`
	GroupId               types.String `tfsdk:"group_id"`
	ServicePlanId         types.String `tfsdk:"service_plan_id"`
	AllotmentLicenseCount types.Int64  `tfsdk:"allotment_license_count"`
	AllotmentDisplayName  types.String `tfsdk:"allotment_display_name"`
}
