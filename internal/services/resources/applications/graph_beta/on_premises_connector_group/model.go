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

// connectorGroupRequest is the connectorGroup create and update body. Only the writable name and region are sent.
type connectorGroupRequest struct {
	Name   *string `json:"name,omitempty"`
	Region *string `json:"region,omitempty"`
}

// connectorGroupResponse is the connectorGroup Graph returns. region and connectorGroupType are decoded as strings,
// not the SDK enums, so values the SDK does not know (such as region "japan") are kept. The live API returns
// isDefault=false for created connector groups and isDefault=true for the imported tenant default.
type connectorGroupResponse struct {
	ID                 *string `json:"id"`
	Name               *string `json:"name"`
	ConnectorGroupType *string `json:"connectorGroupType"`
	IsDefault          *bool   `json:"isDefault"`
	Region             *string `json:"region"`
}
