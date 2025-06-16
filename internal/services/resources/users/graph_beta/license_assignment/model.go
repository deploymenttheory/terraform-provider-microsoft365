// REF: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta&tabs=http

package graphBetaUserLicenseAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserLicenseAssignmentResourceModel represents the Terraform resource model for user license assignment
type UserLicenseAssignmentResourceModel struct {
	ID                types.String                   `tfsdk:"id"`
	UserId            types.String                   `tfsdk:"user_id"`
	UserPrincipalName types.String                   `tfsdk:"user_principal_name"`
	AddLicenses       []AssignedLicenseResourceModel `tfsdk:"add_licenses"`
	RemoveLicenses    types.Set                      `tfsdk:"remove_licenses"`
	AssignedLicenses  types.List                     `tfsdk:"assigned_licenses"`
	Timeouts          timeouts.Value                 `tfsdk:"timeouts"`
}

// AssignedLicenseResourceModel represents a license to be assigned to the user
type AssignedLicenseResourceModel struct {
	SkuId         types.String `tfsdk:"sku_id"`
	DisabledPlans types.Set    `tfsdk:"disabled_plans"`
}

// LicenseDetailsResourceModel represents the current license assignment state (read-only)
type LicenseDetailsResourceModel struct {
	SkuId         types.String `tfsdk:"sku_id"`
	SkuPartNumber types.String `tfsdk:"sku_part_number"`
	ServicePlans  types.List   `tfsdk:"service_plans"`
}

// ServicePlanResourceModel represents individual service plans within a license
type ServicePlanResourceModel struct {
	ServicePlanId      types.String `tfsdk:"service_plan_id"`
	ServicePlanName    types.String `tfsdk:"service_plan_name"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	AppliesTo          types.String `tfsdk:"applies_to"`
}
