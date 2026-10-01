// REF: https://learn.microsoft.com/en-us/graph/api/resources/externalidentitiespolicy?view=graph-rest-beta
package graphBetaExternalIdentitiesPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExternalIdentitiesPolicyResourceModel represents the tenant-wide external identities policy.
type ExternalIdentitiesPolicyResourceModel struct {
	ID                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	Description                       types.String   `tfsdk:"description"`
	AllowExternalIdentitiesToLeave    types.Bool     `tfsdk:"allow_external_identities_to_leave"`
	AllowDeletedIdentitiesDataRemoval types.Bool     `tfsdk:"allow_deleted_identities_data_removal"`
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
}
