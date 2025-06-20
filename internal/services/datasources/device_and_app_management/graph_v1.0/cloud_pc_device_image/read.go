package graphCloudPcDeviceImage

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_v1.0/cloud_pc_device_image"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Read handles the Read operation for the CloudPcDeviceImageDataSource.
func (d *CloudPcDeviceImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state resource.CloudPcDeviceImageResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, resource.ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("Reading assignment filter with display name: %s", state.DisplayName.ValueString()))

	filters := d.client.
		DeviceManagement().
		VirtualEndpoint().
		DeviceImages()
	result, err := filters.Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Assignment Filter",
			fmt.Sprintf("Could not read assignment filter: %s", err),
		)
		return
	}

	var foundFilter models.CloudPcDeviceImageable
	for _, filter := range result.GetValue() {
		if *filter.GetDisplayName() == state.DisplayName.ValueString() {
			foundFilter = filter
			break
		}
	}

	if foundFilter == nil {
		resp.Diagnostics.AddError(
			"Error Reading Assignment Filter Datasource",
			fmt.Sprintf("No assignment filter found with display name: %s", state.DisplayName.ValueString()),
		)
		return
	}

	resource.MapRemoteStateToTerraform(ctx, &state, foundFilter)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
