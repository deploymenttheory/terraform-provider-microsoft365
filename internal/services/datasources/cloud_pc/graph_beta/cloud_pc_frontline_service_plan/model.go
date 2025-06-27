// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcfrontlineserviceplan?view=graph-rest-beta

package cloudPcFrontlineServicePlan

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPcFrontlineServicePlanDataSourceModel represents the Terraform data source model for frontline service plans
type CloudPcFrontlineServicePlanDataSourceModel struct {
	FilterType  types.String                           `tfsdk:"filter_type"`
	FilterValue types.String                           `tfsdk:"filter_value"`
	Items       []CloudPcFrontlineServicePlanItemModel `tfsdk:"items"`
	Timeouts    timeouts.Value                         `tfsdk:"timeouts"`
}

// CloudPcFrontlineServicePlanItemModel represents an individual frontline service plan
type CloudPcFrontlineServicePlanItemModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	TotalCount  types.Int64  `tfsdk:"total_count"`
	UsedCount   types.Int64  `tfsdk:"used_count"`
}
