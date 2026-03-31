package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchOperationalInsightsConnectionResourceModel struct {
	ID                     types.String   `tfsdk:"id"`
	AzureResourceGroupName types.String   `tfsdk:"azure_resource_group_name"`
	AzureSubscriptionId    types.String   `tfsdk:"azure_subscription_id"`
	WorkspaceName          types.String   `tfsdk:"workspace_name"`
	State                  types.String   `tfsdk:"state"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}
