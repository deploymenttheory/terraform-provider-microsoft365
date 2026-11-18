package graphBetaServicePrincipalClaimsPolicyJSON

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServicePrincipalClaimsPolicyJSONResourceModel struct {
	ID                 types.String    `tfsdk:"id"`
	ServicePrincipalID types.String    `tfsdk:"service_principal_id"`
	PolicyID           types.String    `tfsdk:"policy_id"`
	PolicyJSON         PolicyJSONValue `tfsdk:"policy_json"`
	Timeouts           timeouts.Value  `tfsdk:"timeouts"`
}
