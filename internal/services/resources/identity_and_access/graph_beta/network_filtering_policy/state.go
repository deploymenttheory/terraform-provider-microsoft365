package graphBetaNetworkFilteringPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// MapRemoteStateToTerraform maps the base properties of a NetworkFilteringPolicyResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *NetworkFilteringPolicyResourceModel, remoteResource models.FilteringPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Action = convert.GraphToFrameworkEnum(remoteResource.GetAction())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
