package graphBetaApplicationsIpApplicationSegment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *IpApplicationSegmentResourceModel, remoteResource json.RawMessage) {
	var ipApplicationData struct {
		ID              string   `json:"id"`
		DestinationHost string   `json:"destinationHost"`
		DestinationType string   `json:"destinationType"`
		Ports           []string `json:"ports"`
		Protocol        string   `json:"protocol"`
	}

	if err := json.Unmarshal(remoteResource, &ipApplicationData); err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to unmarshal remote resource: %v", err))
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": ipApplicationData.ID,
	})

	data.ID = convert.GraphToFrameworkString(&ipApplicationData.ID)
	data.DestinationHost = convert.GraphToFrameworkString(&ipApplicationData.DestinationHost)
	data.DestinationType = convert.GraphToFrameworkString(&ipApplicationData.DestinationType)
	data.Ports = convert.GraphToFrameworkStringSet(ctx, ipApplicationData.Ports)
	data.Protocol = convert.GraphToFrameworkString(&ipApplicationData.Protocol)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
