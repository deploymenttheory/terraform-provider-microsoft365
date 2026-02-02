// REF: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
package graphBetaApplicationIdentifierUri

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicationIdentifierUriResourceModel struct {
	Id            types.String   `tfsdk:"id"`
	ApplicationID types.String   `tfsdk:"application_id"`
	IdentifierUri types.String   `tfsdk:"identifier_uri"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}
