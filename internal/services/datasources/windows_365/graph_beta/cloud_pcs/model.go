// REF: https://learn.microsoft.com/en-us/graph/api/virtualendpoint-list-cloudpcs?view=graph-rest-beta

package graphBetaCloudPcs

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPcsDataSourceModel represents the Terraform data source model for Cloud PCs
type CloudPcsDataSourceModel struct {
	FilterType  types.String       `tfsdk:"filter_type"`
	FilterValue types.String       `tfsdk:"filter_value"`
	ODataFilter types.String       `tfsdk:"odata_filter"`
	ODataSelect types.String       `tfsdk:"odata_select"`
	ODataTop    types.Int64        `tfsdk:"odata_top"`
	ODataCount  types.Bool         `tfsdk:"odata_count"`
	Items       []CloudPcItemModel `tfsdk:"items"`
	Timeouts    timeouts.Value     `tfsdk:"timeouts"`
}

// CloudPcItemModel represents an individual Cloud PC
type CloudPcItemModel struct {
	ID                       types.String `tfsdk:"id"`
	DisplayName              types.String `tfsdk:"display_name"`
	AADDeviceID              types.String `tfsdk:"aad_device_id"`
	ImageDisplayName         types.String `tfsdk:"image_display_name"`
	ManagedDeviceID          types.String `tfsdk:"managed_device_id"`
	ManagedDeviceName        types.String `tfsdk:"managed_device_name"`
	ProvisioningPolicyID     types.String `tfsdk:"provisioning_policy_id"`
	ProvisioningPolicyName   types.String `tfsdk:"provisioning_policy_name"`
	OnPremisesConnectionName types.String `tfsdk:"on_premises_connection_name"`
	ServicePlanID            types.String `tfsdk:"service_plan_id"`
	ServicePlanName          types.String `tfsdk:"service_plan_name"`
	ServicePlanType          types.String `tfsdk:"service_plan_type"`
	Status                   types.String `tfsdk:"status"`
	UserPrincipalName        types.String `tfsdk:"user_principal_name"`
	LastModifiedDateTime     types.String `tfsdk:"last_modified_date_time"`
	StatusDetailCode         types.String `tfsdk:"status_detail_code"`
	StatusDetailMessage      types.String `tfsdk:"status_detail_message"`
	GracePeriodEndDateTime   types.String `tfsdk:"grace_period_end_date_time"`
	ProvisioningType         types.String `tfsdk:"provisioning_type"`
	DeviceRegionName         types.String `tfsdk:"device_region_name"`
	DiskEncryptionState      types.String `tfsdk:"disk_encryption_state"`
	ProductType              types.String `tfsdk:"product_type"`
	UserAccountType          types.String `tfsdk:"user_account_type"`
	EnableSingleSignOn       types.Bool   `tfsdk:"enable_single_sign_on"`
}
