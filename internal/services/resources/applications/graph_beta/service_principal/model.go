package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalResourceModel describes the Terraform resource data model for a Service Principal
type ServicePrincipalResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	AppID                     types.String   `tfsdk:"app_id"`
	DisplayName               types.String   `tfsdk:"display_name"`
	AccountEnabled            types.Bool     `tfsdk:"account_enabled"`
	AppRoleAssignmentRequired types.Bool     `tfsdk:"app_role_assignment_required"`
	ServicePrincipalType      types.String   `tfsdk:"service_principal_type"`
	ServicePrincipalNames     types.Set      `tfsdk:"service_principal_names"`
	SignInAudience            types.String   `tfsdk:"sign_in_audience"`
	Tags                      types.Set      `tfsdk:"tags"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
}
