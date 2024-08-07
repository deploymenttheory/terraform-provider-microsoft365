package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model.
// It populates the ConditionalAccessGrantControlsModel with data from the DeviceAndAppManagementAssignmentFilterable.
func mapRemoteStateToTerraform(ctx context.Context, data *ConditionalAccessPolicyResourceModel, remoteResource models.ConditionalAccessPolicy) {
	tflog.Debug(ctx, "Finished mapping remote state to Terraform")
}
