// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-owners?view=graph-rest-beta
package graphBetaServicePrincipalOwner

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServicePrincipalOwnerResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	ServicePrincipalID   types.String   `tfsdk:"service_principal_id"`
	OwnerID              types.String   `tfsdk:"owner_id"`
	OwnerObjectType      types.String   `tfsdk:"owner_object_type"`
	OwnerType            types.String   `tfsdk:"owner_type"`
	OwnerDisplayName     types.String   `tfsdk:"owner_display_name"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
