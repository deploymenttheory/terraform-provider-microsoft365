// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpc?view=graph-rest-beta

package graphBetaCloudPC

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPCDataSourceModel defines the data source model
type CloudPCDataSourceModel struct {
	FilterType   types.String                 `tfsdk:"filter_type"`   // Required field to specify how to filter
	FilterValue  types.String                 `tfsdk:"filter_value"`  // Value to filter by (not used for "all" or "odata")
	ODataFilter  types.String                 `tfsdk:"odata_filter"`  // OData filter parameter
	ODataTop     types.Int32                  `tfsdk:"odata_top"`     // OData top parameter for limiting results
	ODataSkip    types.Int32                  `tfsdk:"odata_skip"`    // OData skip parameter for pagination
	ODataSelect  types.String                 `tfsdk:"odata_select"`  // OData select parameter for field selection
	ODataOrderBy types.String                 `tfsdk:"odata_orderby"` // OData orderby parameter for sorting
	Items        []CloudPCItemDataSourceModel `tfsdk:"items"`         // List of cloud PCs that match the filters
	Timeouts     timeouts.Value               `tfsdk:"timeouts"`
}

// CloudPCItemDataSourceModel represents a single cloud PC
type CloudPCItemDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`                          // The unique identifier for the cloud PC
	AadDeviceId              types.String `tfsdk:"aad_device_id"`               // Azure AD device ID of the cloud PC
	DisplayName              types.String `tfsdk:"display_name"`                // The display name of the cloud PC
	ImageDisplayName         types.String `tfsdk:"image_display_name"`          // Name of the OS image that's on the cloud PC
	ManagedDeviceId          types.String `tfsdk:"managed_device_id"`           // The Intune managed device ID of the cloud PC
	ManagedDeviceName        types.String `tfsdk:"managed_device_name"`         // The Intune managed device name of the cloud PC
	ProvisioningPolicyId     types.String `tfsdk:"provisioning_policy_id"`      // The provisioning policy ID of the cloud PC
	ProvisioningPolicyName   types.String `tfsdk:"provisioning_policy_name"`    // The provisioning policy name of the cloud PC
	OnPremisesConnectionName types.String `tfsdk:"on_premises_connection_name"` // The Azure network connection that is applied during the provisioning of cloud PCs
	ServicePlanId            types.String `tfsdk:"service_plan_id"`             // The service plan ID of the cloud PC
	ServicePlanName          types.String `tfsdk:"service_plan_name"`           // The service plan name of the cloud PC
	ServicePlanType          types.String `tfsdk:"service_plan_type"`           // The service plan type of the cloud PC (enterprise, business, unknownFutureValue)
	Status                   types.String `tfsdk:"status"`                      // The status of the cloud PC
	UserPrincipalName        types.String `tfsdk:"user_principal_name"`         // The user principal name (UPN) of the user assigned to the cloud PC
	LastModifiedDateTime     types.String `tfsdk:"last_modified_date_time"`     // The last modified date and time of the cloud PC
	StatusDetails            types.String `tfsdk:"status_details"`              // The details of the cloud PC status
	GracePeriodEndDateTime   types.String `tfsdk:"grace_period_end_date_time"`  // The date and time when the grace period ends and reprovisioning/deprovisioning happens
	ProvisioningType         types.String `tfsdk:"provisioning_type"`           // Specifies the type of license used when provisioning cloud PCs using this policy
	DeviceRegionName         types.String `tfsdk:"device_region_name"`          // The name of the geographical region where the cloud PC is currently provisioned
	DiskEncryptionState      types.String `tfsdk:"disk_encryption_state"`       // The disk encryption applied to the cloud PC
}
