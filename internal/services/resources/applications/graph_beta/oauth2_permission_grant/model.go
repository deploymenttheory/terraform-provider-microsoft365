package graphBetaOAuth2PermissionGrant

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OAuth2PermissionGrantResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	ClientID    types.String   `tfsdk:"client_id"`
	ResourceID  types.String   `tfsdk:"resource_id"`
	ConsentType types.String   `tfsdk:"consent_type"`
	PrincipalID types.String   `tfsdk:"principal_id"`
	Scopes      types.Set      `tfsdk:"scopes"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}
