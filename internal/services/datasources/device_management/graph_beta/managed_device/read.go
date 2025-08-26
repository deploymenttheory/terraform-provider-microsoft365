// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-windowsmanageddevice-list?view=graph-rest-beta
package graphBetaManagedDevice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *ManagedDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ManagedDeviceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s with filter_type: %s", DataSourceName, filterType))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []ManagedDeviceDeviceDataItemModel
	filterValue := object.FilterValue.ValueString()

	if filterType == "id" {
		// Fetch a single managed device by ID
		itemResp, err := d.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		item := MapRemoteStateToDataSource(itemResp)
		filteredItems = append(filteredItems, item)
	} else {
		// Fetch all managed devices
		respList, err := d.client.
			DeviceManagement().
			ManagedDevices().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, managedDevice := range respList.GetValue() {
			item := MapRemoteStateToDataSource(managedDevice)
			switch filterType {
			case "all":
				filteredItems = append(filteredItems, item)
			case "device_name":
				if managedDevice.GetDeviceName() != nil && strings.Contains(
					strings.ToLower(*managedDevice.GetDeviceName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
				}
			case "serial_number":
				if managedDevice.GetSerialNumber() != nil && strings.Contains(
					strings.ToLower(*managedDevice.GetSerialNumber()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
				}
			case "user_id":
				if managedDevice.GetUserId() != nil && strings.Contains(
					strings.ToLower(*managedDevice.GetUserId()),
					strings.ToLower(filterValue)) {
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
	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}
