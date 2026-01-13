// read.go
package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Windows Update Ring data source.
func (d *WindowsUpdateRingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsUpdateRingDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s with filter_type: %s", DataSourceName, filterType))

	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []WindowsUpdateRingModel
	filterValue := object.FilterValue.ValueString()

	// For ID filter, we can make a direct API call
	if filterType == "id" {

		respItem, err := d.client.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		// Verify this is a Windows Update Ring
		if isWindowsUpdateRing(respItem) {
			filteredItems = append(filteredItems, MapRemoteStateToDataSource(respItem))
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Windows Update Ring",
				fmt.Sprintf("The device configuration with ID %s is not a Windows Update Ring", filterValue),
			)
			return
		}
	} else {
		// For all other filters, we need to get all device configurations and filter for Windows Update Rings

		respList, err := d.client.
			DeviceManagement().
			DeviceConfigurations().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, item := range respList.GetValue() {
			// Filter for Windows Update Rings only
			if !isWindowsUpdateRing(item) {
				continue
			}

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))

			case "display_name":
				if item.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*item.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
				}
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}

// isWindowsUpdateRing checks if a device configuration is a Windows Update Ring
func isWindowsUpdateRing(config graphmodels.DeviceConfigurationable) bool {
	odataType := config.GetOdataType()
	return odataType != nil && *odataType == "#microsoft.graph.windowsUpdateForBusinessConfiguration"
}
