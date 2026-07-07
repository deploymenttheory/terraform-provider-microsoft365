package graphBetaApplicationsOnPremisesConnectorGroupAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OnPremisesConnectorGroupAssignmentResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	ApplicationID      types.String   `tfsdk:"application_id"`
	ConnectorGroupID   types.String   `tfsdk:"connector_group_id"`
	ConnectorGroupName types.String   `tfsdk:"connector_group_name"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
