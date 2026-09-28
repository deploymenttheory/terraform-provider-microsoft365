// REF: https://learn.microsoft.com/en-us/graph/api/resources/b2bmanagementpolicy?view=graph-rest-beta
package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// B2bManagementPolicyResourceModel represents the schema for the B2B Management Policy resource.
// The documented description property is intentionally not modelled: Microsoft Graph rejects any
// write that includes it with 404 Request_ResourceNotFound and never returns it on reads.
type B2bManagementPolicyResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	DisplayName           types.String   `tfsdk:"display_name"`
	Definition            types.List     `tfsdk:"definition"`
	IsOrganizationDefault types.Bool     `tfsdk:"is_organization_default"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
