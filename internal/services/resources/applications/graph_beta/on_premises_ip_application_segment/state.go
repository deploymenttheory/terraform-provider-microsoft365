package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(ctx context.Context, data *OnPremisesIpApplicationSegmentResourceModel, remoteResource *ipApplicationSegmentResponse) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.ID).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.ID)
	data.DestinationHost = convert.GraphToFrameworkString(remoteResource.DestinationHost)
	if remoteResource.DestinationType != nil {
		// Preserve Terraform's public schema and Microsoft Learn's documented value
		// ("ipAddress") while accepting the beta API's observed response value
		// ("ip") for single IP address segments.
		// https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
		destinationType := *remoteResource.DestinationType
		if destinationType == "ip" {
			destinationType = "ipAddress"
		}
		data.DestinationType = types.StringValue(destinationType)
	}
	data.Ports = convert.GraphToFrameworkStringSet(ctx, remoteResource.Ports)
	if remoteResource.Protocol != nil {
		// Graph returns the protocols as one comma-separated string, e.g. "tcp,udp".
		var protocols []string
		for part := range strings.SplitSeq(*remoteResource.Protocol, ",") {
			if protocol := strings.TrimSpace(part); protocol != "" {
				protocols = append(protocols, protocol)
			}
		}
		data.Protocol = convert.GraphToFrameworkStringSet(ctx, protocols)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
