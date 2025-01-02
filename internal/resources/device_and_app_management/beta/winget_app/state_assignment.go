package graphBetaWinGetApp

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps the remote mobile app assignments to the Terraform state
func MapRemoteAssignmentStateToTerraform(ctx context.Context, assignment []sharedmodels.MobileAppAssignmentResourceModel, remoteAssignmentsResponse graphmodels.MobileAppAssignmentCollectionResponseable) {
	if remoteAssignmentsResponse == nil || remoteAssignmentsResponse.GetValue() == nil {
		tflog.Debug(ctx, "Remote assignments response is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map assignments from remote resource to Terraform state", map[string]interface{}{
		"assignmentsCount": len(remoteAssignmentsResponse.GetValue()),
	})

	return
}
