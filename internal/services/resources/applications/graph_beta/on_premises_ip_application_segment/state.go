package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *OnPremisesIpApplicationSegmentResourceModel, remoteResource *ipApplicationSegmentResponse) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.id).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.id)
	data.DestinationHost = convert.GraphToFrameworkString(remoteResource.destinationHost)
	if remoteResource.destinationType != nil {
		data.DestinationType = convert.GraphToFrameworkString(helpers.StringPtr(terraformDestinationType(*remoteResource.destinationType)))
	}
	data.Ports = convert.GraphToFrameworkStringSet(ctx, remoteResource.ports)
	if remoteResource.protocol != nil {
		data.Protocol = convert.GraphToFrameworkString(remoteResource.protocol)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func terraformDestinationType(destinationType string) string {
	// Preserve the Terraform schema and public docs value while accepting the
	// different literal returned by the beta API.
	if destinationType == "ip" {
		return "ipAddress"
	}

	return destinationType
}
