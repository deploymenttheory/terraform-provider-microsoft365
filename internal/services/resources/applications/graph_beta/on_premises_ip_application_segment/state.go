package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"
	"strings"

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
		data.Protocol = convert.GraphToFrameworkStringSet(ctx, terraformProtocols(*remoteResource.protocol))
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func terraformDestinationType(destinationType string) string {
	// Preserve Terraform's public schema and Microsoft Learn's documented value
	// ("ipAddress") while accepting the beta API's observed response value
	// ("ip") for single IP address segments.
	// https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
	if destinationType == "ip" {
		return "ipAddress"
	}

	return destinationType
}

func terraformProtocols(protocol string) []string {
	parts := strings.Split(protocol, ",")
	protocols := make([]string, 0, len(parts))
	for _, part := range parts {
		if protocol := strings.TrimSpace(part); protocol != "" {
			protocols = append(protocols, protocol)
		}
	}

	return protocols
}
