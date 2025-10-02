// REF: https://learn.microsoft.com/en-us/graph/api/resources/authenticationcontextclassreference?view=graph-rest-beta
package graphBetaAuthenticationContext

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthenticationContextResourceModel represents the authentication context resource model
type AuthenticationContextResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	DisplayName types.String   `tfsdk:"display_name"`
	Description types.String   `tfsdk:"description"`
	IsAvailable types.Bool     `tfsdk:"is_available"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}
