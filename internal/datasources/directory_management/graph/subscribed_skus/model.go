// REF: https://learn.microsoft.com/en-us/graph/api/subscribedsku-list?view=graph-rest-1.0&tabs=http

package graphSubscribedSkus

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SubscribedSkusDataSourceModel represents the Terraform data source model for subscribed SKUs
type SubscribedSkusDataSourceModel struct {
	ID             types.String   `tfsdk:"id"`
	SkuId          types.String   `tfsdk:"sku_id"`          // Filter by specific SKU ID
	SkuPartNumber  types.String   `tfsdk:"sku_part_number"` // Filter by SKU part number
	AppliesTo      types.String   `tfsdk:"applies_to"`      // Filter by applies to (User/Company)
	SubscribedSkus types.List     `tfsdk:"subscribed_skus"` // List of all SKUs
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

// SubscribedSkuModel represents an individual subscribed SKU
type SubscribedSkuModel struct {
	ID               types.String `tfsdk:"id"`
	AccountId        types.String `tfsdk:"account_id"`
	AccountName      types.String `tfsdk:"account_name"`
	AppliesTo        types.String `tfsdk:"applies_to"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	ConsumedUnits    types.Int64  `tfsdk:"consumed_units"`
	SkuId            types.String `tfsdk:"sku_id"`
	SkuPartNumber    types.String `tfsdk:"sku_part_number"`
	PrepaidUnits     types.Object `tfsdk:"prepaid_units"`
	ServicePlans     types.List   `tfsdk:"service_plans"`
	SubscriptionIds  types.List   `tfsdk:"subscription_ids"`
}

// LicenseUnitsDetailModel represents the prepaid units information
type LicenseUnitsDetailModel struct {
	Enabled   types.Int64 `tfsdk:"enabled"`
	LockedOut types.Int64 `tfsdk:"locked_out"`
	Suspended types.Int64 `tfsdk:"suspended"`
	Warning   types.Int64 `tfsdk:"warning"`
}

// ServicePlanInfoModel represents individual service plans within a SKU
type ServicePlanInfoModel struct {
	ServicePlanId      types.String `tfsdk:"service_plan_id"`
	ServicePlanName    types.String `tfsdk:"service_plan_name"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	AppliesTo          types.String `tfsdk:"applies_to"`
}
