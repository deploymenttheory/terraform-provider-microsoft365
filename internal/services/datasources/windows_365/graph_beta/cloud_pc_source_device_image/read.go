package graphBetaCloudPcSourceDeviceImage

import (
	"context"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *CloudPcSourceDeviceImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPcSourceDeviceImageDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...) // Get config into object
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, "Starting Read method for datasource: "+DataSourceName+" with filter_type: "+filterType)

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	filterValue := object.FilterValue.ValueString()
	var filteredItems []CloudPcSourceDeviceImageItemModel

	result, err := d.client.
		DeviceManagement().
		VirtualEndpoint().
		DeviceImages().
		GetSourceImages().
		GetAsGetSourceImagesGetResponse(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	if result.GetValue() != nil {
		for _, img := range result.GetValue() {
			item := CloudPcSourceDeviceImageItemModel{
				ID:                      convert.GraphToFrameworkString(img.GetId()),
				ResourceId:              convert.GraphToFrameworkString(img.GetResourceId()),
				DisplayName:             convert.GraphToFrameworkString(img.GetDisplayName()),
				SubscriptionId:          convert.GraphToFrameworkString(img.GetSubscriptionId()),
				SubscriptionDisplayName: convert.GraphToFrameworkString(img.GetSubscriptionDisplayName()),
			}

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, item)
			case "id":
				if img.GetId() != nil && *img.GetId() == filterValue {
					filteredItems = append(filteredItems, item)
				}
			case "display_name":
				if img.GetDisplayName() != nil && strings.Contains(strings.ToLower(*img.GetDisplayName()), strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
				}
			}
		}
	}
	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...) // Set state
	if resp.Diagnostics.HasError() {
		return
	}

}
