// REF: https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta&tabs=http

package graphBetaGroupLicenseAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GroupLicenseAssignmentResourceModel represents the Terraform resource model for group license assignment
type GroupLicenseAssignmentResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	GroupId       types.String   `tfsdk:"group_id"`
	DisplayName   types.String   `tfsdk:"display_name"`
	SkuId         types.String   `tfsdk:"sku_id"`
	DisabledPlans types.Set      `tfsdk:"disabled_plans"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}
