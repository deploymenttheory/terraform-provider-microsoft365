// REF: https://learn.microsoft.com/en-us/graph/api/application-list-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/application-post-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/application-delete-owners?view=graph-rest-beta
package graphBetaApplicationOwner

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicationOwnerResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	ApplicationID    types.String   `tfsdk:"application_id"`
	OwnerID          types.String   `tfsdk:"owner_id"`
	OwnerObjectType  types.String   `tfsdk:"owner_object_type"`
	OwnerType        types.String   `tfsdk:"owner_type"`
	OwnerDisplayName types.String   `tfsdk:"owner_display_name"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}
