// REF: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringpolicy?view=graph-rest-beta
package graphBetaFilteringPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FilteringPolicyResourceModel represents the schema for the Filtering Policy resource
//
// Note: The Microsoft Graph API documentation for filtering policies mentions 'state' and 'priority'
// properties in the update documentation. However, these properties are actually used when linking
// policies to security profiles, not as direct properties of the filtering policy resource itself.
// This appears to be a documentation error on Microsoft's side.
//
// Additionally, the 'version' property is not documented in the Microsoft Graph API documentation,
// but it is actually included in the API responses.
type FilteringPolicyResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Version              types.String   `tfsdk:"version"`
	Action               types.String   `tfsdk:"action"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
