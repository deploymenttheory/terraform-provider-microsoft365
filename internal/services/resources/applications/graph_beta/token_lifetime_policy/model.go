// REF: https://learn.microsoft.com/en-us/graph/api/resources/tokenlifetimepolicy?view=graph-rest-beta
package graphBetaApplicationsTokenLifetimePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TokenLifetimePolicyResourceModel represents the schema for the Token Lifetime Policy resource
type TokenLifetimePolicyResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	DisplayName           types.String   `tfsdk:"display_name"`
	Description           types.String   `tfsdk:"description"`
	Definition            types.List     `tfsdk:"definition"`
	IsOrganizationDefault types.Bool     `tfsdk:"is_organization_default"`
	DeletedDateTime       types.String   `tfsdk:"deleted_date_time"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
