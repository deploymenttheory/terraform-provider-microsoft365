package graphBetaApplicationsIpApplicationSegment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *IpApplicationSegmentResourceModel, remoteResource graphmodels.IpApplicationSegmentable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DestinationHost = convert.GraphToFrameworkString(remoteResource.GetDestinationHost())
	data.DestinationType = convert.GraphToFrameworkString(helpers.StringPtr(remoteResource.GetDestinationType().String()))
	data.Ports = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetPorts())
	data.Protocol = convert.GraphToFrameworkString(helpers.StringPtr(remoteResource.GetProtocol().String()))

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
