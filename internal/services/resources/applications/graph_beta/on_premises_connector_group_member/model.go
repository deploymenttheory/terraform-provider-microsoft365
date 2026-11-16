// REF: https://learn.microsoft.com/en-us/graph/api/resources/connector?view=graph-rest-beta
package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OnPremisesConnectorGroupMemberResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	ConnectorID        types.String   `tfsdk:"connector_id"`
	ConnectorGroupID   types.String   `tfsdk:"connector_group_id"`
	ConnectorGroupName types.String   `tfsdk:"connector_group_name"`
	MachineName        types.String   `tfsdk:"machine_name"`
	ExternalIP         types.String   `tfsdk:"external_ip"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
