package graphVersionResourceTemplate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state to the Terraform state
func mapRemoteStateToTerraform(ctx context.Context, data *ResourceTemplateResourceModel, remoteResource graphmodels.DeviceAndAppManagementAssignmentFilterable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = types.StringValue(convert.GraphToFrameworkString
(remoteResource.GetId()))
	// add more fields here as needed. use the helpers from the state package as needed.

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})
}
