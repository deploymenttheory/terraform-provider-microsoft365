// REF: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
package graphBetaNetworkFilteringProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkFilteringProfileResourceModel represents the schema for the Filtering Profile resource
type NetworkFilteringProfileResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	Priority             types.Int64    `tfsdk:"priority"`
	State                types.String   `tfsdk:"state"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Version              types.String   `tfsdk:"version"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
