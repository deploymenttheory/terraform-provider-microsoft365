// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-approleassignedto?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-approleassignedto?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-approleassignedto?view=graph-rest-beta
package graphBetaServicePrincipalAppRoleAssignedTo

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalAppRoleAssignedToResourceModel represents the Terraform resource model for app role assignments
type ServicePrincipalAppRoleAssignedToResourceModel struct {
	ID                             types.String   `tfsdk:"id"`
	ResourceObjectID               types.String   `tfsdk:"resource_object_id"`
	AppRoleID                      types.String   `tfsdk:"app_role_id"`
	TargetServicePrincipalObjectID types.String   `tfsdk:"target_service_principal_object_id"`
	PrincipalType                  types.String   `tfsdk:"principal_type"`
	PrincipalDisplayName           types.String   `tfsdk:"principal_display_name"`
	ResourceDisplayName            types.String   `tfsdk:"resource_display_name"`
	CreatedDateTime                types.String   `tfsdk:"created_date_time"`
	Timeouts                       timeouts.Value `tfsdk:"timeouts"`
}
