// read.go
package graphBetaWindowsQualityUpdateExpeditePolicy

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

// Read handles the Read operation for Windows Quality Update Expedite Policy data source.
func (d *WindowsQualityUpdateExpeditePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsQualityUpdateExpeditePolicyDataSourceModel

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

	var filteredItems []WindowsQualityUpdateExpeditePolicyModel
	filterValue := object.FilterValue.ValueString()

	// For ID filter, we can make a direct API call
	if filterType == "id" {

		respItem, err := d.client.
			DeviceManagement().
			WindowsQualityUpdateProfiles().
			ByWindowsQualityUpdateProfileId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		filteredItems = append(filteredItems, MapRemoteStateToDataSource(respItem))
	} else {
		// For all other filters, we need to get all policies and filter locally

		respList, err := d.client.
			DeviceManagement().
			WindowsQualityUpdateProfiles().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		// Iterate through all profiles and get their expedite settings
		for _, profile := range respList.GetValue() {
			if profile.GetId() == nil {
				continue
			}

			expediteSettings, err := d.client.
				DeviceManagement().
				WindowsQualityUpdateProfiles().
				ByWindowsQualityUpdateProfileId(*profile.GetId()).
				Get(ctx, nil)

			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Error getting expedite settings for profile %s: %s", *profile.GetId(), err))
				continue
			}

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(expediteSettings))

			case "display_name":
				if expediteSettings.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*expediteSettings.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(expediteSettings))
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
