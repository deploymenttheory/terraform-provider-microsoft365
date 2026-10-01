// REF: https://learn.microsoft.com/en-us/graph/api/resources/authenticationflowspolicy?view=graph-rest-beta
package graphBetaAuthenticationFlowsPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthenticationFlowsPolicyResourceModel represents the tenant-wide authentication flows policy.
type AuthenticationFlowsPolicyResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	DisplayName       types.String   `tfsdk:"display_name"`
	Description       types.String   `tfsdk:"description"`
	SelfServiceSignUp types.Object   `tfsdk:"self_service_sign_up"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

// SelfServiceSignUpModel represents self-service sign-up configuration.
type SelfServiceSignUpModel struct {
	IsEnabled types.Bool `tfsdk:"is_enabled"`
}

func selfServiceSignUpAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{"is_enabled": types.BoolType}
}
