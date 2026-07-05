// REF: https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OnPremisesConnectorGroupResourceModel represents an Application Proxy connector group.
type OnPremisesConnectorGroupResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	Name               types.String   `tfsdk:"name"`
	Region             types.String   `tfsdk:"region"`
	ConnectorGroupType types.String   `tfsdk:"connector_group_type"`
	IsDefault          types.Bool     `tfsdk:"is_default"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
