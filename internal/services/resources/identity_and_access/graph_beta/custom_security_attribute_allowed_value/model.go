// https://learn.microsoft.com/en-us/graph/api/resources/allowedvalue?view=graph-rest-beta

package graphBetaCustomSecurityAttributeAllowedValue

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomSecurityAttributeAllowedValueResourceModel struct {
	ID                                  types.String   `tfsdk:"id"`
	CustomSecurityAttributeDefinitionId types.String   `tfsdk:"custom_security_attribute_definition_id"`
	IsActive                            types.Bool     `tfsdk:"is_active"`
	Timeouts                            timeouts.Value `tfsdk:"timeouts"`
}
