// REF: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta&tabs=http

package graphBetaUserLicenseAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserLicenseAssignmentResourceModel represents the Terraform resource model for user license assignment
type UserLicenseAssignmentResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	UserId            types.String   `tfsdk:"user_id"`
	UserPrincipalName types.String   `tfsdk:"user_principal_name"`
	SkuId             types.String   `tfsdk:"sku_id"`
	DisabledPlans     types.Set      `tfsdk:"disabled_plans"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}
